package cloud_gpu_virtual_cluster_volume_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// gpuClusterPreCheck verifies GPU-cluster-specific environment variables are
// set, mirroring cloud_gpu_virtual_cluster's own precheck (not importable
// across packages). Real GPU virtual server creation is expensive and slow,
// so this test only runs when a human opts in with these env vars.
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

func gpuFlavor() string  { return os.Getenv("GCORE_GPU_VIRTUAL_CLUSTER_FLAVOR") }
func gpuImageID() string { return os.Getenv("GCORE_GPU_VIRTUAL_CLUSTER_IMAGE_ID") }
func gpuSSHKey() string  { return os.Getenv("GCORE_GPU_VIRTUAL_CLUSTER_SSH_KEY") }

// Volumes do not exist without a cluster, so this test creates one (its
// servers_settings always includes a root-volume) and then reads the list
// data source (deferred via depends_on) to assert the volume is returned.
//
// Known issue (GCLOUD2-28055): run against real infra, this currently fails
// its implicit post-apply empty-plan check -- the gcore_cloud_gpu_virtual_cluster
// resource leaves several computed fields (created_at, status, updated_at,
// has_pending_changes, servers_ids) without a stable plan once this data
// source's depends_on forces a re-plan. The list data source itself works.
// Confirmed both real clusters created while diagnosing this were torn down
// cleanly. Left enabled but env-gated (see gpuClusterPreCheck below), matching
// this repo's existing convention, so it does not affect default CI runs.
func TestAccCloudGPUVirtualClusterVolumesDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUVirtualClusterVolumesDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_virtual_cluster_volumes.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_virtual_cluster_volumes.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("root-volume")),
				},
			},
		},
	})
}

func testAccCloudGPUVirtualClusterVolumesDataSourceConfig(name string) string {
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
}

data "gcore_cloud_gpu_virtual_cluster_volumes" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  cluster_id = gcore_cloud_gpu_virtual_cluster.test.id
  depends_on = [gcore_cloud_gpu_virtual_cluster.test]
}
`, acctest.ProjectID(), acctest.RegionID(), name, gpuFlavor(), gpuImageID(), gpuSSHKey())
}
