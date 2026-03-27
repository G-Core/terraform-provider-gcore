package cloud_file_share_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCloudFileShareDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_file_share.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_cloud_file_share.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_file_share.test",
						tfjsonpath.New("protocol"), knownvalue.StringExact("NFS")),
					statecheck.ExpectKnownValue("data.gcore_cloud_file_share.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(1)),
					statecheck.CompareValuePairs(
						"gcore_cloud_file_share.test", tfjsonpath.New("id"),
						"data.gcore_cloud_file_share.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_cloud_file_share.test", tfjsonpath.New("name"),
						"data.gcore_cloud_file_share.test", tfjsonpath.New("name"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudFileShareDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  name          = "%[3]s-network"
  create_router = true
  type          = "vxlan"
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

data "gcore_cloud_file_share" "test" {
  project_id     = %[1]s
  region_id      = %[2]s
  file_share_id  = gcore_cloud_file_share.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
