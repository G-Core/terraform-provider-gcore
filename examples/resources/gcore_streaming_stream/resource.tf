resource "gcore_streaming_stream" "example_streaming_stream" {
  name = "name"
  active = true
  auto_record = true
  broadcast_ids = [0]
  cdn_id = 0
  client_entity_data = "client_entity_data"
  client_user_id = 0
  dvr_duration = 0
  dvr_enabled = true
  hls_mpegts_endlist_tag = true
  html_overlay = true
  low_latency_enabled = true
  projection = "regular"
  pull = true
  quality_set_id = 0
  record_type = "origin"
  uri = "uri"
}
