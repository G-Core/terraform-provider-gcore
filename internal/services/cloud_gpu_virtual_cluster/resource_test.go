package cloud_gpu_virtual_cluster_test

import (
	"context"
	"fmt"
	"os"
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
