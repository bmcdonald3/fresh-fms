## Generate and execute E2E tests

### Step outcome

A test script validating the full declarative workflow against a mock implementation of the external system.

### What to give

The resource models (Step 3) and the reconciler logic (Step 5).

### Prompt

Write a validation bash script using `curl` to simulate the user workflow. 
Constraint 1: The script must first stand up a lightweight local mock server (e.g., a Python HTTP server) on a designated port to simulate the external hardware/API interactions defined in our concrete operations.
Constraint 2: The script must create the necessary Fabrica resources by sending POST requests. You must wrap the resource creation `curl` commands with `set +e` and `set -e` to prevent the script from aborting on failure. Capture the HTTP status code and print the full JSON response if the creation fails.
Constraint 3: The script must poll the Fabrica GET endpoints to check if the `status.phase` field transitions to the expected terminal state, proving that the reconciliation loop successfully interacted with the mock server.

### Context

Fabrica generates standard CRUD endpoints for all defined resources. 
- Create: POST /<resource-name-plural> (Requires "apiVersion", "kind", "metadata.name", and "spec")
- Read: GET /<resource-name-plural>/<uid>
- List: GET /<resource-name-plural>

By testing against a local mock server, we validate the complete network and parsing logic of the reconciler without requiring actual hardware or waiting for network timeouts.
