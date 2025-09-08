package gcore

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		CustomizeDiff: validatePolicies,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The WAAP domain ID for which the Policy is configured.",
			},
			"policies": {
				Type:     schema.TypeMap,
				Required: true,
				Description: "A map of policies where each key is a policy ID and the value is a boolean indicating whether the policy is enabled (true) or disabled (false). " +
					"Policy IDs can be obtained from the API endpoint /v1/domains/{domain_id}/rule-sets (the 'rules' field) or you can use the gcore_waap_domain_policy data source.",
				Elem: &schema.Schema{
					Type: schema.TypeBool,
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
	log.Println("[DEBUG] Start WAAP Policy reading")

	client := m.(*Config).WaapClient
	domainID, _ := strconv.Atoi(d.Get("domain_id").(string))

	policiesFromConfig := d.Get("policies").(map[string]interface{})

	// Get policies from API
	policiesFromApi, statusCode, err := getPoliciesFromApi(ctx, client, domainID)
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("") // Resource not found, remove from state
		}
		return err
	}

	// Update state
	for policyID, _ := range policiesFromConfig {
		if apiState, exists := policiesFromApi[policyID]; exists {
			policiesFromConfig[policyID] = apiState
		}
	}

	d.Set("policies", policiesFromConfig)

	log.Printf("[DEBUG] Finish WAAP Policy reading (id=%d)\n", domainID)
	return nil
}

func resourceWaapPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start WAAP Policy updating")

	var diags diag.Diagnostics
	client := m.(*Config).WaapClient
	domainID, _ := strconv.Atoi(d.Get("domain_id").(string))

	// Get policies from API
	policiesFromApi, _, err := getPoliciesFromApi(ctx, client, domainID)
	if err != nil {
		return err
	}

	// Get policies from TF config
	policiesFromConfig := d.Get("policies").(map[string]interface{})

	// Compare policies from API with policies from TF config and create a map of policies to be updated
	policiesToUpdate := make(map[string]bool)
	for policyID, expectedState := range policiesFromConfig {
		apiState, exists := policiesFromApi[policyID]

		if !exists {
			// add warning to diagnostics if policy ID from config does not exist in API
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Policy ID '%s' specified in configuration does not exist in API for domain %d", policyID, domainID),
			})
			continue
		}

		if apiState != expectedState {
			policiesToUpdate[policyID] = expectedState.(bool)
		}
	}

	// Update policies
	for policyID, _ := range policiesToUpdate {
		updateResp, err := client.ToggleDomainPolicyV1DomainsDomainIdPoliciesPolicyIdTogglePatchWithResponse(ctx, domainID, policyID)

		if err != nil {
			return diag.Errorf("Failed to update Policy state: %s", err)
		}

		if updateResp.StatusCode() != http.StatusOK {
			return diag.Errorf("Failed to update Policy state. Status code: %d with error: %s", updateResp.StatusCode(), updateResp.Body)
		}
	}

	log.Printf("[DEBUG] Finish WAAP Policy updating (id=%d)\n", domainID)

	diags = append(diags, resourceWaapPolicyRead(ctx, d, m)...)
	return diags
}

func resourceWaapPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func getPoliciesFromApi(ctx context.Context, waapClient *waap.ClientWithResponses, domainID int) (map[string]interface{}, int, diag.Diagnostics) {
	policiesResp, err := waapClient.GetRuleSetListV1DomainsDomainIdRuleSetsGetWithResponse(ctx, domainID)
	if err != nil {
		return nil, 0, diag.Errorf("Failed to read Policy: %s", err)
	}

	statusCode := policiesResp.StatusCode()

	if statusCode == http.StatusNotFound {
		return nil, statusCode, diag.Diagnostics{
			{Severity: diag.Warning, Summary: fmt.Sprintf("Policy for Domain (%d) was not found, removed from TF state", domainID)},
		}
	}

	if statusCode != http.StatusOK {
		return nil, statusCode, diag.Errorf("Failed to read Policy. Status code: %d with error: %s", policiesResp.StatusCode(), policiesResp.Body)
	}

	// Get flat list of policies from API with their states
	policies := make(map[string]interface{})
	for _, policySet := range *policiesResp.JSON200 {
		for _, policy := range *policySet.Rules {
			policies[policy.Id] = policy.Mode
		}
	}

	return policies, statusCode, nil
}

func validatePolicies(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	domainIDStr := d.Get("domain_id").(string)

	// Skip validation if domain ID has not been set
	if domainIDStr == "" {
		return nil
	}

	domainID, err := strconv.Atoi(domainIDStr)
	if err != nil {
		return fmt.Errorf("invalid domain ID: %s", err)
	}
	policiesFromConfig := d.Get("policies").(map[string]interface{})

	client := m.(*Config).WaapClient
	policiesFromApi, _, diagErr := getPoliciesFromApi(ctx, client, domainID)
	if diagErr != nil {
		return fmt.Errorf("failed to get policies for validation: %v", diagErr)
	}

	// Validate that all specified policy IDs exist
	var nonExistentPolicies []string
	for policyID := range policiesFromConfig {
		if _, exists := policiesFromApi[policyID]; !exists {
			nonExistentPolicies = append(nonExistentPolicies, policyID)
		}
	}
	if len(nonExistentPolicies) > 0 {
		return fmt.Errorf("the following policy IDs do not exist for domain %d: %s", domainID, strings.Join(nonExistentPolicies, ", "))
	}

	return nil
}
