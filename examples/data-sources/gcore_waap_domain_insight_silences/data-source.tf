data "gcore_waap_domain_insight_silences" "example_waap_domain_insight_silences" {
  domain_id = 1
  id = ["182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e", "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"]
  author = "author"
  comment = "comment"
  insight_type = ["string", "string"]
}
