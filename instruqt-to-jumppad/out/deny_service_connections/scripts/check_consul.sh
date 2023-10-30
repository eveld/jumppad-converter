#!/bin/bash
ALLOWED=$(curl -s "http://localhost:8500/v1/connect/intentions/check?source=dashboard&destination=counting" | jq .Allowed)
if [[ $ALLOWED == "true" ]]; then
    fail-message "Communication is still possible from * to *. Did you correctly define the intention?"
    exit 1
fi

exit 0