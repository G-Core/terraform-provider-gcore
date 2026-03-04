# Test 4: Multiple mapping entries
mapping_name = "test-tf-minimal"
mapping_entries = [
  {
    tags  = ["europe", "primary"]
    cidr4 = ["10.0.0.0/8"]
    cidr6 = []
  },
  {
    tags  = ["asia", "secondary"]
    cidr4 = ["172.16.0.0/12"]
    cidr6 = ["2001:db8:1::/48"]
  }
]
