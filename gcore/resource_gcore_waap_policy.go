package gcore

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	waap "github.com/G-Core/gcore-waap-sdk-go"
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
				ForceNew:    true,
				Description: "The WAAP domain ID for which the Policy is configured.",
			},
			"rules": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "A map of rules where each key is a rule ID and the value is a boolean indicating whether the rule is enabled (true) or disabled (false).",
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
				// todo: validate if there are rules with provided IDs for the domain
			},
		},
	}
}

func resourceWaapPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId(d.Get("domain_id").(string))

	return resourceWaapPolicyUpdate(ctx, d, m)
}

func resourceWaapPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start WAAP Policy reading")

	client := m.(*Config).WaapClient
	domainID, _ := strconv.Atoi(d.Get("domain_id").(string))

	rulesFromConfig := d.Get("rules").(map[string]interface{})

	// Get policy rules from API
	rulesFromApi, statusCode, err := getRulesFromApi(ctx, client, domainID)
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("") // Resource not found, remove from state
		}
		return err
	}

	// Update state
	for ruleID, _ := range rulesFromConfig {
		if apiState, exists := rulesFromApi[ruleID]; exists {
			rulesFromConfig[ruleID] = apiState
		}
	}

	d.Set("rules", rulesFromConfig)

	log.Printf("[DEBUG] Finish WAAP Policy reading (id=%s)\n", domainID)
	return nil
}

func resourceWaapPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start WAAP Policy updating")

	client := m.(*Config).WaapClient
	domainID, _ := strconv.Atoi(d.Get("domain_id").(string))

	// Get policy rules from API
	rulesFromApi, _, err := getRulesFromApi(ctx, client, domainID)
	if err != nil {
		return err
	}

	// Get rules from TF config
	rulesFromConfig := d.Get("rules").(map[string]interface{})

	// Compare rules from API with rules from TF config and create a map of rules to be updated
	rulesToUpdate := make(map[string]bool)
	for ruleID, expectedState := range rulesFromConfig {
		apiState, exists := rulesFromApi[ruleID]
		if exists && apiState != expectedState {
			rulesToUpdate[ruleID] = expectedState.(bool)
		}
	}

	// Update rules
	for ruleID, _ := range rulesToUpdate {
		updateResp, err := client.ToggleDomainPolicyV1DomainsDomainIdPoliciesPolicyIdTogglePatchWithResponse(ctx, domainID, ruleID)

		if err != nil {
			return diag.Errorf("Failed to update Policy Rule: %w", err)
		}

		if updateResp.StatusCode() != http.StatusOK {
			return diag.Errorf("Failed to update Policy Rule. Status code: %d with error: %s", updateResp.StatusCode(), updateResp.Body)
		}
	}

	log.Printf("[DEBUG] Finish WAAP Policy updating (id=%s)\n", domainID)
	return resourceWaapPolicyRead(ctx, d, m)
}

func resourceWaapPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func getRulesFromApi(ctx context.Context, waapClient *waap.ClientWithResponses, domainID int) (map[string]interface{}, int, diag.Diagnostics) {
	policiesResp, err := waapClient.GetRuleSetListV1DomainsDomainIdRuleSetsGetWithResponse(ctx, domainID)
	if err != nil {
		return nil, 0, diag.Errorf("Failed to read Policy: %w", err)
	}

	statusCode := policiesResp.StatusCode()

	if statusCode == http.StatusNotFound {
		return nil, statusCode, diag.Diagnostics{
			{Severity: diag.Warning, Summary: fmt.Sprintf("Policy for Domain (%s) was not found, removed from TF state", domainID)},
		}
	}

	if statusCode != http.StatusOK {
		return nil, statusCode, diag.Errorf("Failed to read Policy. Status code: %d with error: %s", policiesResp.StatusCode(), policiesResp.Body)
	}

	// Get flat list of rules from API with their states
	rules := make(map[string]interface{})
	for _, ruleSet := range *policiesResp.JSON200 {
		for _, rule := range *ruleSet.Rules {
			rules[rule.Id] = rule.Mode
		}
	}

	return rules, statusCode, nil
}
