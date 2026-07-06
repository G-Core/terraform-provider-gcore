data "gcore_dns_zone_rrsets" "example_dns_zone_rrsets" {
  zone_name = "zoneName"
  order_by = "order_by"
  order_direction = "asc"
}
