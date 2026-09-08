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

// TestAccCloudInstance_importNoReplace guards the regression QA hit on the first
// pass of this fix: ssh_key_name is response-backed but not computed, so
// UnmarshalComputed never wrote it and it stayed null in state after an import.
// With the unconditional RequiresReplace() it carried, a configured key then
// planned a destroy+recreate of a live instance:
//
//	# gcore_cloud_instance.test must be replaced
//	    + ssh_key_name = "..." # forces replacement
//
// ImportState now resolves it, so the plan after import is an ordinary in-place
// update. It is deliberately NOT asserted empty: import intentionally omits
// volumes[].boot_index, so a config that sets it always leaves that one diff.
// The assertion is therefore only that nothing gets destroyed.
//
// The existing TestAccCloudInstance_basic import step cannot catch this:
// ImportStateVerify only compares attributes and never runs a plan, so
// "forces replacement" is invisible to it. Only an import-block step plans.
func TestAccCloudInstance_importNoReplace(t *testing.T) {
	rName := acctest.RandomName()
	imageID := latestUbuntuImageID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfigWithSSHKey(rName, imageID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("ssh_key_name"), knownvalue.StringExact(rName+"-key")),
				},
			},
			{
				Config:            testAccCloudInstanceConfigWithSSHKey(rName, imageID),
				ResourceName:      "gcore_cloud_instance.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_instance.test", "project_id", "region_id", "id"),
				// boot_index is not imported by design, so the plan is non-empty.
				ExpectNonEmptyPlan: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// The invariant, stated without naming the action: whatever
						// else the plan does, it must not destroy the instance.
						acctest.ExpectResourceNotReplaced("gcore_cloud_instance.test"),
					},
				},
			},
		},
	})
}

// testAccCloudInstanceConfigWithSSHKey pairs the instance with a keypair of its
// own, so the import path exercises the keypair UUID -> name lookup rather than
// depending on a key that happens to exist in the test project.
func testAccCloudInstanceConfigWithSSHKey(name, imageID string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_ssh_key" "test" {
  project_id = %[1]s
  name       = "%[3]s-key"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAuZsAIHRmkZ8q2fwq3itCamNcD+kgb/iggHW+B6xxND %[3]s"
}

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
  project_id   = %[1]s
  region_id    = %[2]s
  name         = %[3]q
  flavor       = "g1-standard-1-2"
  ssh_key_name = gcore_cloud_ssh_key.test.name

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

