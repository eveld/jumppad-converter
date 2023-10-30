COUNTING_SERVICE=$(pgrep counting-service)
if [[ -z $COUNTING_SERVICE ]]; then
    fail-message "Is the counting-service running?"
    exit 1
fi

exit 0
