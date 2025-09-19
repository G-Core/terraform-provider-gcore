provider gcore {
  permanent_api_token = "768660$.............a43f91f"
}

resource "gcore_cdn_client_config" "my_account" {
  utilization_level = 90
}

output "cdn_client" {
  value = {
    id                            = gcore_cdn_client_config.my_account.id
    cname                         = gcore_cdn_client_config.my_account.cname
    created                       = gcore_cdn_client_config.my_account.created
    updated                       = gcore_cdn_client_config.my_account.updated
    utilization_level             = gcore_cdn_client_config.my_account.utilization_level
    use_balancer                  = gcore_cdn_client_config.my_account.use_balancer
    auto_suspend_enabled          = gcore_cdn_client_config.my_account.auto_suspend_enabled
    cdn_resources_rules_max_count = gcore_cdn_client_config.my_account.cdn_resources_rules_max_count

    service = {
      enabled = gcore_cdn_client_config.my_account.service[0].enabled
      status  = gcore_cdn_client_config.my_account.service[0].status
      updated = gcore_cdn_client_config.my_account.service[0].updated
    }
  }
  sensitive = false
}