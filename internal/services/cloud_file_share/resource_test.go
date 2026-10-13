package cloud_file_share_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/G-Core/gcore-go"
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

func TestAccCloudFileShare_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("protocol"), knownvalue.StringExact("NFS")),
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("type_name"), knownvalue.StringExact("standard")),
				},
			},
		},
	})
}

func TestAccCloudFileShare_update(t *testing.T) {
	rName := acctest.RandomName()
	newName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					compareIDSame.AddStateValue(
						"gcore_cloud_file_share.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudFileShareConfigUpdated(newName),
				// size is resized in place through the extend endpoint
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_file_share.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("name"), knownvalue.StringExact(newName)),
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(2)),
					compareIDSame.AddStateValue(
						"gcore_cloud_file_share.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

// TestAccCloudFileShare_resizeOnly changes only size. The plan must not
// send a PATCH: with nothing to rename or retag its body would be empty, and
// the API rejects an empty update.
func TestAccCloudFileShare_resizeOnly(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfigSize(rName, 1),
			},
			{
				Config: testAccCloudFileShareConfigSize(rName, 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_file_share.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(2)),
				},
			},
			{
				Config:   testAccCloudFileShareConfigSize(rName, 2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccCloudFileShare_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cloud_file_share.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_file_share.test", "project_id", "region_id", "id"),
			},
		},
	})
}

