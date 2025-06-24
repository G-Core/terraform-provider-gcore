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

func dataWaapSecurityInsightType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataWaapSecurityInsightTypeRead,
		Description: "Represent WAAP security insight type",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the insight type",
				Required:    true,
			},
		},
	}
}

func dataWaapSecurityInsightTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start reading WAAP Security Insight Types")

	config := m.(*Config)
	client := config.WaapClient

	name := d.Get("name").(string)
	params := waap.GetInsightTypesV1SecurityInsightsTypesGetParams{
		Name: &name,
	}

	result, err := client.GetInsightTypesV1SecurityInsightsTypesGetWithResponse(ctx, &params)
	if err != nil {
		return diag.FromErr(err)
	}

	if result.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to read Insight Types. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	insightType := findInsightTypeByName(*result.JSON200, name)

	if insightType == nil {
		return diag.Errorf("Insight type with name '%s' not found.", name)
	}

	d.SetId(insightType.Slug)

	log.Println("[DEBUG] Finish reading WAAP Security Insight Types")
	return nil
}

func findInsightTypeByName(response waap.PaginatedResponseInsightType, name string) *waap.InsightType {
	for _, insightType := range response.Results {
		if strings.EqualFold(insightType.Name, name) {
			return &insightType
		}
	}
	return nil
}
