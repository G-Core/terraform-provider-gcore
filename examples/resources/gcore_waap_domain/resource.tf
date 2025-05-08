provider gcore {
  permanent_api_token = "768660$.............a43f91f"
  gcore_waap_api = "https://api.cdb-staging.cdn.orange.com/waap"
}

resource "gcore_waap_domain" "domain" {
  name   = "nsf-demo.cdb-staging.cdn.orange.com"
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
}
