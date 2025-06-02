package gcore

import (
	"context"
	"log"
	"net/http"
	"strings"

	waap "github.com/G-Core/gcore-waap-sdk-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWaapTag() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataWaapTagRead,
		Description: "Represent WAAP tag",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the tag.",
				Required:    true,
			},
		},
	}
}

func dataWaapTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start reading WAAP Tags")

	client := m.(*Config).WaapClient

	name := d.Get("name").(string)
	params := waap.GetTagsV1TagsGetParams{
		ReadableName: &name,
	}

	reqEditor := func(ctx context.Context, req *http.Request) error {
		q := req.URL.Query()
		q.Set("show_all", "1")
		req.URL.RawQuery = q.Encode()
		return nil
	}

	result, err := client.GetTagsV1TagsGetWithResponse(ctx, &params, reqEditor)
	if err != nil {
		return diag.FromErr(err)
	}

	if result.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to read Tags. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	tag := findTagByName(*&result.JSON200.Results, name)

	if tag == nil {
		return diag.Errorf("Tag with name '%s' not found.", name)
	}

	d.SetId(tag.Name)

	log.Println("[DEBUG] Finish reading WAAP Tags")

	return nil
}

func findTagByName(tags []waap.AppModelsTagsTag, name string) *waap.AppModelsTagsTag {
	for _, tag := range tags {
		if strings.EqualFold(tag.ReadableName, name) {
			return &tag
		}
	}
	return nil
}
