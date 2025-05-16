package gcore

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWaapPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapPolicyCreate,
		ReadContext:   resourceWaapPolicyRead,
		UpdateContext: resourceWaapPolicyUpdate,
		DeleteContext: resourceWaapPolicyDelete,
		Description:   "Represent WAAP Policy",

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the domain",
			},
			"policies": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ConfigMode:  schema.SchemaConfigModeAttr,
				Description: "Policy Rule List -- You can find the list of available rules using the `data_source_gcore_waap_policy`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Formatted Policy Name in snake_case format.",
						},
						"rules": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "List of Formatted Rule Names with required modes (true/false).",
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

func resourceWaapPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId(d.Get("domain_id").(string))
	return resourceWaapPolicyUpdate(ctx, d, m)
}

func resourceWaapPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	domainID, _ := strconv.Atoi(d.Get("domain_id").(string))

	definedPolicies := make(map[string]map[string]string)

	if v, ok := d.GetOk("policies"); ok {
		rawPolicies := v.([]interface{})

		for _, rawItem := range rawPolicies {
			item := rawItem.(map[string]interface{})

			policyName := item["policy_name"].(string)

			if _, exists := definedPolicies[policyName]; !exists {
				definedPolicies[policyName] = make(map[string]string)
			}

			if rulesRaw, ok := item["rules"]; ok && rulesRaw != nil {
				for ruleName, val := range rulesRaw.(map[string]interface{}) {
					definedPolicies[policyName][ruleName] = fmt.Sprintf("%v", val)
				}
			}
		}
	}

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
		policyName := strings.ReplaceAll(*policy.ResourceSlug, "-", "_")
		rulesMap := make(map[string]interface{})

		if definedRules, ok := definedPolicies[policyName]; ok {
			for _, rule := range *policy.Rules {
				ruleName := formatRuleName(rule.Name)

				if _, shouldInclude := definedRules[ruleName]; shouldInclude {
					rulesMap[ruleName] = strconv.FormatBool(rule.Mode)
				}
			}
		}

		if len(rulesMap) > 0 {
			policyGroup := map[string]interface{}{
				"policy_name": policyName,
				"rules":       rulesMap,
			}
			policies = append(policies, policyGroup)
		}
	}

	if err := d.Set("policies", policies); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set policies: %v", err))
	}

	return nil
}

func resourceWaapPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	domainID, _ := strconv.Atoi(d.Get("domain_id").(string))

	if d.HasChange("policies") {
		definedPolicies := make(map[string]map[string]string)

		if v, ok := d.GetOk("policies"); ok {
			rawPolicies := v.([]interface{})

			for _, rawItem := range rawPolicies {
				item := rawItem.(map[string]interface{})

				policyName := item["policy_name"].(string)

				if _, exists := definedPolicies[policyName]; !exists {
					definedPolicies[policyName] = make(map[string]string)
				}

				if rulesRaw, ok := item["rules"]; ok && rulesRaw != nil {
					for ruleName, val := range rulesRaw.(map[string]interface{}) {
						definedPolicies[policyName][ruleName] = fmt.Sprintf("%v", val)
					}
				}
			}
		}

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

		var apiPolicies []map[string]interface{}

		for _, policy := range *policiesResp.JSON200 {
			policyName := strings.ReplaceAll(*policy.ResourceSlug, "-", "_")
			rulesMap := make(map[string]interface{})

			for _, rule := range *policy.Rules {
				ruleName := formatRuleName(rule.Name)
				rulesMap[ruleName] = map[string]interface{}{
					"mode": strconv.FormatBool(rule.Mode),
					"id":   rule.Id,
				}
			}

			if len(rulesMap) > 0 {
				policyGroup := map[string]interface{}{
					"policy_name": policyName,
					"rules":       rulesMap,
				}
				apiPolicies = append(apiPolicies, policyGroup)
			}
		}

		for definedPolicy, definedRules := range definedPolicies {
			for _, apiPolicy := range apiPolicies {
				policyName := apiPolicy["policy_name"].(string)
				if policyName != definedPolicy {
					continue
				}

				apiRules, _ := apiPolicy["rules"].(map[string]interface{})

				for ruleName, definedMode := range definedRules {
					ruleData := apiRules[ruleName].(map[string]interface{})
					apiMode := ruleData["mode"].(string)

					if definedMode != apiMode {
						ruleID := ruleData["id"].(string)

						updateResp, err := client.ToggleDomainPolicyV1DomainsDomainIdPoliciesPolicyIdTogglePatchWithResponse(
							context.Background(),
							domainID,
							ruleID,
						)

						if err != nil || updateResp.StatusCode() != 200 {
							return diag.FromErr(fmt.Errorf("failed to update policy rule %s. Status code: %d with error: %v", ruleName, updateResp.StatusCode(), err))
						}
					}
				}

			}
		}
	}

	return resourceWaapPolicyRead(ctx, d, m)
}

func resourceWaapPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
