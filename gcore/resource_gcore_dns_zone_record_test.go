//go:build !cloud
// +build !cloud

package gcore

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDnsZoneRecord(t *testing.T) {
	random := time.Now().Nanosecond()
	zone := "kokizzu.neuroops.link"
	rrSetName := strings.ReplaceAll(zone, ".", "_")
	subDomain := fmt.Sprintf("key%d", random)
	name := fmt.Sprintf("%s_%s", subDomain, rrSetName)
	fullDomain := subDomain + "." + zone

	resourceName := fmt.Sprintf("%s.%s", DNSZoneRecordResource, name)

	templateCreate := func() string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  zone = "%s"
  domain = "%s"
  type = "TXT"
  ttl = 210

  filter {
    type = "geodistance"
    strict = true
  }
  
  filter {
    type = "first_n"
    limit = 1
    strict = false
  }

  resource_record {
    content  = "1234"
    enabled = true
    
    meta {
      latlong = [52.367,4.9041]
	  asn = [12345]
	  ip = ["1.1.1.1"]
	  notes = "notes"
	  continents = ["europe"]
	  countries = ["pl"]
	  default = true
  	}
  }
}
		`, DNSZoneRecordResource, name, zone, fullDomain)
	}
	fmt.Println(templateCreate())

	templateUpdate := func() string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  zone = "%s"
  domain = "%s"
  type = "TXT"
  ttl = 120

  resource_record {
    content  = "12345"
    
    meta {
      latlong = [52.367,4.9041]
	  ip = ["1.1.2.2"]
	  notes = "notes"
	  continents = ["na"]
	  countries = ["us"]
	  default = false
  	}
  }
}
		`, DNSZoneRecordResource, name, zone, fullDomain)
	}
	fmt.Println(templateUpdate())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_PERMANENT_TOKEN_VAR, GCORE_DNS_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaDomain, fullDomain),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaType, "TXT"),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaTTL, "210"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaFilter, DNSZoneRecordSchemaFilterType),
						"geodistance"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaFilter, DNSZoneRecordSchemaFilterStrict),
						"true"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.1.%s", DNSZoneRecordSchemaFilter, DNSZoneRecordSchemaFilterType),
						"first_n"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.1.%s", DNSZoneRecordSchemaFilter, DNSZoneRecordSchemaFilterLimit),
						"1"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.1.%s", DNSZoneRecordSchemaFilter, DNSZoneRecordSchemaFilterStrict),
						"false"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaContent),
						"1234"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaEnabled),
						"true"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaLatLong),
						"52.367"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.1",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaLatLong),
						"4.9041"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaAsn),
						"12345"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaIP),
						"1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaNotes),
						"notes"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaContinents),
						"europe"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaCountries),
						"pl"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaDefault),
						"true"),
				),
			},
			{
				Config: templateUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaDomain, fullDomain),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaType, "TXT"),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaTTL, "120"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaContent),
						"12345"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaLatLong),
						"52.367"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.1",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaLatLong),
						"4.9041"),
					resource.TestCheckNoResourceAttr(resourceName, fmt.Sprintf("%s.0.%s.0.%s.0",
						DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaAsn)),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaIP),
						"1.1.2.2"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaNotes),
						"notes"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaContinents),
						"na"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaCountries),
						"us"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaDefault),
						"false"),
				),
			},
		},
	})
}

