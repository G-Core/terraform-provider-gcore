# Test 6: Clear cidr4 (empty array)
mapping_name = "test-tf-minimal"
mapping_entries = [
  {
    tags  = ["europe", "primary"]
    cidr4 = []
    cidr6 = []
  }
]
