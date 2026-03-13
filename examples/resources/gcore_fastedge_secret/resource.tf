resource "gcore_fastedge_secret" "example_fastedge_secret" {
  name = "name"
  comment = "comment"
  secret_slots = [{
    slot = 1704067200
    value = "P@ssw0rd123!"
  }]
}
