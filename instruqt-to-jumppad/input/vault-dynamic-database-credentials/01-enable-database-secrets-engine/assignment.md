---
slug: enable-database-secrets-engine
id: jyvln5ec7w07
type: challenge
title: Enable the Database Secrets Engine
teaser: |
  Enable the Database secrets engine on the Vault server.
notes:
- type: text
  contents: |-
    Secrets engines are Vault plugins that store, generate, or encrypt data. Secrets engines are incredibly flexible, so it is easiest to think about them in terms of their function.

    Vault's Database secrets engine dynamically generates credentials for many databases.

    To learn more, see these links:
    https://www.vaultproject.io/docs/secrets/databases/
    https://www.vaultproject.io/docs/secrets/databases/mysql-maria/
tabs:
- title: Vault CLI
  type: terminal
  hostname: vault-mysql-server
- title: Vault UI
  type: service
  hostname: vault-mysql-server
  port: 8200
- title: Editor
  type: code
  hostname: vault-mysql-server
difficulty: basic
timelimit: 600
---
The Database secrets engine generates credentials dynamically for various databases.  In this track, we are using an instance of MySQL that is running on the same VM as the Vault server itself.

The Database credentials are time-bound and are automatically revoked when the Vault lease expires. The credentials can also be revoked at any time.

All secrets engines must be enabled before they can be used. Check which secrets engines are currently enabled.
```
vault secrets list
```

Note that the Database secrets engine is not enabled. Please enable it at the path "lob_a/workshop/database".
```
vault secrets enable -path=lob_a/workshop/database database
```

If you like, you can login to the Vault UI with the Vault token `root` and see that the database secrets engine has been enabled.