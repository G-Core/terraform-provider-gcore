package cloud_instance_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func testAccCheckCloudInstanceDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_instance" {
			continue
		}

		// project_id/region_id are optional path params: a config that omits
		// them leaves the attributes null in state and the provider falls back
		// to GCORE_CLOUD_PROJECT_ID/GCORE_CLOUD_REGION_ID. Do the same here.
		projectAttr := rs.Primary.Attributes["project_id"]
		if projectAttr == "" {
			projectAttr = acctest.ProjectID()
		}
		projectID, err := strconv.ParseInt(projectAttr, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing project_id: %w", err)
		}
		regionAttr := rs.Primary.Attributes["region_id"]
		if regionAttr == "" {
			regionAttr = acctest.RegionID()
		}
		regionID, err := strconv.ParseInt(regionAttr, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing region_id: %w", err)
		}

		_, err = client.Cloud.Instances.Get(
			context.Background(),
			rs.Primary.ID,
			cloud.InstanceGetParams{
				ProjectID: param.NewOpt(projectID),
				RegionID:  param.NewOpt(regionID),
			},
		)
		if err == nil {
			return fmt.Errorf("cloud instance %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking instance deletion: %w", err)
		}
	}
	return nil
}

func TestAccCloudInstance_basic(t *testing.T) {
	rName := acctest.RandomName()
	imageID := latestUbuntuImageID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfig(rName, imageID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("status"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:      "gcore_cloud_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password_wo",
					"password_wo_version",
					"allow_app_ports",
					"user_data",
					"volumes",
				},
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_instance.test", "project_id", "region_id", "id"),
			},
		},
	})
}

func TestAccCloudInstance_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()
	imageID := latestUbuntuImageID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfig(rName, imageID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				Config: testAccCloudInstanceConfig(rNameUpdated, imageID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
				},
			},
		},
	})
}

func testAccCloudInstanceConfig(name, imageID string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "boot" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-vol"
  size       = 10
  type_name  = "ssd_hiiops"
  source     = "image"
  image_id   = %[4]q
}

resource "gcore_cloud_instance" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  flavor     = "g1-standard-1-2"

  volumes = [
    {
      volume_id  = gcore_cloud_volume.boot.id
      boot_index = 0
    }
  ]

  interfaces = [
    {
      type = "external"
    }
  ]
}`, acctest.ProjectID(), acctest.RegionID(), name, imageID)
}

// Reproduces the import-then-apply failure: importing sets concrete
// project_id/region_id, the config omits them, so the next apply invokes
// Update with a diff no handler tracks; every bare Computed attribute is
// planned unknown and must be resolved by the unconditional refresh.
func TestAccCloudInstance_importThenApply(t *testing.T) {
	rName := acctest.RandomName()
	imageID := latestUbuntuImageID(t)
	config := testAccCloudInstanceConfigNoProjectRegion(rName, imageID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			{Config: config},
			{
				ResourceName:       "gcore_cloud_instance.test",
				ImportState:        true,
				ImportStatePersist: true,
				// The imported state deliberately differs from the config's
				// state: project_id/region_id come back concrete, the config
				// leaves them null. That difference is the point of the test.
				ImportStateVerify: false,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["gcore_cloud_instance.test"]
					if !ok {
						return "", fmt.Errorf("resource not found in state")
					}
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), rs.Primary.ID), nil
				},
			},
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Guards against this step going vacuous: the imported
						// project_id/region_id must produce an in-place update.
						plancheck.ExpectResourceAction("gcore_cloud_instance.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("addresses"), knownvalue.NotNull()),
				},
			},
		},
	})
}

// Same as testAccCloudInstanceConfig but the instance block omits
// project_id/region_id so the provider falls back to the environment. The
// volume keeps them so its own lifecycle is unaffected.
func testAccCloudInstanceConfigNoProjectRegion(name, imageID string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "boot" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-vol"
  size       = 10
  type_name  = "ssd_hiiops"
  source     = "image"
  image_id   = %[4]q
}

resource "gcore_cloud_instance" "test" {
  name   = %[3]q
  flavor = "g1-standard-1-2"

  volumes = [
    {
      volume_id  = gcore_cloud_volume.boot.id
      boot_index = 0
    }
  ]

  interfaces = [
    {
      type = "external"
    }
  ]
}`, acctest.ProjectID(), acctest.RegionID(), name, imageID)
}

