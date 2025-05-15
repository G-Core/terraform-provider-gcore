package gcore

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWaapPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWaapPoliciesRead,
		Description: "Represent WAAP Policies",
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Description: "The ID of the domain",
				Required:    true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rules": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceWaapPoliciesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	domainID, _ := strconv.Atoi(d.Get("domain_id").(string))

	policiesResp, err := client.GetRuleSetListV1DomainsDomainIdRuleSetsGetWithResponse(
		context.Background(),
		domainID,
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting domain policies: %v", err))
	}

	if policiesResp.JSON200 == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("no policies found for domain ID %d", domainID))
	}

	var policies []map[string]interface{}

	for _, policy := range *policiesResp.JSON200 {
		resourceSlug := strings.ReplaceAll(*policy.ResourceSlug, "-", "_")
		rulesMap := make(map[string]interface{})

		for _, rule := range *policy.Rules {
			ruleName := formatRuleName(rule.Name)
			rulesMap[ruleName] = rule.Id
		}

		if len(rulesMap) > 0 {
			policyGroup := map[string]interface{}{
				"resource_slug": resourceSlug,
				"rules":         rulesMap,
			}
			policies = append(policies, policyGroup)
		}
	}

	if err := d.Set("policies", policies); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set policies: %v", err))
	}

	d.SetId(strconv.Itoa(domainID))

	return nil
}

func formatRuleName(ruleName string) string {
	re := regexp.MustCompile(`\(([^)]+)\)`)
	matches := re.FindStringSubmatch(ruleName)

	if len(matches) > 1 {
		return strings.ToLower(matches[1])
	}

	formatted := strings.ToLower(ruleName)
	replacements := map[string]string{
		"'s": "s",
		".":  "_",
		"-":  "_",
		" ":  "_",
	}

	for old, new := range replacements {
		formatted = strings.ReplaceAll(formatted, old, new)
	}

	return formatted
}
