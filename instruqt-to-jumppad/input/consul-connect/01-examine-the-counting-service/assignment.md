---
slug: examine-the-counting-service
id: ryomnkvhbe6o
type: challenge
title: Examine the counting service
teaser: Taking a look at how services and proxies are configured in Consul Connect.
notes:
- type: text
  contents: |-
    We're **already running** Consul for you, with the **Consul web UI** running on port 8500.

    To start, you'll see a **red X** next to the counting and dashboard services since **neither are running** (so both are unhealthy).

    On the **next slide** we will explain what Consul Connect is and how it works.
- type: video
  url: https://www.youtube.com/embed/8T8t4-hQY74?autoplay=0&rel=0
tabs:
- title: Consul UI
  type: service
  hostname: consul
  port: 8500
- title: Terminal
  type: terminal
  hostname: consul
- title: Editor
  type: code
  hostname: consul
  path: /etc/consul.d
difficulty: basic
timelimit: 300
---
Take a look at the **contents of the counting service configuration**.

---

The counting service configuration is stored in a file located at /etc/consul.d/counting.json.

You can do this run the following command in the **Terminal** tab:

```
cat /etc/consul.d/counting.json
```