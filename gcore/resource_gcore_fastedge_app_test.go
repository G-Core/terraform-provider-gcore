package gcore

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

func TestAddAppFromExistingBinary(t *testing.T) {
	baseAppJson := `
		"status": 1,
		"binary": 314,
		"env": {
			"key": "value"
		},
		"rsp_headers": {
			"key1": "value1"
		},
		"debug": true,
		"comment": "test application"
	`
	baseApp := sdk.App{
		Name:    ptr("test-app"),
		Binary:  ptr[int64](314),
		Comment: ptr("test application"),
		Status:  ptr(1),
		Debug:   ptr(true),
		Env: ptr(map[string]string{
			"key": "value",
		}),
		RspHeaders: ptr(map[string]string{
			"key1": "value1",
		}),
	}
	baseTfConfig := `resource "gcore_fastedge_app" "test" {
		binary = 314
		status = "enabled"
		env = {
			"key" = "value"
		}
		rsp_headers = {
			"key1" = "value1"
		}
		debug = true
		comment = "test application"
	`
	updatedApp := baseApp
	updatedApp.Name = ptr("test-app1")
	addAppMock := &mockSdk{
		t: t,
		getApp: []mockParams{
			{
				expectId:  42,
				retStatus: http.StatusOK,
				retBody:   `{"id": 42, "name": "test-app", ` + baseAppJson + `}`,
			},
		},
		addApp: []mockParams{
			{
				expectPayload: baseApp,
				retStatus:     http.StatusOK,
				retBody:       `{"id": 42, "name": "test-app", ` + baseAppJson + `}`,
			},
		},
		updApp: []mockParams{
			{
				expectId:      42,
				expectPayload: updatedApp,
				retStatus:     http.StatusOK,
				retBody:       `{"id": 42, "name": "test-app1", ` + baseAppJson + `}`,
			},
		},
		delApp: []mockParams{
			{
				expectId:  42,
				retStatus: http.StatusNoContent,
			},
		},
	}
	providers := map[string]func() (*schema.Provider, error){
		"gcore": func() (*schema.Provider, error) {
			return &schema.Provider{
				ResourcesMap: map[string]*schema.Resource{
					"gcore_fastedge_binary": resourceFastEdgeBinary(),
					"gcore_fastedge_app":    resourceFastEdgeApp(),
				},
				ConfigureContextFunc: func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
					config := Config{FastEdgeClient: &sdk.ClientWithResponses{ClientInterface: addAppMock}}
					return &config, nil
				},
			}, nil
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: providers,
		IsUnitTest:        true,
		Steps: []resource.TestStep{
			{
				Config: baseTfConfig + `name = "test-app"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_app.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "name", "test-app"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "binary", "314"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "status", "enabled"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "env.key", "value"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "rsp_headers.key1", "value1"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "comment", "test application"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "debug", "true"),
				),
			},
			/*			{
						Config: baseTfConfig + `name = "test-app"
						}`,
						Check: resource.ComposeTestCheckFunc(
							testAccCheckResourceExists("gcore_fastedge_app.test"),
						),
					}, */
		},
	})
	compare(t, addAppMock.getAppCount, 1)
	compare(t, addAppMock.addAppCount, 1)
	compare(t, addAppMock.updAppCount, 0)
	compare(t, addAppMock.delAppCount, 1)
}
