## Define Fabrica resources

### Step outcome

Resource models split into spec and status fields based on the concrete operations.

### What to give

The workflow mapping from Step 2.

### Prompt

Using the workflow mapping, identify the required API resources. For each resource, generate the Go structs defining the schema. You must follow the Kubernetes-style resource pattern by splitting the data into two components:
1. `Spec`: The desired state provided by the user (containing the required inputs for the concrete operations).
2. `Status`: The observed state managed by the system in the background (containing the outputs, status tracking, and error messages).

### Context

Fabrica resources strictly separate data into 'Spec' and 'Status'.
1. Spec: The desired state provided by the user. Must use 'validate' struct tags for input validation.
2. Status: The observed state managed exclusively by the system's reconciliation loop.

Example Resource Implementation:

```go
type UserSpec struct {
    Email string `json:"email" validate:"required,email"`
    Role  string `json:"role" validate:"oneof=admin user guest"`
    ParentTeamUID string `json:"parentTeamUid,omitempty"`
}

type UserStatus struct {
    Phase      string     `json:"phase,omitempty"`
    Message    string     `json:"message,omitempty"`
    LastLogin  *time.Time `json:"lastLogin,omitempty"`
}
```
