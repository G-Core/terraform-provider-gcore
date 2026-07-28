package cloud_registry_artifact_test

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

// Artifacts are created implicitly by pushing a tagged image to a repository,
// which is impractical to do from an acceptance test. This test instead
// verifies the data source against a freshly created registry with a
// repository name that does not exist: an empty items list is itself the
// expected, correctly-handled result (confirmed against the live API:
// listing artifacts under a nonexistent repository returns count=0, not an
// error).
func TestAccCloudRegistryArtifactsDataSource_basic(t *testing.T) {
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
				Config: testAccCloudRegistryArtifactsDataSourceConfig(registry.ID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_registry_artifacts.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(0)),
				},
			},
		},
	})
}

func testAccCloudRegistryArtifactsDataSourceConfig(registryID int64) string {
	return fmt.Sprintf(`
data "gcore_cloud_registry_artifacts" "test" {
  project_id      = %[1]s
  region_id       = %[2]s
  registry_id     = %[3]d
  repository_name = "nonexistent-repo"
}
`, acctest.ProjectID(), acctest.RegionID(), registryID)
}
