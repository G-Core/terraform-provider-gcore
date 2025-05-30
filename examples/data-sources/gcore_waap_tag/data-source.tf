provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_waap_tag" "abnormal_traffic_volume_tag" {
  name = "Abnormal Traffic Volume"
}

output "abnormal_traffic_volume_tag_id" {
  value = data.gcore_waap_tag.abnormal_traffic_volume_tag.id
}
