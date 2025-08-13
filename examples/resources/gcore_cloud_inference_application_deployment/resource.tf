resource "gcore_cloud_inference_application_deployment" "example_cloud_inference_application_deployment" {
  project_id = 1
  application_name = "demo-app"
  components_configuration = {
    model = {
      exposed = true
      flavor = "inference-16vcpu-232gib-1xh100-80gb"
      scale = {
        max = 1
        min = 1
      }
      parameter_overrides = {
        foo = {
          value = "value"
        }
      }
    }
  }
  name = "name"
  regions = [1, 2]
  api_keys = ["key1", "key2"]
}
