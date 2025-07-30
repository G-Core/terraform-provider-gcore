resource "gcore_streaming_restream" "example_streaming_restream" {
  restream = {
    active = true
    client_user_id = 10
    live = true
    name = "first restream"
    stream_id = 20
    uri = "rtmp://a.rtmp.youtube.com/live/k17a-13s8"
  }
}
