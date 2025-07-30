resource "gcore_fastedge_kv_store" "example_fastedge_kv_store" {
  byod = {
    prefix = "prefix"
    url = "url"
  }
  comment = "comment"
}
