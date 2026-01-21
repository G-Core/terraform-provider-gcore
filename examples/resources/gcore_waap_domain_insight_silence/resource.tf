resource "gcore_waap_domain_insight_silence" "example_waap_domain_insight_silence" {
  domain_id = 1
  author = "author"
  comment = "comment"
  insight_type = "26f1klzn5713-56bincal4ca-60zz1k91s4"
  labels = {
    foo = "string"
  }
  expire_at = "2019-12-27T18:11:19.117Z"
}
