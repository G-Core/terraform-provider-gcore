package gcore

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

func TestAddBinary(t *testing.T) {
	checksum, _ := fileChecksum(os.Args[0])
	addBinMock := &mockSdk{
		t: t,
		getBin: []mockParams{
			{
				expectId:  42,
				retStatus: http.StatusOK,
				retBody: `{
					"id": 42,
					"checksum": "` + checksum + `"
				}`,
			},
		},
		addBin: []mockParams{
			{
				retStatus: http.StatusOK,
				retBody: `{
					"id": 42,
					"checksum": "` + checksum + `"
				}`,
			},
		},
		delBin: []mockParams{
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
					config := Config{FastEdgeClient: &sdk.ClientWithResponses{ClientInterface: addBinMock}}
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
				Config: `resource "gcore_fastedge_binary" "test" {	filename = "` + os.Args[0] + `"}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gcore_fastedge_binary.test", "checksum", checksum),
				),
			},
		},
	})

	compare(t, addBinMock.getBinCount, 1)
	compare(t, addBinMock.addBinCount, 1)
	compare(t, addBinMock.delBinCount, 1)
}
