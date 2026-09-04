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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// k8sClusterVersion is the Kubernetes version used in tests.
// It should be a supported version in the target region.
const k8sClusterVersion = "v1.33.13"

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
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(0).AtMapKey("security_group_ids"),
						knownvalue.ListSizeExact(0)),
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

// mapLacksKey is a knownvalue.Check asserting that a map value does not
// contain the given key. Useful when the API adds system-managed entries, so
// exact/size checks are too strict.
type mapLacksKey struct {
	key string
}

func (c mapLacksKey) CheckValue(other any) error {
	otherVal, ok := other.(map[string]any)
	if !ok {
		return fmt.Errorf("expected map[string]any for mapLacksKey check, got: %T", other)
	}
	if _, found := otherVal[c.key]; found {
		return fmt.Errorf("expected key %q to be absent, but it is present (value: %v)", c.key, otherVal[c.key])
	}
	return nil
}

func (c mapLacksKey) String() string {
	return fmt.Sprintf("map lacking key %q", c.key)
}

// testAccCapturePoolIDs returns a TestCheckFunc that fetches the cluster's
// pools from the API and stores their UUIDs keyed by pool name into dst.
func testAccCapturePoolIDs(resourceName string, dst map[string]string) resource.TestCheckFunc {
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

		pools, err := client.Cloud.K8S.Clusters.Pools.List(
			context.Background(),
			rs.Primary.Attributes["name"],
			cloud.K8SClusterPoolListParams{
				ProjectID: param.NewOpt(projectID),
				RegionID:  param.NewOpt(regionID),
			},
		)
		if err != nil {
			return fmt.Errorf("error listing cluster pools: %w", err)
		}

		for k := range dst {
			delete(dst, k)
		}
		for _, pool := range pools.Results {
			dst[pool.Name] = pool.ID
		}
		return nil
	}
}

// TestAccCloudK8SCluster_insertPoolMiddle verifies that inserting a new pool
// in the middle of the pools list does not recreate the existing pools
// (GCLOUD2-20597). Pool changes are correlated by name, not list index.
func TestAccCloudK8SCluster_insertPoolMiddle(t *testing.T) {
	rName := acctest.RandomName()
	compareClusterIDSame := statecheck.CompareValue(compare.ValuesSame())

	poolIDsBefore := map[string]string{}
	poolIDsAfter := map[string]string{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudK8SClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudK8SClusterConfigTwoPools(rName),
				Check:  testAccCapturePoolIDs("gcore_cloud_k8s_cluster.test", poolIDsBefore),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools"), knownvalue.ListSizeExact(2)),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("first-pool")),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(1).AtMapKey("name"),
						knownvalue.StringExact("second-pool")),
					compareClusterIDSame.AddStateValue("gcore_cloud_k8s_cluster.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccCloudK8SClusterConfigThreePoolsMiddle(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// The cluster itself must be updated in place, never replaced.
						plancheck.ExpectResourceAction("gcore_cloud_k8s_cluster.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCapturePoolIDs("gcore_cloud_k8s_cluster.test", poolIDsAfter),
					func(_ *terraform.State) error {
						for _, name := range []string{"first-pool", "second-pool"} {
							before, after := poolIDsBefore[name], poolIDsAfter[name]
							if before == "" || after == "" {
								return fmt.Errorf("pool %q missing before (%q) or after (%q) the update", name, before, after)
							}
							if before != after {
								return fmt.Errorf("pool %q was recreated: UUID changed from %s to %s", name, before, after)
							}
						}
						if poolIDsAfter["middle-pool"] == "" {
							return fmt.Errorf("middle-pool was not created")
						}
						return nil
					},
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools"), knownvalue.ListSizeExact(3)),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("first-pool")),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(1).AtMapKey("name"),
						knownvalue.StringExact("middle-pool")),
					// The inserted pool must not inherit values from the pool
					// that previously occupied its list index (second-pool).
					// The API may add system-managed labels, so only assert
					// the inherited key is absent.
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(1).AtMapKey("labels"),
						mapLacksKey{key: "pool-role"}),
					// min_node_count was omitted for middle-pool; the provider
					// must default it to 1, not send 0.
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(1).AtMapKey("min_node_count"),
						knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(2).AtMapKey("name"),
						knownvalue.StringExact("second-pool")),
					// The displaced pool keeps its own labels.
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(2).AtMapKey("labels"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"pool-role": knownvalue.StringExact("secondary"),
						})),
					compareClusterIDSame.AddStateValue("gcore_cloud_k8s_cluster.test", tfjsonpath.New("id")),
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

