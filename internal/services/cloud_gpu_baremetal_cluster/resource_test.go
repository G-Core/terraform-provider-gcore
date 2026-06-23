package cloud_gpu_baremetal_cluster_test

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

func gpuBmClusterPreCheck(t *testing.T) {
	t.Helper()
	acctest.PreCheck(t)
	if v := os.Getenv("GCORE_GPU_BAREMETAL_CLUSTER_FLAVOR"); v == "" {
		t.Skip("GCORE_GPU_BAREMETAL_CLUSTER_FLAVOR must be set for GPU baremetal cluster acceptance tests")
	}
	if v := os.Getenv("GCORE_GPU_BAREMETAL_CLUSTER_IMAGE_ID"); v == "" {
		t.Skip("GCORE_GPU_BAREMETAL_CLUSTER_IMAGE_ID must be set for GPU baremetal cluster acceptance tests")
	}
	if v := os.Getenv("GCORE_GPU_BAREMETAL_CLUSTER_SSH_KEY"); v == "" {
		t.Skip("GCORE_GPU_BAREMETAL_CLUSTER_SSH_KEY must be set for GPU baremetal cluster acceptance tests")
	}
}

func gpuBmFlavor() string {
	return os.Getenv("GCORE_GPU_BAREMETAL_CLUSTER_FLAVOR")
}

func gpuBmImageID() string {
	return os.Getenv("GCORE_GPU_BAREMETAL_CLUSTER_IMAGE_ID")
}

func gpuBmSSHKey() string {
	return os.Getenv("GCORE_GPU_BAREMETAL_CLUSTER_SSH_KEY")
}

func testAccCheckCloudGPUBaremetalClusterDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_gpu_baremetal_cluster" {
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

		_, err = client.Cloud.GPUBaremetal.Clusters.Get(context.Background(), rs.Primary.ID, cloud.GPUBaremetalClusterGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("GPU baremetal cluster %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking GPU baremetal cluster deletion: %w", err)
		}
	}
	return nil
}

func TestAccCloudGPUBaremetalCluster_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuBmClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUBaremetalClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUBaremetalClusterConfig(rName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("servers_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("flavor"), knownvalue.StringExact(gpuBmFlavor())),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("servers_ids"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func TestAccCloudGPUBaremetalCluster_update(t *testing.T) {
	rName := acctest.RandomName()
	newName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuBmClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUBaremetalClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUBaremetalClusterConfig(rName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("servers_count"), knownvalue.Int64Exact(1)),
					compareIDSame.AddStateValue(
						"gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudGPUBaremetalClusterConfig(newName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(newName)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("servers_count"), knownvalue.Int64Exact(1)),
					compareIDSame.AddStateValue(
						"gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func TestAccCloudGPUBaremetalCluster_tags(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuBmClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUBaremetalClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUBaremetalClusterConfigWithTags(rName, "env", "test"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("tags").AtMapKey("env"), knownvalue.StringExact("test")),
				},
			},
			{
				Config: testAccCloudGPUBaremetalClusterConfigWithTags(rName, "env", "production"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("tags").AtMapKey("env"), knownvalue.StringExact("production")),
				},
			},
		},
	})
}

func TestAccCloudGPUBaremetalCluster_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuBmClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUBaremetalClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUBaremetalClusterConfig(rName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cloud_gpu_baremetal_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_gpu_baremetal_cluster.test", "project_id", "region_id", "id"),
				ImportStateVerifyIgnore: []string{
					"servers_settings.credentials",
					"servers_settings.user_data",
				},
			},
		},
	})
}

func testAccCloudGPUBaremetalClusterConfig(name string, serversCount int) string {
	return fmt.Sprintf(`
resource "gcore_cloud_gpu_baremetal_cluster" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  name          = %[3]q
  flavor        = %[4]q
  image_id      = %[5]q
  servers_count = %[6]d

  servers_settings = {
    interfaces = [{
      type      = "external"
      ip_family = "ipv4"
    }]
    credentials = {
      ssh_key_name = %[7]q
    }
  }
}`, acctest.ProjectID(), acctest.RegionID(), name, gpuBmFlavor(), gpuBmImageID(), serversCount, gpuBmSSHKey())
}

func testAccCloudGPUBaremetalClusterConfigWithTags(name, tagKey, tagValue string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_gpu_baremetal_cluster" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  name          = %[3]q
  flavor        = %[4]q
  image_id      = %[5]q
  servers_count = 1

  servers_settings = {
    interfaces = [{
      type      = "external"
      ip_family = "ipv4"
    }]
    credentials = {
      ssh_key_name = %[6]q
    }
  }

  tags = {
    %[7]s = %[8]q
  }
}`, acctest.ProjectID(), acctest.RegionID(), name, gpuBmFlavor(), gpuBmImageID(), gpuBmSSHKey(), tagKey, tagValue)
}
