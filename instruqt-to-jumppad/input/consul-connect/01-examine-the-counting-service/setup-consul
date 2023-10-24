#!/bin/bash
until $(curl -o /dev/null -s --head --fail http://localhost:8500); do
    sleep 1
done

cat <<EOF > /etc/consul.d/counting.json
{
  "service": {
    "name": "counting",
    "port": 9003,
    "connect": {
      "proxy": {
      }
    },
    "check": {
      "id": "counting-check",
      "http": "http://localhost:9003/health",
      "method": "GET",
      "interval": "1s",
      "timeout": "1s"
    }
  }
}
EOF

cat <<EOF > /etc/consul.d/dashboard.json
{
  "service": {
    "name": "dashboard",
    "port": 9002,
    "connect": {
      "proxy": {
        "config": {
          "upstreams": [
            {
              "destination_name": "counting",
              "local_bind_port": 9001
            }
          ]
        }
      }
    },
    "check": {
      "id": "dashboard-check",
      "http": "http://localhost:9002/health",
      "method": "GET",
      "interval": "1s",
      "timeout": "1s"
    }
  }
}
EOF

consul reload