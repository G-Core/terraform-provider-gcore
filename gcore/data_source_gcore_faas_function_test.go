//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFaaSFunctionDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, faasPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	nsName := "test-ns"
	if err := createTestNamespace(client, nsName); err != nil {
		t.Fatal(err)
	}
	defer faas.DeleteNamespace(client, nsName)

	opts := faas.CreateFunctionOpts{
		Name:    "testname",
		Envs:    map[string]string{},
		Runtime: "go1.16.6",
		Timeout: 5,
		Flavor:  "80mCPU-128MB",
		Autoscaling: faas.FunctionAutoscaling{
			MinInstances: pointer.To(1),
			MaxInstances: pointer.To(2),
		},
		CodeText: `
		package kubeless

import (
        "github.com/kubeless/kubeless/pkg/functions"
)

func Run(evt functions.Event, ctx functions.Context) (string, error) {
        return "Hello World!", nil
}
`,
		MainMethod: "main",
	}

	res, err := faas.CreateFunction(client, nsName, opts).Extract()
	if err != nil {
		t.Fatal(err)
	}

	taskID := res.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, FloatingIPCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		return nil, nil
	})

	f, err := faas.GetFunction(client, nsName, opts.Name).Extract()
	if err != nil {
		t.Fatal(err)
	}

	fullName := "data.gcore_faas_function.acctest"
	tpl := func(ns, f string) string {
		return fmt.Sprintf(`
			data "gcore_faas_function" "acctest" {
			  %s
              %s
              name = "%s"
              namespace = "%s"
			}
		`, projectInfo(), regionInfo(), f, ns)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(nsName, opts.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "id", funcID(f.Name, nsName)),
					resource.TestCheckResourceAttr(fullName, "name", f.Name),
					resource.TestCheckResourceAttr(fullName, "runtime", f.Runtime),
				),
			},
		},
	})
}
