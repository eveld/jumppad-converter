#!/bin/bash
curl -L https://github.com/hashicorp/katakoda/raw/master/consul-connect/assets/bin/counting-service -o /usr/local/bin/counting-service
chmod +x /usr/local/bin/counting-service

curl -L https://github.com/hashicorp/katakoda/raw/master/consul-connect/assets/bin/dashboard-service -o /usr/local/bin/dashboard-service
chmod +x /usr/local/bin/dashboard-service