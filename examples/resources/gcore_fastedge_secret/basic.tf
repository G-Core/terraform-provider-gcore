# Create a FastEdge secret with multiple slots
resource "gcore_fastedge_secret" "example" {
  name    = "my-edge-secret"
  comment = "API keys for edge functions"

  secret_slots = [
    {
      slot  = 0
      value = var.secret_slot_0
    },
    {
      slot  = 1
      value = var.secret_slot_1
    },
  ]
}
