resource "track" "consul_connect" {
  title       = "Getting started with Consul Connect"
  owner       = "instruqt"
  teaser      = "Learn the basics of Consul Connect"
  description = file("description.mdx")

  tags       = ["hashicorp", "consul connect"]
  developers = ["ade@instruqt.com"]
  icon       = "https://storage.googleapis.com/instruqt-frontend/assets/hashicorp/tracks/consul.png"

  challenges = {
    resource.challenge.examine_the_counting_service,
    resource.challenge.service_configuration,
    resource.challenge.start_the_counting_service,
    resource.challenge.start_the_dashboard_service,
    resource.challenge.deny_service_connections,
    resource.challenge.allow_dashboard_connection
  }
}

resource "challenge" "examine_the_counting_service" {
  title      = "Examine the counting service"
  teaser     = "Taking a look at how services and proxies are configured in Consul Connect."
  assignment = file("examine_the_counting_service/assignment.mdx")

  tabs = {
    resource.tab.terminal_consul,
    resource.tab.code_consul,
    resource.tab.service_consul_8500
  }

  note {
    type = "text"
    note = file("examine_the_counting_service/notes/note_0.mdx")
  }

  note {
    type = "video"
    url  = "https://www.youtube.com/embed/8T8t4-hQY74?autoplay=0&rel=0"
  }

  check {
    target = "consul"
    script = file("examine_the_counting_service/scripts/check_consul.sh")
  }

  setup {
    target = "consul"
    script = file("examine_the_counting_service/scripts/setup_consul.sh")
  }
}

resource "challenge" "service_configuration" {
  title      = "Service configuration"
  teaser     = "The service configuration files describe the service and configure the Consul Connect proxy."
  assignment = file("service_configuration/assignment.mdx")

  tabs = {
    resource.tab.terminal_consul
  }

  note {
    type = "text"
    note = file("service_configuration/notes/note_0.mdx")
  }
}

resource "challenge" "start_the_counting_service" {
  title      = "Starting the counting service"
  teaser     = "Starting the backend service which we will connect through via the Consul Connect proxy."
  assignment = file("start_the_counting_service/assignment.mdx")

  tabs = {
    resource.tab.terminal_consul,
    resource.tab.service_consul_9003,
    resource.tab.service_consul_8500
  }

  note {
    type = "text"
    note = file("start_the_counting_service/notes/note_0.mdx")
  }

  check {
    target = "consul"
    script = file("start_the_counting_service/scripts/check_consul.sh")
  }

  setup {
    target = "consul"
    script = file("start_the_counting_service/scripts/setup_consul.sh")
  }
}

resource "challenge" "start_the_dashboard_service" {
  title      = "Start the dashboard service"
  teaser     = "Starting the frontend service that will connect to the backend via the Consul Connect proxy."
  assignment = file("start_the_dashboard_service/assignment.mdx")

  tabs = {
    resource.tab.service_consul_8500,
    resource.tab.terminal_consul,
    resource.tab.service_consul_9002
  }

  note {
    type = "text"
    note = file("start_the_dashboard_service/notes/note_0.mdx")
  }

  check {
    target = "consul"
    script = file("start_the_dashboard_service/scripts/check_consul.sh")
  }

  setup {
    target = "consul"
    script = file("start_the_dashboard_service/scripts/setup_consul.sh")
  }
}

resource "challenge" "deny_service_connections" {
  title      = "Deny service connections"
  teaser     = "Denying all connections to services by default."
  assignment = file("deny_service_connections/assignment.mdx")

  tabs = {
    resource.tab.service_consul_8500,
    resource.tab.service_consul_9002,
    resource.tab.terminal_consul
  }

  note {
    type = "text"
    note = file("deny_service_connections/notes/note_0.mdx")
  }

  check {
    target = "consul"
    script = file("deny_service_connections/scripts/check_consul.sh")
  }

  setup {
    target = "consul"
    script = file("deny_service_connections/scripts/setup_consul.sh")
  }
}

resource "challenge" "allow_dashboard_connection" {
  title      = "Allow dashboard connection"
  teaser     = "Allow the frontend service to connect to the backend service."
  assignment = file("allow_dashboard_connection/assignment.mdx")

  tabs = {
    resource.tab.service_consul_8500,
    resource.tab.service_consul_9002
  }

  note {
    type = "text"
    note = file("allow_dashboard_connection/notes/note_0.mdx")
  }

  check {
    target = "consul"
    script = file("allow_dashboard_connection/scripts/check_consul.sh")
  }
}

resource "tab" "terminal_consul" {
  type     = "terminal"
  title    = "Terminal"
  hostname = "consul"
}

resource "tab" "code_consul" {
  type     = "code"
  title    = "Editor"
  hostname = "consul"
  path     = "/etc/consul.d"
}

resource "tab" "service_consul_9003" {
  type     = "service"
  title    = "Counting"
  hostname = "consul"
  port     = 9003
}

resource "tab" "service_consul_9002" {
  type     = "service"
  title    = "Dashboard UI"
  hostname = "consul"
  port     = 9002
}

resource "tab" "service_consul_8500" {
  type     = "service"
  title    = "Consul UI"
  hostname = "consul"
  path     = "/ui/dc1/intentions"
  port     = 8500
}

