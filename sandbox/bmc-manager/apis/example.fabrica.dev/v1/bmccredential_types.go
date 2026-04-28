package v1

import "github.com/openchami/fabrica/pkg/fabrica"

type BMCCredential struct {
	APIVersion string              `json:"apiVersion" validate:"required"`
	Kind       string              `json:"kind" validate:"required"`
	Metadata   fabrica.Metadata    `json:"metadata"`
	Spec       BMCCredentialSpec   `json:"spec"`
	Status     BMCCredentialStatus `json:"status,omitempty"`
}

type BMCCredentialSpec struct {
	BMCAddress            string `json:"bmcAddress" validate:"required,ip"`
	AuthorizationUsername string `json:"authorizationUsername" validate:"required"`
	AuthorizationPassword string `json:"authorizationPassword" validate:"required"`
	TargetUsername        string `json:"targetUsername" validate:"required"`
	DesiredPassword       string `json:"desiredPassword" validate:"required"`
}

type BMCCredentialStatus struct {
	Phase      string `json:"phase,omitempty"`
	Message    string `json:"message,omitempty"`
	AccountURI string `json:"accountUri,omitempty"`
}

func (r *BMCCredential) GetKind() string { return "BMCCredential" }
func (r *BMCCredential) GetName() string { return r.Metadata.Name }
func (r *BMCCredential) GetUID() string { return r.Metadata.UID }
func (r *BMCCredential) IsHub() {}
