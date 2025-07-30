resource "gcore_streaming_video_subtitle" "example_streaming_video_subtitle" {
  video_id = 0
  body = {
    language = "language"
    name = "German (AI-generated)"
    vtt = <<EOT
    WEBVTT

    1
    00:00:07.154 --> 00:00:12.736
    Wir haben 100 Millionen registrierte Benutzer oder aktive Benutzer, die mindestens einmal pro Woche spielen.

    2
    00:00:13.236 --> 00:00:20.198
    Wir haben vielleicht 80 oder 100.000, die auf einem bestimmten Cluster spielen.
    EOT
    auto_transcribe_audio_language = "auto"
    auto_translate_subtitles_language = "default"
  }
}
