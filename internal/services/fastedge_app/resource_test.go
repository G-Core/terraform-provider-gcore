package fastedge_app_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

var minimalWasm = []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}

func createTestWasmFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.wasm")
	if err := os.WriteFile(path, minimalWasm, 0644); err != nil {
		t.Fatalf("failed to write test wasm file: %s", err)
	}
	return path
}

func TestAccFastedgeApp_basic(t *testing.T) {
	rName := acctest.RandomName()
	wasmPath := createTestWasmFile(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFastedgeAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFastedgeAppConfig(wasmPath, rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_fastedge_app.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_fastedge_app.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_fastedge_app.test",
						tfjsonpath.New("status"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("gcore_fastedge_app.test",
						tfjsonpath.New("binary"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccFastedgeApp_update(t *testing.T) {
	rName := acctest.RandomName()
	updatedName := acctest.RandomName()
	wasmPath := createTestWasmFile(t)

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFastedgeAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFastedgeAppConfig(wasmPath, rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_fastedge_app.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_fastedge_app.test",
						tfjsonpath.New("status"), knownvalue.Int64Exact(1)),
					compareIDSame.AddStateValue(
						"gcore_fastedge_app.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccFastedgeAppConfigUpdated(wasmPath, updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_fastedge_app.test",
						tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("gcore_fastedge_app.test",
						tfjsonpath.New("status"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("gcore_fastedge_app.test",
						tfjsonpath.New("env"), knownvalue.MapSizeExact(1)),
					statecheck.ExpectKnownValue("gcore_fastedge_app.test",
						tfjsonpath.New("env").AtMapKey("TEST_KEY"), knownvalue.StringExact("test_value")),
					compareIDSame.AddStateValue(
						"gcore_fastedge_app.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func TestAccFastedgeApp_import(t *testing.T) {
	rName := acctest.RandomName()
	wasmPath := createTestWasmFile(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFastedgeAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFastedgeAppConfig(wasmPath, rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_fastedge_app.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:    "gcore_fastedge_app.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateVerifyIgnore: []string{
					// binary is computed from template by API, may not match config exactly
					"binary",
					// status may change after import
					"status",
				},
			},
		},
	})
}

// TestAccFastedgeApp_unknownBinary verifies that plan succeeds when binary is
// unknown during plan (e.g., from a resource reference). Regression test for
// GCLOUD2-24005.
func TestAccFastedgeApp_unknownBinary(t *testing.T) {
	rName := acctest.RandomName()
	wasmPath := createTestWasmFile(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// binary = gcore_fastedge_binary.test.id is unknown during plan
				Config:             testAccFastedgeAppConfig(wasmPath, rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckFastedgeAppDestroy(s *terraform.State) error {
	return acctest.CheckResourceDestroyed(s, "gcore_fastedge_app", func(client *gcore.Client, id string) error {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse app ID: %w", err)
		}
		_, err = client.Fastedge.Apps.Get(context.Background(), idInt)
		return err
	})
}

func testAccFastedgeAppConfig(wasmPath, name string) string {
	return fmt.Sprintf(`
resource "gcore_fastedge_binary" "test" {
  filename = %[1]q
}

resource "gcore_fastedge_app" "test" {
  binary = gcore_fastedge_binary.test.id
  name   = %[2]q
  status = 1
}`, wasmPath, name)
}

func testAccFastedgeAppConfigUpdated(wasmPath, name string) string {
	return fmt.Sprintf(`
resource "gcore_fastedge_binary" "test" {
  filename = %[1]q
}

resource "gcore_fastedge_app" "test" {
  binary = gcore_fastedge_binary.test.id
  name   = %[2]q
  status = 1
  env = {
    "TEST_KEY" = "test_value"
  }
}`, wasmPath, name)
}