// Reproduces the failure seen on the apply that follows a terraform import.
//
// project_id/region_id are Optional-only path params, so importing (which
// writes both concretely from the import ID) and then applying a config that
// omits them leaves state and config disagreeing on exactly those two
// attributes. That is a plain in-place update that none of Update()'s
// handlers act on, and the same diff is produced deterministically by moving
// from a config that sets them to one that does not - no import step needed,
// which keeps the test off ImportStatePersist's "resource already managed by
// Terraform" path.
//
// On such a diff every bare Computed attribute is planned unknown, and the
// unconditional post-update refresh is the only thing that resolves them.
// Without it the apply fails with "Provider returned invalid result object
// after apply".
func TestAccCloudInstance_pathOnlyDiffResolvesComputed(t *testing.T) {
	rName := acctest.RandomName()
	imageID := latestUbuntuImageID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			// State ends up with concrete project_id/region_id, as an import would leave it.
			{Config: testAccCloudInstanceConfig(rName, imageID)},
			{
				Config: testAccCloudInstanceConfigNoProjectRegion(rName, imageID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Guards against this step going vacuous: dropping
						// project_id/region_id must still produce an in-place
						// update, not a no-op and not a replacement.
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

// The adoption-window tests below share one instance shape: no create-only
// attribute set, and - unlike testAccCloudInstanceConfig - no boot_index on
// the volume. Import intentionally omits volumes[].boot_index, so a config
// that sets it can never produce an empty post-import plan, and these tests
// need one: the no-op apply between import and the day-two edit is the whole
// point. (If real infra turns out to reject a create without boot_index,
// swap the no-op apply step for a `RefreshState: true` step - it persists the
// same marker advance - and give the create config its boot_index back.)

// testAccCloudInstanceAdoptionVolumeBlock is the boot volume shared by every
// step of the adoption-window tests, kept identical across them so the volume
// itself is never replaced mid-test.
func testAccCloudInstanceAdoptionVolumeBlock(name, imageID string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "boot" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-vol"
  size       = 10
  type_name  = "ssd_hiiops"
  source     = "image"
  image_id   = %[4]q
}`, acctest.ProjectID(), acctest.RegionID(), name, imageID)
}

// testAccCloudInstanceConfigNoCreateOnly declares the instance with none of
// the create-only attributes, so an import writes the adoption marker but the
// following applies have nothing to adopt.
func testAccCloudInstanceConfigNoCreateOnly(name, imageID string) string {
	return testAccCloudInstanceAdoptionVolumeBlock(name, imageID) + fmt.Sprintf(`

resource "gcore_cloud_instance" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  flavor     = "g1-standard-1-2"

  volumes = [
    {
      volume_id = gcore_cloud_volume.boot.id
    }
  ]

  interfaces = [
    {
      type = "external"
    }
  ]
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

// testAccCloudInstanceConfigNoCreateOnlyForgotten keeps the boot volume of
// testAccCloudInstanceConfigNoCreateOnly and replaces the instance block with
// a removed block (Terraform >= 1.7): its state entry goes away, the live
// instance stays, and a later step can legally import it again. Necessary
// because `terraform import` refuses an address that is already in state.
func testAccCloudInstanceConfigNoCreateOnlyForgotten(name, imageID string) string {
	return testAccCloudInstanceAdoptionVolumeBlock(name, imageID) + `

removed {
  from = gcore_cloud_instance.test

  lifecycle {
    destroy = false
  }
}`
}

// testAccCloudInstanceConfigAdoptsCreateOnly is the no-create-only config plus
// user_data and username - the two attributes whose only legitimate paths are
// "adopted right after import" and "replacement". The values never reach the
// API on the adoption path, so they only have to be well-formed.
func testAccCloudInstanceConfigAdoptsCreateOnly(name, imageID string) string {
	return testAccCloudInstanceAdoptionVolumeBlock(name, imageID) + fmt.Sprintf(`

resource "gcore_cloud_instance" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  flavor     = "g1-standard-1-2"
  user_data  = "Zm9v"
  username   = "gcore-qa"

  volumes = [
    {
      volume_id = gcore_cloud_volume.boot.id
    }
  ]

  interfaces = [
    {
      type = "external"
    }
  ]
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

// TestAccCloudInstance_importAdoptionApplies covers the half of the window the
// import-block test cannot reach: the apply that actually adopts. Command-mode
// import with ImportStatePersist puts the imported state into the working
// state, so later steps really plan and apply against it.
//
// The standalone plan-only step in the middle is deliberate: standalone plans
// never persist private state, so any number of them must leave the window
// open for the apply that follows. A lifecycle that retired the marker on
// plans would pass every other step and fail exactly there.
func TestAccCloudInstance_importAdoptionApplies(t *testing.T) {
	rName := acctest.RandomName()
	imageID := latestUbuntuImageID(t)
	var instanceID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfigNoCreateOnly(rName, imageID),
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_instance.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_instance.test not found in state")
					}
					instanceID = rs.Primary.ID
					return nil
				},
			},
			{
				// Forget the instance without destroying it, so the import
				// below targets an address absent from state.
				Config: testAccCloudInstanceConfigNoCreateOnlyForgotten(rName, imageID),
			},
			{
				Config:             testAccCloudInstanceConfigNoCreateOnly(rName, imageID),
				ResourceName:       "gcore_cloud_instance.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), instanceID), nil
				},
			},
			{
				// A standalone plan of the adopting config: in-place update
				// planned, window not consumed - plans persist nothing.
				Config:             testAccCloudInstanceConfigAdoptsCreateOnly(rName, imageID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectResourceNotReplaced("gcore_cloud_instance.test"),
					},
				},
			},
			{
				// The adopting apply: user_data/username land in state without
				// the instance being replaced, and the marker is retired.
				Config: testAccCloudInstanceConfigAdoptsCreateOnly(rName, imageID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_instance.test",
							plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("user_data"), knownvalue.StringExact("Zm9v")),
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("username"), knownvalue.StringExact("gcore-qa")),
				},
			},
			{
				// Adoption is stable: with the values in state, the same
				// config plans nothing at all.
				Config:   testAccCloudInstanceConfigAdoptsCreateOnly(rName, imageID),
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudInstance_importThenNoopApplyThenSetForcesReplace is the reported
// stale-marker defect end to end: import with a config that sets no create-only
// attribute, run a no-op apply, then set user_data on day two. The no-op apply
// consumed the adoption window - Terraform persists its refresh even though
// nothing changed - so the day-two set must plan the replacement that actually
// delivers the value to the instance, not a silent in-place adoption that only
// rewrites state.
//
// Before the marker lifecycle this planned an in-place update: nothing ever
// retired the marker, because a config with nothing to adopt never runs Update.
func TestAccCloudInstance_importThenNoopApplyThenSetForcesReplace(t *testing.T) {
	rName := acctest.RandomName()
	imageID := latestUbuntuImageID(t)
	var instanceID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfigNoCreateOnly(rName, imageID),
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_instance.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_instance.test not found in state")
					}
					instanceID = rs.Primary.ID
					return nil
				},
			},
			{
				// Forget the instance without destroying it (Terraform >= 1.7).
				Config: testAccCloudInstanceConfigNoCreateOnlyForgotten(rName, imageID),
			},
			{
				Config:             testAccCloudInstanceConfigNoCreateOnly(rName, imageID),
				ResourceName:       "gcore_cloud_instance.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), instanceID), nil
				},
			},
			{
				// The no-op apply. Its refresh is persisted - marker advances
				// armed -> spent - and the framework's implicit post-apply
				// empty-plan check doubles as proof the step really was a
				// no-op, which is what makes the window consumption silent.
				Config: testAccCloudInstanceConfigNoCreateOnly(rName, imageID),
			},
			{
				// Day two: the window is spent, so setting user_data must plan
				// a replacement - the only action that applies the value.
				// Plan-only: actually replacing the instance proves nothing
				// more and burns minutes of real infra.
				Config:             testAccCloudInstanceConfigAdoptsCreateOnly(rName, imageID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_instance.test",
							plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}
