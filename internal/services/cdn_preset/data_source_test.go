package cdn_preset_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/G-Core/gcore-go/cdn"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCDNPresetDataSource_basic(t *testing.T) {
	presetID := anyPresetID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNPresetDataSourceConfig(presetID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cdn_preset.test",
						tfjsonpath.New("id"), knownvalue.Int64Exact(presetID)),
					statecheck.ExpectKnownValue("data.gcore_cdn_preset.test",
						tfjsonpath.New("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cdn_preset.test",
						tfjsonpath.New("object_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cdn_preset.test",
						tfjsonpath.New("service"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cdn_preset.test",
						tfjsonpath.New("preset_settings"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCDNPresetsDataSource_basic(t *testing.T) {
	// Presets are account-scoped and read-only, so the list is only meaningful
	// on an account that has at least one.
	anyPresetID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "gcore_cdn_presets" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cdn_presets.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cdn_presets.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cdn_presets.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("object_type"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCDNPresetsDataSource_maxItems(t *testing.T) {
	anyPresetID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "gcore_cdn_presets" "test" {
  max_items = 1
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cdn_presets.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func testAccCDNPresetDataSourceConfig(presetID int64) string {
	return fmt.Sprintf(`
data "gcore_cdn_preset" "test" {
  preset_id = %[1]d
}`, presetID)
}

// anyPresetID returns the ID of any preset on the account, skipping the test
// when there is none. Presets are read-only, so a test cannot create one.
func anyPresetID(t *testing.T) int64 {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	client, err := acctest.NewGcoreClient()
	if err != nil {
		t.Fatal(err)
	}

	presets, err := client.CDN.Presets.List(context.Background(), cdn.PresetListParams{})
	if err != nil {
		t.Fatalf("failed to list CDN presets: %s", err)
	}
	if len(presets.Results) == 0 {
		t.Skip("no CDN presets available on this account")
	}

	return presets.Results[0].ID
}
