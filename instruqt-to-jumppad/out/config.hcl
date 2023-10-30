resource "network" "main" {
  subnet = "10.0.5.0/16"
}

resource "container" "consul" {
  image {
    name = "gcr.io/instruqt/consul-connect"
  }

  network {
    id = resource.network.main.id
  }

  port {
    local  = 8300
    remote = 8300
    host   = 8300
  }

  port {
    local  = 8301
    remote = 8301
    host   = 8301
  }

  port {
    local  = 8302
    remote = 8302
    host   = 8302
  }

  port {
    local  = 8500
    remote = 8500
    host   = 8500
  }

  port {
    local  = 8600
    remote = 8600
    host   = 8600
  }

  port {
    local  = 9002
    remote = 9002
    host   = 9002
  }

  port {
    local  = 9003
    remote = 9003
    host   = 9003
  }

  environment = {
    SHELL = "/bin/bash"
  }

  resources {
    memory = 128
  }
}

