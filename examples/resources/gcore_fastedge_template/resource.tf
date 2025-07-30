resource "gcore_fastedge_template" "example_fastedge_template" {
  binary_id = 0
  name = "name"
  owned = true
  params = [{
    data_type = "string"
    mandatory = true
    name = "name"
    descr = "descr"
  }]
  long_descr = "long_descr"
  short_descr = "short_descr"
}
