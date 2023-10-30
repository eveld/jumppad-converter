#!/bin/bash
COMMAND=$(cat /root/.bash_history | grep "^cat" | grep "counting.json" | wc -l)
if [ $COMMAND -eq 0 ]; then
  fail-message "Did you take a look at /etc/consul.d/counting.json?"
  exit 1
fi

exit 0
