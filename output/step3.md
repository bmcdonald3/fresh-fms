The `BMCCredential` resource schema maps the Phase A inputs into the `Spec` and the Phase B execution outputs into the `Status`. 

```go
// BMCCredentialSpec represents the user's desired state (Phase A).
type BMCCredentialSpec struct {
	BMCAddress            string `json:"bmcAddress" validate:"required,ip"`
	AuthorizationUsername string `json:"authorizationUsername" validate:"required"`
	AuthorizationPassword string `json:"authorizationPassword" validate:"required"`
	TargetUsername        string `json:"targetUsername" validate:"required"`
	DesiredPassword       string `json:"desiredPassword" validate:"required"`
}

// BMCCredentialStatus represents the system's observed state (Phase B).
type BMCCredentialStatus struct {
	Phase      string `json:"phase" validate:"oneof=Pending Discovering Updating Ready Error"`
	Message    string `json:"message,omitempty"`
	AccountURI string `json:"accountUri,omitempty"`
}
```

### Schema Details

**Spec (Desired State)**
* `BMCAddress` requires a valid IP address to ensure the background loop attempts to connect to a properly formatted endpoint.
* The `Authorization` and `Target` credential fields are marked as required. Fabrica will reject the user's API payload immediately if any of these variables are omitted, preventing the creation of an un-actionable event.

**Status (Observed State)**
* `Phase`: Tracks the controller's progress through the asynchronous operations. It transitions from `Pending` (event received), to `Discovering` (GET request), to `Updating` (PATCH request), and finally lands on `Ready` or `Error`.
* `AccountURI`: Persists the discovered `@odata.id` (e.g., `/redfish/v1/AccountService/Accounts/3`) for the target user. Storing this in the status prevents the controller from needing to execute the GET request discovery phase again during future reconciliation checks.
* `Message`: Captures system outputs, specifically the HTTP response codes (`200 OK`, `401 Unauthorized`, `400 Bad Request`) or timeouts returned during Phase B execution.
