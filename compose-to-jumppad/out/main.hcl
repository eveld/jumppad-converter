resource "network" "elastic" {
  subnet = "10.0.5.0/16"
}

resource "container" "kibana" {
  depends_on = ["resource.container.elasticsearch"]

  image {
    name = "kibana:7.16.1"
  }

  network {
    id      = resource.network.elastic.id
    aliases = ["kib"]
  }

  port {
    local  = 5601
    remote = 5601
    host   = 5601
  }
}

resource "container" "elasticsearch" {
  image {
    name = "elasticsearch:7.16.1"
  }

  network {
    id      = resource.network.elastic.id
    aliases = ["es"]
  }

  port {
    local  = 9200
    remote = 9200
    host   = 9200
  }

  port {
    local  = 9300
    remote = 9300
    host   = 9300
  }

  environment = {
    ES_JAVA_OPTS     = "-Xms512m -Xmx512m"
    "discovery.type" = "single-node"
  }

  health_check {
    timeout = "60s"

    exec {
      script = "curl --silent --fail localhost:9200/_cluster/health || exit 1"
    }
  }
}

resource "container" "logstash" {
  depends_on = ["resource.container.elasticsearch"]

  image {
    name = "logstash:7.16.1"
  }

  network {
    id      = resource.network.elastic.id
    aliases = ["log"]
  }

  command = ["logstash", "-f", "/usr/share/logstash/pipeline/logstash-nginx.config"]

  port {
    local  = 5000
    remote = 5000
    host   = 5000
  }

  port {
    protocol = "udp"

    local  = 5000
    remote = 5000
    host   = 5000
  }

  port {
    local  = 5044
    remote = 5044
    host   = 5044
  }

  port {
    local  = 9600
    remote = 9600
    host   = 9600
  }

  environment = {
    LS_JAVA_OPTS           = "-Xms512m -Xmx512m"
    "discovery.seed_hosts" = "logstash"
  }

  volume {
    source      = "/home/erik/code/jumppad/compose-to-jumppad/examples/elasticsearch-logstash-kibana/logstash/pipeline/logstash-nginx.config"
    destination = "/usr/share/logstash/pipeline/logstash-nginx.config"
  }

  volume {
    source      = "/home/erik/code/jumppad/compose-to-jumppad/examples/elasticsearch-logstash-kibana/logstash/nginx.log"
    destination = "/home/nginx.log"
  }
}

