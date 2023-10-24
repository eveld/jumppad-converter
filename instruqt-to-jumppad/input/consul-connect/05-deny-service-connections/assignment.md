---
slug: deny-service-connections
id: nzgajsahsxo0
type: challenge
title: Deny service connections
teaser: Denying all connections to services by default.
notes:
- type: text
  contents: |-
    Let's secure these services by defining **intentions** in Consul.

    Intentions are the way that we describe which services should be able to connect to other services.
tabs:
- title: Consul UI
  type: service
  hostname: consul
  path: /ui/dc1/intentions
  port: 8500
- title: Dashboard
  type: service
  hostname: consul
  port: 9002
- title: Terminal
  type: terminal
  hostname: consul
difficulty: basic
timelimit: 600
---
**Create an intention** that **denies all** communication between all services.

---

This is achieved be creating an intention from *** (All Services)** to *** (All Services)** with a value of **deny**.

You can optionally add a description such as "Deny all communication by default".

Then save the newly created intention.

Since existing proxies will not be terminated when a deny rule is created, we must restart the dashboard-service.

Restart it by typing

```
killall dashboard-service
```

in the Terminal tab. We will then restart it for you.

If you look at the **Dashboard** and you will see that it cannot reach the backend counting-service because all connections are denied.