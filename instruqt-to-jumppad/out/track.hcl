resource "track" "vault_dynamic_database_credentials" {
  title       = "Vault Dynamic Database Credentials"
  owner       = "instruqt"
  teaser      = "Generate dynamic credentials for a MySQL database from Vault."
  description = <<-EOF
    The Vault Database secrets engine lets you generate dynamic, time-bound credentials for many different databases.
    
    In this track, you will do this for a MySQL database that is running on the same server as the Vault server itself.
  EOF
  tags        = ["dynamic-secrets", "vault", "database"]
  developers  = ["wietse@instruqt.com"]
  icon        = "https://storage.googleapis.com/instruqt-hashicorp-tracks/logo/vault.png"

  challenges = {
    resource.challenge.enable_database_secrets_engine,
    resource.challenge.configure_database_secrets_engine,
    resource.challenge.generate_database_creds,
    resource.challenge.renew_revoke_database_creds
  }
}

resource "challenge" "enable_database_secrets_engine" {
  title      = "Enable the Database Secrets Engine"
  teaser     = "Enable the Database secrets engine on the Vault server."
  assignment = file("assignments/enable_database_secrets_engine.mdx")

  tabs = {
    resource.tab.enable_database_secrets_engine_0,
    resource.tab.enable_database_secrets_engine_1,
    resource.tab.enable_database_secrets_engine_2
  }

  note {
    type     = "text"
    contents = <<-EOF
      Secrets engines are Vault plugins that store, generate, or encrypt data. Secrets engines are incredibly flexible, so it is easiest to think about them in terms of their function.
      
      Vault's Database secrets engine dynamically generates credentials for many databases.
      
      To learn more, see these links:
      https://www.vaultproject.io/docs/secrets/databases/
      https://www.vaultproject.io/docs/secrets/databases/mysql-maria/
    EOF
  }

  check {
    target = "vault-mysql-server"
    script = "out/scripts/enable_database_secrets_engine/check_vault-mysql-server.sh"
  }

  setup {
    target = "vault-mysql-server"
    script = "out/scripts/enable_database_secrets_engine/setup_vault-mysql-server.sh"
  }

  solve {
    target = "vault-mysql-server"
    script = "out/scripts/enable_database_secrets_engine/solve_vault-mysql-server.sh"
  }
}

resource "tab" "code_vault-mysql-server" {
  type     = "terminal"
  title    = "Vault CLI"
  hostname = "vault-mysql-server"
}

resource "tab" "code_vault-mysql-server" {
  type     = "service"
  title    = "Vault UI"
  hostname = "vault-mysql-server"
  port     = 8200
}

resource "tab" "code_vault-mysql-server" {
  type     = "code"
  title    = "Editor"
  hostname = "vault-mysql-server"
}

resource "challenge" "configure_database_secrets_engine" {
  title      = "Configure the Database Secrets Engine"
  teaser     = "Configure the Database secrets engine on the Vault server."
  assignment = file("assignments/configure_database_secrets_engine.mdx")

  tabs = {
    resource.tab.configure_database_secrets_engine_0,
    resource.tab.configure_database_secrets_engine_1
  }

  note {
    type     = "text"
    contents = <<-EOF
      Vault's Database secrets engine dynamically generates credentials (username and password) for many databases.
      
      In this challenge, you will configure the database secrets engine you enabled in the previous challenge on the path `lob_a/workshop/database` to work with the local instance of the MySQL database. We use a specific path rather than the default "database" to illustrate that multiple instances of the database secrets engine could be configured for different lines of business that might each have multiple databases.
    EOF
  }

  note {
    type     = "text"
    contents = <<-EOF
      We will configure a connection and two roles for the database. The roles will allow dynamic generation of credentials with different lifetimes.
      
      The first role, "workshop-app-long", will generate credentials initially valid for 1 hour with a maximum lifetime of 24 hours. The second role, "workshop-app", will generate credentials initially valid for 3 minutes with a maximum lifetime of 6 minutes.
      
      To learn more, see these links:
      https://www.vaultproject.io/docs/secrets/databases/
      https://www.vaultproject.io/docs/secrets/databases/mysql-maria/
      https://www.vaultproject.io/api/secret/databases/
      https://www.vaultproject.io/api/secret/databases/mysql-maria/
    EOF
  }

  check {
    target = "vault-mysql-server"
    script = "out/scripts/configure_database_secrets_engine/check_vault-mysql-server.sh"
  }

  solve {
    target = "vault-mysql-server"
    script = "out/scripts/configure_database_secrets_engine/solve_vault-mysql-server.sh"
  }
}

