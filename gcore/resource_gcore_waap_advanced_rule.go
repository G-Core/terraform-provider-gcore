package gcore

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	waap "github.com/G-Core/gcore-waap-sdk-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// This also used for the WAAP Custom Rule resource
var (
	waapActionSchema = &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		MaxItems:    1,
		Description: "The action that the rule takes when triggered.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allow": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "The WAAP allows the request.",
					ExactlyOneOf: []string{
						"action.0.allow",
						"action.0.block",
						"action.0.captcha",
						"action.0.handshake",
						"action.0.monitor",
						"action.0.tag",
					},
				},
				"block": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "The WAAP blocks the request.",
					ExactlyOneOf: []string{
						"action.0.allow",
						"action.0.block",
						"action.0.captcha",
						"action.0.handshake",
						"action.0.monitor",
						"action.0.tag",
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"status_code": {
								Type:     schema.TypeInt,
								Optional: true,
								Default:  403,
								Description: "A custom HTTP status code that the WAAP returns if a rule blocks a request. " +
									"It must be one of these values {403, 405, 418, 429}. Default is 403.",
								ValidateFunc: validation.IntInSlice([]int{403, 405, 418, 429}),
							},
							"action_duration": {
								Type:     schema.TypeString,
								Optional: true,
								Description: "How long a rule's block action will apply to subsequent requests. " +
									"Can be specified in seconds or by using a numeral followed by 's', 'm', 'h', or 'd' " +
									"to represent time format (seconds, minutes, hours, or days). Example: 12h. Must match the pattern ^[0-9]*[smhd]?$",
								ValidateFunc: validation.StringMatch(
									regexp.MustCompile(`^[0-9]+[smhd]?$`),
									"Must be a number optionally followed by 's', 'm', 'h', or 'd' (e.g., 60, 5m, 12h, 1d)",
								),
							},
						},
					},
				},
				"captcha": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "The WAAP requires the user to solve a CAPTCHA challenge.",
					ExactlyOneOf: []string{
						"action.0.allow",
						"action.0.block",
						"action.0.captcha",
						"action.0.handshake",
						"action.0.monitor",
						"action.0.tag",
					},
				},
				"handshake": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "The WAAP performs automatic browser validation.",
					ExactlyOneOf: []string{
						"action.0.allow",
						"action.0.block",
						"action.0.captcha",
						"action.0.handshake",
						"action.0.monitor",
						"action.0.tag",
					},
				},
				"monitor": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "The WAAP monitors the request but took no action.",
					ExactlyOneOf: []string{
						"action.0.allow",
						"action.0.block",
						"action.0.captcha",
						"action.0.handshake",
						"action.0.monitor",
						"action.0.tag",
					},
				},
				"tag": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "The WAAP tags the request.",
					ExactlyOneOf: []string{
						"action.0.allow",
						"action.0.block",
						"action.0.captcha",
						"action.0.handshake",
						"action.0.monitor",
						"action.0.tag",
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"tags": {
								Type:        schema.TypeList,
								Required:    true,
								Description: "The list of user defined tags to tag the request with.",
								MinItems:    1,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
			},
		},
	}
)

func resourceWaapAdvancedRule() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
				domainIdStr, ruleId, err := resourceAdvancedRuleImportParseId(d.Id())
				if err != nil {
					return nil, err
				}

				domainId, err := strconv.ParseInt(domainIdStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("unexpected format of domain_id (%s), expected number", domainIdStr)
				}

				d.Set("domain_id", domainId)
				d.SetId(ruleId)

				return []*schema.ResourceData{d}, nil
			},
		},
		CreateContext: resourceWaapAdvancedRuleCreate,
		ReadContext:   resourceWaapAdvancedRuleRead,
		UpdateContext: resourceWaapAdvancedRuleUpdate,
		DeleteContext: resourceWaapAdvancedRuleDelete,
		Description:   "Represent Advanced Rules for a specific WAAP domain",

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The WAAP domain ID for which the Advanced Rule is configured.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name assigned to the rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description assigned to the rule.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the rule is enabled.",
			},
			"action": waapActionSchema,
			"source": {
				Type:     schema.TypeString,
				Required: true,
				Description: "A CEL syntax expression that contains the rule's conditions. " +
					"Allowed objects are: request, whois, session, response, tags, user_defined_tags, user_agent, client_data. " +
					"More info can be found here: https://gcore.com/docs/waap/waap-rules/advanced-rules",
			},
			"phase": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The WAAP request/response phase for applying the rule. " +
					"The 'access' phase is responsible for modifying the request before it is sent to the origin server. " +
					"The 'header_filter' phase is responsible for modifying the HTTP headers of a response before they are sent back to the client." +
					"The 'body_filter' phase is responsible for modifying the body of a response before it is sent back to the client. Default is 'access'.",
				Default: "access",
				ValidateFunc: validation.StringInSlice([]string{
					"access",
					"header_filter",
					"body_filter",
				}, false),
			},
		},
	}
}

func resourceWaapAdvancedRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start WAAP Advanced Rule creating")

	client := m.(*Config).WaapClient

	req := waap.AdvancedRule{
		Name:    d.Get("name").(string),
		Enabled: d.Get("enabled").(bool),
		Source:  d.Get("source").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		req.Description = &description
	}

	if v, ok := d.GetOk("phase"); ok {
		phase := waap.AdvancedRulePhase(v.(string))
		req.Phase = &phase
	}

	if v, ok := d.GetOk("action"); ok {
		if action := getWaapActionPayload(v); action != nil {
			req.Action = *action
		}
	}

	result, err := client.CreateAdvancedRuleV1DomainsDomainIdAdvancedRulesPostWithResponse(ctx, d.Get("domain_id").(int), req)

	if err != nil {
		return diag.Errorf("Failed to create Advanced Rule: %w", err)
	}

	if result.StatusCode() != http.StatusCreated {
		return diag.Errorf("Failed to create Advanced Rule. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	d.SetId(fmt.Sprintf("%d", result.JSON201.Id))

	log.Printf("[DEBUG] Finish WAAP Advanced Rule creating (id=%s)\n", d.Id())

	return resourceWaapAdvancedRuleRead(ctx, d, m)
}

func resourceWaapAdvancedRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start WAAP Advanced Rule reading")

	client := m.(*Config).WaapClient

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.GetAdvancedRuleV1DomainsDomainIdAdvancedRulesRuleIdGetWithResponse(ctx, d.Get("domain_id").(int), ruleID)
	if err != nil {
		return diag.Errorf("Failed to read Advanced Rule: %w", err)
	}

	if result.StatusCode() == http.StatusNotFound {
		d.SetId("") // Resource not found, remove from state
		return diag.Diagnostics{
			{Severity: diag.Warning, Summary: fmt.Sprintf("Advanced Rule (%s) was not found, removed from TF state", ruleID)},
		}
	}

	if result.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to read Advanced Rule. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	d.Set("name", result.JSON200.Name)
	d.Set("description", result.JSON200.Description)
	d.Set("enabled", result.JSON200.Enabled)
	d.Set("source", result.JSON200.Source)
	d.Set("phase", result.JSON200.Phase)
	d.Set("action", readWaapActionFromResponse(result.JSON200.Action))

	log.Printf("[DEBUG] Finish WAAP Advanced Rule reading (id=%s)\n", ruleID)
	return nil
}

func resourceWaapAdvancedRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Start WAAP Advanced Rule updating (id=%s)\n", d.Id())

	client := m.(*Config).WaapClient

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	enabled := d.Get("enabled").(bool)
	source := d.Get("source").(string)
	req := waap.UpdateAdvancedRule{
		Name:    &name,
		Enabled: &enabled,
		Source:  &source,
	}

	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		req.Description = &description
	}

	if v, ok := d.GetOk("phase"); ok {
		phase := waap.UpdateAdvancedRulePhase(v.(string))
		req.Phase = &phase
	}

	if d.HasChange("action") {
		if action := getWaapActionPayload(d.Get("action")); action != nil {
			req.Action = action
		}
	}

	result, err := client.UpdateAdvancedRuleV1DomainsDomainIdAdvancedRulesRuleIdPatchWithResponse(ctx, d.Get("domain_id").(int), ruleID, req)

	if err != nil {
		return diag.Errorf("Failed to update Advanced Rule: %w", err)
	}

	if result.StatusCode() != http.StatusNoContent {
		return diag.Errorf("Failed to update Advanced Rule. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	log.Printf("[DEBUG] Finish WAAP Advanced Rule updating (id=%s)", ruleID)
	return resourceWaapAdvancedRuleRead(ctx, d, m)
}

func resourceWaapAdvancedRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Start WAAP Advanced Rule deleting (id=%s)\n", d.Id())

	client := m.(*Config).WaapClient

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.DeleteAdvancedRuleV1DomainsDomainIdAdvancedRulesRuleIdDeleteWithResponse(ctx, d.Get("domain_id").(int), ruleID)
	if err != nil {
		return diag.FromErr(err)
	}

	if result.StatusCode() != http.StatusNoContent {
		return diag.Errorf("Failed to delete Advanced Rule. Status code: %d with error: %s", result.StatusCode(), result.Body)
	}

	log.Printf("[DEBUG] Finish WAAP Advanced Rule deleting (id=%s)\n", d.Id())
	d.SetId("")

	return nil
}

// This also used for the WAAP Custom Rule resource
func getWaapActionPayload(actionRaw any) *waap.CustomerRuleActionInput {
	actions := actionRaw.([]interface{})

	if len(actions) > 0 && actions[0] != nil {
		var action waap.CustomerRuleActionInput
		actionsMap := actions[0].(map[string]interface{})
		emptyMap := map[string]interface{}{}

		if v, exists := actionsMap["allow"]; exists && v.(bool) {
			action.Allow = &emptyMap
		}

		if v, exists := actionsMap["block"].([]interface{}); exists && len(v) > 0 {
			blockAction := waap.RuleBlockAction{}
			action.Block = &blockAction

			if v[0] != nil {
				blockMap := v[0].(map[string]interface{})

				if v, ok := blockMap["status_code"]; ok {
					val := waap.RuleBlockStatusCode(v.(int))
					blockAction.StatusCode = &val
				}
				if v, ok := blockMap["action_duration"]; ok {
					val := v.(string)
					blockAction.ActionDuration = &val
				}
			}
		}

		if v, exists := actionsMap["captcha"]; exists && v.(bool) {
			action.Captcha = &emptyMap
		}

		if v, exists := actionsMap["handshake"]; exists && v.(bool) {
			action.Handshake = &emptyMap
		}

		if v, exists := actionsMap["monitor"]; exists && v.(bool) {
			action.Monitor = &emptyMap
		}

		if v, exists := actionsMap["tag"].([]interface{}); exists && len(v) > 0 {
			tagMap := v[0].(map[string]interface{})
			action.Tag = &waap.RuleTagAction{
				Tags: convertStringList(tagMap["tags"].([]interface{})),
			}
		}

		return &action
	}

	return nil
}

func readWaapActionFromResponse(action waap.CustomerRuleActionOutput) []interface{} {
	actionMap := map[string]interface{}{
		"allow":     action.Allow != nil,
		"captcha":   action.Captcha != nil,
		"handshake": action.Handshake != nil,
		"monitor":   action.Monitor != nil,
	}

	if action.Block != nil {
		actionMap["block"] = []interface{}{
			map[string]interface{}{
				"status_code":     action.Block.StatusCode,
				"action_duration": action.Block.ActionDuration,
			},
		}
	} else {
		actionMap["block"] = []interface{}{}
	}

	if action.Tag != nil {
		actionMap["tag"] = []interface{}{
			map[string]interface{}{
				"tags": action.Tag.Tags,
			},
		}
	} else {
		actionMap["tag"] = []interface{}{}
	}

	return []interface{}{actionMap}
}

func resourceAdvancedRuleImportParseId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected domain_id:rule_id", id)
	}

	return parts[0], parts[1], nil
}
