package fastedge_binary_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccFastedgeBinaryDataSource_basic(t *testing.T) {
	wasmPath := createTestWasmFile(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFastedgeBinaryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFastedgeBinaryDataSourceConfig(wasmPath),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						"gcore_fastedge_binary.test", tfjsonpath.New("id"),
						"data.gcore_fastedge_binary.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_fastedge_binary.test", tfjsonpath.New("checksum"),
						"data.gcore_fastedge_binary.test", tfjsonpath.New("checksum"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_fastedge_binary.test", tfjsonpath.New("status"),
						"data.gcore_fastedge_binary.test", tfjsonpath.New("status"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccFastedgeBinaryDataSourceConfig(wasmPath string) string {
	return testAccFastedgeBinaryConfig(wasmPath) + `

data "gcore_fastedge_binary" "test" {
  id = gcore_fastedge_binary.test.id
}`
}
