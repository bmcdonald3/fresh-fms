## Map to Fabrica workflow

### Step outcome

A declarative workflow mapping that separates user intent from background execution.

### What to give

The defined inputs and outputs from Step 1.

### Prompt

Using the inputs and outputs we defined, map this process into a declarative Fabrica workflow. Differentiate the workflow into two distinct phases:
Phase A: What desired state the user submits (based on the inputs).
Phase B: What the background reconciliation loop must do asynchronously to execute the commands and track the outputs.

### Context

This service is built using Fabrica, a framework for generating Kubernetes-style declarative APIs. 
Architecture Overview:
1. Declarative Design: Users do not trigger imperative actions. Instead, they declare a "Desired State" by creating or updating a JSON resource via a REST API.
2. Asynchronous Processing: Creating or modifying a resource publishes a CloudEvent.
3. Reconciliation Controller: A background worker receives the event, compares the "Current State" of the system to the "Desired State", and executes the concrete operations to align them.