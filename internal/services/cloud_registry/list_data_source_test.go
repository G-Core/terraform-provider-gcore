package cloud_registry_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// There is no gcore_cloud_registry Terraform resource (registries are exposed
// as data sources only), so the prerequisite registry is created and torn
// down directly via the Go SDK rather than via an HCL resource block.
func TestAccCloudRegistriesDataSource_basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	rName := acctest.RandomName()

	client, err := acctest.NewGcoreClient()
	if err != nil {
		t.Fatal(err)
	}
	projectID, err := strconv.ParseInt(acctest.ProjectID(), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	regionID, err := strconv.ParseInt(acctest.RegionID(), 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	registry, err := client.Cloud.Registries.New(context.Background(), cloud.RegistryNewParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
		Name:      rName,
	})
	if err != nil {
		t.Fatalf("error creating prerequisite registry: %s", err)
	}
	t.Cleanup(func() {
		err := client.Cloud.Registries.Delete(context.Background(), registry.ID, cloud.RegistryDeleteParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			t.Logf("error deleting prerequisite registry %d: %s", registry.ID, err)
		}
	})

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRegistriesDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_registries.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_registries.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccCloudRegistriesDataSourceConfig() string {
	return fmt.Sprintf(`
data "gcore_cloud_registries" "test" {
  project_id = %[1]s
  region_id  = %[2]s
}
`, acctest.ProjectID(), acctest.RegionID())
}