// note: when testing, set GCORE_DNS_API=https://api.gcore.com/dns
func TestAccDnsZoneRecordSvcbHttps(t *testing.T) {
	random := time.Now().Nanosecond()
	zone := "kokizzu.neuroops.link"
	subDomain := fmt.Sprintf("key%d", random)
	name := strings.ReplaceAll(fmt.Sprintf("%s_%s", subDomain, zone), `.`, `_`)
	fullDomain := subDomain + "." + zone

	resourceName := fmt.Sprintf("%s.%s", DNSZoneRecordResource, name)

	content := `1 . alpn="h3,h2" port=1443 ipv4hint=10.0.0.1 ech=AEn+DQBFKwAgACABWIHUGj4u+PIggYXcR5JF0gYk3dCRioBW8uJq9H4mKAAIAAEAAQABAANAEnB1YmxpYy50bHMtZWNoLmRldgAA`

	templateCreate := func() string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  zone = "%s"
  domain = "%s"
  type = "HTTPS"
  ttl = 120

  resource_record {
    content = <<EOT
%s
EOT
    enabled = true
  }
}
		`, DNSZoneRecordResource, name, zone, fullDomain, content)
	}
	fmt.Println(templateCreate())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_PERMANENT_TOKEN_VAR, GCORE_DNS_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaDomain, fullDomain),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaType, "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaTTL, "120"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaContent), content+"\n"),
				),
			},
		},
	})
}

// note: when testing, set GCORE_DNS_API=https://api.gcore.com/dns
func TestAccDnsZoneRecordRRSetMeta(t *testing.T) {
	// checks for rrset.meta failover/healthchecks and geodns_link

	random := time.Now().Nanosecond()
	const zone = "kokizzu.neuroops.link"
	subDomain := fmt.Sprintf("key%d", random)
	name := strings.ReplaceAll(fmt.Sprintf("%s_%s", subDomain, zone), `.`, `_`)
	fullDomain := subDomain + "." + zone

	resourceName := fmt.Sprintf("%s.%s", DNSZoneRecordResource, name)

	content := `127.0.0.1`

	templateCreate := func() string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  zone = "%s"
  domain = "%s"
  type = "A"
  ttl = 120

  resource_record {
    content = "%s"
    enabled = true
  }

  filter {
    type = "is_healthy"
    limit = 0
    strict = false
  }

  meta {
    healthchecks {
      frequency = 300
      host = "%s"
      http_status_code = 200
      method = "GET"
      port = 80
      protocol = "HTTP"
      regexp = ""
      timeout = 10
      tls = false
      url = "/"
    }
    geodns_link = "test.com"
  }
}
		`, DNSZoneRecordResource, name, zone, fullDomain, content, zone)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_PERMANENT_TOKEN_VAR, GCORE_DNS_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaDomain, fullDomain),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaType, "A"),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaTTL, "120"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaContent), content),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaHealthchecks, DNSZoneRRSetSchemaMetaFailoverFrequency), "300"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaHealthchecks, DNSZoneRRSetSchemaMetaFailoverHost), zone),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaHealthchecks, DNSZoneRRSetSchemaMetaFailoverHTTPStatusCode), "200"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaHealthchecks, DNSZoneRRSetSchemaMetaFailoverMethod), "GET"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaHealthchecks, DNSZoneRRSetSchemaMetaFailoverPort), "80"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaHealthchecks, DNSZoneRRSetSchemaMetaFailoverProtocol), "HTTP"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaHealthchecks, DNSZoneRRSetSchemaMetaFailoverRegexp), ""),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaHealthchecks, DNSZoneRRSetSchemaMetaFailoverTimeout), "10"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaHealthchecks, DNSZoneRRSetSchemaMetaFailoverTLS), "false"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaHealthchecks, DNSZoneRRSetSchemaMetaFailoverURL), "/"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaGeodnsLink), "test.com"),
				),
			},
		},
	})
}

