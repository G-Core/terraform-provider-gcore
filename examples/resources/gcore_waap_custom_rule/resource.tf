provider gcore {
    permanent_api_token = "768660$.............a43f91f"
}

data "gcore_waap_tag" "proxy_network" {
  name = "Proxy Network"
}

resource "gcore_cdn_resource" "cdn_resource" {
    cname  = "api.example.com"
    origin = "origin.example.com"
    options {
        waap { value = true }
    }
}

resource "gcore_waap_domain" "domain" {
    name   = gcore_cdn_resource.cdn_resource.cname
    status = "monitor"
}

resource "gcore_waap_custom_rule" "custom_rule" {
    domain_id = gcore_waap_domain.domain.id
    name = "Custom Rule"
    enabled = true
    action {
        block {
            status_code = 403
            action_duration = "5m"
        }
    }
    conditions {
        ip {
            negation    = true
            ip_address = "192.168.0.6"
        }
        http_method {
            http_method  = "POST"
        }
        tags {
            tags = [data.gcore_waap_tag.proxy_network.id]
        }
    }
}
