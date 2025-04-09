package gcore

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

var baseTemplateJson string = `
	"binary_id": 314,
	"short_descr": "short description",
	"long_descr": "long description",
	"params": [
		{
			"name": "param1",
			"data_type": "string",
			"mandatory": true,
			"descr": "param1 description"
		}
`
var extraParamJson string = `,
		{
			"name": "param2",
			"data_type": "number",
			"mandatory": false,
			"descr": "param2 description"
		}
`
var baseTemplate sdk.Template = sdk.Template{
	Name:       "test-template",
	BinaryId:   314,
	ShortDescr: ptr("short description"),
	LongDescr:  ptr("long description"),
	Params: []sdk.TemplateParam{
		{
			Name:      "param1",
			DataType:  "string",
			Mandatory: true,
			Descr:     ptr("param1 description"),
		},
	},
}
var baseTfFastEdgeTemplConfig string = `resource "gcore_fastedge_template" "test" {
	binary = 314
	short_descr = "short description"
	long_descr = "long description"
	param {
		name = "param1"
		type = "string"
		mandatory = true
		descr = "param1 description"
	}
`
var extraParamTf string = `param {
		name = "param2"
		type = "number"
		descr = "param2 description"
	}
`

func TestFastEdgeTemplate_basic(t *testing.T) {
	updatedTemplate := baseTemplate
	updatedTemplate.Name = "test-template1"
	/*	updatedTemplate.Params = append(updatedTemplate.Params, sdk.TemplateParam{
		Name:      "param2",
		DataType:  "number",
		Mandatory: false,
		Descr:     ptr("param2 description"),
	}) */

	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetTemplate": {
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-template", ` + baseTemplateJson + `]}`,
					},
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-template", ` + baseTemplateJson + `]}`,
					},
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-template1", ` + baseTemplateJson /*+ extraParamJson*/ + `]}`,
					},
				},
			},
			"AddTemplate": {
				params: []mockParams{
					{
						expectPayload: baseTemplate,
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, "name": "test-template", ` + baseTemplateJson + `]}`,
					},
				},
			},
			"UpdateTemplate": {
				params: []mockParams{
					{
						expectId:      42,
						expectPayload: updatedTemplate,
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, "name": "test-template1", ` + baseTemplateJson /* + extraParamJson*/ + `]}`,
					},
				},
			},
			"DelTemplate": {
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
				Config: baseTfFastEdgeTemplConfig + `name = "test-template"
				}`,
				Check: resource.ComposeTestCheckFunc(
				//					testAccCheckResourceExists("gcore_fastedge_template.test"),
				/*					resource.TestCheckResourceAttr("gcore_fastedge_app.test", "id", "42"),
									resource.TestCheckResourceAttr("gcore_fastedge_app.test", "name", "test-app"),
									resource.TestCheckResourceAttr("gcore_fastedge_app.test", "binary", "314"),
									resource.TestCheckResourceAttr("gcore_fastedge_app.test", "status", "enabled"),
									resource.TestCheckResourceAttr("gcore_fastedge_app.test", "env.key", "value"),
									resource.TestCheckResourceAttr("gcore_fastedge_app.test", "rsp_headers.key1", "value1"),
									resource.TestCheckResourceAttr("gcore_fastedge_app.test", "comment", "test application"),
									resource.TestCheckResourceAttr("gcore_fastedge_app.test", "debug", "true"), */
				),
			},
			{ // update resource
				Config: baseTfFastEdgeTemplConfig + /*extraParamTf +*/ `name = "test-template1"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_template.test"),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}
