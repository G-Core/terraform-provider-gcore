//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v1/metadata"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVolume(t *testing.T) {
	type Params struct {
		Name        string
		Size        int
		Type        string
		Source      string
		SnapshotID  string
		ImageID     string
		MetadataMap string
	}

	create := Params{
		Name: "test",
		Size: 1,
		Type: "standard",
		MetadataMap: `{
			key1 = "val1"
			key2 = "val2"
		}`,
	}

	update := Params{
		Name: "test2",
		Size: 2,
		Type: "ssd_hiiops",
		MetadataMap: `{
			key3 = "val3"
		}`,
	}

	fullName := "gcore_volume.acctest"
	importStateIDPrefix := fmt.Sprintf("%s:%s:", os.Getenv("TEST_PROJECT_ID"), os.Getenv("TEST_REGION_ID"))

	VolumeTemplate := func(params *Params) string {
		additional := fmt.Sprintf("%s\n        %s", regionInfo(), projectInfo())
		if params.SnapshotID != "" {
			additional += fmt.Sprintf(`%s        snapshot_id = "%s"`, "\n", params.SnapshotID)
		}
		if params.ImageID != "" {
			additional += fmt.Sprintf(`%s        image_id = "%s"`, "\n", params.ImageID)
		}

		template := fmt.Sprintf(`
		resource "gcore_volume" "acctest" {
			name = "%s"
			size = %d
			type_name = "%s"
			%s
 			metadata_map = %s
		`, params.Name, params.Size, params.Type, additional, params.MetadataMap)

		return template + "\n}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: VolumeTemplate(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "size", strconv.Itoa(create.Size)),
					resource.TestCheckResourceAttr(fullName, "type_name", create.Type),
					resource.TestCheckResourceAttr(fullName, "name", create.Name),
					testAccCheckMetadata(fullName, true, map[string]interface{}{
						"key1": "val1",
						"key2": "val2",
					}),
				),
			},
			{
				Config: VolumeTemplate(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "size", strconv.Itoa(update.Size)),
					resource.TestCheckResourceAttr(fullName, "type_name", update.Type),
					resource.TestCheckResourceAttr(fullName, "name", update.Name),
					testAccCheckMetadata(fullName, true, map[string]interface{}{
						"key3": "val3",
					}),
					testAccCheckMetadata(fullName, false, map[string]interface{}{
						"key1": "val1",
					}),
					testAccCheckMetadata(fullName, false, map[string]interface{}{
						"key2": "val2",
					})),
			},
			{
				ImportStateIdPrefix: importStateIDPrefix,
				ResourceName:        fullName,
				ImportState:         true,
			},
		},
	})
}

func testAccVolumeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, volumesPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_volume" {
			continue
		}

		_, err := networks.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("volume still exists")
		}
	}

	return nil
}

func setNonStateMeta(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, volumesPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_volume" {
			continue
		}

		resourceID := rs.Primary.ID
		nonStateMeta := map[string]string{"key4": "val4"}
		err = metadata.MetadataReplace(client, resourceID, nonStateMeta).Err
		if err != nil {
			return err
		}
	}

	return nil
}
