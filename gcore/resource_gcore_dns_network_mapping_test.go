package gcore

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDNSNetworkMapping(t *testing.T) {
	random := time.Now().Nanosecond()
	mappingName := fmt.Sprintf("test-mapping-%d", random)
	resourceName := fmt.Sprintf("%s.%s", DNSNetworkMappingResource, "test")

	template := func(name, mappings string) string {
		return fmt.Sprintf(`
resource "%s" "test" {
  name = "%s"
  %s
}
		`, DNSNetworkMappingResource, name, mappings)
	}

	mappingsCreate := `
  mapping {
    tags  = ["dev", "dc1"]
    cidr4 = ["10.0.0.0/24"]
    cidr6 = ["2001:db8:1::/48"]
  }
  mapping {
    tags  = ["prod", "dc2"]
    cidr4 = ["10.1.0.0/24"]
  }
`
	mappingsUpdate := `
  mapping {
    tags  = ["dev", "dc1", "updated"]
    cidr4 = ["10.0.0.0/24", "10.0.1.0/24"]
    cidr6 = ["2001:db8:2::/48"]
  }
`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// PreCheck logic from other tests if needed, e.g., testAccPreCheckVars
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Network Mapping creation
				Config: template(mappingName, mappingsCreate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaName, mappingName),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".#", "2"),
					// Check first mapping block
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".0."+DNSNetworkMappingSchemaTags+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".0."+DNSNetworkMappingSchemaTags+".0", "dev"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".0."+DNSNetworkMappingSchemaCidr4+".#", "1"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".0."+DNSNetworkMappingSchemaCidr4+".0", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".0."+DNSNetworkMappingSchemaCidr6+".#", "1"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".0."+DNSNetworkMappingSchemaCidr6+".0", "2001:db8:1::/48"),
					// Check second mapping block
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".1."+DNSNetworkMappingSchemaTags+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".1."+DNSNetworkMappingSchemaTags+".0", "prod"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".1."+DNSNetworkMappingSchemaCidr4+".#", "1"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".1."+DNSNetworkMappingSchemaCidr4+".0", "10.1.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".1."+DNSNetworkMappingSchemaCidr6+".#", "0"),
				),
			},
			{
				// Test Network Mapping update
				Config: template(mappingName, mappingsUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaName, mappingName),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".#", "1"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".0."+DNSNetworkMappingSchemaTags+".#", "3"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".0."+DNSNetworkMappingSchemaTags+".2", "updated"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".0."+DNSNetworkMappingSchemaCidr4+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".0."+DNSNetworkMappingSchemaCidr4+".1", "10.0.1.0/24"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".0."+DNSNetworkMappingSchemaCidr6+".#", "1"),
					resource.TestCheckResourceAttr(resourceName, DNSNetworkMappingSchemaMapping+".0."+DNSNetworkMappingSchemaCidr6+".0", "2001:db8:2::/48"),
				),
			},
		},
	})
}
