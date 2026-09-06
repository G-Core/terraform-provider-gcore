package cloud_gpu_baremetal_cluster_test

import (
	"context"
	"encoding/base64"
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

// TestAccCloudGPUBaremetalCluster_userData covers the servers-settings branch of Update:
// changing user_data patches the cluster template through
// PATCH /v3/gpu/baremetal/{project_id}/{region_id}/clusters/{cluster_id} and then rolls the
// change out to the running servers through POST .../apply_settings.
//
// user_data is exempt from RequiresReplaceOnConfigChange, so the change is applied in place -
// the id check below fails the test if the resource is replaced instead. Applying settings
// re-images every server in the cluster, so this test is slow and consumes GPU bare metal quota.
func TestAccCloudGPUBaremetalCluster_userData(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { gpuBmClusterPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUBaremetalClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUBaremetalClusterConfigWithUserData(rName, "#!/bin/bash\necho initial\n"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("servers_settings").AtMapKey("user_data"),
						knownvalue.StringExact(base64.StdEncoding.EncodeToString([]byte("#!/bin/bash\necho initial\n")))),
					compareIDSame.AddStateValue(
						"gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudGPUBaremetalClusterConfigWithUserData(rName, "#!/bin/bash\necho updated\n"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster.test",
						tfjsonpath.New("servers_settings").AtMapKey("user_data"),
						knownvalue.StringExact(base64.StdEncoding.EncodeToString([]byte("#!/bin/bash\necho updated\n")))),
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

// testAccCloudGPUBaremetalClusterConfigWithUserData mirrors
// testAccCloudGPUBaremetalClusterConfig with a user_data script. The attribute carries the
// Base64-encoded script, matching what Read writes back into state.
func testAccCloudGPUBaremetalClusterConfigWithUserData(name, userData string) string {
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
    user_data = %[7]q
  }
}`, acctest.ProjectID(), acctest.RegionID(), name, gpuBmFlavor(), gpuBmImageID(), gpuBmSSHKey(),
		base64.StdEncoding.EncodeToString([]byte(userData)))
}

// Synthetic IDs for the plan-only validation cases. These must be literals, not
// references to real resources: an unresolved reference is unknown at plan time, and
// the validator deliberately does not judge unknown values.
const (
	fakeNetworkID = "11111111-1111-1111-1111-111111111111"
	fakeSubnetID  = "22222222-2222-2222-2222-222222222222"
	fakeImageID   = "33333333-3333-3333-3333-333333333333"
)

// TestAccCloudGPUBaremetalCluster_interfaceVariantValidation exercises the per-variant
// field contract through a real provider binary. Every case fails during plan, so no
// infrastructure is created and no GPU bare-metal quota is consumed - which is what
// makes this affordable for a resource that has no positive lifecycle coverage of the
// subnet-interface path.
//
// This resource shares the interface contract with cloud_gpu_virtual_cluster; see
// GCLOUD2-28207 for why the contract has to be enforced provider-side.
func TestAccCloudGPUBaremetalCluster_interfaceVariantValidation(t *testing.T) {
	testCases := map[string]struct {
		iface     string
		wantError string
	}{
		"ip_family on subnet": {
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
						Config:      testAccCloudGPUBaremetalClusterConfigWithInterface(tc.iface),
						PlanOnly:    true,
						ExpectError: regexp.MustCompile(regexp.QuoteMeta(tc.wantError)),
					},
				},
			})
		})
	}
}

// testAccCloudGPUBaremetalClusterConfigWithInterface builds a fully synthetic cluster
// carrying a single interface block. It never gets applied, so the flavor, image and
// key values only need to satisfy the other schema validators.
func testAccCloudGPUBaremetalClusterConfigWithInterface(ifaceBlock string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_gpu_baremetal_cluster" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  name          = "tf-test-variant-validation"
  flavor        = "tf-test-placeholder-flavor"
  image_id      = %[4]q
  servers_count = 1

  servers_settings = {
    interfaces = [
      %[3]s
    ]
    credentials = {
      ssh_key_name = "tf-test-placeholder-key"
    }
  }
}`, acctest.ProjectID(), acctest.RegionID(), ifaceBlock, fakeImageID)
}
