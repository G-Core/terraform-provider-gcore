resource "gcore_fastedge_secret" "example_fastedge_secret" {
  name = "name"
  comment = "comment"
  secret_slots = [{
    slot = 0
    value = "value"
  }]
}
