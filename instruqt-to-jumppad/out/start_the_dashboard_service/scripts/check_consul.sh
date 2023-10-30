#!/bin/bash
DASHBOARD_SERVICE=$(pgrep dashboard-service)
if [[ -z $DASHBOARD_SERVICE ]]; then
    fail-message "Is the dashboard-service running?"
    exit 1
fi

exit 0