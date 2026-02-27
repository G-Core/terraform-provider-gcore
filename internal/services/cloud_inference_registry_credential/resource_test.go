package cloud_inference_registry_credential_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/stainless-sdks/gcore-terraform/internal/acctest"
)

func TestAccCloudInferenceRegistryCredential_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInferenceRegistryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInferenceRegistryCredentialConfig(rName, "testuser", "testpassword1", 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_inference_registry_credential.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_inference_registry_credential.test",
						tfjsonpath.New("id"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_inference_registry_credential.test",
						tfjsonpath.New("registry_url"), knownvalue.StringExact("https://registry.example.com")),
					statecheck.ExpectKnownValue("gcore_cloud_inference_registry_credential.test",
						tfjsonpath.New("username"), knownvalue.StringExact("testuser")),
				},
			},
		},
	})
}

func TestAccCloudInferenceRegistryCredential_update(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInferenceRegistryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInferenceRegistryCredentialConfig(rName, "user1", "password1", 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_inference_registry_credential.test",
						tfjsonpath.New("username"), knownvalue.StringExact("user1")),
				},
			},
			{
				Config: testAccCloudInferenceRegistryCredentialConfig(rName, "user2", "password2", 2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_inference_registry_credential.test",
						tfjsonpath.New("username"), knownvalue.StringExact("user2")),
				},
			},
		},
	})
}

func TestAccCloudInferenceRegistryCredential_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInferenceRegistryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInferenceRegistryCredentialConfig(rName, "testuser", "testpassword", 1),
			},
			{
				ResourceName:      "gcore_cloud_inference_registry_credential.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_inference_registry_credential.test", "project_id", "id"),
				ImportStateVerify: true,
				// password_wo and password_wo_version are write-only, not in imported state
				ImportStateVerifyIgnore: []string{"password_wo", "password_wo_version"},
			},
		},
	})
}

func testAccCheckCloudInferenceRegistryCredentialDestroy(s *terraform.State) error {
	return acctest.CheckResourceDestroyed(s, "gcore_cloud_inference_registry_credential", func(client *gcore.Client, id string) error {
		// Get project_id from environment
		projectIDStr := acctest.ProjectID()
		projectID, _ := strconv.ParseInt(projectIDStr, 10, 64)

		_, err := client.Cloud.Inference.RegistryCredentials.Get(context.Background(), id, cloud.InferenceRegistryCredentialGetParams{
			ProjectID: param.NewOpt(projectID),
		})
		return err
	})
}

func testAccCloudInferenceRegistryCredentialConfig(name, username, password string, passwordVersion int) string {
	return fmt.Sprintf(`
resource "gcore_cloud_inference_registry_credential" "test" {
  project_id          = %[1]s
  name                = %[2]q
  registry_url        = "https://registry.example.com"
  username            = %[3]q
  password_wo         = %[4]q
  password_wo_version = %[5]d
}`, acctest.ProjectID(), name, username, password, passwordVersion)
}
