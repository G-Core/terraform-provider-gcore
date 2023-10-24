//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/aiimages"
	ai "github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/ais"
	"github.com/G-Core/gcorelabscloud-go/gcore/keypair/v2/keypairs"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
)

const (
	AIImagesPoint = "ai/images"
)

func TestAccAIClusterResource(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	kpClient, err := CreateTestClient(cfg.Provider, keypairsPoint, versionPointV2)
	if err != nil {
		t.Fatal(err)
	}

	ProjectID, err := strconv.Atoi(os.Getenv("TEST_PROJECT_ID"))
	if err != nil {
		t.Fatal(err)
	}

	kpOpts := keypairs.CreateOpts{
		Name:      testAIClusterKeypair,
		PublicKey: pkTest,
		ProjectID: ProjectID,
	}
	keyPair, err := keypairs.Create(kpClient, kpOpts).Extract()
	if err != nil {
		t.Fatal(err)
	}
	defer keypairs.Delete(kpClient, keyPair.ID)

	sgClient, err := CreateTestClient(cfg.Provider, securityGroupPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	sgList, err := securitygroups.ListAll(sgClient, securitygroups.ListOpts{})
	if err != nil {
		t.Fatal(err)
	}
	var defaultSgID string
	for _, sg := range sgList {
		if sg.Name == "default" {
			defaultSgID = sg.ID
			break
		}
	}

	aiImageClient, err := CreateTestClient(cfg.Provider, AIImagesPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}
	opts := aiimages.AIImageListOpts{}
	aiImages, err := aiimages.ListAll(aiImageClient, opts)
	if err != nil {
		t.Fatal(err)
	}

	testAIImage := aiImages[0].ID

	fullName := "gcore_ai_cluster.acctest"

	AiClusterResourceTemplate := fmt.Sprintf(`
	resource "gcore_ai_cluster" "acctest" {
		%[1]s
		%[2]s
		flavor = "%[3]s"
		image_id = "%[4]s"
		cluster_name = "acctest"
		keypair_name = "%[5]s"
		cluster_status = "ACTIVE"
		cluster_metadata = {
		  qqqq = "quk"
		  zuk  = "muk"
		}
		volume {
		  source     = "image"
		  image_id = "%[4]s"
		  volume_type = "standard"
		  size = 20
		}
		interface {
		  type = "external"
		}
		security_group {
		  id = "%[6]s"
		}
	}
	`,
		regionInfo(),
		projectInfo(),
		testAIClusterFlavor,
		testAIImage,
		keyPair.Name,
		defaultSgID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccAIClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: AiClusterResourceTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "cluster_name", "acctest"),
				),
			},
		},
	})
}

func testAccAIClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, AIClusterPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_ai_cluster" {
			continue
		}

		_, err := ai.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ai cluster still exists")
		}
	}
	return nil
}
