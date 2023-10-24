---
slug: allow-dashboard-connection
id: g9oxvkpdb6np
type: challenge
title: Allow dashboard connection
teaser: Allow the frontend service to connect to the backend service.
notes:
- type: text
  contents: |-
    For the final step in this tutorial, let's enable communication between the dashboard-service and the counting-service.

    This intention will **allow communication** from the source dashboard service to the destination counting service.
tabs:
- title: Consul UI
  type: service
  hostname: consul
  path: /ui/dc1/intentions
  port: 8500
- title: Dashboard UI
  type: service
  hostname: consul
  port: 9002
difficulty: basic
timelimit: 300
---
**Add an intention that allows the dashboard service to connect to the counting service**.

---

You can do this by selecting **dashboard** in the first pulldown and **counting** in the second pulldown. Then selecting **allow** below. Optionally adding a description.

Then save the new intention.

Finally, view the **Dashboard** tab which will automatically discover the new connection (no refresh is needed for this websockets-driven application). It should not take more than a few seconds.

You will not need to restart any services since intentions which allow connectivity will take effect dynamically.