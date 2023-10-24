---
slug: generate-database-creds
id: fuact4dhnief
type: challenge
title: Generate and Use Dynamic Database Credentials
teaser: |
  Generate and use dynamic database credentials for the MySQL database.
notes:
- type: text
  contents: |-
    In this challenge, you will dynamically generate credentials (username and password) against the two roles you configured in the previous challenge.

    You will then connect to the MySQL server with the credentials generated against the shorter duration role, "workshop-app". You will also validate that Vault deletes the credentials from the MySQL server after 3 minutes.

    To learn more, see these links:
    https://www.vaultproject.io/docs/secrets/databases/mysql-maria/
    https://www.vaultproject.io/docs/secrets/databases/#usage
    https://www.vaultproject.io/api/secret/databases/#generate-credentials
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
Now that you have configured the database secrets engine with a connection and two roles for the MySQL database, you can dynamically generate short-lived credentials against the roles and use them to connect to the database.

First, generate credentials for the role, "workshop-app-long", using a `curl` call against the Vault HTTP API and pipe the results to `jq` to make the JSON returned by the API easier to read:
```
curl --header "X-Vault-Token: root" "http://localhost:8200/v1/lob_a/workshop/database/creds/workshop-app-long" | jq
```
You should get something like this:<br>
```
{
  "request_id": "7e4480d3-b8c5-8e21-d45c-29d3440cf01d",
  "lease_id": "lob_a/workshop/database/creds/workshop-app-long/QOsXNeNDruux0QeMVWAF34J3",
  "renewable": true,
  "lease_duration": 3600,
  "data": {
    "password": "A1a-yK2qJ18gRHqUbwJo",
    "username": "v-token-workshop-a-5SEQ8ZLwULeLJ"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```
In these results, you see a several things, including `lease_id`, `username`, and `password`. The first is used if you want to renew or revoke the credentials (as you will do in the next challenge).  The username and password are used to connect to the database.  Note that `renewable` has the value `true`, indicating that the lifetime of the credentials can be extended using Vault's `sys/leases/renew` API endpoint.

You can also generate credentials against the same role with the Vault CLI:
```
vault read lob_a/workshop/database/creds/workshop-app-long
```
This should return something like:<br>
```
Key                Value
---                -----
lease_id           lob_a/workshop/database/creds/workshop-app-long/JeUGIL2xD6BzXSneqity8UmF
lease_duration     1h
lease_renewable    true
password           A1a-zy4ENaf2kwpzGk9t
username           v-token-workshop-a-DM0BJ3eMlMhbf
```

Next, generate credentials against the shorter role, "workshop-app", using the Vault CLI:
```
vault read lob_a/workshop/database/creds/workshop-app
```
This should return something like:<br>
```
Key                Value
---                -----
lease_id           lob_a/workshop/database/creds/workshop-app/t3i85CnEjMlenWbvmJux8SI6
lease_duration     3m
lease_renewable    true
password           A1a-ksrpiyz4tRKmxsRI
username           v-token-workshop-a-kVYT30h6l3e1y
```

Now, use the last set of credentials to connect to the local MySQL server with a command like this:<br>
`mysql -u <username> -p`<br>
Replace `<user_name>` with the one generated from the previous command and provide the generated password when prompted. Then press the '<enter>' or `<return>` key on your keyboard.

You should then see text like that below and be given a `mysql>` prompt:<br>
`
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 7
Server version: 5.7.28 MySQL Community Server (GPL)

mysql>
`<br>

Verify that you can see two of the databases on the MySQL server by running this command:
```
show databases;
```
You should see this:
```
+--------------------+
| Database           |
+--------------------+
| information_schema |
| my_app             |
+--------------------+
2 rows in set (0.00 sec)
```

If you're curious about the my_app database, you can also run the following SQL commands:
```
use my_app;
show tables;
describe customers;
select first_name, last_name from customers;
```
But running these is not required to complete the challenge. If you get any errors while running those commands, it is probable that your credentials have expired.

Logout of the MySQL server by typing `\q` at the `mysql>` prompt. This should return you to the `root@vault-mysql-server:~#` prompt.

At least 3 minutes after you generated the credentials, try to connect to the MySQL server again, using the same username and password as before:<br>
`mysql -u <username> -p`<br>

You should get an error like `ERROR 1045 (28000): Access denied for user 'v-token-workshop-a-tUrh1Z6u5GwKn'@'localhost' (using password: YES)`. If not, type `\q` to log out of MySQL and then try again.

This shows that Vault deleted the credentials from the MySQL database when the lease of the credentials expired after 3 minutes.

In the next challenge, you will learn how to renew and revoke database credentials.