// note: when testing, set GCORE_DNS_API=https://api.gcore.com/dns
func TestAccDnsZoneRecordRRSetMetaFailover(t *testing.T) {
	// checks for rrset.meta failover field (in addition to healthchecks)

	random := time.Now().Nanosecond()
	const zone = "kokizzu.neuroops.link"
	subDomain := fmt.Sprintf("key-failover-%d", random)
	name := strings.ReplaceAll(fmt.Sprintf("%s_%s", subDomain, zone), `.`, `_`)
	fullDomain := subDomain + "." + zone

	resourceName := fmt.Sprintf("%s.%s", DNSZoneRecordResource, name)

	content := `127.0.0.1`

	templateCreate := func() string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  zone = "%s"
  domain = "%s"
  type = "A"
  ttl = 120

  resource_record {
    content = "%s"
    enabled = true
  }

  filter {
    type = "is_healthy"
    limit = 0
    strict = false
  }

  meta {
    failover {
      frequency = 300
      host = "%s"
      http_status_code = 200
      method = "GET"
      port = 80
      protocol = "HTTP"
      regexp = ""
      timeout = 10
      tls = false
      url = "/"
    }
    geodns_link = "test.com"
  }
}
		`, DNSZoneRecordResource, name, zone, fullDomain, content, zone)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_PERMANENT_TOKEN_VAR, GCORE_DNS_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaDomain, fullDomain),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaType, "A"),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaTTL, "120"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaContent), content),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaFailover, DNSZoneRRSetSchemaMetaFailoverFrequency), "300"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaFailover, DNSZoneRRSetSchemaMetaFailoverHost), zone),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaFailover, DNSZoneRRSetSchemaMetaFailoverHTTPStatusCode), "200"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaFailover, DNSZoneRRSetSchemaMetaFailoverMethod), "GET"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaFailover, DNSZoneRRSetSchemaMetaFailoverPort), "80"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaFailover, DNSZoneRRSetSchemaMetaFailoverProtocol), "HTTP"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaFailover, DNSZoneRRSetSchemaMetaFailoverRegexp), ""),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaFailover, DNSZoneRRSetSchemaMetaFailoverTimeout), "10"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaFailover, DNSZoneRRSetSchemaMetaFailoverTLS), "false"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaFailover, DNSZoneRRSetSchemaMetaFailoverURL), "/"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s",
							DNSZoneRRSetSchemaMeta, DNSZoneRRSetSchemaMetaGeodnsLink), "test.com"),
				),
			},
		},
	})
}

// note: when testing, set GCORE_DNS_API=https://api.gcore.com/dns
func TestAccDnsZoneRecordNetworkMapping(t *testing.T) {
	random := time.Now().Nanosecond()
	zone := "kokizzu.neuroops.link"
	subDomain := fmt.Sprintf("key-nm-%d", random)
	fullDomain := subDomain + "." + zone
	mappingName := fmt.Sprintf("test-mapping-%d", random)

	const mappingResource = "gcore_dns_network_mapping"
	mappingResourceName := fmt.Sprintf("%s.network_mapping", mappingResource)
	recordResourceName := fmt.Sprintf("%s.rrset_example", DNSZoneRecordResource)

	template := func(val1, val2 int) string {
		return fmt.Sprintf(`
resource "%s" "network_mapping" {
  name = "%s"
  mapping {
    tags = ["development", "test"]
    cidr4 = ["10.0.0.0/16", "10.1.0.0/16"]
    cidr6 = ["fd00::/8"]
  }
}

resource "%s" "rrset_example" {
  zone   = "%s"
  domain = "%s"
  type   = "CNAME"
  ttl    = 120

  filter {
    type = "network_mapping"
    strict = false
  }

  meta {
    cidr_mapping = %s.network_mapping.name
  }

  resource_record {
    content = "ezt.cdb-staging.cdn.orange.com"
    meta {
      cidr_labels = {
        (%s.network_mapping.mapping[0].tags[0]) = %d
        (%s.network_mapping.mapping[0].tags[1]) = %d
      }
    }
    enabled = true
  }
}
`, mappingResource, mappingName, DNSZoneRecordResource, zone, fullDomain, mappingResource, mappingResource, val1, mappingResource, val2)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_PERMANENT_TOKEN_VAR, GCORE_DNS_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: template(10, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(recordResourceName),
					testAccCheckResourceExists(mappingResourceName),
					resource.TestCheckResourceAttr(recordResourceName, "filter.0.type", "network_mapping"),
					resource.TestCheckResourceAttr(recordResourceName, "meta.0.cidr_mapping", mappingName),
					resource.TestCheckResourceAttr(recordResourceName, "resource_record.0.meta.0.cidr_labels.development", "10"),
					resource.TestCheckResourceAttr(recordResourceName, "resource_record.0.meta.0.cidr_labels.test", "20"),
				),
			},
			{
				Config: template(50, 75),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(recordResourceName),
					testAccCheckResourceExists(mappingResourceName),
					resource.TestCheckResourceAttr(recordResourceName, "resource_record.0.meta.0.cidr_labels.development", "50"),
					resource.TestCheckResourceAttr(recordResourceName, "resource_record.0.meta.0.cidr_labels.test", "75"),
				),
			},
		},
	})
}
