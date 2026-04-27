# Fabrica Auto-Generator Pipeline

This repository contains an autonomous AI-driven pipeline designed to scaffold, implement, and test event-driven, Kubernetes-style declarative APIs using the Fabrica framework. 

Instead of manually defining schemas and writing boilerplate controller logic, this pipeline allows a user to input the concrete manual operations (e.g., standard `curl` commands or scripts) required to interact with an external system. The orchestrator script utilizes a Large Language Model (LLM) to deduce the required data models, generate the Go code, and validate the logic against a mock server.

## Architecture and Workflow

The generation pipeline executes in six sequential steps. The orchestration script handles the state and context window management between these steps to ensure architectural consistency.

1. **Define concrete operations:** The user provides the raw manual commands and their successful outputs.
2. **Map to Fabrica workflow:** The LLM internally maps the inputs and outputs into a declarative "Desired State" (Phase A) and an asynchronous "Observed State" execution loop (Phase B).
3. **Define Fabrica resources:** The LLM generates the Go structs for the `Spec` (user inputs) and `Status` (system outputs and state tracking) schemas.
4. **Bootstrap service:** A generated bash script initializes the Fabrica project, adds the resources, and runs the code generators.
5. **Implement Reconciler Logic:** The LLM writes the Go business logic into the generated `pkg/reconcilers/` worker files, including idempotency checks and status updates.
6. **Generate and execute E2E tests:** The orchestrator generates a test script that stands up a local mock server and validates the full declarative workflow.

## Error Correction

The pipeline includes an autonomous self-correction loop. If Step 4 (compilation and project scaffolding) or Step 6 (end-to-end testing) fails and returns a non-zero exit code, the orchestrator captures the standard error output. It automatically feeds this error log back into the LLM context and attempts to rewrite the script or Go code, retrying the execution up to three times.

## Prerequisites

To run the pipeline locally, your environment must have the following installed and accessible in your system PATH:

* **Go:** Version 1.21 or higher.
* **Fabrica CLI:** Installed and configured.
* **Python:** Version 3.8 or higher.
* **Python OpenAI Package:** Installed via `pip install openai`.

## Usage

1. Clone this repository and navigate to the root directory.
2. Set your OpenAI API key in your terminal session:
   ```bash
   export OPENAI_API_KEY="your-api-key-here"
   ```
3. Execute the orchestrator script:
   ```bash
   python3 orchestrator.py
   ```
4. When prompted, input the exact manual commands and expected outputs you want the service to automate. 
5. Type `EOF` on a new line and press Enter to start the pipeline.

### Example Input Payload

```text
Initial target IP: 192.168.1.100
Current credentials: admin / password

Step 1: Retrieve account details
curl -s -k -u admin:password -X GET https://192.168.1.100/redfish/v1/AccountService/Accounts

Successful Output:
{
  "Members": [
    {
      "@odata.id": "/redfish/v1/AccountService/Accounts/1"
    }
  ]
}

Step 2: Update the password using the discovered URI
curl -i -k -u admin:password -X PATCH https://192.168.1.100/redfish/v1/AccountService/Accounts/1 -H "Content-Type: application/json" -d '{"Password":"new_secure_password"}'

Successful Output:
HTTP/1.1 200 OK
```