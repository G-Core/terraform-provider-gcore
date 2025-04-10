resource "gcore_fastedge_binary" "test_binary" {
  filename = "test.wasm"
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
