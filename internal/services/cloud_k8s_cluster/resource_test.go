package cloud_k8s_cluster_test

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

// k8sClusterVersion is the Kubernetes version used in tests.
// It should be a supported version in the target region.
const k8sClusterVersion = "v1.32.10"

// k8sClusterFlavorID is the VM flavor ID used for worker node pools in tests.
// This is a small VM flavor suitable for testing.
const k8sClusterFlavorID = "g1-standard-2-4"

func testAccCheckCloudK8SClusterDestroy(s *terraform.State) error {
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

func TestAccCloudK8SCluster_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudK8SClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudK8SClusterConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("status"), knownvalue.StringExact("Provisioned")),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("version"), knownvalue.StringExact(k8sClusterVersion)),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func TestAccCloudK8SCluster_updatePool(t *testing.T) {
	rName := acctest.RandomName()
	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudK8SClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudK8SClusterConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(0).AtMapKey("min_node_count"),
						knownvalue.Int64Exact(1)),
					compareIDSame.AddStateValue("gcore_cloud_k8s_cluster.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccCloudK8SClusterConfigWithNodeCount(rName, 2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(0).AtMapKey("min_node_count"),
						knownvalue.Int64Exact(2)),
					compareIDSame.AddStateValue("gcore_cloud_k8s_cluster.test", tfjsonpath.New("id")),
				},
			},
		},
	})
}

func TestAccCloudK8SCluster_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudK8SClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudK8SClusterConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cloud_k8s_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// keypair is not returned by the API in the cluster object
					"keypair",
				},
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_k8s_cluster.test", "project_id", "region_id", "name"),
			},
		},
	})
}

func testAccCloudK8SClusterConfig(name string) string {
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
}`, acctest.ProjectID(), acctest.RegionID(), name, k8sClusterVersion, k8sClusterFlavorID)
}

func testAccCloudK8SClusterConfigWithNodeCount(name string, minNodeCount int) string {
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
      min_node_count     = %[6]d
      max_node_count     = 3
      boot_volume_size   = 50
      boot_volume_type   = "standard"
      servergroup_policy = "soft-anti-affinity"
      is_public_ipv4     = true
    }
  ]
}`, acctest.ProjectID(), acctest.RegionID(), name, k8sClusterVersion, k8sClusterFlavorID, minNodeCount)
}
