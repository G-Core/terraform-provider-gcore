package cloud_gpu_virtual_cluster_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
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

// gpuClusterPreCheck verifies GPU-cluster-specific environment variables are set.
func gpuClusterPreCheck(t *testing.T) {
	t.Helper()
	acctest.PreCheck(t)
	if v := os.Getenv("GCORE_GPU_VIRTUAL_CLUSTER_FLAVOR"); v == "" {
		t.Skip("GCORE_GPU_VIRTUAL_CLUSTER_FLAVOR must be set for GPU virtual cluster acceptance tests")
	}
	if v := os.Getenv("GCORE_GPU_VIRTUAL_CLUSTER_IMAGE_ID"); v == "" {
		t.Skip("GCORE_GPU_VIRTUAL_CLUSTER_IMAGE_ID must be set for GPU virtual cluster acceptance tests")
	}
	if v := os.Getenv("GCORE_GPU_VIRTUAL_CLUSTER_SSH_KEY"); v == "" {
		t.Skip("GCORE_GPU_VIRTUAL_CLUSTER_SSH_KEY must be set for GPU virtual cluster acceptance tests")
	}
}

// gpuFlavor returns the GPU cluster flavor from environment variable.
func gpuFlavor() string {
	return os.Getenv("GCORE_GPU_VIRTUAL_CLUSTER_FLAVOR")
}

// gpuImageID returns the GPU cluster image ID from environment variable.
func gpuImageID() string {
	return os.Getenv("GCORE_GPU_VIRTUAL_CLUSTER_IMAGE_ID")
}

// gpuSSHKey returns the SSH key name for GPU clusters from environment variable.
func gpuSSHKey() string {
	return os.Getenv("GCORE_GPU_VIRTUAL_CLUSTER_SSH_KEY")
}

