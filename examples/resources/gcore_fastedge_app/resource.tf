resource "gcore_fastedge_binary" "test_binary" {
  filename = "test.wasm"
}

resource "gcore_fastedge_app" "app_from_binary" {
  status = "enabled"
  name = "terraform-test1"
  comment = "Terraform test app 1"
  binary = gcore_fastedge_binary.test_binary.id
  env = {
    "foo" = "bar"
  }
}

resource "gcore_fastedge_template" "test_template" {
  name = "terraform_test_template"
  binary = gcore_fastedge_binary.test_binary.id
  short_descr = "short description"
  long_descr = "long description"
  param {
      name = "foo"
      type = "string"
      mandatory = true
      descr = "Parameter foo"
  }
  param {
      name = "bar"
      type = "number"
      descr = "Parameter bar"
  }
}

resource "gcore_fastedge_app" "app_from_template" {
  status = "enabled"
  name = "terraform-test2"
  comment = "Terraform test app 2"
  template = gcore_fastedge_template.test_template.id
  env = {
    "foo" = "foo_value"
    "bar" = 123
  }
}