// TestAccCloudFileShare_importNoReplacement imports a file share whose
// configuration sets network and checks that the plan leaves it in place.
// Before the API returned network, import left it empty and planned a destroy
// and recreate that wiped the share's data.
func TestAccCloudFileShare_importNoReplacement(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfig(rName),
			},
			{
				ResourceName:      "gcore_cloud_file_share.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_file_share.test", "project_id", "region_id", "id"),
				ImportStateKind:   resource.ImportBlockWithID,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_file_share.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("gcore_cloud_file_share.test",
							tfjsonpath.New("network").AtMapKey("network_id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("gcore_cloud_file_share.test",
							tfjsonpath.New("network").AtMapKey("subnet_id"), knownvalue.NotNull()),
					},
				},
			},
			{
				// network is computed, so dropping the block keeps the share's
				// network and plans no change. A PlanOnly step fails on any diff.
				Config:   testAccCloudFileShareConfigNoNetwork(rName),
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudFileShare_autoSubnet creates a file share with only network_id
// set, so the API picks the subnet. The chosen subnet is read back into state
// and later plans stay empty.
func TestAccCloudFileShare_autoSubnet(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfigNetworkOnly(rName),
				Check:  testAccWaitForNetworkSubnet("gcore_cloud_network.test"),
			},
			{
				Config: testAccCloudFileShareConfigAutoSubnet(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						"gcore_cloud_file_share.test", tfjsonpath.New("network").AtMapKey("subnet_id"),
						"gcore_cloud_network_subnet.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
			{
				Config: testAccCloudFileShareConfigAutoSubnet(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				ResourceName:      "gcore_cloud_file_share.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_file_share.test", "project_id", "region_id", "id"),
				ImportStateKind:   resource.ImportBlockWithID,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_file_share.test", plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// TestAccCloudFileShare_changeNetwork moves a file share to another network.
// The update endpoint cannot change it, so the share must be replaced.
func TestAccCloudFileShare_changeNetwork(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfigTwoNetworks(rName, "test"),
			},
			{
				Config:             testAccCloudFileShareConfigTwoNetworks(rName, "other"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					// PlanOnly steps run their checks here, against the plan.
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_file_share.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

// TestAccCloudFileShare_vastWithoutNetwork creates a VAST share, which takes no
// network on create although the API returns one. The returned network must
// not make later plans replace the share.
func TestAccCloudFileShare_vastWithoutNetwork(t *testing.T) {
	regionID := os.Getenv("GCORE_FILE_SHARE_VAST_REGION_ID")
	if regionID == "" {
		t.Skip("GCORE_FILE_SHARE_VAST_REGION_ID must be set to a region with VAST file shares")
	}
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfigVast(rName, regionID),
				// The API assigns the VAST network itself. It must be read back,
				// or the empty plan below would pass for the wrong reason.
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("type_name"), knownvalue.StringExact("vast")),
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("network").AtMapKey("network_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("network").AtMapKey("subnet_id"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccCloudFileShareConfigVast(rName, regionID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAccCloudFileShare_vastSettingsResizeOnly changes only the size of a VAST
// share that sets some share_settings. The share_settings it leaves unset must
// keep their state values: if they plan as unknown, the update sends a PATCH
// with no name or tags, which the API rejects, and the resize never runs.
func TestAccCloudFileShare_vastSettingsResizeOnly(t *testing.T) {
	regionID := os.Getenv("GCORE_FILE_SHARE_VAST_REGION_ID")
	if regionID == "" {
		t.Skip("GCORE_FILE_SHARE_VAST_REGION_ID must be set to a region with VAST file shares")
	}
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfigVastSettings(rName, regionID, 1),
			},
			{
				Config: testAccCloudFileShareConfigVastSettings(rName, regionID, 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_file_share.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("gcore_cloud_file_share.test",
							tfjsonpath.New("share_settings").AtMapKey("allowed_characters"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue("gcore_cloud_file_share.test",
							tfjsonpath.New("share_settings").AtMapKey("path_length"), knownvalue.NotNull()),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(2)),
				},
			},
			{
				Config: testAccCloudFileShareConfigVastSettings(rName, regionID, 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAccCloudFileShare_ddnNoDrift checks that an unchanged DDN share plans no
// update, after both create and import. Any planned update would send an empty
// PATCH, which the API rejects.
func TestAccCloudFileShare_ddnNoDrift(t *testing.T) {
	regionID := os.Getenv("GCORE_FILE_SHARE_DDN_REGION_ID")
	if regionID == "" {
		t.Skip("GCORE_FILE_SHARE_DDN_REGION_ID must be set to a region with DDN file shares")
	}
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfigDDN(rName, regionID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("type_name"), knownvalue.StringExact("ddn")),
				},
			},
			{
				Config: testAccCloudFileShareConfigDDN(rName, regionID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				ResourceName:      "gcore_cloud_file_share.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_file_share.test", "project_id", "region_id", "id"),
				ImportStateKind:   resource.ImportBlockWithID,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_file_share.test", plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// TestAccCloudFileShare_networkRequiresStandard checks the config validator.
// Only standard shares take a network, and a configured network block still
// needs network_id although the schema makes it computed. Validation fails
// before any API call, so these plan-only steps create nothing.
func TestAccCloudFileShare_networkRequiresStandard(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudFileShareConfigVastWithNetwork(rName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`network can only be provided when type_name is 'standard'`),
			},
			{
				Config:      testAccCloudFileShareConfigNetworkWithoutID(rName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`network_id is required when network is provided`),
			},
		},
	})
}

// testAccCloudFileShareLegacyVersion is the last release where network is not
// read back from the API, so state it wrote holds a null network.subnet_id when
// the configuration omits it.
const testAccCloudFileShareLegacyVersion = "2.0.0-rc.2"

// TestAccCloudFileShare_upgradeAutoSubnet checks that state written by the
// legacy release, with network.subnet_id null, gets the subnet the API chose
// on refresh and plans no change.
//
// Run it without a CLI dev override for G-Core/gcore. An override would
// replace the legacy release with the local binary. The namespace must match
// the provider address the legacy release stores in state and the lock file,
// and setting it rules out a parallel test.
func TestAccCloudFileShare_upgradeAutoSubnet(t *testing.T) {
	t.Setenv(resource.EnvTfAccProviderNamespace, "g-core")
	rName := acctest.RandomName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		CheckDestroy: testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"gcore": {Source: "G-Core/gcore", VersionConstraint: testAccCloudFileShareLegacyVersion},
				},
				Config: testAccCloudFileShareConfigNetworkOnly(rName),
				Check:  testAccWaitForNetworkSubnet("gcore_cloud_network.test"),
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"gcore": {Source: "G-Core/gcore", VersionConstraint: testAccCloudFileShareLegacyVersion},
				},
				Config: testAccCloudFileShareConfigAutoSubnet(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("network").AtMapKey("subnet_id"), knownvalue.Null()),
				},
			},
			{
				ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
				Config:                   testAccCloudFileShareConfigAutoSubnet(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						"gcore_cloud_file_share.test", tfjsonpath.New("network").AtMapKey("subnet_id"),
						"gcore_cloud_network_subnet.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCheckCloudFileShareDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_file_share" {
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

		_, err = client.Cloud.FileShares.Get(context.Background(), rs.Primary.ID, cloud.FileShareGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("file share %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking file share deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudFileShareConfig(name string) string {
	return testAccCloudFileShareConfigSize(name, 1)
}

func testAccCloudFileShareConfigSize(name string, size int) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  network_id = gcore_cloud_network.test.id
  name       = "%[3]s-subnet"
  cidr       = "10.0.0.0/24"
}

resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "NFS"
  size       = %[4]d
  type_name  = "standard"
  network = {
    network_id = gcore_cloud_network.test.id
    subnet_id  = gcore_cloud_network_subnet.test.id
  }
}`, acctest.ProjectID(), acctest.RegionID(), name, size)
}

// testAccWaitForNetworkSubnet waits until the API lists a subnet on the
// network. Right after the subnet is created, a file share create that omits
// subnet_id can still fail with "No subnets found for network".
func testAccWaitForNetworkSubnet(resourceName string) resource.TestCheckFunc {
	return acctest.CheckResourceExists(resourceName, func(client *gcore.Client, id string) error {
		projectID, err := strconv.ParseInt(acctest.ProjectID(), 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing project_id: %w", err)
		}
		regionID, err := strconv.ParseInt(acctest.RegionID(), 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing region_id: %w", err)
		}

		deadline := time.Now().Add(2 * time.Minute)
		for {
			network, err := client.Cloud.Networks.Get(context.Background(), id, cloud.NetworkGetParams{
				ProjectID: param.NewOpt(projectID),
				RegionID:  param.NewOpt(regionID),
			})
			if err != nil {
				return fmt.Errorf("error getting network %s: %w", id, err)
			}
			if len(network.Subnets) > 0 {
				return nil
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("network %s has no subnets after 2 minutes", id)
			}
			time.Sleep(5 * time.Second)
		}
	})
}

// testAccCloudFileShareConfigNetworkOnly creates the network and subnet that
// testAccCloudFileShareConfigAutoSubnet uses, without the file share.
func testAccCloudFileShareConfigNetworkOnly(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  network_id = gcore_cloud_network.test.id
  name       = "%[3]s-subnet"
  cidr       = "10.0.0.0/24"
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

// testAccCloudFileShareConfigAutoSubnet omits subnet_id so the API picks the
// network's only subnet. depends_on makes sure that subnet exists first.
func testAccCloudFileShareConfigAutoSubnet(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  network_id = gcore_cloud_network.test.id
  name       = "%[3]s-subnet"
  cidr       = "10.0.0.0/24"
}

resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "NFS"
  size       = 1
  type_name  = "standard"
  network = {
    network_id = gcore_cloud_network.test.id
  }

  depends_on = [gcore_cloud_network_subnet.test]
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

// testAccCloudFileShareConfigNoNetwork keeps the network and subnet but drops
// the network block from the file share.
func testAccCloudFileShareConfigNoNetwork(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  network_id = gcore_cloud_network.test.id
  name       = "%[3]s-subnet"
  cidr       = "10.0.0.0/24"
}

resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "NFS"
  size       = 1
  type_name  = "standard"
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

// testAccCloudFileShareConfigTwoNetworks creates two networks with a subnet
// each and puts the file share on the one named by shareNetwork ("test" or
// "other"), so switching it changes only the file share.
func testAccCloudFileShareConfigTwoNetworks(name, shareNetwork string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  network_id = gcore_cloud_network.test.id
  name       = "%[3]s-subnet"
  cidr       = "10.0.0.0/24"
}

resource "gcore_cloud_network" "other" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-network-other"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "other" {
  project_id = %[1]s
  region_id  = %[2]s
  network_id = gcore_cloud_network.other.id
  name       = "%[3]s-subnet-other"
  cidr       = "10.0.1.0/24"
}

resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "NFS"
  size       = 1
  type_name  = "standard"
  network = {
    network_id = gcore_cloud_network.%[4]s.id
    subnet_id  = gcore_cloud_network_subnet.%[4]s.id
  }
}`, acctest.ProjectID(), acctest.RegionID(), name, shareNetwork)
}

// testAccCloudFileShareConfigVastWithNetwork sets network on a VAST share,
// which the validator rejects. The network ID is never looked up.
func testAccCloudFileShareConfigVastWithNetwork(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "NFS"
  size       = 1
  type_name  = "vast"
  network = {
    network_id = "00000000-0000-4000-8000-000000000000"
  }
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

// testAccCloudFileShareConfigNetworkWithoutID sets an empty network block on a
// standard share, which the validator rejects.
func testAccCloudFileShareConfigNetworkWithoutID(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "NFS"
  size       = 1
  type_name  = "standard"
  network    = {}
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

// testAccCloudFileShareConfigVast creates a VAST share, which takes no network.
func testAccCloudFileShareConfigVast(name, regionID string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "NFS"
  size       = 1
  type_name  = "vast"
}`, acctest.ProjectID(), regionID, name)
}

// testAccCloudFileShareConfigVastSettings creates a VAST share that sets only
// some share_settings and leaves the rest to the API.
func testAccCloudFileShareConfigVastSettings(name, regionID string, size int) string {
	return fmt.Sprintf(`
resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "NFS"
  size       = %[4]d
  type_name  = "vast"
  share_settings = {
    root_squash  = false
    trash_access = true
  }
}`, acctest.ProjectID(), regionID, name, size)
}

// testAccCloudFileShareConfigDDN creates a DDN share, which takes no network
// and no share_settings.
func testAccCloudFileShareConfigDDN(name, regionID string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "LUSTRE"
  size       = 1
  type_name  = "ddn"
}`, acctest.ProjectID(), regionID, name)
}

func testAccCloudFileShareConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  network_id = gcore_cloud_network.test.id
  name       = "%[3]s-subnet"
  cidr       = "10.0.0.0/24"
}

resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "NFS"
  size       = 2
  type_name  = "standard"
  network = {
    network_id = gcore_cloud_network.test.id
    subnet_id  = gcore_cloud_network_subnet.test.id
  }
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
