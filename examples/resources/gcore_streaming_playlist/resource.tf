resource "gcore_streaming_playlist" "example_streaming_playlist" {
  active = true
  ad_id = 0
  client_id = 0
  client_user_id = 0
  countdown = true
  hls_cmaf_url = "hls_cmaf_url"
  hls_url = "hls_url"
  iframe_url = "iframe_url"
  loop = true
  name = "name"
  player_id = 0
  playlist_type = "live"
  start_time = "start_time"
  video_ids = [0]
}
