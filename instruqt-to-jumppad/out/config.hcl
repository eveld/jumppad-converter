resource "network" "main" {
  subnet = "10.0.5.0/16"
}

resource "vm" "vault-mysql-server" {
  config {
    arch = "x86_64"
  }

  image = "instruqt-hashicorp/vault-1-8-3-with-mysql-and-python-web-app"

  resources {
    cpu    = 1
    memory = 4096
  }

  network {
    id = resource.network.main.id
  }

  environment = {
    MYSQL_ENDPOINT = "localhost:3306"
    MYSQL_HOST     = "localhost"
    MYSQL_PASSWORD = "sJ2w*8NX"
    MYSQL_PORT     = "3306"
    SHELL          = "/bin/bash -l"
    VAULT_ADDR     = "http://localhost:8200"
    VAULT_TOKEN    = "root"
  }
}

