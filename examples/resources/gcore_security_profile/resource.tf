resource "gcore_security_profile" "example_security_profile" {
  fields = [{
    base_field = 1
    field_value = {
      key = "bar"
    }
  }]
  profile_template = 1
  ip_address = "123.43.2.10"
  site = "ED"
}
