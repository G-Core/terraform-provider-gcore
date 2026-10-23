# Tag requests from hosting providers in selected countries, for use in other rules
resource "gcore_waap_domain_custom_rule" "tag_hosting" {
  domain_id = 12345
  name      = "Tag hosting traffic"
  enabled   = true

  action = {
    tag = {
      tags = ["hosting"]
    }
  }

  conditions = [
    { owner_types = { owner_types = ["HOSTING_SERVICES"] } },
    { country = { country_code = ["CN", "RU"] } },
  ]
}
