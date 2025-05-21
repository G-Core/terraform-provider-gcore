package gcore

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWaapDomainPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataWaapDomainPolicyRead,
		Description: "Represent WAAP domain policy",
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The WAAP domain ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Domain Policy",
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The policy group name to which the policy belongs",
			},
		},
	}
}

func dataWaapDomainPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start reading WAAP Domain Policies")

	client := m.(*Config).WaapClient

	domainID := d.Get("domain_id").(int)
	policyName := d.Get("name").(string)
	policyGroup := d.Get("group").(string)

	// Get policies from API
	policiesResp, err := client.GetRuleSetListV1DomainsDomainIdRuleSetsGetWithResponse(ctx, domainID)
	if err != nil {
		return diag.FromErr(err)
	}

	if policiesResp.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to read Domain Policies. Status code: %d with error: %s", policiesResp.StatusCode(), policiesResp.Body)
	}

	// Find the policy by name and group
	for _, policySet := range *policiesResp.JSON200 {
		for _, policy := range *policySet.Rules {
			if strings.EqualFold(policy.Group, policyGroup) && strings.EqualFold(policy.Name, policyName) {
				d.SetId(policy.Id)
				break
			}
		}
	}

	if d.Id() == "" {
		return diag.Errorf("Domain Policy with name '%s' and group '%s' not found in domain ID %d", policyName, policyGroup, domainID)
	}

	log.Println("[DEBUG] Finish reading WAAP Domain Policies")
	return nil
}
