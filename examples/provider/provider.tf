terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = "2.0.0-alpha.1"
    }
  }
}

# Configure the Gcore provider.
# The API key can also be set via the GCORE_API_KEY environment variable.
provider "gcore" {
  api_key = "251$d3361.............1b35f26d8"
}
