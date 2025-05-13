provider gcore {
  permanent_api_token = "768660$.............a43f91f"
}

resource "gcore_cdn_resource" "example" {
  cname  = "api.example.com"
  origin = "origin.example.com"
  options {
    waap { value = true }
  }
}

resource "gcore_waap_domain" "domain" {
  name   = gcore_cdn_resource.example.cname
  status = "monitor"

  settings {
    ddos {
      global_threshold     = 2000
      burst_threshold      = 1000
    }
    
    api {
      api_urls = [
        "https://api.example.com/v1",
        "https://api.example.com/v2"
      ]
    }
  }

  policies {
    protocol_validation {
      invalid_user_agent      = false
      unknown_user_agent      = false
    }
    core_waf_owasp_top_threats{
      sql_injection           = false
      xss                     = false
    }
  }
}
