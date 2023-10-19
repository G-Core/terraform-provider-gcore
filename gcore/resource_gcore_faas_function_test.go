//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFaaSFunction(t *testing.T) {

	type Params struct {
		Name         string
		Description  string
		CodeText     string
		Timeout      int
		MaxInstances int
		MinInstances int
		Key          string
		EnableAPIKey bool
	}

	create := Params{
		Name:        "f-name",
		Description: "description",
		CodeText: `
		package kubeless

import (
        "github.com/kubeless/kubeless/pkg/functions"
)

func Run(evt functions.Event, ctx functions.Context) (string, error) {
        return "Hello World!!", nil
}`,
		Timeout:      5,
		MaxInstances: 2,
		MinInstances: 1,
		EnableAPIKey: true,
		Key:          "\"test-key\"",
	}

	update := Params{
		Name:        "f-name",
		Description: "changed description",
		CodeText: `
		package kubeless

import (
        "github.com/kubeless/kubeless/pkg/functions"
)

func Run(evt functions.Event, ctx functions.Context) (string, error) {
        return "Hello World!", nil
}`,
		Timeout:      6,
		MaxInstances: 3,
		MinInstances: 1,
		EnableAPIKey: false,
		Key:          "",
	}

	fullName := "gcore_faas_function.acctest"

	tpl := func(params *Params) string {
		template := fmt.Sprintf(`
		resource "gcore_faas_key" "acctestkey" {
			name = "test-key"
			description = "test apikey"
			%s
			%s
        }

		resource "gcore_faas_namespace" "acctest1" {
			name = "test"
			description = "test"
			%s
			%s
		}

		resource "gcore_faas_function" "acctest" {
			name = "%s"
			description = "%s"
			namespace = gcore_faas_namespace.acctest1.name
			runtime = "go1.16.6"
			code_text = <<EOF
%s
EOF
			timeout = %d
			max_instances = %d
			min_instances = %d
			%s
			%s
			flavor = "80mCPU-128MB"
			main_method = "main"
			keys = [%s]
			enable_api_key = %t
		}`, regionInfo(), projectInfo(), regionInfo(), projectInfo(), params.Name,
			params.Description, params.CodeText,
			params.Timeout, params.MaxInstances, params.MinInstances,
			regionInfo(), projectInfo(), params.Key, params.EnableAPIKey)

		return template
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      CheckDestroyFaaSFunction,
		Steps: []resource.TestStep{
			{
				Config: tpl(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", create.Name),
					resource.TestCheckResourceAttr(fullName, "description", create.Description),
					resource.TestCheckResourceAttr(fullName, "code_text", create.CodeText+"\n"),
					resource.TestCheckResourceAttr(fullName, "timeout", strconv.Itoa(create.Timeout)),
					resource.TestCheckResourceAttr(fullName, "max_instances", strconv.Itoa(create.MaxInstances)),
					resource.TestCheckResourceAttr(fullName, "min_instances", strconv.Itoa(create.MinInstances)),
					resource.TestCheckResourceAttr(fullName, "keys.0", "test-key"),
					resource.TestCheckResourceAttr(fullName, "enable_api_key", "true"),
				),
			},
			{
				Config: tpl(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", update.Name),
					resource.TestCheckResourceAttr(fullName, "description", update.Description),
					resource.TestCheckResourceAttr(fullName, "code_text", update.CodeText+"\n"),
					resource.TestCheckResourceAttr(fullName, "timeout", strconv.Itoa(update.Timeout)),
					resource.TestCheckResourceAttr(fullName, "max_instances", strconv.Itoa(update.MaxInstances)),
					resource.TestCheckResourceAttr(fullName, "min_instances", strconv.Itoa(update.MinInstances)),
					resource.TestCheckNoResourceAttr(fullName, "keys.0"),
					resource.TestCheckResourceAttr(fullName, "enable_api_key", "false"),
				),
			},
		},
	})
}

func CheckDestroyFaaSFunction(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, faasPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_faas_function" {
			continue
		}

		nsName := rs.Primary.Attributes["namespace"]
		fName := rs.Primary.Attributes["name"]
		_, err := faas.GetFunction(client, nsName, fName).Extract()
		if err == nil {
			return fmt.Errorf("faas function still exists")
		}
	}

	return nil
}
