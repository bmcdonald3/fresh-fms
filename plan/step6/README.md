## Generate and execute E2E tests

### Step outcome

A test script validating the full declarative workflow against a mock implementation of the external system.

### What to give

The resource models (Step 3) and the reconciler logic (Step 5).

### Prompt

Write a validation bash script using `curl` to simulate the user workflow. 
Constraint 1: The script must start the Fabrica server on port 8085 and terminate it upon completion. It can be ran from the root of the project with `go run ./cmd/server serve --port=8085 --database-url="file:data.db?cache=shared&_fk=1"`.
Constraint 1: The script must stand up a lightweight local mock server (e.g., a Python HTTP server, if applicable) on a designated port to simulate the external hardware/API interactions defined in our concrete operations. Here is an example of a working post request for a different service:
curl -s -X POST http://localhost:8085/bmccredentials \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion": "fabrica.dev/v1",
    "kind": "BMCCredential",
    "metadata": {
      "name": "real-bmc-task"
    },
    "spec": {
      "bmcAddress": "172.24.0.3",
      "authorizationUsername": "root",
      "authorizationPassword": "initial0",
      "targetUsername": "root",
      "desiredPassword": "initial1"
    }
  }'
Constraint 2: The script must create the necessary Fabrica resources by sending POST requests. You must wrap the resource creation `curl` commands with `set +e` and `set -e` to prevent the script from aborting on failure. Capture the HTTP status code and print the full JSON response if the creation fails.
Constraint 3: The script must poll the Fabrica GET endpoints to check if the `status.phase` field transitions to the expected terminal state, proving that the reconciliation loop successfully interacted with the mock server.

### Context

Fabrica generates standard CRUD endpoints for all defined resources. 
- Create: POST /<resource-name-plural> (Requires "apiVersion", "kind", "metadata.name", and "spec")
- Read: GET /<resource-name-plural>/<uid>
- List: GET /<resource-name-plural>
