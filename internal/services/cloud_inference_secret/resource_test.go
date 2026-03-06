package cloud_inference_secret_test

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

func TestAccCloudInferenceSecret_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInferenceSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInferenceSecretConfig(rName, "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_inference_secret.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_inference_secret.test",
						tfjsonpath.New("id"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_inference_secret.test",
						tfjsonpath.New("type"), knownvalue.StringExact("aws-iam")),
				},
			},
		},
	})
}

func TestAccCloudInferenceSecret_update(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInferenceSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInferenceSecretConfig(rName, "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_inference_secret.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				// Update with new credentials (increment version)
				Config: testAccCloudInferenceSecretConfig(rName, "AKIAI44QH8DHBEXAMPLE", "je7MtGbClwBF/2Zp9Utk/h3yCo8nvbEXAMPLEKEY", 2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_inference_secret.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_inference_secret.test",
						tfjsonpath.New("data_wo_version"), knownvalue.Int64Exact(2)),
				},
			},
		},
	})
}

func TestAccCloudInferenceSecret_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInferenceSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInferenceSecretConfig(rName, "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", 1),
			},
			{
				ResourceName:      "gcore_cloud_inference_secret.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_inference_secret.test", "project_id", "id"),
				ImportStateVerify: true,
				// data block and data_wo_version are write-only, not in imported state
				ImportStateVerifyIgnore: []string{"data", "data_wo_version"},
			},
		},
	})
}

func testAccCheckCloudInferenceSecretDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_inference_secret" {
			continue
		}

		projectID, _ := strconv.ParseInt(rs.Primary.Attributes["project_id"], 10, 64)

		_, err := client.Cloud.Inference.Secrets.Get(context.Background(), rs.Primary.ID, cloud.InferenceSecretGetParams{
			ProjectID: param.NewOpt(projectID),
		})

		if err == nil {
			return fmt.Errorf("inference secret %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking inference secret deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudInferenceSecretConfig(name, accessKeyID, secretAccessKey string, dataVersion int) string {
	return fmt.Sprintf(`
resource "gcore_cloud_inference_secret" "test" {
  project_id      = %[1]s
  name            = %[2]q
  type            = "aws-iam"
  data_wo_version = %[5]d
  data = {
    aws_access_key_id_wo     = %[3]q
    aws_secret_access_key_wo = %[4]q
  }
}`, acctest.ProjectID(), name, accessKeyID, secretAccessKey, dataVersion)
}
