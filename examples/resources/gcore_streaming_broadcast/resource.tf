resource "gcore_streaming_broadcast" "example_streaming_broadcast" {
  broadcast = {
    name = "Broadcast"
    ad_id = 1
    custom_iframe_url = ""
    pending_message = "pending_message"
    player_id = 14
    poster = "poster"
    share_url = ""
    show_dvr_after_finish = true
    status = "live"
    stream_ids = [10]
  }
}
