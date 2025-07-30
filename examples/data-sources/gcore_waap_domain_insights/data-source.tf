data "gcore_waap_domain_insights" "example_waap_domain_insights" {
  domain_id = 1
  id = ["182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e", "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"]
  description = "description"
  insight_type = ["string", "string"]
  ordering = "id"
  status = ["OPEN", "ACKED"]
}
