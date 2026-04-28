#!/bin/bash
set -euo pipefail

PROJECT_NAME="bmc-manager"
MODULE_NAME="github.com/example/bmc-manager"
GROUP="example.fabrica.dev"
API_VERSION="v1"
API_DIR="apis/$GROUP/$API_VERSION"

rm -rf $PROJECT_NAME

fabrica init $PROJECT_NAME --module $MODULE_NAME --storage-type ent --db sqlite --events --events-bus memory --reconcile

cd $PROJECT_NAME

fabrica add resource BMCCredential

cat << 'EOF' > $API_DIR/bmccredential_types.go
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
	Phase      string `json:"phase" validate:"oneof=Pending Discovering Updating Ready Error"`
	Message    string `json:"message,omitempty"`
	AccountURI string `json:"accountUri,omitempty"`
}

func (r *BMCCredential) GetKind() string { return "BMCCredential" }
func (r *BMCCredential) GetName() string { return r.Metadata.Name }
func (r *BMCCredential) GetUID() string { return r.Metadata.UID }
func (r *BMCCredential) IsHub() {}
EOF

fabrica generate

go mod tidy