func testAccCloudK8SClusterConfigTwoPools(name string) string {
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
      name               = "first-pool"
      flavor_id          = %[5]q
      min_node_count     = 1
      max_node_count     = 2
      boot_volume_size   = 50
      boot_volume_type   = "standard"
      servergroup_policy = "soft-anti-affinity"
      is_public_ipv4     = true
    },
    {
      name               = "second-pool"
      flavor_id          = %[5]q
      min_node_count     = 1
      max_node_count     = 2
      boot_volume_size   = 50
      boot_volume_type   = "standard"
      servergroup_policy = "soft-anti-affinity"
      is_public_ipv4     = true
      labels = {
        "pool-role" = "secondary"
      }
    }
  ]
}`, acctest.ProjectID(), acctest.RegionID(), name, k8sClusterVersion, k8sClusterFlavorID)
}

func testAccCloudK8SClusterConfigThreePoolsMiddle(name string) string {
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
      name               = "first-pool"
      flavor_id          = %[5]q
      min_node_count     = 1
      max_node_count     = 2
      boot_volume_size   = 50
      boot_volume_type   = "standard"
      servergroup_policy = "soft-anti-affinity"
      is_public_ipv4     = true
    },
    {
      name               = "middle-pool"
      flavor_id          = %[5]q
      max_node_count     = 2
      boot_volume_size   = 50
      boot_volume_type   = "standard"
      servergroup_policy = "soft-anti-affinity"
      is_public_ipv4     = true
    },
    {
      name               = "second-pool"
      flavor_id          = %[5]q
      min_node_count     = 1
      max_node_count     = 2
      boot_volume_size   = 50
      boot_volume_type   = "standard"
      servergroup_policy = "soft-anti-affinity"
      is_public_ipv4     = true
      labels = {
        "pool-role" = "secondary"
      }
    }
  ]
}`, acctest.ProjectID(), acctest.RegionID(), name, k8sClusterVersion, k8sClusterFlavorID)
}

// testAccCloudK8SClusterConfigNamedPools renders a cluster with two pools
// whose names, labels and max_node_count are parameterized, letting tests
// swap names/attributes between steps.
func testAccCloudK8SClusterConfigNamedPools(name, pool1Name, pool1Env string, pool1Max int, pool2Name, pool2Env string, pool2Max int) string {
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
      name               = %[6]q
      flavor_id          = %[5]q
      min_node_count     = 1
      max_node_count     = %[8]d
      boot_volume_size   = 50
      boot_volume_type   = "standard"
      servergroup_policy = "soft-anti-affinity"
      is_public_ipv4     = true
      labels = {
        "env" = %[7]q
      }
    },
    {
      name               = %[9]q
      flavor_id          = %[5]q
      min_node_count     = 1
      max_node_count     = %[11]d
      boot_volume_size   = 50
      boot_volume_type   = "standard"
      servergroup_policy = "soft-anti-affinity"
      is_public_ipv4     = true
      labels = {
        "env" = %[10]q
      }
    }
  ]
}`, acctest.ProjectID(), acctest.RegionID(), name, k8sClusterVersion, k8sClusterFlavorID,
		pool1Name, pool1Env, pool1Max, pool2Name, pool2Env, pool2Max)
}

// TestAccCloudK8SCluster_swapPoolNames reproduces GCLOUD2-20597 QA issue 1:
// swapping pool names in config (attributes staying at their positions) must
// produce a valid plan and be applied as in-place updates of both pools by
// name — no recreation, no cross-pool value mixing.
func TestAccCloudK8SCluster_swapPoolNames(t *testing.T) {
	rName := acctest.RandomName()

	poolIDsBefore := map[string]string{}
	poolIDsAfter := map[string]string{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudK8SClusterDestroy,
		Steps: []resource.TestStep{
			{
				// pool-a{env=prod,max=2}, pool-b{env=staging,max=3}
				Config: testAccCloudK8SClusterConfigNamedPools(rName,
					"pool-a", "prod", 2, "pool-b", "staging", 3),
				Check: testAccCapturePoolIDs("gcore_cloud_k8s_cluster.test", poolIDsBefore),
			},
			{
				// Names swapped: pool-b{env=prod,max=2}, pool-a{env=staging,max=3}
				Config: testAccCloudK8SClusterConfigNamedPools(rName,
					"pool-b", "prod", 2, "pool-a", "staging", 3),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_k8s_cluster.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCapturePoolIDs("gcore_cloud_k8s_cluster.test", poolIDsAfter),
					func(_ *terraform.State) error {
						// Both pools must survive as in-place updates by name.
						for _, name := range []string{"pool-a", "pool-b"} {
							before, after := poolIDsBefore[name], poolIDsAfter[name]
							if before == "" || after == "" {
								return fmt.Errorf("pool %q missing before (%q) or after (%q) the swap", name, before, after)
							}
							if before != after {
								return fmt.Errorf("pool %q was recreated: UUID changed from %s to %s", name, before, after)
							}
						}
						return nil
					},
				),
				ConfigStateChecks: []statecheck.StateCheck{
					// Attribute ownership follows config: index 0 is pool-b
					// with prod/max=2, index 1 is pool-a with staging/max=3.
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("pool-b")),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(0).AtMapKey("labels").AtMapKey("env"),
						knownvalue.StringExact("prod")),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(0).AtMapKey("max_node_count"),
						knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(1).AtMapKey("name"),
						knownvalue.StringExact("pool-a")),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(1).AtMapKey("labels").AtMapKey("env"),
						knownvalue.StringExact("staging")),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(1).AtMapKey("max_node_count"),
						knownvalue.Int64Exact(3)),
				},
			},
		},
	})
}

func testAccCloudK8SClusterConfigInsertUpdateStep(name string, poolAMax int, withMiddle bool) string {
	middle := ""
	if withMiddle {
		middle = fmt.Sprintf(`
    {
      name               = "middle-pool"
      flavor_id          = %[1]q
      max_node_count     = 2
      boot_volume_size   = 50
      boot_volume_type   = "standard"
      servergroup_policy = "soft-anti-affinity"
      is_public_ipv4     = true
    },`, k8sClusterFlavorID)
	}
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
      name               = "pool-a"
      flavor_id          = %[5]q
      min_node_count     = 1
      max_node_count     = %[6]d
      boot_volume_size   = 50
      boot_volume_type   = "standard"
      servergroup_policy = "soft-anti-affinity"
      is_public_ipv4     = true
      labels             = {}
    },%[7]s
    {
      name               = "pool-b"
      flavor_id          = %[5]q
      min_node_count     = 1
      max_node_count     = 3
      boot_volume_size   = 50
      boot_volume_type   = "standard"
      servergroup_policy = "soft-anti-affinity"
      is_public_ipv4     = true
      labels = {
        "env" = "staging"
      }
      taints = {
        "dedicated" = "true:NoSchedule"
      }
    }
  ]
}`, acctest.ProjectID(), acctest.RegionID(), name, k8sClusterVersion, k8sClusterFlavorID, poolAMax, middle)
}

