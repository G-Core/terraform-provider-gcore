package cloud_registry_user_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
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

// There is no gcore_cloud_registry_user Terraform resource, so both the
// prerequisite registry and the registry user are created and torn down
// directly via the Go SDK rather than via HCL resource blocks.
func TestAccCloudRegistryUsersDataSource_basic(t *testing.T) {
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

	// Registry user names must be lowercase letters/numbers, max 16 chars.
	userName := "tftestuser"
	user, err := client.Cloud.Registries.Users.New(context.Background(), registry.ID, cloud.RegistryUserNewParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
		Name:      userName,
		Duration:  1,
	})
	if err != nil {
		t.Fatalf("error creating prerequisite registry user: %s", err)
	}
	t.Cleanup(func() {
		err := client.Cloud.Registries.Users.Delete(context.Background(), user.ID, cloud.RegistryUserDeleteParams{
			ProjectID:  param.NewOpt(projectID),
			RegionID:   param.NewOpt(regionID),
			RegistryID: registry.ID,
		})
		if err != nil {
			t.Logf("error deleting prerequisite registry user %d: %s", user.ID, err)
		}
	})

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRegistryUsersDataSourceConfig(registry.ID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_registry_users.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_registry_users.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"),
						// The API namespaces the requested name with a project/registry-scoped
						// prefix (e.g. "r_<client>-<project>-<region>-<registry>+<name>"), so
						// match the suffix rather than the exact requested name.
						knownvalue.StringRegexp(regexp.MustCompile(regexp.QuoteMeta("+"+userName)+"$"))),
				},
			},
		},
	})
}

func testAccCloudRegistryUsersDataSourceConfig(registryID int64) string {
	return fmt.Sprintf(`
data "gcore_cloud_registry_users" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  registry_id = %[3]d
}
`, acctest.ProjectID(), acctest.RegionID(), registryID)
}
