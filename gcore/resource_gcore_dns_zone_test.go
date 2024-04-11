//go:build !cloud
// +build !cloud

package gcore

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDnsZone(t *testing.T) {

	random := time.Now().Nanosecond()
	name := fmt.Sprintf("terraformtestkey%d", random)
	zone := name + ".com"
	zoneRenamed := strconv.Itoa(random) + "." + zone
	resourceName := fmt.Sprintf("%s.%s", DNSZoneResource, name)

	template := func(zoneName string, enableDNSSec bool) string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
  dnssec = %t	
}
		`, DNSZoneResource, name, zoneName, enableDNSSec)
	}

	templateNoDnssec := func(zoneName string) string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
}
		`, DNSZoneResource, name, zoneName)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_DNS_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test DNS Zone creation
				Config: templateNoDnssec(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSZoneSchemaName, zone),
				),
			},
			{
				// Test change not affecting DNS Zone, since dnssec by default disabled
				Config: template(zone, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSZoneSchemaName, zone),
				),
			},
			{
				// Test change affecting DNS Zone, since dnssec enabled
				Config: template(zone, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSZoneSchemaName, zone),
				),
			},
			{
				// Test DNS Zone deleted and recreated with new name
				Config: template(zoneRenamed, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSZoneSchemaName, zoneRenamed),
				),
			},
		},
	})

}
