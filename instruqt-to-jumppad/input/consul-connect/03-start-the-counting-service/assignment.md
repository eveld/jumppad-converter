---
slug: start-the-counting-service
id: iofyw71drrbk
type: challenge
title: Starting the counting service
teaser: Starting the backend service which we will connect through via the Consul
  Connect proxy.
notes:
- type: text
  contents: |-
    Now let's start some services.

    In the following challenges we will then let them communicate via Consul Connect.
tabs:
- title: Consul UI
  type: service
  hostname: consul
  port: 8500
- title: Terminal
  type: terminal
  hostname: consul
- title: Counting
  type: service
  hostname: consul
  port: 9003
difficulty: basic
timelimit: 300
---
**Start the counting service**, specifying PORT as an environment variable.

---

You can do this by running the following command in the **Terminal** tab:

```
PORT=9003 counting-service
```

You can view the output of the counting service in the **Counting** tab at the top of the screen.
It's a simple JSON API that returns a number.

If you refresh the Consul Web UI in the **Consul UI** tab, you'll notice that the counting service now shows as healthy.