run_tc4 = true
tc4_name = "test-cdn-og-tc4-updated-skill"
tc4_proxy_next_upstream = ["error", "timeout", "http_500"]
tc4_sources = [
  {
    source  = "93.184.216.41"
    enabled = true
    backup  = false
  }
]
