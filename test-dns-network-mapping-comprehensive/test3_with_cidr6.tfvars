# Test 3: With cidr4 and cidr6 values
mapping_name = "test-tf-minimal"
mapping_entries = [
  {
    tags  = ["test-tag-1"]
    cidr4 = ["192.168.0.0/16", "10.0.0.0/8"]
    cidr6 = ["2001:db8::/32"]
  }
]
