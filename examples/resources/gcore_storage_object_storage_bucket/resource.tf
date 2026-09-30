resource "gcore_storage_object_storage_bucket" "example_storage_object_storage_bucket" {
  storage_id = 0
  name = "my-new-bucket"
  cors = {
    allowed_origins = ["https://example.com"]
  }
  storage_object_storage_bucket_lifecycle = {
    expiration_days = 30
  }
  policy = {
    is_public = true
  }
}