func TestAccCloudInstance_tags(t *testing.T) {
	rName := acctest.RandomName()
	imageID := latestUbuntuImageID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfigWithTags(rName, imageID, map[string]string{
					"env":  "test",
					"team": "platform",
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("tags"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"env":  knownvalue.StringExact("test"),
							"team": knownvalue.StringExact("platform"),
						})),
				},
			},
			{
				Config: testAccCloudInstanceConfigWithTags(rName, imageID, map[string]string{
					"env":     "staging",
					"team":    "platform",
					"version": "2",
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("tags"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"env":     knownvalue.StringExact("staging"),
							"team":    knownvalue.StringExact("platform"),
							"version": knownvalue.StringExact("2"),
						})),
				},
			},
			// Step 3: remove tags (verify JSON Merge Patch sends null for removed keys)
			{
				Config: testAccCloudInstanceConfigWithTags(rName, imageID, map[string]string{
					"env": "staging",
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("tags"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"env": knownvalue.StringExact("staging"),
						})),
				},
			},
		},
	})
}

func testAccCloudInstanceConfigWithTags(name, imageID string, tags map[string]string) string {
	tagLines := ""
	for k, v := range tags {
		tagLines += fmt.Sprintf("    %s = %q\n", k, v)
	}
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "boot" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-vol"
  size       = 10
  type_name  = "ssd_hiiops"
  source     = "image"
  image_id   = %[4]q
}

