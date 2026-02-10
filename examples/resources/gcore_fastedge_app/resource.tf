resource "gcore_fastedge_app" "example_fastedge_app" {
  binary = 0
  comment = "comment"
  debug = true
  env = {
    var1 = "value1"
    var2 = "value2"
  }
  log = "kafka"
  name = "name"
  rsp_headers = {
    header1 = "value1"
    header2 = "value2"
  }
  secrets = {
    foo = {
      id = 0
    }
  }
  status = 0
  stores = {
    foo = {
      id = 0
    }
  }
  template = 0
}
