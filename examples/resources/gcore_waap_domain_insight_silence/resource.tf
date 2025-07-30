resource "gcore_waap_domain_insight_silence" "example_waap_domain_insight_silence" {
  domain_id = 1
  author = "author"
  comment = "comment"
  insight_type = "insight_type"
  labels = {
    foo = "string"
  }
  expire_at = "2019-12-27T18:11:19.117Z"
}
