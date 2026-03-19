//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func checkSubnetAttrs(resourceName string, opts *subnets.CreateOpts) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s.Empty() == true {
			return fmt.Errorf("State not updated")
		}

		checksStore := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceName, "name", opts.Name),
			resource.TestCheckResourceAttr(resourceName, "cidr", opts.CIDR.String()),
			resource.TestCheckResourceAttr(resourceName, "enable_dhcp", strconv.FormatBool(*opts.EnableDHCP)),
			resource.TestCheckResourceAttr(resourceName, "dns_nameservers.#", strconv.Itoa(len(opts.DNSNameservers))),
			resource.TestCheckResourceAttr(resourceName, "host_routes.#", strconv.Itoa(len(opts.HostRoutes))),
		}

		if opts.GatewayIP == nil && opts.EnableDHCP != nil && !*opts.EnableDHCP {
			checksStore = append(checksStore, resource.TestCheckResourceAttr(resourceName, "gateway_ip", "disable"))
		}

		for i, hr := range opts.HostRoutes {
			checksStore = append(checksStore,
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`host_routes.%d.destination`, i), hr.Destination.String()),
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`host_routes.%d.nexthop`, i), hr.NextHop.String()),
			)
		}

		return resource.ComposeTestCheckFunc(checksStore...)(s)
	}
}

func TestAccSubnet(t *testing.T) {

	var dst1, dst2, cidr gcorecloud.CIDR

	_, netIPNet, _ := net.ParseCIDR("10.0.3.0/24")
	dst1.IP = netIPNet.IP
	dst1.Mask = netIPNet.Mask

	_, netIPNet, _ = net.ParseCIDR("10.0.4.0/24")
	dst2.IP = netIPNet.IP
	dst2.Mask = netIPNet.Mask

	_, netIPNet, _ = net.ParseCIDR("192.168.10.0/24")
	cidr.IP = netIPNet.IP
	cidr.Mask = netIPNet.Mask

	enableDHCPTrue := true
	enableDHCPFalse := false

	createFixt := subnets.CreateOpts{
		Name:           "create_subnet",
		CIDR:           cidr,
		DNSNameservers: []net.IP{net.ParseIP("8.8.4.4"), net.ParseIP("1.1.1.1")},
		EnableDHCP:     &enableDHCPTrue,
		HostRoutes: []subnets.HostRoute{
			{
				Destination: dst1,
				NextHop:     net.ParseIP("192.168.10.1"),
			},
			{
				Destination: dst2,
				NextHop:     net.ParseIP("192.168.10.1"),
			},
		},
	}

	gateway := net.ParseIP("disable")

	updateFixt := subnets.CreateOpts{
		Name:           "update_subnet",
		CIDR:           cidr,
		DNSNameservers: make([]net.IP, 0),
		EnableDHCP:     &enableDHCPFalse,
		HostRoutes:     make([]subnets.HostRoute, 0),
		GatewayIP:      &gateway,
	}

	type Params struct {
		Name        string
		CIDR        string
		DNS         []string
		HRoutes     []map[string]string
		DHCP        string
		Gateway     string
		MetadataMap string
	}

	create := Params{
		Name: "create_subnet",
		CIDR: "192.168.10.0/24",
		DNS:  []string{"8.8.4.4", "1.1.1.1"},
		HRoutes: []map[string]string{
			{"destination": "10.0.3.0/24", "nexthop": "192.168.10.1"},
			{"destination": "10.0.4.0/24", "nexthop": "192.168.10.1"},
		},
		MetadataMap: `{
				key1 = "val1"
				key2 = "val2"
		}`,
	}

	update := Params{
		Name:    "update_subnet",
		CIDR:    "192.168.10.0/24",
		DHCP:    "false",
		DNS:     []string{},
		HRoutes: []map[string]string{},
		Gateway: "disable",
		MetadataMap: `{
				key3 = "val3"
	  	}`,
	}

	SubnetTemplate := func(params *Params) string {
		template := `
		locals {
	    	dns_nameservers = [`

		for i := range params.DNS {
			template += fmt.Sprintf(`"%s",`, params.DNS[i])
		}

		template += fmt.Sprint(`]
			host_routes = [`)

		for i := range params.HRoutes {
			template += fmt.Sprintf(`
			{
              destination = "%s"
              nexthop = "%s"
			},`, params.HRoutes[i]["destination"], params.HRoutes[i]["nexthop"])
		}

		template += fmt.Sprintf(`]
        	}

		resource "gcore_network" "acctest" {
			name = "create_network"
  			mtu = 1450
  			type = "vxlan"
			create_router = false
			%[1]s
			%[2]s
		}

		resource "gcore_subnet" "acctest" {
			name = "%s"
  			cidr = "%s"
  			network_id = gcore_network.acctest.id
			dns_nameservers = local.dns_nameservers
			connect_to_network_router = false
            dynamic host_routes {
				iterator = hr
				for_each = local.host_routes
				  content {
					destination = hr.value.destination
					nexthop = hr.value.nexthop
				  }
			  }	
			metadata_map = %s
            %[1]s
			%[2]s
	

		`, regionInfo(), projectInfo(), params.Name, params.CIDR, params.MetadataMap)

		if params.DHCP != "" {
			template += fmt.Sprintf("enable_dhcp = %s\n", params.DHCP)
		}

		if params.Gateway != "" {
			template += fmt.Sprintf(`gateway_ip = "%s"`, params.Gateway)
		}

		return template + "\n}"
	}

	fullName := "gcore_subnet.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: SubnetTemplate(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					checkSubnetAttrs(fullName, &createFixt),
					testAccCheckMetadata(fullName, true, map[string]interface{}{
						"key1": "val1",
						"key2": "val2",
					}),
				),
			},
			{
				Config: SubnetTemplate(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					checkSubnetAttrs(fullName, &updateFixt),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key3": "val3",
					}),
					testAccCheckMetadata(fullName, false, map[string]string{
						"key1": "val1",
					}),
					testAccCheckMetadata(fullName, false, map[string]string{
						"key2": "val2",
					}),
				),
			},
		},
	})
}

func testAccSubnetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, subnetPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_subnet" {
			continue
		}

		_, err := subnets.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}
