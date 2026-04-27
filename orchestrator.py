import os
import subprocess
from openai import OpenAI

client = OpenAI()

# ---------------------------------------------------------------------------
# Prompt Templates
# ---------------------------------------------------------------------------

STEP_1_PROMPT = """
I am designing a new service. I will provide the manual operations (e.g., API calls, CLI commands, or scripts) and their successful outputs required to interact with an external system. Please analyze these operations and output a precise list of the required inputs (variables the user must provide) and the expected outputs so that they can be used to determine the expected Fabrica workflow.
"""

STEP_2_PROMPT = """
Using the inputs and outputs defined below, map this process into a declarative Fabrica workflow. Differentiate the workflow into two distinct phases:
Phase A: What desired state the user submits (based on the inputs).
Phase B: What the background reconciliation loop must do asynchronously to execute the commands and track the outputs.

Input/Output mapping:
"""

STEP_3_PROMPT = """
Using the workflow mapping provided, identify the required API resources. For each resource, generate the Go structs defining the schema. You must follow the Kubernetes-style resource pattern by splitting the data into two components:
1. `Spec`: The desired state provided by the user (containing the required inputs for the concrete operations).
2. `Status`: The observed state managed by the system in the background.
Output ONLY the raw Go code block.

Workflow Mapping:
"""

STEP_4_PROMPT = """
Generate a bash bootstrap script to initialize the Fabrica project. The script must use `fabrica init --events --reconcile`, followed by `fabrica add resource` for the identified resources. Write the custom Go structs into the `apis/` directory overriding the defaults, then run `fabrica generate`. 
Output ONLY the raw bash script.

Go Structs:
"""

STEP_5_PROMPT = """
Generate the Go code for the reconcilers. The logic must execute the concrete operations identified initially. Ensure idempotency checks, progressive status updates, and explicit `r.Client.Update` calls. 
Output ONLY the raw Go code block.

Go Structs and Workflow:
"""

STEP_6_PROMPT = """
Write a validation bash script using `curl` to simulate the user workflow.
Constraint 1: Start a lightweight local mock server (e.g., Python HTTP server) to simulate the external system.
Constraint 2: Create Fabrica resources via POST requests. Use `set +e` and `set -e` to prevent aborts on failure.
Constraint 3: Poll the Fabrica GET endpoints to check if `status.phase` transitions to the expected terminal state.
Output ONLY the raw bash script.

Reconciler Code:
"""

# ---------------------------------------------------------------------------
# Core Execution Logic
# ---------------------------------------------------------------------------

def call_llm(system_prompt, user_content):
    response = client.chat.completions.create(
        model="gpt-4o", # Replace with your preferred model
        messages=[
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": user_content}
        ],
        temperature=0.2
    )
    return response.choices[0].message.content.strip("` \n").removeprefix("bash\n").removeprefix("go\n")

def execute_script(script_content, filename):
    with open(filename, "w") as f:
        f.write(script_content)
    
    os.chmod(filename, 0o755)
    
    result = subprocess.run(
        [f"./{filename}"], 
        capture_output=True, 
        text=True
    )
    return result.returncode, result.stdout, result.stderr

def run_with_auto_correction(llm_prompt, input_data, script_name, max_retries=3):
    current_input = input_data
    
    for attempt in range(max_retries):
        print(f"Generating {script_name} (Attempt {attempt + 1})...")
        script_content = call_llm(llm_prompt, current_input)
        
        print(f"Executing {script_name}...")
        returncode, stdout, stderr = execute_script(script_content, script_name)
        
        if returncode == 0:
            print(f"Execution successful for {script_name}.")
            return script_content
            
        print(f"Execution failed. Exit code: {returncode}. Initiating self-correction.")
        
        # Append the error log to the prompt for the next iteration
        current_input = f"""
        Previous output generated an error. 
        Original Request Data: {input_data}
        
        Execution Error Log:
        {stderr}
        {stdout}
        
        Please correct the script to resolve these errors.
        """
        
    raise Exception(f"Failed to execute {script_name} after {max_retries} attempts.")

# ---------------------------------------------------------------------------
# Orchestration Pipeline
# ---------------------------------------------------------------------------

def main():
    print("Enter the concrete operations (Phase A Input). Type 'EOF' on a new line when finished:")
    lines = []
    while True:
        line = input()
        if line.strip() == "EOF":
            break
        lines.append(line)
    
    raw_commands = "\n".join(lines)
    state = {}

    print("\n[Step 1] Analyzing concrete operations...")
    state["step1_io"] = call_llm(STEP_1_PROMPT, raw_commands)
    
    print("[Step 2] Mapping to Fabrica workflow (Hidden CoT)...")
    state["step2_workflow"] = call_llm(STEP_2_PROMPT, state["step1_io"])

    print("[Step 3] Generating Go Structs...")
    state["step3_structs"] = call_llm(STEP_3_PROMPT, state["step2_workflow"])

    print("[Step 4] Bootstrapping project...")
    state["step4_script"] = run_with_auto_correction(STEP_4_PROMPT, state["step3_structs"], "bootstrap.sh")

    print("[Step 5] Generating Reconciler Logic...")
    step5_input = f"Structs:\n{state['step3_structs']}\n\nWorkflow:\n{state['step2_workflow']}"
    state["step5_reconciler"] = call_llm(STEP_5_PROMPT, step5_input)
    # Note: To auto-correct Step 5, the script would need to write the file into the generated Fabrica directory and run `go build`

    print("[Step 6] Generating and executing E2E Tests...")
    state["step6_test"] = run_with_auto_correction(STEP_6_PROMPT, state["step5_reconciler"], "e2e_test.sh")

    print("\nPipeline complete. Review the generated scripts in your directory.")

if __name__ == "__main__":
    main()
