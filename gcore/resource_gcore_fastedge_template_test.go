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
			"data_type": "enum",
			"mandatory": true,
			"descr": "param1 description",
			"metadata": "{\"allowed\":[\"a\",\"b\"]}"
		},
		{
			"name": "param2",
			"data_type": "bool",
			"mandatory": false,
			"descr": "param2 description"
		}
	]
}
`

var baseTemplate sdk.Template = sdk.Template{
	Name:       "test-template",
	BinaryId:   314,
	ShortDescr: ptr("short description"),
	LongDescr:  ptr("long description"),
	Params: []sdk.TemplateParam{
		{
			Name:      "param2",
			DataType:  "bool",
			Mandatory: false,
			Descr:     ptr("param2 description"),
		},
		{
			Name:      "param1",
			DataType:  "enum",
			Mandatory: true,
			Descr:     ptr("param1 description"),
			Metadata:  ptr(`{"allowed":["a","b"]}`),
		},
	},
}
var baseTfFastEdgeTemplConfig string = `resource "gcore_fastedge_template" "test" {
	binary = 314
	short_descr = "short description"
	long_descr = "long description"
	param {
		name = "param1"
		type = "enum"
		mandatory = true
		descr = "param1 description"
		metadata = "{\"allowed\":[\"a\",\"b\"]}"
	}
	param {
		name = "param2"
		type = "bool"
		descr = "param2 description"
	}
`

func TestFastEdgeTemplate_basic(t *testing.T) {
	updatedTemplate := baseTemplate
	updatedTemplate.Name = "test-template1"

	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetTemplate": {
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-template", ` + baseTemplateJson,
					},
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-template", ` + baseTemplateJson,
					},
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-template1", ` + baseTemplateJson,
					},
				},
			},
			"AddTemplate": {
				params: []mockParams{
					{
						expectPayload: baseTemplate,
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, "name": "test-template", ` + baseTemplateJson,
					},
				},
			},
			"UpdateTemplate": {
				params: []mockParams{
					{
						expectId:      42,
						expectPayload: updatedTemplate,
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, "name": "test-template1", ` + baseTemplateJson,
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
					testAccCheckResourceExists("gcore_fastedge_template.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "name", "test-template"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "binary", "314"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "short_descr", "short description"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "long_descr", "long description"),
					resource.TestCheckTypeSetElemNestedAttrs("gcore_fastedge_template.test", "param.*", map[string]string{
						"name":      "param1",
						"type":      "enum",
						"mandatory": "true",
						"descr":     "param1 description",
						"metadata":  `{"allowed":["a","b"]}`,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("gcore_fastedge_template.test", "param.*", map[string]string{
						"name":      "param2",
						"type":      "bool",
						"mandatory": "false",
						"descr":     "param2 description",
					}),
				),
			},
			{ // update resource
				Config: baseTfFastEdgeTemplConfig + `name = "test-template1"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_template.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "name", "test-template1"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "binary", "314"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "short_descr", "short description"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "long_descr", "long description"),
					resource.TestCheckTypeSetElemNestedAttrs("gcore_fastedge_template.test", "param.*", map[string]string{
						"name":      "param1",
						"type":      "enum",
						"mandatory": "true",
						"descr":     "param1 description",
						"metadata":  `{"allowed":["a","b"]}`,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("gcore_fastedge_template.test", "param.*", map[string]string{
						"name":      "param2",
						"type":      "bool",
						"mandatory": "false",
						"descr":     "param2 description",
					}),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}

func TestFastEdgeTemplate_disappear(t *testing.T) {
	updatedTemplate := baseTemplate
	updatedTemplate.Name = "test-template1"

	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetTemplate": {
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-template", ` + baseTemplateJson,
					},
					{
						expectId:  42,
						retStatus: http.StatusNotFound, // resource disappeared from the backend
					},
					{
						expectId:  43,
						retStatus: http.StatusOK,
						retBody:   `{"id": 43, "name": "test-template1", ` + baseTemplateJson,
					},
				},
			},
			"AddTemplate": {
				params: []mockParams{
					{
						expectPayload: baseTemplate,
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, "name": "test-template", ` + baseTemplateJson,
					},
					{
						expectPayload: updatedTemplate,
						retStatus:     http.StatusOK,
						retBody:       `{"id": 43, "name": "test-template1", ` + baseTemplateJson,
					},
				},
			},
			"DelTemplate": {
				params: []mockParams{
					{
						expectId:  43,
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
					testAccCheckResourceExists("gcore_fastedge_template.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "name", "test-template"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "binary", "314"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "short_descr", "short description"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "long_descr", "long description"),
					resource.TestCheckTypeSetElemNestedAttrs("gcore_fastedge_template.test", "param.*", map[string]string{
						"name":      "param1",
						"type":      "enum",
						"mandatory": "true",
						"descr":     "param1 description",
						"metadata":  `{"allowed":["a","b"]}`,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("gcore_fastedge_template.test", "param.*", map[string]string{
						"name":      "param2",
						"type":      "bool",
						"mandatory": "false",
						"descr":     "param2 description",
					}),
				),
			},
			{ // resource disappeared - re-create
				Config: baseTfFastEdgeTemplConfig + `name = "test-template1"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_template.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "id", "43"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "name", "test-template1"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "binary", "314"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "short_descr", "short description"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "long_descr", "long description"),
					resource.TestCheckTypeSetElemNestedAttrs("gcore_fastedge_template.test", "param.*", map[string]string{
						"name":      "param1",
						"type":      "enum",
						"mandatory": "true",
						"descr":     "param1 description",
						"metadata":  `{"allowed":["a","b"]}`,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("gcore_fastedge_template.test", "param.*", map[string]string{
						"name":      "param2",
						"type":      "bool",
						"mandatory": "false",
						"descr":     "param2 description",
					}),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}

func TestFastEdgeTemplate_import(t *testing.T) {
	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetTemplate": {
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "name": "test-template", ` + baseTemplateJson,
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
				Config: baseTfFastEdgeTemplConfig + `name = "test-template"
				}`,
				ImportState:   true,
				ImportStateId: "42",
				ResourceName:  "gcore_fastedge_template.test",
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_template.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "name", "test-template"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "binary", "314"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "short_descr", "short description"),
					resource.TestCheckResourceAttr("gcore_fastedge_template.test", "long_descr", "long description"),
					resource.TestCheckTypeSetElemNestedAttrs("gcore_fastedge_template.test", "param.*", map[string]string{
						"name":      "param1",
						"type":      "enum",
						"mandatory": "true",
						"descr":     "param1 description",
						"metadata":  `{"allowed":["a","b"]}`,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("gcore_fastedge_template.test", "param.*", map[string]string{
						"name":      "param2",
						"type":      "bool",
						"mandatory": "false",
						"descr":     "param2 description",
					}),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}
