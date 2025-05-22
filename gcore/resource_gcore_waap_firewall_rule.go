package gcore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	waap "github.com/G-Core/gcore-waap-sdk-go"
)

func resourceWaapFirewallRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapFirewallRuleCreate,
		ReadContext:   resourceWaapFirewallRuleRead,
		UpdateContext: resourceWaapFirewallRuleUpdate,
		DeleteContext: resourceWaapFirewallRuleDelete,
		Description:   "Represent WAAP Firewall Rule",

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
				domainIdStr, ruleId, err := resourceFirewallRuleImportParseId(d.Id())
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

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the domain",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the firewall rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the firewall rule",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the rule is enabled",
			},
			"action": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The action that the rule takes when triggered",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The WAAP allows the request",
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
							Description: "The WAAP blocks the request",
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
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "A custom HTTP status code that the WAAP returns if a rule blocks a request",
									},
									"action_duration": {
										Type:        schema.TypeString,
										Optional:    true,
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
					},
				},
			},
			"conditions": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The condition required for the WAAP engine to trigger the rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "IP address condition",
							ExactlyOneOf: []string{
								"conditions.0.ip",
								"conditions.0.ip_range",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:      false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									},
									"ip_address": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "A single IPv4 or IPv6 address to match",
									},
								},
							},
						},
						"ip_range": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "IP range condition",
							ExactlyOneOf: []string{
								"conditions.0.ip",
								"conditions.0.ip_range",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:      false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									},
									"lower_bound": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The lower bound IPv4 or IPv6 address to match against",
									},
									"upper_bound": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The upper bound IPv4 or IPv6 address to match against",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceWaapFirewallRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	domainID := d.Get("domain_id").(int)

	reqBody := waap.CreateFirewallRuleV1DomainsDomainIdFirewallRulesPostJSONRequestBody{
		Name:    d.Get("name").(string),
		Enabled: d.Get("enabled").(bool),
	}

	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		reqBody.Description = &description
	}

	if v, ok := d.GetOk("action"); ok {
		reqBody.Action = parseFirewallActionBlock(v.([]interface{}), false).(waap.FirewallRuleActionInput)
	}

	if v, ok := d.GetOk("conditions"); ok {
		reqBody.Conditions = parseFirewallConditionBlock(v.([]interface{}))
	}

	resp, err := client.CreateFirewallRuleV1DomainsDomainIdFirewallRulesPostWithResponse(ctx, domainID, reqBody)
	if err != nil || resp.StatusCode() != http.StatusCreated {
		return diag.FromErr(fmt.Errorf("error creating WAAP firewall rule with status code: %d, response: %s",
			resp.StatusCode(), string(resp.Body)))
	}

	d.SetId(strconv.Itoa(resp.JSON201.Id))

	return resourceWaapFirewallRuleRead(ctx, d, m)
}

func resourceWaapFirewallRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	domainID := d.Get("domain_id").(int)
	ruleID, _ := strconv.Atoi(d.Id())

	resp, err := client.GetFirewallRuleV1DomainsDomainIdFirewallRulesRuleIdGetWithResponse(ctx, domainID, ruleID)
	if err != nil {
		return diag.Errorf("Failed to read Firewall Rule: %w", err)
	}
	if resp.StatusCode() == http.StatusNotFound {
		d.SetId("") // Resource not found, remove from state
		return diag.Diagnostics{
			{Severity: diag.Warning, Summary: fmt.Sprintf("Firewall Rule (%s) was not found, removed from TF state", ruleID)},
		}
	}
	if resp.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to read Firewall Rule. Status code: %d with error: %s", resp.StatusCode(), resp.Body)
	}

	d.Set("name", resp.JSON200.Name)
	d.Set("description", resp.JSON200.Description)
	d.Set("enabled", resp.JSON200.Enabled)

	actionMap := map[string]interface{}{
		"allow": resp.JSON200.Action.Allow != nil,
	}
	if resp.JSON200.Action.Block != nil {
		actionMap["block"] = []interface{}{
			map[string]interface{}{
				"status_code":     resp.JSON200.Action.Block.StatusCode,
				"action_duration": resp.JSON200.Action.Block.ActionDuration,
			},
		}
	} else {
		actionMap["block"] = []interface{}{}
	}
	d.Set("action", []interface{}{actionMap})

	conditionMap := map[string]interface{}{}
	if resp.JSON200.Conditions[0].Ip != nil {
		ipAddress := marshalStructToJSONString(resp.JSON200.Conditions[0].Ip.IpAddress)
		conditionMap["ip"] = []interface{}{
			map[string]interface{}{
				"ip_address": ipAddress,
				"negation":   resp.JSON200.Conditions[0].Ip.Negation,
			},
		}
	} else {
		conditionMap["ip"] = []interface{}{}
	}

	if resp.JSON200.Conditions[0].IpRange != nil {
		lowerBound := marshalStructToJSONString(resp.JSON200.Conditions[0].IpRange.LowerBound)
		upperBound := marshalStructToJSONString(resp.JSON200.Conditions[0].IpRange.UpperBound)
		conditionMap["ip_range"] = []interface{}{
			map[string]interface{}{
				"lower_bound": lowerBound,
				"upper_bound": upperBound,
				"negation":    resp.JSON200.Conditions[0].IpRange.Negation,
			},
		}
	} else {
		conditionMap["ip_range"] = []interface{}{}
	}
	d.Set("conditions", []interface{}{conditionMap})

	return nil
}

func resourceWaapFirewallRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	domainID := d.Get("domain_id").(int)
	ruleID, _ := strconv.Atoi(d.Id())

	name := d.Get("name").(string)
	enabled := d.Get("enabled").(bool)

	reqBody := waap.UpdateFirewallRuleV1DomainsDomainIdFirewallRulesRuleIdPatchJSONRequestBody{
		Name:    &name,
		Enabled: &enabled,
	}

	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		reqBody.Description = &description
	}

	if d.HasChange("action") {
		if v, ok := d.GetOk("action"); ok {
			updateAction := parseFirewallActionBlock(v.([]interface{}), true).(waap.CustomerRuleActionInput)
			reqBody.Action = &updateAction
		}
	}

	if d.HasChange("conditions") {
		if v, ok := d.GetOk("conditions"); ok {
			conditions := parseFirewallConditionBlock(v.([]interface{}))
			reqBody.Conditions = &conditions
		}
	}

	resp, err := client.UpdateFirewallRuleV1DomainsDomainIdFirewallRulesRuleIdPatchWithResponse(ctx, domainID, ruleID, reqBody)
	if err != nil || (resp.StatusCode() != http.StatusNoContent) {
		return diag.FromErr(fmt.Errorf("error updating WAAP firewall rule '%s', status code: %d, Error: %s",
			name, resp.StatusCode(), string(resp.Body)))
	}

	return resourceWaapFirewallRuleRead(ctx, d, m)
}

func resourceWaapFirewallRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	domainID := d.Get("domain_id").(int)
	ruleID, _ := strconv.Atoi(d.Id())

	resp, err := client.DeleteFirewallRuleV1DomainsDomainIdFirewallRulesRuleIdDeleteWithResponse(ctx, domainID, ruleID)
	if err != nil || (resp.StatusCode() != http.StatusNoContent) {
		return diag.FromErr(fmt.Errorf("error deleting WAAP firewall rule ID '%s', status code: %d, Error: %s",
			ruleID, resp.StatusCode(), resp.Body))
	}

	d.SetId("")
	return nil
}

