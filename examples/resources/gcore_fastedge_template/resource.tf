resource "gcore_fastedge_template" "example_fastedge_template" {
  binary_id = 12345
  name = "api-gateway-template"
  owned = true
  params = [{
    data_type = "string"
    mandatory = true
    name = "api_key"
    descr = "API key for external service authentication"
    metadata = "metadata"
  }]
  long_descr = "Complete API gateway solution with JWT authentication, rate limiting, and request transformation capabilities."
  short_descr = "HTTP API gateway with authentication"
}
