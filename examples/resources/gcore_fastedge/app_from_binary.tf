resource "gcore_fastedge_binary" "test_binary" {
  filename = "test.wasm"
}

resource "gcore_fastedge_app" "test_app" {
  status = "enabled"
  name = "terraform-test"
  comment = "Terraform test app"
  binary = gcore_fastedge_binary.test_binary.id
  env = {
    "foo" = "bar"
  }
}
