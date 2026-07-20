package cloud_volume_snapshot_test

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

func TestAccCloudVolumeSnapshotDataSource_byID(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeSnapshotDataSourceConfigByID(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_volume_snapshot.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_cloud_volume_snapshot.test", tfjsonpath.New("id"),
						"data.gcore_cloud_volume_snapshot.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_cloud_volume_snapshot.test", tfjsonpath.New("volume_id"),
						"data.gcore_cloud_volume_snapshot.test", tfjsonpath.New("volume_id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_cloud_volume_snapshot.test", tfjsonpath.New("size"),
						"data.gcore_cloud_volume_snapshot.test", tfjsonpath.New("size"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_cloud_volume_snapshot.test", tfjsonpath.New("status"),
						"data.gcore_cloud_volume_snapshot.test", tfjsonpath.New("status"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestAccCloudVolumeSnapshotsDataSource_byVolumeID(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeSnapshotsDataSourceConfigByVolumeID(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					// The list is filtered by the freshly created volume, so it
					// contains exactly the one snapshot from this config.
					statecheck.ExpectKnownValue("data.gcore_cloud_volume_snapshots.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
					statecheck.CompareValuePairs(
						"gcore_cloud_volume_snapshot.test", tfjsonpath.New("id"),
						"data.gcore_cloud_volume_snapshots.test", tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"),
						compare.ValuesSame(),
					),
					statecheck.ExpectKnownValue("data.gcore_cloud_volume_snapshots.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(rName)),
				},
			},
		},
	})
}

func testAccCloudVolumeSnapshotDataSourceConfigByID(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  name       = %[3]q
  size       = 1
}

resource "gcore_cloud_volume_snapshot" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  volume_id  = gcore_cloud_volume.test.id
  name       = %[3]q
}

data "gcore_cloud_volume_snapshot" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  snapshot_id = gcore_cloud_volume_snapshot.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudVolumeSnapshotsDataSourceConfigByVolumeID(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  name       = %[3]q
  size       = 1
}

resource "gcore_cloud_volume_snapshot" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  volume_id  = gcore_cloud_volume.test.id
  name       = %[3]q
}

data "gcore_cloud_volume_snapshots" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  volume_id  = gcore_cloud_volume_snapshot.test.volume_id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
