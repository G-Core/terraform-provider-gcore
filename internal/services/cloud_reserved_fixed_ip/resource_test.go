package cloud_reserved_fixed_ip_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func testAccCheckCloudReservedFixedIPDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_reserved_fixed_ip" {
			continue
		}

		projectID, err := strconv.ParseInt(rs.Primary.Attributes["project_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing project_id: %w", err)
		}
		regionID, err := strconv.ParseInt(rs.Primary.Attributes["region_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing region_id: %w", err)
		}

		_, err = client.Cloud.ReservedFixedIPs.Get(context.Background(), rs.Primary.Attributes["port_id"], cloud.ReservedFixedIPGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("reserved fixed IP %s still exists", rs.Primary.Attributes["port_id"])
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking reserved fixed IP deletion: %w", err)
		}
	}
	return nil
}

func TestAccCloudReservedFixedIP_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudReservedFixedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudReservedFixedIPConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("port_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("is_external"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("fixed_ip_address"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("name"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCloudReservedFixedIP_update(t *testing.T) {
	comparePortIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudReservedFixedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudReservedFixedIPConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("port_id"), knownvalue.NotNull()),
					comparePortIDSame.AddStateValue(
						"gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("port_id"),
					),
				},
			},
			{
				Config: testAccCloudReservedFixedIPConfigIsVIP(true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("is_vip"), knownvalue.Bool(true)),
					comparePortIDSame.AddStateValue(
						"gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("port_id"),
					),
				},
			},
			// Set is_vip back to false so the resource can be destroyed
			{
				Config: testAccCloudReservedFixedIPConfigIsVIP(false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("is_vip"), knownvalue.Bool(false)),
					comparePortIDSame.AddStateValue(
						"gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("port_id"),
					),
				},
			},
		},
	})
}

func TestAccCloudReservedFixedIP_import(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudReservedFixedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudReservedFixedIPConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("port_id"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:      "gcore_cloud_reserved_fixed_ip.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_reserved_fixed_ip.test", "project_id", "region_id", "port_id"),
				ImportStateKind:   resource.ImportBlockWithID,
			},
		},
	})
}

func testAccCloudReservedFixedIPConfig() string {
	return fmt.Sprintf(`
resource "gcore_cloud_reserved_fixed_ip" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  type       = "external"
}`, acctest.ProjectID(), acctest.RegionID())
}

func testAccCloudReservedFixedIPConfigIsVIP(isVIP bool) string {
	return fmt.Sprintf(`
resource "gcore_cloud_reserved_fixed_ip" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  type       = "external"
  is_vip     = %[3]t
}`, acctest.ProjectID(), acctest.RegionID(), isVIP)
}
