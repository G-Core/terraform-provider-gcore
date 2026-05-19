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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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

// TestAccCloudReservedFixedIP_importSubnet verifies that importing a
// subnet-type reserved fixed IP produces a clean plan with no drift on
// type, subnet_id, network_id, or is_vip. Uses is_vip=false in config
// to also cover the is_vip field population during import.
func TestAccCloudReservedFixedIP_importSubnet(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudReservedFixedIPDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with type=subnet, is_vip=false
			{
				Config: testAccCloudReservedFixedIPConfigSubnetIsVip(rName, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("port_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("type"), knownvalue.StringExact("subnet")),
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("is_vip"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("subnet_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("network_id"), knownvalue.NotNull()),
				},
			},
			// Step 2: Import. Expect clean plan — no drift on type,
			// subnet_id, network_id, or is_vip.
			{
				ResourceName:      "gcore_cloud_reserved_fixed_ip.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_reserved_fixed_ip.test", "project_id", "region_id", "port_id"),
				ImportStateKind:   resource.ImportBlockWithID,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_reserved_fixed_ip.test", plancheck.ResourceActionNoop),
						// type must be inferred as "subnet" from the API response
						plancheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
							tfjsonpath.New("type"), knownvalue.StringExact("subnet")),
						// subnet_id must be preserved, not "known after apply"
						plancheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
							tfjsonpath.New("subnet_id"), knownvalue.NotNull()),
						// network_id must be preserved
						plancheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
							tfjsonpath.New("network_id"), knownvalue.NotNull()),
						// is_vip must be populated from API (false), not null.
						// Guards against the GCLOUD2-24268 2026-04-15 regression
						// where is_vip was stored as null after import and
						// triggered a spurious PATCH.
						plancheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
							tfjsonpath.New("is_vip"), knownvalue.Bool(false)),
					},
				},
			},
		},
	})
}

// TestAccCloudReservedFixedIP_importAnySubnet verifies that importing an
// any_subnet-type reserved fixed IP:
//   - triggers an in-place update (not replacement) because the API cannot
//     distinguish "subnet" from "any_subnet" and import infers "subnet"
//   - preserves subnet_id and network_id (no drift on these fields)
//   - produces a fully clean plan after one reconciliation apply
func TestAccCloudReservedFixedIP_importAnySubnet(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudReservedFixedIPDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with type=any_subnet, only network_id in config
			{
				Config: testAccCloudReservedFixedIPConfigAnySubnet(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("port_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("type"), knownvalue.StringExact("any_subnet")),
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("subnet_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("network_id"), knownvalue.NotNull()),
				},
			},
			// Step 2: Import. The API infers type="subnet" but the config
			// says "any_subnet", so the plan shows an in-place update.
			// subnet_id and network_id must be preserved (not "known after apply").
			{
				ResourceName:       "gcore_cloud_reserved_fixed_ip.test",
				ImportState:        true,
				ImportStateIdFunc:  acctest.BuildImportID("gcore_cloud_reserved_fixed_ip.test", "project_id", "region_id", "port_id"),
				ImportStateKind:    resource.ImportBlockWithID,
				ExpectNonEmptyPlan: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Update (not Replace) — the type plan modifier must
						// treat subnet/any_subnet as equivalent.
						plancheck.ExpectResourceAction("gcore_cloud_reserved_fixed_ip.test", plancheck.ResourceActionUpdate),
						// Config value "any_subnet" should be adopted into
						// the plan without triggering replacement.
						plancheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
							tfjsonpath.New("type"), knownvalue.StringExact("any_subnet")),
						plancheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
							tfjsonpath.New("subnet_id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
							tfjsonpath.New("network_id"), knownvalue.NotNull()),
						// is_vip must be populated from API (false), not null.
						plancheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
							tfjsonpath.New("is_vip"), knownvalue.Bool(false)),
					},
				},
			},
			// Step 3: After reconciliation apply, plan must be fully clean.
			{
				Config: testAccCloudReservedFixedIPConfigAnySubnet(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// TestAccCloudReservedFixedIP_importExternal verifies that importing an
// external-type reserved fixed IP produces a fully clean plan (no drift,
// no forced replacement). This guards against the original GCLOUD2-24268
// regression where `type` was lost after import and `RequiresReplace`
// caused destroy+recreate.
func TestAccCloudReservedFixedIP_importExternal(t *testing.T) {
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
					statecheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("type"), knownvalue.StringExact("external")),
				},
			},
			// Import must be a noop: type inferred as "external", is_vip
			// populated from the API, no replacement.
			{
				ResourceName:      "gcore_cloud_reserved_fixed_ip.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_reserved_fixed_ip.test", "project_id", "region_id", "port_id"),
				ImportStateKind:   resource.ImportBlockWithID,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_reserved_fixed_ip.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
							tfjsonpath.New("type"), knownvalue.StringExact("external")),
						plancheck.ExpectKnownValue("gcore_cloud_reserved_fixed_ip.test",
							tfjsonpath.New("is_vip"), knownvalue.Bool(false)),
					},
				},
			},
			// After import, re-apply the config and confirm the plan is
			// fully empty — catches any lingering drift that the import
			// plan check might miss.
			{
				Config: testAccCloudReservedFixedIPConfig(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func testAccCloudReservedFixedIPConfigSubnetIsVip(name string, isVip bool) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "tf-test-net-%[3]s"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = "tf-test-subnet-%[3]s"
  network_id  = gcore_cloud_network.test.id
  cidr        = "10.0.0.0/24"
  enable_dhcp = true
}

resource "gcore_cloud_reserved_fixed_ip" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  type       = "subnet"
  subnet_id  = gcore_cloud_network_subnet.test.id
  is_vip     = %[4]t
}`, acctest.ProjectID(), acctest.RegionID(), name, isVip)
}

func testAccCloudReservedFixedIPConfigAnySubnet(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "tf-test-net-%[3]s"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = "tf-test-subnet-%[3]s"
  network_id  = gcore_cloud_network.test.id
  cidr        = "10.0.0.0/24"
  enable_dhcp = true
}

resource "gcore_cloud_reserved_fixed_ip" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  type       = "any_subnet"
  network_id = gcore_cloud_network.test.id

  depends_on = [gcore_cloud_network_subnet.test]
}`, acctest.ProjectID(), acctest.RegionID(), name)
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
