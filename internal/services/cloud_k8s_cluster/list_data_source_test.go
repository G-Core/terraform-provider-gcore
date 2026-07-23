package cloud_k8s_cluster_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// Clusters do not exist by default, so this test creates one and then reads the
// list data source (deferred via depends_on) to assert it is returned.
func TestAccCloudK8SClustersDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudK8SClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudK8SClustersDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_clusters.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_clusters.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_clusters.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("status"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccCloudK8SClustersDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%[1]s

data "gcore_cloud_k8s_clusters" "test" {
  project_id = %[2]s
  region_id  = %[3]s
  depends_on = [gcore_cloud_k8s_cluster.test]
}
`, testAccCloudK8SClusterConfig(name), acctest.ProjectID(), acctest.RegionID())
}