func parseFirewallActionBlock(actionCfg []interface{}, isUpdate bool) interface{} {
	var createStruct waap.FirewallRuleActionInput

	if len(actionCfg) > 0 && actionCfg[0] != nil {
		actionsMap := actionCfg[0].(map[string]interface{})

		if v, exists := actionsMap["allow"]; exists && v.(bool) {
			createStruct.Allow = &waap.RuleAllowAction{}
		}

		if v, exists := actionsMap["block"].([]interface{}); exists && len(v) > 0 {
			blockAction := waap.RuleBlockAction{}
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
			createStruct.Block = &blockAction
		}

		if isUpdate {
			var updateStruct waap.CustomerRuleActionInput
			updateStruct.Allow = createStruct.Allow
			updateStruct.Block = createStruct.Block

			if v, exists := actionsMap["captcha"]; exists && v.(bool) {
				updateStruct.Captcha = &waap.RuleCaptchaAction{}
			}

			if v, exists := actionsMap["handshake"]; exists && v.(bool) {
				updateStruct.Handshake = &waap.RuleHandshakeAction{}
			}

			if v, exists := actionsMap["monitor"]; exists && v.(bool) {
				updateStruct.Monitor = &waap.RuleMonitorAction{}
			}

			if v, exists := actionsMap["tag"].([]interface{}); exists && len(v) > 0 {
				tagMap := v[0].(map[string]interface{})

				updateStruct.Tag = &waap.RuleTagAction{
					Tags: convertStringList(tagMap["tags"].([]interface{})),
				}
			}
			return updateStruct
		}
	}

	return createStruct
}

func parseFirewallConditionBlock(conditionCfg []interface{}) []waap.FirewallRuleCondition {
	conditionStruct := waap.FirewallRuleCondition{}

	if len(conditionCfg) > 0 && conditionCfg[0] != nil {
		conditionMap := conditionCfg[0].(map[string]interface{})

		if v, exists := conditionMap["ip"].([]interface{}); exists && len(v) > 0 {
			ipStruct := waap.IpCondition{}
			if v[0] != nil {
				ipMap := v[0].(map[string]interface{})

				if v, ok := ipMap["ip_address"]; ok {
					var ipAddress waap.IpCondition_IpAddress
					unmarshalStringToJSONStruct(v.(string), &ipAddress)
					ipStruct.IpAddress = ipAddress
				}

				if v, exists := ipMap["negation"]; exists && v.(bool) {
					negation := v.(bool)
					ipStruct.Negation = &negation
				}
			}
			conditionStruct.Ip = &ipStruct
		}

		if v, exists := conditionMap["ip_range"].([]interface{}); exists && len(v) > 0 {
			ipRangeStruct := waap.IpRangeCondition{}
			if v[0] != nil {
				ipRangeMap := v[0].(map[string]interface{})

				if v, ok := ipRangeMap["lower_bound"]; ok {
					var lowerBound waap.IpRangeCondition_LowerBound
					unmarshalStringToJSONStruct(v.(string), &lowerBound)
					ipRangeStruct.LowerBound = lowerBound
				}

				if v, ok := ipRangeMap["upper_bound"]; ok {
					var upperBound waap.IpRangeCondition_UpperBound
					unmarshalStringToJSONStruct(v.(string), &upperBound)
					ipRangeStruct.UpperBound = upperBound
				}

				if v, exists := ipRangeMap["negation"]; exists && v.(bool) {
					negation := v.(bool)
					ipRangeStruct.Negation = &negation
				}
			}
			conditionStruct.IpRange = &ipRangeStruct
		}
	}

	return []waap.FirewallRuleCondition{conditionStruct}
}

func unmarshalStringToJSONStruct[T any](addr string, target *T) {
	raw := json.RawMessage(`"` + addr + `"`)
	json.Unmarshal(raw, target)
}

func marshalStructToJSONString[T any](input T) string {
	jsonBytes, _ := json.Marshal(input)
	ipAddr := string(jsonBytes)
	ipAddr = strings.ReplaceAll(ipAddr, "\"", "")
	return ipAddr
}

func resourceFirewallRuleImportParseId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected domain_id:rule_id", id)
	}

	return parts[0], parts[1], nil
}