resource "tab" "service_vault-mysql-server_8200" {
  type     = "terminal"
  title    = "Vault CLI"
  hostname = "vault-mysql-server"
}

resource "tab" "service_vault-mysql-server_8200" {
  type     = "service"
  title    = "Vault UI"
  hostname = "vault-mysql-server"
  port     = 8200
}

resource "challenge" "generate_database_creds" {
  title      = "Generate and Use Dynamic Database Credentials"
  teaser     = "Generate and use dynamic database credentials for the MySQL database."
  assignment = file("assignments/generate_database_creds.mdx")

  tabs = {
    resource.tab.generate_database_creds_0,
    resource.tab.generate_database_creds_1
  }

  note {
    type     = "text"
    contents = <<-EOF
      In this challenge, you will dynamically generate credentials (username and password) against the two roles you configured in the previous challenge.
      
      You will then connect to the MySQL server with the credentials generated against the shorter duration role, "workshop-app". You will also validate that Vault deletes the credentials from the MySQL server after 3 minutes.
      
      To learn more, see these links:
      https://www.vaultproject.io/docs/secrets/databases/mysql-maria/
      https://www.vaultproject.io/docs/secrets/databases/#usage
      https://www.vaultproject.io/api/secret/databases/#generate-credentials
    EOF
  }

  check {
    target = "vault-mysql-server"
    script = "out/scripts/generate_database_creds/check_vault-mysql-server.sh"
  }

  setup {
    target = "vault-mysql-server"
    script = "out/scripts/generate_database_creds/setup_vault-mysql-server.sh"
  }

  solve {
    target = "vault-mysql-server"
    script = "out/scripts/generate_database_creds/solve_vault-mysql-server.sh"
  }
}

resource "tab" "service_vault-mysql-server_8200" {
  type     = "terminal"
  title    = "Vault CLI"
  hostname = "vault-mysql-server"
}

resource "tab" "service_vault-mysql-server_8200" {
  type     = "service"
  title    = "Vault UI"
  hostname = "vault-mysql-server"
  port     = 8200
}

resource "challenge" "renew_revoke_database_creds" {
  title      = "Renew and Revoke Database Credentials"
  teaser     = "Renew and revoke database credentials for the MySQL database."
  assignment = file("assignments/renew_revoke_database_creds.mdx")

  tabs = {
    resource.tab.renew_revoke_database_creds_0,
    resource.tab.renew_revoke_database_creds_1
  }

  note {
    type     = "text"
    contents = <<-EOF
      In this challenge, you will learn how to renew and revoke credentials generated by Vault's database secrets engine.
      
      You will see that it is possible to extend the lifetime of generated credentials when they have not yet expired by renewing them. You will also see that they cannot be renewed beyond the `max_ttl` of the role against which the credentials were generated.
      
      To learn more, see these links:
      https://www.vaultproject.io/api/system/leases/#renew-lease
      https://www.vaultproject.io/api/system/leases/#revoke-lease
    EOF
  }

  check {
    target = "vault-mysql-server"
    script = "out/scripts/renew_revoke_database_creds/check_vault-mysql-server.sh"
  }

  solve {
    target = "vault-mysql-server"
    script = "out/scripts/renew_revoke_database_creds/solve_vault-mysql-server.sh"
  }
}

resource "tab" "service_vault-mysql-server_8200" {
  type     = "terminal"
  title    = "Vault CLI"
  hostname = "vault-mysql-server"
}

resource "tab" "service_vault-mysql-server_8200" {
  type     = "service"
  title    = "Vault UI"
  hostname = "vault-mysql-server"
  port     = 8200
}

