package fastedge_binary_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// Minimal valid WebAssembly binary (magic number + version)
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

func TestAccFastedgeBinary_basic(t *testing.T) {
	wasmPath := createTestWasmFile(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFastedgeBinaryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFastedgeBinaryConfig(wasmPath),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_fastedge_binary.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_fastedge_binary.test",
						tfjsonpath.New("checksum"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_fastedge_binary.test",
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_fastedge_binary.test",
						tfjsonpath.New("filename"), knownvalue.StringExact(wasmPath)),
				},
			},
		},
	})
}

func TestAccFastedgeBinary_import(t *testing.T) {
	wasmPath := createTestWasmFile(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFastedgeBinaryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFastedgeBinaryConfig(wasmPath),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_fastedge_binary.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:            "gcore_fastedge_binary.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportBlockWithID,
				ImportStateVerifyIgnore: []string{"filename"},
			},
		},
	})
}

func testAccCheckFastedgeBinaryDestroy(s *terraform.State) error {
	return acctest.CheckResourceDestroyed(s, "gcore_fastedge_binary", func(client *gcore.Client, id string) error {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse binary ID: %w", err)
		}
		_, err = client.Fastedge.Binaries.Get(context.Background(), idInt)
		return err
	})
}

func testAccFastedgeBinaryConfig(wasmPath string) string {
	return fmt.Sprintf(`
resource "gcore_fastedge_binary" "test" {
  filename = %[1]q
}`, wasmPath)
}
