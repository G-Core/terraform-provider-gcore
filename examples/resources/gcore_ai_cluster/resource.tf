provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}
resource "gcore_ai_cluster" "cluster1" {
  region_id = 10
  project_id = 15
  flavor = "g2a-ai-fake-v1pod-8"
  image_id = "f6aa6e75-ab88-4c19-889d-79133366cb83"
  cluster_name      = "cluster1"
  keypair_name = "front"

  volume {
    source     = "image"
    image_id = "f6aa6e75-ab88-4c19-889d-79133366cb83"
    volume_type = "standard"
    size = 20
  }
  
  interface {
    type = "external" 
  }

  security_group {
    id = "4c74142d-9374-4aa6-b11b-43469b66f746"
  }

  cluster_metadata = {
    meta = "meta1"
  }
}  
