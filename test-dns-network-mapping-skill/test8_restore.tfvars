mapping = [
  {
    tags  = ["dev", "dc1"]
    cidr4 = ["10.0.0.0/24"]
    cidr6 = ["2001:db8:1::/48"]
  },
  {
    tags  = ["prod", "dc2"]
    cidr4 = ["10.1.0.0/24"]
  }
]
