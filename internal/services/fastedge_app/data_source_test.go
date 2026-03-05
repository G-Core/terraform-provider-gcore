package fastedge_app_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/stainless-sdks/gcore-terraform/internal/acctest"
)

func TestAccFastedgeAppDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()
	wasmPath := createTestWasmFile(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFastedgeAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFastedgeAppDataSourceConfig(wasmPath, rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_fastedge_app.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_fastedge_app.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_fastedge_app.test",
						tfjsonpath.New("status"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccFastedgeAppDataSourceConfig(wasmPath, name string) string {
	return fmt.Sprintf(`
resource "gcore_fastedge_binary" "test" {
  filename = %[1]q
}

resource "gcore_fastedge_app" "test" {
  binary = gcore_fastedge_binary.test.id
  name   = %[2]q
  status = 1
}

data "gcore_fastedge_app" "test" {
  id = gcore_fastedge_app.test.id
}`, wasmPath, name)
}
