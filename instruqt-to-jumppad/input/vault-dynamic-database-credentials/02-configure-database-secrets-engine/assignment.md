---
slug: configure-database-secrets-engine
id: wmcfwbiihbaq
type: challenge
title: Configure the Database Secrets Engine
teaser: |
  Configure the Database secrets engine on the Vault server.
notes:
- type: text
  contents: |-
    Vault's Database secrets engine dynamically generates credentials (username and password) for many databases.

    In this challenge, you will configure the database secrets engine you enabled in the previous challenge on the path `lob_a/workshop/database` to work with the local instance of the MySQL database. We use a specific path rather than the default "database" to illustrate that multiple instances of the database secrets engine could be configured for different lines of business that might each have multiple databases.
- type: text
  contents: |-
    We will configure a connection and two roles for the database. The roles will allow dynamic generation of credentials with different lifetimes.

    The first role, "workshop-app-long", will generate credentials initially valid for 1 hour with a maximum lifetime of 24 hours. The second role, "workshop-app", will generate credentials initially valid for 3 minutes with a maximum lifetime of 6 minutes.

    To learn more, see these links:
    https://www.vaultproject.io/docs/secrets/databases/
    https://www.vaultproject.io/docs/secrets/databases/mysql-maria/
    https://www.vaultproject.io/api/secret/databases/
    https://www.vaultproject.io/api/secret/databases/mysql-maria/
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
All secrets engines must be configured before they can be used.

We first need to configure the database secrets engine to use the MySQL database plugin and valid connection information. We are configuring a database connection called "wsmysqldatabase" that is allowed to use two roles that we will create below.
```
vault write lob_a/workshop/database/config/wsmysqldatabase \
  plugin_name=mysql-database-plugin \
  connection_url="{{username}}:{{password}}@tcp(localhost:3306)/" \
  allowed_roles="workshop-app","workshop-app-long" \
  username="hashicorp" \
  password="Password123"
```
This will not return anything if successful.

Note that the username and password are templated in the "connection_url" string, getting their values from the "username" and "password" fields.  We do this so that reading the path "lob_a/workshop/database/config/wsmysqldatabase" will not show them.

To test this, try running this command:
```
vault read lob_a/workshop/database/config/wsmysqldatabase
```
You will not see the password.

We used the initial MySQL username "hashicorp" and password "Password123" above. Validate that you can login to the MySQL server with this command:
```
mysql -u hashicorp -pPassword123
```
You should be given a `mysql>` prompt.

Logout of the MySQL server by typing `\q` at the `mysql>` prompt. This should return you to the `root@vault-mysql-server:~#` prompt.

We can make the configuration of the database secrets engine even more secure by rotating the root credentials (actually just the password) that we passed into the configuration.  We do this by running this command:
```
vault write -force lob_a/workshop/database/rotate-root/wsmysqldatabase
```
This should return "Success! Data written to: lob_a/workshop/database/rotate-root/wsmysqldatabase".

Now, if you try to login to the MySQL server with the same command given above, it should fail and give you the message "ERROR 1045 (28000): Access denied for user 'hashicorp'@'localhost' (using password: YES)". Please verify that:
```
mysql -u hashicorp -pPassword123
```

Note: You should **not** use the actual `root` user of the MySQL database (despite the reference to "root credentials"); instead, create a separate user with sufficient privileges to create users and to change its own password.

Now, you should create the first of the two roles we will be using, "workshop-app-long", which generates credentials with an initial lease of 1 hour that can be renewed for up to 24 hours.
```
vault write lob_a/workshop/database/roles/workshop-app-long \
  db_name=wsmysqldatabase \
  creation_statements="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';GRANT ALL ON my_app.* TO '{{name}}'@'%';" \
  default_ttl="1h" \
  max_ttl="24h"
```
This should return "Success! Data written to: lob_a/workshop/database/roles/workshop-app-long".

And then create the second role, "workshop-app" which has shorter default and max leases of 3 minutes and 6 minutes. (These are intentionally set long enough so that you can use the credentials generated for the role to connect to the database but also see them expire in the next challenge.)
```
vault write lob_a/workshop/database/roles/workshop-app \
  db_name=wsmysqldatabase \
  creation_statements="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';GRANT ALL ON my_app.* TO '{{name}}'@'%';" \
  default_ttl="3m" \
  max_ttl="6m"
```
This should return "Success! Data written to: lob_a/workshop/database/roles/workshop-app".

The database secrets engine is now configured to talk to the MySQL server and is allowed to create users with two different roles. In the next challenge, you'll generate credentials (username and password) for these roles.