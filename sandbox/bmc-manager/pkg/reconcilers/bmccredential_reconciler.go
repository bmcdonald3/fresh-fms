package reconcilers

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	v1 "github.com/example/bmc-manager/apis/example.fabrica.dev/v1"
)

// reconcileBMCCredential executes the background operations for the Desired State.
func (r *BMCCredentialReconciler) reconcileBMCCredential(ctx context.Context, res *v1.BMCCredential) error {
	// 1. Idempotency Check
	if res.Status.Phase == "Ready" {
		return nil
	}

	// Initialize HTTP Client with insecure TLS for BMC connections
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr, Timeout: 15 * time.Second}

	// 2. Current State Discovery Phase
	if res.Status.AccountURI == "" {
		res.Status.Phase = "Discovering"
		if err := r.Client.Update(ctx, res); err != nil {
			return fmt.Errorf("failed to save Discovering phase: %w", err)
		}

		uri, err := discoverAccountURI(ctx, httpClient, res)
		if err != nil {
			return r.handleError(ctx, res, fmt.Errorf("discovery phase failed: %w", err))
		}
		
		res.Status.AccountURI = uri
		if err := r.Client.Update(ctx, res); err != nil {
			return fmt.Errorf("failed to save discovered AccountURI: %w", err)
		}
	}

	// 3. State Alignment Execution Phase
	res.Status.Phase = "Updating"
	if err := r.Client.Update(ctx, res); err != nil {
		return fmt.Errorf("failed to save Updating phase: %w", err)
	}

	err := updateBMCPassword(ctx, httpClient, res)
	if err != nil {
		return r.handleError(ctx, res, fmt.Errorf("execution phase failed: %w", err))
	}

	// 4. Observed State Synchronization
	res.Status.Phase = "Ready"
	res.Status.Message = "Credential update synchronized successfully"
	if err := r.Client.Update(ctx, res); err != nil {
		return fmt.Errorf("failed to save Ready phase: %w", err)
	}

	return nil
}

// handleError is a helper to persist error states before returning to the event queue
func (r *BMCCredentialReconciler) handleError(ctx context.Context, res *v1.BMCCredential, err error) error {
	res.Status.Phase = "Error"
	res.Status.Message = err.Error()
	_ = r.Client.Update(ctx, res) // Ignore update error to prioritize returning the execution error
	return err
}

// discoverAccountURI queries the Redfish AccountService to find the precise URI for the target username
func discoverAccountURI(ctx context.Context, client *http.Client, res *v1.BMCCredential) (string, error) {
	collectionURL := fmt.Sprintf("https://%s/redfish/v1/AccountService/Accounts", res.Spec.BMCAddress)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, collectionURL, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(res.Spec.AuthorizationUsername, res.Spec.AuthorizationPassword)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code requesting accounts collection: %d", resp.StatusCode)
	}

	var collection struct {
		Members []struct {
			ODataID string `json:"@odata.id"`
		} `json:"Members"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&collection); err != nil {
		return "", err
	}

	// Iterate through the returned accounts to find the specific username match
	for _, member := range collection.Members {
		memberURL := fmt.Sprintf("https://%s%s", res.Spec.BMCAddress, member.ODataID)
		mReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, memberURL, nil)
		mReq.SetBasicAuth(res.Spec.AuthorizationUsername, res.Spec.AuthorizationPassword)

		mResp, mErr := client.Do(mReq)
		if mErr != nil {
			continue
		}

		var account struct {
			UserName string `json:"UserName"`
		}
		json.NewDecoder(mResp.Body).Decode(&account)
		mResp.Body.Close()

		if account.UserName == res.Spec.TargetUsername {
			return member.ODataID, nil
		}
	}

	return "", fmt.Errorf("target account '%s' not found on BMC", res.Spec.TargetUsername)
}

// updateBMCPassword executes the HTTP PATCH to align the desired credentials
func updateBMCPassword(ctx context.Context, client *http.Client, res *v1.BMCCredential) error {
	fmt.Printf("DEBUG: Discovered URI for %s is %s\n", res.Spec.TargetUsername, res.Status.AccountURI)
	updateURL := fmt.Sprintf("https://%s%s", res.Spec.BMCAddress, res.Status.AccountURI)

	payload := map[string]string{
		"Password": res.Spec.DesiredPassword,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, updateURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.SetBasicAuth(res.Spec.AuthorizationUsername, res.Spec.AuthorizationPassword)
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("DEBUG: Sending PATCH to %s with body: %s\n", updateURL, string(body))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var bodyBytes []byte
    	bodyBytes, _ = io.ReadAll(resp.Body) 
    	return fmt.Errorf("rejected by BMC with HTTP %d: %s", resp.StatusCode, string(bodyBytes))
}
	}

	return nil
}