func TestAccCloudGPUVirtualCluster_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUVirtualClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUVirtualClusterConfig(rName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("servers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("flavor"), knownvalue.StringExact(gpuFlavor())),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("servers_ids"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func TestAccCloudGPUVirtualCluster_update(t *testing.T) {
	rName := acctest.RandomName()
	newName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUVirtualClusterDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create
			{
				Config: testAccCloudGPUVirtualClusterConfig(rName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("servers_count"), knownvalue.Int64Exact(1)),
					compareIDSame.AddStateValue(
						"gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 2: Update name
			{
				Config: testAccCloudGPUVirtualClusterConfig(newName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(newName)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("servers_count"), knownvalue.Int64Exact(1)),
					// verify same resource (in-place update)
					compareIDSame.AddStateValue(
						"gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 3: Re-plan same config - guards against perpetual diffs on
			// volatile computed attributes (updated_at/status/has_pending_changes)
			// after an in-place update (GCLOUD2-28055).
			{
				Config:   testAccCloudGPUVirtualClusterConfig(newName, 1),
				PlanOnly: true,
			},
		},
	})
}

func TestAccCloudGPUVirtualCluster_resize(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUVirtualClusterDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with 1 server
			{
				Config: testAccCloudGPUVirtualClusterConfig(rName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("servers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("servers_ids"), knownvalue.ListSizeExact(1)),
					compareIDSame.AddStateValue(
						"gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 2: Resize to 2 servers
			{
				Config: testAccCloudGPUVirtualClusterConfig(rName, 2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("servers_count"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("servers_ids"), knownvalue.ListSizeExact(2)),
					// verify same resource (in-place resize, not replacement)
					compareIDSame.AddStateValue(
						"gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 3: Re-plan same config - guards against perpetual diffs on
			// servers_ids and volatile computed attributes after a resize
			// (GCLOUD2-28055).
			{
				Config:   testAccCloudGPUVirtualClusterConfig(rName, 2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccCloudGPUVirtualCluster_tags(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUVirtualClusterDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with tags
			{
				Config: testAccCloudGPUVirtualClusterConfigWithTags(rName, "env", "test"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("tags").AtMapKey("env"), knownvalue.StringExact("test")),
				},
			},
			// Step 2: Update tags
			{
				Config: testAccCloudGPUVirtualClusterConfigWithTags(rName, "env", "production"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("tags").AtMapKey("env"), knownvalue.StringExact("production")),
				},
			},
		},
	})
}

func TestAccCloudGPUVirtualCluster_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUVirtualClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUVirtualClusterConfig(rName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cloud_gpu_virtual_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_gpu_virtual_cluster.test", "project_id", "region_id", "id"),
				// credentials and user_data are write-only / no_refresh; source is no_refresh
				ImportStateVerifyIgnore: []string{
					"servers_settings.credentials",
					"servers_settings.user_data",
					"servers_settings.volumes.0.source",
				},
			},
		},
	})
}

// TestAccCloudGPUVirtualCluster_dependentDataSourceNoDrift is a regression test
// for GCLOUD2-28055: a data source referencing the cluster via depends_on forces
// Terraform to re-plan the resource even though its config is unchanged. Computed
// attributes without state-preserving plan modifiers (created_at, updated_at,
// status, has_pending_changes, servers_ids) then reverted to "(known after
// apply)", leaving a non-empty post-apply plan. Each test step implicitly
// asserts the post-apply plan is empty.
func TestAccCloudGPUVirtualCluster_dependentDataSourceNoDrift(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUVirtualClusterDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create the cluster alone
			{
				Config: testAccCloudGPUVirtualClusterConfig(rName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("status"), knownvalue.StringExact("active")),
				},
			},
			// Step 2: Same resource config plus a dependent data source. The
			// implicit post-apply empty-plan check fails if any computed
			// attribute reverts to unknown.
			{
				Config: testAccCloudGPUVirtualClusterConfigWithDataSource(rName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_virtual_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
		},
	})
}

// testAccCheckCloudGPUVirtualClusterDestroy verifies the GPU virtual cluster is deleted.
func testAccCheckCloudGPUVirtualClusterDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_gpu_virtual_cluster" {
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

		_, err = client.Cloud.GPUVirtual.Clusters.Get(context.Background(), rs.Primary.ID, cloud.GPUVirtualClusterGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("GPU virtual cluster %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking GPU virtual cluster deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudGPUVirtualClusterConfig(name string, serversCount int) string {
	return fmt.Sprintf(`
resource "gcore_cloud_gpu_virtual_cluster" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  name          = %[3]q
  flavor        = %[4]q
  servers_count = %[5]d

  servers_settings = {
    interfaces = [{
      type      = "external"
      ip_family = "ipv4"
    }]
    volumes = [{
      name       = "root-volume"
      size       = 50
      type       = "ssd_hiiops"
      source     = "image"
      image_id   = %[6]q
      boot_index = 0
    }]
    credentials = {
      ssh_key_name = %[7]q
    }
  }
}`, acctest.ProjectID(), acctest.RegionID(), name, gpuFlavor(), serversCount, gpuImageID(), gpuSSHKey())
}

func testAccCloudGPUVirtualClusterConfigWithDataSource(name string, serversCount int) string {
	return testAccCloudGPUVirtualClusterConfig(name, serversCount) + fmt.Sprintf(`

data "gcore_cloud_gpu_virtual_cluster" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  cluster_id = gcore_cloud_gpu_virtual_cluster.test.id

  depends_on = [gcore_cloud_gpu_virtual_cluster.test]
}`, acctest.ProjectID(), acctest.RegionID())
}

func testAccCloudGPUVirtualClusterConfigWithTags(name, tagKey, tagValue string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_gpu_virtual_cluster" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  name          = %[3]q
  flavor        = %[4]q
  servers_count = 1

  servers_settings = {
    interfaces = [{
      type      = "external"
      ip_family = "ipv4"
    }]
    volumes = [{
      name       = "root-volume"
      size       = 50
      type       = "ssd_hiiops"
      source     = "image"
      image_id   = %[5]q
      boot_index = 0
    }]
    credentials = {
      ssh_key_name = %[6]q
    }
  }

  tags = {
    %[7]s = %[8]q
  }
}`, acctest.ProjectID(), acctest.RegionID(), name, gpuFlavor(), gpuImageID(), gpuSSHKey(), tagKey, tagValue)
}

// Synthetic IDs for the plan-only validation cases. These must be literals, not
// references to real resources: an unresolved reference is unknown at plan time, and
// the validator deliberately does not judge unknown values.
const (
	fakeNetworkID = "11111111-1111-1111-1111-111111111111"
	fakeSubnetID  = "22222222-2222-2222-2222-222222222222"
	fakeImageID   = "33333333-3333-3333-3333-333333333333"
)

// TestAccCloudGPUVirtualCluster_interfaceVariants covers all three interface variants
// on one cluster with ip_family omitted throughout, then re-plans to prove none of them
// drift. This is the regression test for GCLOUD2-28207, where a "subnet" interface
// forced a full cluster replacement on every plan.
//
// The existing tests in this file only ever create "external" interfaces with ip_family
// set explicitly, which is why the bug shipped.
//
// If the API rejects mixing variants (or this interface count) on a single cluster,
// split this into one test per variant keeping the same assertions.
func TestAccCloudGPUVirtualCluster_interfaceVariants(t *testing.T) {
	rName := acctest.RandomName()
	const clusterAddr = "gcore_cloud_gpu_virtual_cluster.test"

	ifacePath := func(index int, attribute string) tfjsonpath.Path {
		return tfjsonpath.New("servers_settings").AtMapKey("interfaces").AtSliceIndex(index).AtMapKey(attribute)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUVirtualClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUVirtualClusterConfigInterfaceVariants(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(clusterAddr,
						tfjsonpath.New("status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue(clusterAddr,
						tfjsonpath.New("servers_settings").AtMapKey("interfaces"), knownvalue.ListSizeExact(3)),

					// external and any_subnet: the API returns ip_family, so state is populated
					// and the plan modifier's non-null path applies.
					statecheck.ExpectKnownValue(clusterAddr, ifacePath(0, "ip_family"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(clusterAddr, ifacePath(2, "ip_family"), knownvalue.NotNull()),
					// subnet: the API does not return ip_family, so state stays null. Without
					// StringUseStateForUnknownInclNull this re-plans as unknown on every plan.
					statecheck.ExpectKnownValue(clusterAddr, ifacePath(1, "ip_family"), knownvalue.Null()),

					// The computed default security group the ticket blamed. It lands on every
					// interface and must not, by itself, produce a diff.
					statecheck.ExpectKnownValue(clusterAddr, ifacePath(0, "security_groups"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(clusterAddr, ifacePath(1, "security_groups"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(clusterAddr, ifacePath(2, "security_groups"), knownvalue.NotNull()),
				},
			},
			{
				// Byte-identical config. ExpectEmptyPlan cannot go on the creation step, so the
				// no-op assertion needs a second step.
				Config: testAccCloudGPUVirtualClusterConfigInterfaceVariants(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(clusterAddr, plancheck.ResourceActionNoop),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// TestAccCloudGPUVirtualCluster_interfaceVariantValidation exercises the per-variant
// field contract through a real provider binary. Every case fails during plan, so no
// infrastructure is created and no GPU quota is consumed.
func TestAccCloudGPUVirtualCluster_interfaceVariantValidation(t *testing.T) {
	testCases := map[string]struct {
		iface     string
		wantError string
	}{
		"ip_family on subnet": { // GCLOUD2-28207
			iface: `{
        type       = "subnet"
        network_id = "` + fakeNetworkID + `"
        subnet_id  = "` + fakeSubnetID + `"
        ip_family  = "ipv4"
      }`,
			wantError: `'ip_family' is not supported when type = "subnet"`,
		},
		// Same rule, uppercase discriminator: the schema validates `type`
		// case-insensitively, so the validator must normalize before matching.
		"ip_family on uppercase SUBNET": {
			iface: `{
        type       = "SUBNET"
        network_id = "` + fakeNetworkID + `"
        subnet_id  = "` + fakeSubnetID + `"
        ip_family  = "ipv4"
      }`,
			wantError: `'ip_family' is not supported when type = "subnet"`,
		},
		"network_id on external": {
			iface: `{
        type       = "external"
        network_id = "` + fakeNetworkID + `"
      }`,
			wantError: `'network_id' is not supported when type = "external"`,
		},
		"subnet_id on external": {
			iface: `{
        type      = "external"
        subnet_id = "` + fakeSubnetID + `"
      }`,
			wantError: `'subnet_id' is not supported when type = "external"`,
		},
		"floating_ip on external": {
			iface: `{
        type        = "external"
        floating_ip = { source = "new" }
      }`,
			wantError: `'floating_ip' is not supported when type = "external"`,
		},
		"subnet_id on any_subnet": {
			iface: `{
        type       = "any_subnet"
        network_id = "` + fakeNetworkID + `"
        subnet_id  = "` + fakeSubnetID + `"
      }`,
			wantError: `'subnet_id' is not supported when type = "any_subnet"`,
		},
		"subnet missing network_id": {
			iface: `{
        type      = "subnet"
        subnet_id = "` + fakeSubnetID + `"
      }`,
			wantError: `'network_id' is required when type = "subnet"`,
		},
		"subnet missing subnet_id": {
			iface: `{
        type       = "subnet"
        network_id = "` + fakeNetworkID + `"
      }`,
			wantError: `'subnet_id' is required when type = "subnet"`,
		},
		"any_subnet missing network_id": {
			iface: `{
        type = "any_subnet"
      }`,
			wantError: `'network_id' is required when type = "any_subnet"`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:                 func() { acctest.PreCheck(t) },
				ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      testAccCloudGPUVirtualClusterConfigWithInterface(tc.iface),
						PlanOnly:    true,
						ExpectError: regexp.MustCompile(regexp.QuoteMeta(tc.wantError)),
					},
				},
			})
		})
	}
}

func testAccCloudGPUVirtualClusterConfigInterfaceVariants(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "pinned" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-pinned-net"
}

resource "gcore_cloud_network_subnet" "pinned" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-pinned-subnet"
  network_id = gcore_cloud_network.pinned.id
  cidr       = "192.168.90.0/24"
}

resource "gcore_cloud_network" "any" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-any-net"
}

resource "gcore_cloud_network_subnet" "any" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-any-subnet"
  network_id = gcore_cloud_network.any.id
  cidr       = "192.168.91.0/24"
}

resource "gcore_cloud_gpu_virtual_cluster" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  name          = %[3]q
  flavor        = %[4]q
  servers_count = 1

  servers_settings = {
    interfaces = [
      {
        type = "external"
      },
      {
        type       = "subnet"
        network_id = gcore_cloud_network.pinned.id
        subnet_id  = gcore_cloud_network_subnet.pinned.id
      },
      {
        type       = "any_subnet"
        network_id = gcore_cloud_network.any.id
      },
    ]
    volumes = [{
      name       = "root-volume"
      size       = 50
      type       = "ssd_hiiops"
      source     = "image"
      image_id   = %[5]q
      boot_index = 0
    }]
    credentials = {
      ssh_key_name = %[6]q
    }
  }

  # The any_subnet interface only references the network, so nothing would otherwise
  # order it after the subnet it needs to pick from.
  depends_on = [gcore_cloud_network_subnet.any]
}`, acctest.ProjectID(), acctest.RegionID(), name, gpuFlavor(), gpuImageID(), gpuSSHKey())
}

// testAccCloudGPUVirtualClusterConfigWithInterface builds a fully synthetic cluster
// carrying a single interface block. It never gets applied, so the flavor, image and
// key values only need to satisfy the other schema validators.
func testAccCloudGPUVirtualClusterConfigWithInterface(ifaceBlock string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_gpu_virtual_cluster" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  name          = "tf-test-variant-validation"
  flavor        = "tf-test-placeholder-flavor"
  servers_count = 1

  servers_settings = {
    interfaces = [
      %[3]s
    ]
    volumes = [{
      name       = "root-volume"
      size       = 50
      type       = "ssd_hiiops"
      source     = "image"
      image_id   = %[4]q
      boot_index = 0
    }]
    credentials = {
      ssh_key_name = "tf-test-placeholder-key"
    }
  }
}`, acctest.ProjectID(), acctest.RegionID(), ifaceBlock, fakeImageID)
}
