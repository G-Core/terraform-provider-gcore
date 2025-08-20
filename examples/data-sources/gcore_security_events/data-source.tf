data "gcore_security_events" "example_security_events" {
  alert_type = "ddos_alert"
  targeted_ip_addresses = "targeted_ip_addresses"
}
