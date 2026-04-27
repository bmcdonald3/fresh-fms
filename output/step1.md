The declarative workflow in Fabrica will execute in five distinct stages, bridging the configuration intent with the concrete Redfish operations.

**Stage 1: Desired State Declaration**
The user submits a declarative JSON payload to the Fabrica REST API. This payload defines the target BMC and the intended credential configuration.
Payload parameters:
* `bmc_address`: 172.24.0.3
* `current_username`: root
* `current_password`: initial0
* `target_username`: root
* `desired_password`: initial1

**Stage 2: Event Publication**
Upon accepting the JSON payload, Fabrica writes the "Desired State" to its datastore and publishes an asynchronous CloudEvent (e.g., `io.fabrica.bmc.credentials.updated`) to the message bus.

**Stage 3: Reconciliation - Discovery Phase**
A background reconciliation controller receives the CloudEvent and initiates a session with the BMC. Because Redfish account URIs vary by manufacturer and account creation order, the controller must query the AccountService to map the `target_username` to its specific Account ID.
Command executed by the controller:
`curl -s -k -u root:initial0 -X GET https://172.24.0.3/redfish/v1/AccountService/Accounts`

The controller parses the returned JSON payload to locate the member object where `"UserName"` equals `"root"`. In this scenario, it parses the data to extract the specific target URI: `/redfish/v1/AccountService/Accounts/3`.

**Stage 4: Reconciliation - Execution Phase**
Using the discovered URI, the controller executes the HTTP PATCH request to align the BMC's "Current State" with the Fabrica "Desired State" by applying the `desired_password`.
Command executed by the controller:
`curl -i -k -u root:initial0 -X PATCH https://172.24.0.3/redfish/v1/AccountService/Accounts/3 -H "Content-Type: application/json" -d '{"Password":"initial1"}'`

The controller evaluates the HTTP response header. It expects a `200 OK` or `204 No Content` status code to verify that the password complexity policies were met and the change was committed by the BMC. 
Expected system output:
`HTTP/1.1 200 OK`

**Stage 5: Observed State Update**
Following a successful `200 OK` response, the reconciliation controller registers the operation as complete. It updates the BMC's status sub-resource in the Fabrica datastore, marking the "Observed State" as synchronized with the "Desired State". All subsequent reconciliation loops or monitoring checks initiated by the Fabrica service will transition to authenticating with the new `initial1` credentials.