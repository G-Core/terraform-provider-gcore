resource "gcore_streaming_stream_overlay" "example_streaming_stream_overlay" {
  stream_id = 0
  body = [{
    url = "http://domain.com/myoverlay1.html"
    height = 40
    stretch = false
    width = 120
    x = 30
    y = 30
  }]
}
