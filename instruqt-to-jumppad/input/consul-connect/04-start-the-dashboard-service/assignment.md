---
slug: start-the-dashboard-service
id: djsdgizqs6sy
type: challenge
title: Start the dashboard service
teaser: Starting the frontend service that will connect to the backend via the Consul
  Connect proxy.
notes:
- type: text
  contents: In this challenge you'll start the frontend dashboard-service and it will
    communicate to the counting service through an **encrypted proxy**.
tabs:
- title: Consul UI
  type: service
  hostname: consul
  port: 8500
- title: Terminal
  type: terminal
  hostname: consul
- title: 'Dashboard '
  type: service
  hostname: consul
  port: 9002
difficulty: basic
timelimit: 300
---
**Now start the dashboard service**, specifying PORT as an environment variable as in the previous challenge.

---

Consul is configured to look for the **dashboard-service** on port **9002**.
The configuration for this service is located at `/etc/consul.d/dashboard.json`.

You can do this by running the following command in the **Terminal** tab:

```
PORT=9002 dashboard-service
```

When started correctly, you can view the demo dashboard application in the **Dashboard** tab.
And the health check will return alive in the **Consul UI** tab.