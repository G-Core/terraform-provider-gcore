resource "gcore_streaming_video" "example_streaming_video" {
  video = {
    name = "IBC 2025 - International Broadcasting Convention"
    auto_transcribe_audio_language = "auto"
    auto_translate_subtitles_language = "disable"
    client_user_id = 10
    clip_duration_seconds = 60
    clip_start_seconds = 137
    custom_iframe_url = "custom_iframe_url"
    description = "We look forward to welcoming you at IBC2025, which will take place 12-15 September 2025."
    directory_id = 800
    origin_http_headers = "Authorization: Bearer ..."
    origin_url = "https://www.googleapis.com/drive/v3/files/...?alt=media"
    poster = "data:image/jpeg;base64,/9j/4AA...qf/2Q=="
    priority = 0
    projection = "regular"
    quality_set_id = 0
    remote_poster_url = "remote_poster_url"
    remove_poster = true
    screenshot_id = -1
    share_url = "share_url"
    source_bitrate_limit = true
  }
}
