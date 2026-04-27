## Define concrete operations

### Step outcome

A precise list of the exact commands and their successful outputs required to interact with the external system.

### What to give

The manual steps, CLI commands, or HTTP requests a user would execute to achieve the desired outcome, along with the output of those commands completing successfully.

### Prompt

I am designing a new service. I will provide the manual operations (e.g., API calls, CLI commands, or scripts) and their successful outputs required to interact with an external system. Please analyze these operations and output a precise list of the required inputs (variables the user must provide) and the expected outputs so that they can be used to determine the expected Fabrica workflow. Ask me for clarification if any parameters or return values are ambiguous.

### Context

Before designing the declarative API, we must understand the concrete execution steps. The inputs identified here will eventually form the user's "Desired State," and the outputs will be used to track the "Observed State" in the system's background workers.