// TestAccCloudK8SCluster_insertPoolMiddleWithUpdate reproduces GCLOUD2-20597
// QA issue 2: inserting a pool mid-list while simultaneously updating an
// existing pool's attribute. pool-a must not inherit pool-b's taints, the
// apply must be consistent, and a follow-up plan must be empty.
func TestAccCloudK8SCluster_insertPoolMiddleWithUpdate(t *testing.T) {
	rName := acctest.RandomName()

	poolIDsBefore := map[string]string{}
	poolIDsAfter := map[string]string{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudK8SClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudK8SClusterConfigInsertUpdateStep(rName, 2, false),
				Check:  testAccCapturePoolIDs("gcore_cloud_k8s_cluster.test", poolIDsBefore),
			},
			{
				// Insert middle-pool AND change pool-a max_node_count 2 -> 3.
				Config: testAccCloudK8SClusterConfigInsertUpdateStep(rName, 3, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_k8s_cluster.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCapturePoolIDs("gcore_cloud_k8s_cluster.test", poolIDsAfter),
					func(_ *terraform.State) error {
						for _, name := range []string{"pool-a", "pool-b"} {
							before, after := poolIDsBefore[name], poolIDsAfter[name]
							if before == "" || after == "" {
								return fmt.Errorf("pool %q missing before (%q) or after (%q) the update", name, before, after)
							}
							if before != after {
								return fmt.Errorf("pool %q was recreated: UUID changed from %s to %s", name, before, after)
							}
						}
						if poolIDsAfter["middle-pool"] == "" {
							return fmt.Errorf("middle-pool was not created")
						}
						return nil
					},
				),
				ConfigStateChecks: []statecheck.StateCheck{
					// pool-a: updated in place, no taints inherited from pool-b.
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("pool-a")),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(0).AtMapKey("max_node_count"),
						knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(0).AtMapKey("taints"),
						mapLacksKey{key: "dedicated"}),
					// middle-pool: inserted, no inherited labels/taints.
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(1).AtMapKey("name"),
						knownvalue.StringExact("middle-pool")),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(1).AtMapKey("taints"),
						mapLacksKey{key: "dedicated"}),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(1).AtMapKey("labels"),
						mapLacksKey{key: "env"}),
					// pool-b: keeps its own labels and taints.
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(2).AtMapKey("name"),
						knownvalue.StringExact("pool-b")),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(2).AtMapKey("labels").AtMapKey("env"),
						knownvalue.StringExact("staging")),
					statecheck.ExpectKnownValue("gcore_cloud_k8s_cluster.test",
						tfjsonpath.New("pools").AtSliceIndex(2).AtMapKey("taints").AtMapKey("dedicated"),
						knownvalue.StringExact("true:NoSchedule")),
				},
			},
		},
	})
}
