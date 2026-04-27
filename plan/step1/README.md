## Define concrete operations

### Step outcome

A precise list of the exact commands, API calls, or scripts required to interact with the external system.

### What to give

The manual steps, CLI commands, or HTTP requests a user would execute to achieve the desired outcome without the Fabrica service.

### Prompt

I am designing a new service. I will provide the manual operations (e.g., API calls, CLI commands, or scripts) required to interact with an external system. Please analyze these operations and output a precise list of the required inputs (variables the user must provide) and the expected outputs (status codes, return payloads, or state changes the system will produce). Ask me for clarification if any parameters or return values are ambiguous.

### Context

This service will be built using Fabrica, a framework for generating Kubernetes-style declarative APIs. 
Architecture Overview:
1. Declarative Design: Users do not trigger imperative actions. Instead, they declare a "Desired State" by creating or updating a JSON resource via a REST API.
2. Asynchronous Processing: Creating or modifying a resource publishes a CloudEvent (e.g., io.fabrica.resource.created).
3. Reconciliation Controller: A background worker receives the event, compares the "Current State" of the system to the "Desired State", and executes business logic to align them.

Before designing the declarative API, we must understand the concrete execution steps. The inputs identified here will eventually form the user's "Desired State," and the outputs will be used to track the "Observed State" in the system's background workers.
