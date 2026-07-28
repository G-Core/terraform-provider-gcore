package cloud_k8s_cluster_pool_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// k8sClusterVersion/k8sClusterFlavorID mirror the constants used by the
// cloud_k8s_cluster package's own tests; Go test helpers are not importable
// across packages, so the cluster resource config below is duplicated here
// rather than shared.
const (
	k8sClusterVersion  = "v1.33.10"
	k8sClusterFlavorID = "g1-standard-2-4"
)

// Pools do not exist without a cluster, so this test creates a cluster (with
// its default pool) and then reads the list data source (deferred via
// depends_on) to assert the pool is returned.
func TestAccCloudK8SClusterPoolsDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudK8SClusterPoolClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudK8SClusterPoolsDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_cluster_pools.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_cluster_pools.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("default-pool")),
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_cluster_pools.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_cluster_pools.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("status"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccCheckCloudK8SClusterPoolClusterDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_k8s_cluster" {
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

		clusterName := rs.Primary.Attributes["name"]

		_, err = client.Cloud.K8S.Clusters.Get(
			context.Background(),
			clusterName,
			cloud.K8SClusterGetParams{
				ProjectID: param.NewOpt(projectID),
				RegionID:  param.NewOpt(regionID),
			},
		)
		if err == nil {
			return fmt.Errorf("k8s cluster %s still exists", clusterName)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking k8s cluster deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudK8SClusterPoolsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_ssh_key" "test" {
  project_id = %[1]s
  name       = %[3]q
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFo9Yv8MkEsmt8mUm4gMRpOuIiLPkJfSdmR3lKA3GsOG tf-test@example.com"
}

resource "gcore_cloud_k8s_cluster" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  keypair    = gcore_cloud_ssh_key.test.name
  version    = %[4]q

  pools = [
    {
      name               = "default-pool"
      flavor_id          = %[5]q
      min_node_count     = 1
      max_node_count     = 2
      boot_volume_size   = 50
      boot_volume_type   = "standard"
      servergroup_policy = "soft-anti-affinity"
      is_public_ipv4     = true
    }
  ]
}

data "gcore_cloud_k8s_cluster_pools" "test" {
  project_id   = %[1]s
  region_id    = %[2]s
  cluster_name = gcore_cloud_k8s_cluster.test.name
  depends_on   = [gcore_cloud_k8s_cluster.test]
}
`, acctest.ProjectID(), acctest.RegionID(), name, k8sClusterVersion, k8sClusterFlavorID)
}
