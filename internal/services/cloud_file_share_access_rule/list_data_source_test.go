package cloud_file_share_access_rule_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// Access rules require a file share to attach to, so provision a network,
// subnet, file share and access rule first, then read the list data source
// (scoped by file_share_id, deferred via depends_on).

func TestAccCloudFileShareAccessRulesDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareAccessRulesDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_file_share_access_rules.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_file_share_access_rules.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("access_level"), knownvalue.StringExact("rw")),
					statecheck.ExpectKnownValue("data.gcore_cloud_file_share_access_rules.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("access_to"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccCloudFileShareAccessRulesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  network_id = gcore_cloud_network.test.id
  name       = "%[3]s-subnet"
  cidr       = "10.0.0.0/24"
}

resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "NFS"
  size       = 1
  type_name  = "standard"
  network = {
    network_id = gcore_cloud_network.test.id
    subnet_id  = gcore_cloud_network_subnet.test.id
  }
}

resource "gcore_cloud_file_share_access_rule" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  file_share_id = gcore_cloud_file_share.test.id
  access_mode   = "rw"
  ip_address    = "10.0.0.0/24"
}

data "gcore_cloud_file_share_access_rules" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  file_share_id = gcore_cloud_file_share.test.id
  depends_on    = [gcore_cloud_file_share_access_rule.test]
}
`, acctest.ProjectID(), acctest.RegionID(), name)
}
