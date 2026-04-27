Phase A: Desired State Declaration
The user interacts exclusively with the Fabrica REST API by submitting a JSON payload that defines the intended configuration. This submission updates the datastore and initiates the CloudEvent. 

The declarative payload contains the required inputs:
* `bmc_address`: 172.24.0.3
* `authorization_username`: root
* `authorization_password`: initial0
* `target_username`: root
* `desired_password`: initial1

Phase B: Asynchronous Reconciliation Loop
The background worker receives the generated CloudEvent and executes a sequence of operations to align the BMC's actual state with the declared desired state.

1.  **Current State Discovery:** The controller authenticates against the BMC at `172.24.0.3` using the provided `authorization_username` and `authorization_password`. It issues an HTTP GET request to the `/redfish/v1/AccountService/Accounts` endpoint.
2.  **Payload Parsing:** The worker processes the returned JSON payload to locate the specific member object where the `UserName` property matches the `target_username` ("root"). It extracts the corresponding `@odata.id` URI (e.g., `/redfish/v1/AccountService/Accounts/3`).
3.  **State Alignment Execution:** The worker issues an HTTP PATCH request to the discovered URI, injecting the `desired_password` ("initial1") into the request body to modify the account.
4.  **Observed State Synchronization:** The controller evaluates the HTTP response. Upon receiving a 200 OK or 204 No Content status code, it confirms the state change is complete. The worker then updates the `status` sub-resource of the Fabrica object in the datastore to register the observed state as synchronized. All subsequent automated interactions with the BMC will utilize the updated credential payload.