resource "gcore_cloud_instance" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  flavor     = "g1-standard-1-2"

  volumes = [
    {
      volume_id  = gcore_cloud_volume.boot.id
      boot_index = 0
    }
  ]

  interfaces = [
    {
      type = "external"
    }
  ]

  tags = {
%[5]s  }
}`, acctest.ProjectID(), acctest.RegionID(), name, imageID, tagLines)
}

// TestAccCloudInstance_addInterfaceOnly appends a second subnet interface to an existing
// instance without changing the name or tags. The appended interface's computed port_id
// and ip_address must be planned as unknown ("known after apply"), then resolve to real
// values post-apply, and a re-plan of the identical config must be empty.
func TestAccCloudInstance_addInterfaceOnly(t *testing.T) {
	rName := acctest.RandomName()
	imageID := latestUbuntuImageID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			// Step 1: single subnet interface.
			{
				Config: testAccCloudInstanceConfigInterfaces(rName, rName, imageID, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
				Check: testAccCheckInstanceInterfaceCount("gcore_cloud_instance.test", 1),
			},
			// Step 2: same name, append a second subnet interface.
			{
				Config: testAccCloudInstanceConfigInterfaces(rName, rName, imageID, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("gcore_cloud_instance.test",
							tfjsonpath.New("interfaces").AtSliceIndex(1).AtMapKey("port_id")),
						plancheck.ExpectUnknownValue("gcore_cloud_instance.test",
							tfjsonpath.New("interfaces").AtSliceIndex(1).AtMapKey("ip_address")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("gcore_cloud_instance.test", "name", rName),
					resource.TestCheckResourceAttrSet("gcore_cloud_instance.test", "interfaces.1.port_id"),
					resource.TestCheckResourceAttrSet("gcore_cloud_instance.test", "interfaces.1.ip_address"),
					testAccCheckInstanceInterfaceCount("gcore_cloud_instance.test", 2),
				),
			},
			// Step 3: identical config must produce an empty plan.
			{
				Config: testAccCloudInstanceConfigInterfaces(rName, rName, imageID, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// TestAccCloudInstance_addInterfaceWithNameChange appends a second subnet interface while
// also renaming the instance in the same step. The rename exercises the PATCH path while
// the appended interface's computed fields must still plan as unknown and resolve cleanly.
func TestAccCloudInstance_addInterfaceWithNameChange(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := rName + "-upd"
	imageID := latestUbuntuImageID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			// Step 1: single subnet interface.
			{
				Config: testAccCloudInstanceConfigInterfaces(rName, rName, imageID, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
				Check: testAccCheckInstanceInterfaceCount("gcore_cloud_instance.test", 1),
			},
			// Step 2: rename AND append a second subnet interface.
			{
				Config: testAccCloudInstanceConfigInterfaces(rName, rNameUpdated, imageID, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("gcore_cloud_instance.test",
							tfjsonpath.New("interfaces").AtSliceIndex(1).AtMapKey("port_id")),
						plancheck.ExpectUnknownValue("gcore_cloud_instance.test",
							tfjsonpath.New("interfaces").AtSliceIndex(1).AtMapKey("ip_address")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("gcore_cloud_instance.test", "name", rNameUpdated),
					resource.TestCheckResourceAttrSet("gcore_cloud_instance.test", "interfaces.1.port_id"),
					resource.TestCheckResourceAttrSet("gcore_cloud_instance.test", "interfaces.1.ip_address"),
					testAccCheckInstanceInterfaceCount("gcore_cloud_instance.test", 2),
				),
			},
			// Step 3: identical config must produce an empty plan.
			{
				Config: testAccCloudInstanceConfigInterfaces(rName, rNameUpdated, imageID, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// testAccCheckInstanceInterfaceCount asserts, via the API, that the instance has exactly
// the expected number of interfaces. This guards against a silent duplicate-attach where
// the state looks right but the API attached more (or fewer) interfaces than intended.
func testAccCheckInstanceInterfaceCount(resourceName string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		client, err := acctest.NewGcoreClient()
		if err != nil {
			return err
		}

		projectID, err := strconv.ParseInt(rs.Primary.Attributes["project_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing project_id: %w", err)
		}
		regionID, err := strconv.ParseInt(rs.Primary.Attributes["region_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing region_id: %w", err)
		}

		page, err := client.Cloud.Instances.Interfaces.List(
			context.Background(),
			rs.Primary.ID,
			cloud.InstanceInterfaceListParams{
				ProjectID: param.NewOpt(projectID),
				RegionID:  param.NewOpt(regionID),
			},
		)
		if err != nil {
			return fmt.Errorf("listing instance interfaces: %w", err)
		}

		if got := len(page.Results); got != expected {
			return fmt.Errorf("expected %d interfaces on instance %s, got %d", expected, rs.Primary.ID, got)
		}
		return nil
	}
}

// testAccCloudInstanceConfigInterfaces builds an instance backed by a boot volume and two
// networks with one subnet each. The instance always has a subnet interface on network 1;
// when secondInterface is true a second subnet interface on network 2 is appended. The
// second interface deliberately lives on a SEPARATE network: attaching an interface to a
// subnet of a network the instance is already connected to silently no-ops on the API side
// (the attach task finishes without error and without creating a port). prefix names the
// supporting resources (kept stable across steps) while instanceName sets the instance
// name (varied to exercise renames).
func testAccCloudInstanceConfigInterfaces(prefix, instanceName, imageID string, secondInterface bool) string {
	secondIfaceBlock := ""
	if secondInterface {
		secondIfaceBlock = `,
    {
      type       = "subnet"
      network_id = gcore_cloud_network.test2.id
      subnet_id  = gcore_cloud_network_subnet.test2.id
    }`
	}

	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-net"
}

resource "gcore_cloud_network" "test2" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-net2"
}

resource "gcore_cloud_network_subnet" "test1" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-subnet1"
  network_id = gcore_cloud_network.test.id
  cidr       = "192.168.10.0/24"
}

resource "gcore_cloud_network_subnet" "test2" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-subnet2"
  network_id = gcore_cloud_network.test2.id
  cidr       = "192.168.20.0/24"
}

resource "gcore_cloud_volume" "boot" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-vol"
  size       = 10
  type_name  = "ssd_hiiops"
  source     = "image"
  image_id   = %[5]q
}

resource "gcore_cloud_instance" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[4]q
  flavor     = "g1-standard-1-2"

  volumes = [
    {
      volume_id  = gcore_cloud_volume.boot.id
      boot_index = 0
    }
  ]

  interfaces = [
    {
      type       = "subnet"
      network_id = gcore_cloud_network.test.id
      subnet_id  = gcore_cloud_network_subnet.test1.id
    }%[6]s
  ]
}`, acctest.ProjectID(), acctest.RegionID(), prefix, instanceName, imageID, secondIfaceBlock)
}
