#!/bin/bash

# Check if we can connect to the counting service.
ALLOWED=$(curl -s "http://localhost:8500/v1/connect/intentions/check?source=dashboard&destination=counting" | jq .Allowed)
if [[ $ALLOWED == "false" ]]; then
    fail-message "The dashboard-service can not connect to the counting-service."
    exit 1
fi

exit 0