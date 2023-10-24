---
slug: service-configuration
id: 4dk8ylh4kb1b
type: quiz
title: Service configuration
teaser: The service configuration files describe the service and configure the Consul
  Connect proxy.
notes:
- type: text
  contents: |-
    The counting.json service configuration file contains three important settings:

    1. **Consul will look for a service running on port 9003**. It will advertise that as the counting service. On a properly configured node, this can be reached as counting.service.consul through DNS.

    2. **A blank proxy is defined**. This enables proxy communication through Consul Connect but doesn't define any connections right away.

    3. **A health check examines the local /health endpoint** every second to determine whether the service is healthy and can be exposed to other services.
tabs:
- title: Terminal
  type: terminal
  hostname: consul
answers:
- Consul will look for a service running on port 9003
- An envoy proxy is defined
- A blank proxy is defined
- A health check examines the local /health endpoint every second
- A tcp health check examines port 9003 every second
solution:
- 0
- 2
- 3
difficulty: basic
timelimit: 600
---
Which **three** important settings does the **counting.json** configuration file contain?