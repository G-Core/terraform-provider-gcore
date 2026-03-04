# Test 5: Remove one mapping entry (from 2 back to 1)
mapping_name = "test-tf-minimal"
mapping_entries = [
  {
    tags  = ["europe", "primary"]
    cidr4 = ["10.0.0.0/8"]
    cidr6 = []
  }
]
