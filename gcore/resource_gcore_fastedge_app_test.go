package gcore

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

var baseAppJson string = `
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
var baseApp sdk.App = sdk.App{
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
var baseTfFastEdgeAppConfig string = `resource "gcore_fastedge_app" "test" {
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

func TestFastEdgeApp_basic(t *testing.T) {
	updatedApp := baseApp
	updatedApp.Name = ptr("test-app1")
	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetApp": {
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-app", ` + baseAppJson + `}`,
					},
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-app", ` + baseAppJson + `}`,
					},
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-app1", ` + baseAppJson + `}`,
					},
				},
			},
			"AddApp": {
				params: []mockParams{
					{
						expectPayload: baseApp,
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, "name": "test-app", ` + baseAppJson + `}`,
					},
				},
			},
			"UpdateApp": {
				params: []mockParams{
					{
						expectId:      42,
						expectPayload: sdk.UpdateAppJSONRequestBody{App: updatedApp},
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, "name": "test-app1", ` + baseAppJson + `}`,
					},
				},
			},
			"DelApp": {
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusNoContent,
					},
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: fastedgeMockProvider(mock),
		IsUnitTest:        true,
		Steps: []resource.TestStep{
			{ // create resource
				Config: baseTfFastEdgeAppConfig + `name = "test-app"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_app.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "name", "test-app"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "binary", "314"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "status", "enabled"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "env.key", "value"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "rsp_headers.key1", "value1"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "comment", "test application"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "debug", "true"),
				),
			},
			{ // update resource
				Config: baseTfFastEdgeAppConfig + `name = "test-app1"
						}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_app.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "name", "test-app1"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "binary", "314"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "status", "enabled"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "env.key", "value"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "rsp_headers.key1", "value1"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "comment", "test application"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "debug", "true"),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}

func TestFastEdgeApp_disappear(t *testing.T) {
	updatedApp := baseApp
	updatedApp.Name = ptr("test-app1")
	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetApp": {
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-app", ` + baseAppJson + `}`,
					},
					{
						expectId:  42,
						retStatus: http.StatusNotFound, // resource disappeared from the backend
					},
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-app1", ` + baseAppJson + `}`,
					},
				},
			},
			"AddApp": {
				params: []mockParams{
					{
						expectPayload: baseApp,
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, "name": "test-app", ` + baseAppJson + `}`,
					},
					{
						expectPayload: updatedApp,
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, "name": "test-app1", ` + baseAppJson + `}`,
					},
				},
			},
			"DelApp": {
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusNoContent,
					},
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: fastedgeMockProvider(mock),
		IsUnitTest:        true,
		Steps: []resource.TestStep{
			{ // create resource
				Config: baseTfFastEdgeAppConfig + `name = "test-app"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_app.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "name", "test-app"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "binary", "314"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "status", "enabled"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "env.key", "value"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "rsp_headers.key1", "value1"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "comment", "test application"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "debug", "true"),
				),
			},
			{ // resource disappeared - re-create
				Config: baseTfFastEdgeAppConfig + `name = "test-app1"
						}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_app.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "name", "test-app1"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "binary", "314"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "status", "enabled"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "env.key", "value"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "rsp_headers.key1", "value1"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "comment", "test application"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "debug", "true"),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}

func TestFastEdgeApp_import(t *testing.T) {
	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetApp": {
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-app", ` + baseAppJson + `}`,
					},
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: fastedgeMockProvider(mock),
		IsUnitTest:        true,
		Steps: []resource.TestStep{
			{ // import resource
				Config: baseTfFastEdgeAppConfig + `name = "test-app"
				}`,
				ImportState:   true,
				ImportStateId: "42",
				ResourceName:  "gcore_fastedge_app.test",
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_app.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "name", "test-app"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "binary", "314"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "status", "enabled"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "env.key", "value"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "rsp_headers.key1", "value1"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "comment", "test application"),
					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "debug", "true"),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}
