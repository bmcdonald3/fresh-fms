## Bootstrap service

### Step outcome

A script to initialize the Fabrica project with all resources.

### What to give

Defined Go structs from Step 3.

### Prompt

Generate a bash bootstrap script to initialize the Fabrica project. The script must use the `fabrica init` command to create the project with event-driven reconciliation enabled. It must then use `fabrica add resource` for each of our identified resources. Following the resource creation, the script must write our custom Go structs into the generated files in the `apis/` directory, overriding the default types. Finally, the script must run `fabrica generate`.

### Context

Sample bootstrap file:

```bash
#!/bin/bash
set -euo pipefail

PROJECT_NAME="fms"
MODULE_NAME="github.com/bmcdonald3/fms"
GROUP="example.fabrica.dev"
API_VERSION="v1"
API_DIR="apis/$GROUP/$API_VERSION"

rm -rf $PROJECT_NAME

fabrica init $PROJECT_NAME --module $MODULE_NAME --storage-type ent --db sqlite --events --events-bus memory --reconcile

cd $PROJECT_NAME

for res in DeviceProfile FirmwareProfile; do fabrica add resource $res; done

cat << 'EOF' > $API_DIR/deviceprofile_types.go
package v1
import "[github.com/openchami/fabrica/pkg/fabrica](https://github.com/openchami/fabrica/pkg/fabrica)"
type DeviceProfile struct {
	APIVersion string `json:"apiVersion" validate:"required"`
	Kind       string `json:"kind" validate:"required"`
	Metadata   fabrica.Metadata `json:"metadata"`
	Spec       DeviceProfileSpec `json:"spec"`
	Status     DeviceProfileStatus `json:"status,omitempty"`
}
type DeviceProfileSpec struct {
	Manufacturer string `json:"manufacturer" validate:"required"`
}
type DeviceProfileStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
}
func (r *DeviceProfile) GetKind() string { return "DeviceProfile" }
func (r *DeviceProfile) GetName() string { return r.Metadata.Name }
func (r *DeviceProfile) GetUID() string { return r.Metadata.UID }
func (r *DeviceProfile) IsHub() {}
EOF

fabrica generate
go mod tidy
```