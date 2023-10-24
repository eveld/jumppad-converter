---
slug: renew-revoke-database-creds
id: qlxcp8dsaxzv
type: challenge
title: Renew and Revoke Database Credentials
teaser: |
  Renew and revoke database credentials for the MySQL database.
notes:
- type: text
  contents: |-
    In this challenge, you will learn how to renew and revoke credentials generated by Vault's database secrets engine.

    You will see that it is possible to extend the lifetime of generated credentials when they have not yet expired by renewing them. You will also see that they cannot be renewed beyond the `max_ttl` of the role against which the credentials were generated.

    To learn more, see these links:
    https://www.vaultproject.io/api/system/leases/#renew-lease
    https://www.vaultproject.io/api/system/leases/#revoke-lease
tabs:
- title: Vault CLI
  type: terminal
  hostname: vault-mysql-server
- title: Vault UI
  type: service
  hostname: vault-mysql-server
  port: 8200
difficulty: basic
timelimit: 1200
---
In addition to using Vault's database secrets engine to generate credentials for databases, you can also use it to extend their lifetime or revoke them.

First, generate new credentials against the shorter role, "workshop-app", using the Vault CLI:
```
vault read lob_a/workshop/database/creds/workshop-app
```
This should return something like:<br>
```
Key                Value
---                -----
lease_id           lob_a/workshop/database/creds/workshop-app/t3i85CnEjMlenWbvmJux8SI6
lease_duration     1m
lease_renewable    true
password           A1a-ksrpiyz4tRKmxsRI
username           v-token-workshop-a-kVYT30h6l3e1y
```

The lease on credentials returned by the database secrets engine can be manually renewed with a call like this:<br>
`vault write sys/leases/renew lease_id="<lease_id>" increment="120"`<br>
where you should replace `<lease_id>` with the lease_id returned by the previous command. In this case, we are extending the life of the credentials by 2 minutes.

This command should return something like this:<br>
```
Key                Value
---                -----
lease_id           lob_a/workshop/database/creds/workshop-app/jr0to1AfDqE2eiPl2GzYShzR
lease_duration     2m
lease_renewable    true
```

Now, examine the current lease with a command like this:<br>
`vault write sys/leases/lookup lease_id="<lease_id>"`<br>
where you should replace `<lease_id>` with the lease_id returned by either of the last two commands.  This should return something like:<br>
```
Key             Value
---             -----
expire_time     2019-12-12T17:52:41.267656422Z
id              lob_a/workshop/database/creds/workshop-app/5PfygQTgMTwJNCEVqujwaVLS
issue_time      2019-12-12T17:49:41.267656019Z
last_renewal    <nil>
renewable       true
ttl             1m45s
```

The `ttl` will tell you the remaining time to live of the lease and the credentials. When the lease expires, Vault will delete the credentials from MySQL.

Extending the lease will only work if the lease has not yet expired. Additionally, the lease on the credentials cannot be extended beyond the original time of their creation plus the duration given by the `max_ttl` parameter of the role.  If either of these conditions apply, you will get an error.

For instance, if you try to lookup a lease that has already expired, you will get an `invalid lease` error. If you try to extend the lease with an increment of 600 seconds (10 minutes), you will see an error like:<br>
```
WARNING! The following warnings were returned from Vault:
    * TTL of "10m0s" exceeded the effective max_ttl of "2m17s";
    TTL value is capped accordingly
```

Finally, let's explore how you can revoke database credentials.  First, generate a new set of credentials:
```
vault read lob_a/workshop/database/creds/workshop-app
```

Then revoke the credentials:<br>
`vault write sys/leases/revoke lease_id="<lease_id>"`<br>
replacing `<lease_id>` with the one returned with the generated credentials. You should see "Success! Data written to: sys/leases/revoke" returned.

Try to login to the MySQL server with the revoked credentials:<br>
`mysql -u <username> -p`<br>
replacing `<username>` with the generated username and providing the generated password when prompted. You should see a mesage including "ERROR 1045 (28000): Access denied for user".

Congratulations on finishing the Vault Dynamic Database Credentials track.  We recommend that you explore the [Vault Transit Engine](https://instruqt.com/hashicorp/tracks/vault-transit-engine) track next.