#!/bin/bash
set -euo pipefail

openssl req -new -x509 -keyout bmc.pem -out bmc.pem -days 365 -nodes -subj "/CN=localhost" 2>/dev/null

cat << 'EOF' > mock_bmc.py
import http.server
import ssl

class MockBMC(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/redfish/v1/AccountService/Accounts":
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'{"Members": [{"@odata.id": "/redfish/v1/AccountService/Accounts/3"}]}')
        elif self.path == "/redfish/v1/AccountService/Accounts/3":
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'{"UserName": "root"}')
        else:
            self.send_response(404)
            self.end_headers()
            
    def do_PATCH(self):
        if self.path == "/redfish/v1/AccountService/Accounts/3":
            self.send_response(200)
            self.end_headers()
        else:
            self.send_response(404)
            self.end_headers()

httpd = http.server.HTTPServer(('127.0.0.1', 443), MockBMC)
context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
context.load_cert_chain('bmc.pem')
httpd.socket = context.wrap_socket(httpd.socket, server_side=True)
httpd.serve_forever()
EOF

sudo python3 mock_bmc.py &
MOCK_PID=$!
sleep 2

PAYLOAD='{"apiVersion":"v1","kind":"BMCCredential","metadata":{"name":"test-bmc"},"spec":{"bmcAddress":"127.0.0.1","authorizationUsername":"root","authorizationPassword":"initial0","targetUsername":"root","desiredPassword":"initial1"}}'

set +e
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST http://127.0.0.1:8080/bmccredentials -H "Content-Type: application/json" -d "$PAYLOAD")
STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
set -e

if [ "$STATUS" -ne 201 ] && [ "$STATUS" -ne 200 ]; then
    echo "Resource creation failed. HTTP Status: $STATUS"
    echo "$BODY"
    sudo kill $MOCK_PID
    rm bmc.pem mock_bmc.py
    exit 1
fi

UID_VAL=$(echo "$BODY" | jq -r .metadata.uid)

for i in {1..30}; do
    PHASE=$(curl -s http://127.0.0.1:8080/bmccredentials/$UID_VAL | jq -r .status.phase)
    
    if [ "$PHASE" == "Ready" ]; then
        echo "Reconciliation successful. Phase transitioned to Ready."
        break
    fi
    
    if [ "$PHASE" == "Error" ]; then
        echo "Reconciliation failed."
        curl -s http://127.0.0.1:8080/bmccredentials/$UID_VAL | jq .status
        break
    fi
    
    sleep 2
done

sudo kill $MOCK_PID
rm bmc.pem mock_bmc.py