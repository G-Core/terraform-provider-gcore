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

func resourceWaapCustomRules() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceCustomRulesCreate,
		ReadContext:   resourceCustomRulesRead,
		UpdateContext: resourceCustomRulesUpdate,
		DeleteContext: resourceCustomRulesDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of the rule.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the rule is enabled.",
			},
			"action": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// "captcha":   {Type: schema.TypeMap, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"allow":     {Type: schema.TypeList, MaxItems: 0, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{}}},
						"captcha":   {Type: schema.TypeList, MaxItems: 0, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{}}},
						"handshake": {Type: schema.TypeList, MaxItems: 0, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{}}},
						"monitor":   {Type: schema.TypeList, MaxItems: 0, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{}}},
						"tag": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tags": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"block": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status_code": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"action_duration": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"conditions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip":                     conditionField("ip_address"),
						"ip_range":               conditionRangeFields("lower_bound", "upper_bound"),
						"url":                    conditionMatch("url"),
						"user_agent":             conditionMatch("user_agent"),
						"header":                 conditionHeader("header", "value"),
						"header_exists":          conditionExists("header"),
						"response_header":        conditionHeader("header", "value"),
						"response_header_exists": conditionExists("header"),
						"http_method":            conditionExists("http_method"),
						"file_extension":         conditionList("file_extension"),
						"content_type":           conditionList("content_type"),
						"country":                conditionList("country_code"),
						"organization":           conditionMatch("organization"),
						"request_rate": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation":         {Type: schema.TypeBool, Optional: true, Default: false},
									"ips":              listString("List of IP addresses."),
									"http_methods":     listString("List of HTTP methods."),
									"path_pattern":     stringField("Path pattern."),
									"requests":         intField("Request count."),
									"time":             intField("Time window."),
									"user_defined_tag": stringField("User-defined tag."),
								},
							},
						},
						"owner_types": conditionList("owner_types"),
						"tags":        conditionList("tags"),
						"session_request_count": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation":      {Type: schema.TypeBool, Optional: true, Default: false},
									"request_count": {Type: schema.TypeInt, Required: true},
								},
							},
						},
						"user_defined_tags": conditionList("tags"),
					},
				},
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the created rule.",
			},
			"rule_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the rule.",
			},
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type ParsedInput struct {
	Action      *waap.CustomerRuleActionInput
	Conditions  *[]waap.CustomRuleConditionInput
	Description *string
	Enabled     *bool
	Name        *string
	RuleType    *string
}

func parseInput(d *schema.ResourceData) (ParsedInput, error) {

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	rule_type := d.Get("rule_type").(string)
	actionInput := d.Get("action").([]interface{})
	conditionsInput := d.Get("conditions").([]interface{})

	log.Printf("[DEBUG] Input: name=%s, description=%s, enabled=%v, type=%+v", name, description, enabled, rule_type)
	log.Printf("[DEBUG] Input: action=%+v", actionInput)
	log.Printf("[DEBUG] Input: conditions=%+v", conditionsInput)

	actionRequest := waap.CustomerRuleActionInput{}

	if len(actionInput) > 0 {
		for key, value := range actionInput[0].(map[string]interface{}) {
			v := value.([]interface{})
			if len(v) != 0 {
				if key == "tag" {
					tags := []string{}
					for _, tag := range v[0].(map[string]interface{})["tags"].([]interface{}) {
						tags = append(tags, tag.(string))
					}

					actionRequest.Tag = &waap.RuleTagAction{
						Tags: tags,
					}
				}

				if key == "block" {
					blockRequest := &waap.RuleBlockAction{}
					if v[0] != nil {
						block := v[0].(map[string]interface{})           // max 1 item
						duration, _ := block["action_duration"].(string) // check if it's a string
						code := waap.RuleBlockStatusCode(block["status_code"].(int))
						blockRequest.ActionDuration = &duration
						blockRequest.StatusCode = &code
					}
					actionRequest.Block = blockRequest
				}

				if key == "allow" {
					actionRequest.Allow = &waap.RuleAllowAction{}
				}

				if key == "captcha" {
					actionRequest.Captcha = &waap.RuleCaptchaAction{}
				}

				if key == "handshake" {
					actionRequest.Handshake = &waap.RuleHandshakeAction{}
				}

				if key == "monitor" {
					actionRequest.Monitor = &waap.RuleMonitorAction{}
				}
			}
		}
	}

	log.Printf("[DEBUG] Parsed: Action: %+v", actionRequest)

	conditions := make([]waap.CustomRuleConditionInput, 0)

	for _, condition := range conditionsInput {

		for key, value := range condition.(map[string]interface{}) {

			for _, item := range value.([]interface{}) {

				if key == "ip" {

					conditionRequest := waap.CustomRuleConditionInput{}
					ip_obj := item.(map[string]interface{})

					var ipAddress waap.IpCondition_IpAddress
					ipAddress.FromIpConditionIpAddress1(ip_obj["ip_address"].(string))

					negation := ip_obj["negation"].(bool)
					ipCondition := waap.IpCondition{
						IpAddress: ipAddress,
						Negation:  &negation,
					}

					conditionRequest.Ip = &ipCondition
					conditions = append(conditions, conditionRequest)
				}

				if key == "http_method" {
					conditionRequest := waap.CustomRuleConditionInput{}
					httpMethodObj := item.(map[string]interface{})

					method, _ := httpMethodObj["http_method"].(string)
					negation := httpMethodObj["negation"].(bool)

					httpMethod := waap.HttpMethodCondition{
						HttpMethod: waap.HTTPMethod(method),
						Negation:   &negation,
					}

					conditionRequest.HttpMethod = &httpMethod
					conditions = append(conditions, conditionRequest)
				}

				if key == "ip_range" {
					conditionRequest := waap.CustomRuleConditionInput{}
					ipRangeObj := item.(map[string]interface{})
					lowerBound := ipRangeObj["lower_bound"].(string)
					upperBound := ipRangeObj["upper_bound"].(string)
					negation := ipRangeObj["negation"].(bool)

					var lowerObj waap.IpRangeCondition_LowerBound
					lowerObj.FromIpRangeConditionLowerBound0(lowerBound)

					var upperObj waap.IpRangeCondition_UpperBound
					upperObj.FromIpRangeConditionUpperBound0(upperBound)

					ipRange := waap.IpRangeCondition{
						LowerBound: lowerObj,
						UpperBound: upperObj,
						Negation:   &negation,
					}

					conditionRequest.IpRange = &ipRange
					conditions = append(conditions, conditionRequest)
				}

				if key == "url" {
					conditionRequest := waap.CustomRuleConditionInput{}
					urlObj := item.(map[string]interface{})

					url := urlObj["url"].(string)
					negation := urlObj["negation"].(bool)
					match_type := waap.UrlConditionMatchType(urlObj["match_type"].(string))

					urlCondition := waap.UrlCondition{
						Url:       url,
						Negation:  &negation,
						MatchType: &match_type,
					}

					conditionRequest.Url = &urlCondition
					conditions = append(conditions, conditionRequest)
				}

				if key == "user_agent" {
					conditionRequest := waap.CustomRuleConditionInput{}
					agentObj := item.(map[string]interface{})

					userAgent := agentObj["user_agent"].(string)
					negation := agentObj["negation"].(bool)
					match_type := waap.UserAgentConditionMatchType(agentObj["match_type"].(string))

					agentCondition := waap.UserAgentCondition{
						UserAgent: userAgent,
						Negation:  &negation,
						MatchType: &match_type,
					}

					conditionRequest.UserAgent = &agentCondition
					conditions = append(conditions, conditionRequest)
				}

				if key == "header" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})

					header := obj["header"].(string)
					value := obj["value"].(string)
					negation := obj["negation"].(bool)
					match_type := waap.HeaderConditionMatchType(obj["match_type"].(string))

					condition := waap.HeaderCondition{
						Header:    header,
						Value:     value,
						Negation:  &negation,
						MatchType: &match_type,
					}

					conditionRequest.Header = &condition
					conditions = append(conditions, conditionRequest)
				}

				if key == "header_exists" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})

					header := obj["header"].(string)
					negation := obj["negation"].(bool)

					condition := waap.HeaderExistsCondition{
						Header:   header,
						Negation: &negation,
					}

					conditionRequest.HeaderExists = &condition
					conditions = append(conditions, conditionRequest)
				}

				if key == "response_header" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})

					header := obj["header"].(string)
					value := obj["value"].(string)
					negation := obj["negation"].(bool)
					match_type := waap.ResponseHeaderConditionMatchType(obj["match_type"].(string))

					condition := waap.ResponseHeaderCondition{
						Header:    header,
						Value:     value,
						Negation:  &negation,
						MatchType: &match_type,
					}

					conditionRequest.ResponseHeader = &condition
					conditions = append(conditions, conditionRequest)
				}

				if key == "response_header_exists" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})

					header := obj["header"].(string)
					negation := obj["negation"].(bool)

					condition := waap.ResponseHeaderExistsCondition{
						Header:   header,
						Negation: &negation,
					}

					conditionRequest.ResponseHeaderExists = &condition
					conditions = append(conditions, conditionRequest)
				}

				if key == "file_extension" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})
					negation := obj["negation"].(bool)
					extensions := []string{}
					for _, ext := range obj["file_extension"].([]interface{}) {
						extensions = append(extensions, ext.(string))
					}
					condition := waap.FileExtensionCondition{
						FileExtension: extensions,
						Negation:      &negation,
					}
					conditionRequest.FileExtension = &condition
					conditions = append(conditions, conditionRequest)
				}

				if key == "content_type" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})
					negation := obj["negation"].(bool)
					contentTypes := []string{}
					for _, ct := range obj["content_type"].([]interface{}) {
						contentTypes = append(contentTypes, ct.(string))
					}
					condition := waap.ContentTypeCondition{
						ContentType: contentTypes,
						Negation:    &negation,
					}
					conditionRequest.ContentType = &condition
					conditions = append(conditions, conditionRequest)
				}

				if key == "country" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})
					negation := obj["negation"].(bool)
					countries := []string{}
					for _, country := range obj["country"].([]interface{}) {
						countries = append(countries, country.(string))
					}
					condition := waap.CountryCondition{
						CountryCode: countries,
						Negation:    &negation,
					}
					conditionRequest.Country = &condition
					conditions = append(conditions, conditionRequest)
				}

				if key == "organization" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})
					negation := obj["negation"].(bool)
					organization := obj["organization"].(string)
					condition := waap.OrganizationCondition{
						Organization: organization,
						Negation:     &negation,
					}
					conditionRequest.Organization = &condition
					conditions = append(conditions, conditionRequest)
				}

				if key == "request_rate" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})
					negation := obj["negation"].(bool)
					pattern := obj["path_pattern"].(string)
					requests := obj["requests"].(int)
					time := obj["time"].(int)
					userDefinedTag := obj["user_defined_tag"].(string)

					ips := make([]waap.RequestRateCondition_Ips_Item, 0)
					for _, ip := range obj["ips"].([]interface{}) {
						var ipAddress waap.RequestRateCondition_Ips_Item
						ipAddress.FromRequestRateConditionIps0(ip.(string))
						ips = append(ips, ipAddress)
					}

					methods := make([]waap.HTTPMethod, 0)
					for _, method := range obj["http_methods"].([]interface{}) {
						methods = append(methods, waap.HTTPMethod(method.(string)))
					}

					condition := waap.RequestRateCondition{
						HttpMethods:    &methods,
						Negation:       &negation,
						PathPattern:    pattern,
						Requests:       requests,
						Time:           time,
						UserDefinedTag: &userDefinedTag,
						Ips:            &ips,
					}

					conditionRequest.RequestRate = &condition
					conditions = append(conditions, conditionRequest)
				}

				if key == "owner_types" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})
					negation := obj["negation"].(bool)
					ownerTypes := []waap.OwnerTypesConditionOwnerTypes{}
					for _, ownerType := range obj["owner_types"].([]interface{}) {
						ownerTypes = append(ownerTypes, waap.OwnerTypesConditionOwnerTypes(ownerType.(string)))
					}
					condition := waap.OwnerTypesCondition{
						OwnerTypes: &ownerTypes,
						Negation:   &negation,
					}
					conditionRequest.OwnerTypes = &condition
					conditions = append(conditions, conditionRequest)
				}

				if key == "tags" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})
					negation := obj["negation"].(bool)
					tags := []string{}
					for _, tag := range obj["tags"].([]interface{}) {
						tags = append(tags, tag.(string))
					}
					condition := waap.TagsCondition{
						Tags:     tags,
						Negation: &negation,
					}
					conditionRequest.Tags = &condition
					conditions = append(conditions, conditionRequest)
				}

				if key == "session_request_count" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})
					negation := obj["negation"].(bool)
					requestCount := obj["request_count"].(int)
					condition := waap.SessionRequestCountCondition{
						RequestCount: requestCount,
						Negation:     &negation,
					}
					conditionRequest.SessionRequestCount = &condition
					conditions = append(conditions, conditionRequest)
				}

				if key == "user_defined_tags" {
					conditionRequest := waap.CustomRuleConditionInput{}
					obj := item.(map[string]interface{})
					negation := obj["negation"].(bool)
					tags := []string{}
					for _, tag := range obj["tags"].([]interface{}) {
						tags = append(tags, tag.(string))
					}
					condition := waap.UserDefinedTagsCondition{
						Tags:     tags,
						Negation: &negation,
					}
					conditionRequest.UserDefinedTags = &condition
					conditions = append(conditions, conditionRequest)
				}

				fmt.Printf("[DEBUG] Condition Key: %+v, Value: %+v\n", key, item)
			}

		}
	}

	log.Printf("[DEBUG] Parsed: Conditions: %+v", conditions)

	parsed := ParsedInput{
		Action:      &actionRequest,
		Conditions:  &conditions,
		Description: &description,
		Enabled:     &enabled,
		Name:        &name,
		RuleType:    &rule_type,
	}

	return parsed, nil
}

func resourceCustomRulesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Println("[DEBUG] Creating security rule")

	client := m.(*Config).WaapClient
	domainID, _ := strconv.Atoi(d.Get("domain_id").(string))
	parsed, _ := parseInput(d)

	createRequestParams := waap.CreateCustomRuleV1DomainsDomainIdCustomRulesPostJSONRequestBody{
		Name:        *parsed.Name,
		Description: parsed.Description,
		Enabled:     *parsed.Enabled,
		Action:      *parsed.Action,
		Conditions:  *parsed.Conditions,
	}

	resp, err := client.CreateCustomRuleV1DomainsDomainIdCustomRulesPostWithResponse(ctx, domainID, createRequestParams)

	if err != nil {
		return diag.Errorf("Error while creating rule '%+v'", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return diag.Errorf("Failed to create custom rule. Status code: %d with error: %s", resp.StatusCode(), resp.Body)
	}

	log.Printf("[DEBUG] Parsed: Response: %+v, err %+v", resp.JSON201, err)

	d.SetId(strconv.Itoa(resp.JSON201.Id))
	return nil
}

func resourceCustomRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Config).WaapClient
	domainID, _ := strconv.Atoi(d.Get("domain_id").(string))
	ruleId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Failed to load rule ID %s", err)
	}

	log.Printf("[DEBUG] Reading custom rule: ID: %+v | Domain ID: %+v", ruleId, domainID)

	resp, err := client.GetCustomRuleV1DomainsDomainIdCustomRulesRuleIdGetWithResponse(ctx, domainID, ruleId)

	if err != nil {
		return diag.Errorf("Error while updating rule '%+v'", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to Update custom rule. Status code: %d with error: %s", resp.StatusCode(), string(resp.Body))
	}

	d.Set("name", resp.JSON200.Name)

	if resp.JSON200.Description != nil {
		d.Set("description", resp.JSON200.Description)
	}

	d.Set("enabled", resp.JSON200.Enabled)

	// jsonBytes, _ := json.Marshal(resp.JSON200.Action)

	// log.Printf("[DEBUG] Actionsss: ID: %+v", string(jsonBytes))
	// log.Printf("[DEBUG] YYYYY: ID: %+v", d.Get("action"))
	// log.Printf("[DEBUG] YYYYY: ID: %+v", d.Get("conditions"))

	// Handle actions
	actionMap := map[string]interface{}{}

	// Handle action.block
	if resp.JSON200.Action.Block != nil {
		blockMap := map[string]interface{}{}

		if resp.JSON200.Action.Block.StatusCode != nil {
			blockMap["status_code"] = *resp.JSON200.Action.Block.StatusCode
		}
		if resp.JSON200.Action.Block.ActionDuration != nil {
			blockMap["action_duration"] = *resp.JSON200.Action.Block.ActionDuration
		}

		actionMap["block"] = []interface{}{blockMap}
	}

	if resp.JSON200.Action.Allow != nil {
		actionMap["allow"] = []interface{}{nil}
	} else {
		actionMap["allow"] = []interface{}{}
	}

	if resp.JSON200.Action.Captcha != nil {
		actionMap["captcha"] = []interface{}{nil}
	} else {
		actionMap["captcha"] = []interface{}{}
	}

	if resp.JSON200.Action.Handshake != nil {
		actionMap["handshake"] = []interface{}{nil}
	} else {
		actionMap["handshake"] = []interface{}{}
	}

	if resp.JSON200.Action.Monitor != nil {
		actionMap["monitor"] = []interface{}{nil}
	} else {
		actionMap["monitor"] = []interface{}{}
	}

	if resp.JSON200.Action.Tag != nil {
		obj := []map[string]interface{}{}
		obj = append(obj, map[string]interface{}{
			"tags": resp.JSON200.Action.Tag.Tags,
		})
		actionMap["tag"] = obj
	} else {
		actionMap["tag"] = []interface{}{}
	}

	d.Set("action", []interface{}{actionMap})

	if resp.JSON200.Conditions != nil {

		var conditionsList []map[string]interface{}

		for _, condition := range resp.JSON200.Conditions {

			conditionMap := map[string]interface{}{}

			if condition.Ip != nil {
				ipMap := map[string]interface{}{}
				ipMap["ip"] = condition.Ip.IpAddress
				ipMap["negation"] = condition.Ip.Negation
				conditionMap["ip"] = []interface{}{ipMap}
			} else {
				conditionMap["ip"] = []interface{}{}
			}

			if condition.IpRange != nil {
				ipRangeMap := map[string]interface{}{}
				ipRangeMap["lower_bound"] = condition.IpRange.LowerBound
				ipRangeMap["upper_bound"] = condition.IpRange.UpperBound
				ipRangeMap["negation"] = condition.IpRange.Negation
				conditionMap["ip_range"] = []interface{}{ipRangeMap}
			} else {
				conditionMap["ip_range"] = []interface{}{}
			}

			if condition.Url != nil {
				urlMap := map[string]interface{}{}
				urlMap["url"] = condition.Url.Url
				urlMap["negation"] = condition.Url.Negation
				urlMap["match_type"] = condition.Url.MatchType
				conditionMap["url"] = []interface{}{urlMap}
			} else {
				conditionMap["url"] = []interface{}{}
			}

			if condition.UserAgent != nil {
				userAgentMap := map[string]interface{}{}
				userAgentMap["user_agent"] = condition.UserAgent.UserAgent
				userAgentMap["match_type"] = condition.UserAgent.MatchType
				userAgentMap["negation"] = condition.UserAgent.Negation
				conditionMap["user_agnet"] = []interface{}{userAgentMap}
			} else {
				conditionMap["user_agnet"] = []interface{}{}
			}

			if condition.Header != nil {
				headerMap := map[string]interface{}{}
				headerMap["header"] = condition.Header.Header
				headerMap["value"] = condition.Header.Value
				headerMap["negation"] = condition.Header.Negation
				headerMap["match_type"] = condition.Header.MatchType
				conditionMap["header"] = []interface{}{headerMap}
			} else {
				conditionMap["header"] = []interface{}{}
			}

			if condition.HeaderExists != nil {
				headerExistsMap := map[string]interface{}{}
				headerExistsMap["header"] = condition.HeaderExists.Header
				headerExistsMap["negation"] = condition.HeaderExists.Negation
				conditionMap["header_exists"] = []interface{}{headerExistsMap}
			} else {
				conditionMap["header_exists"] = []interface{}{}
			}

			if condition.ResponseHeader != nil {
				respHeaderMap := map[string]interface{}{}
				respHeaderMap["header"] = condition.ResponseHeader.Header
				respHeaderMap["value"] = condition.ResponseHeader.Value
				respHeaderMap["negation"] = condition.ResponseHeader.Negation
				respHeaderMap["match_type"] = condition.ResponseHeader.MatchType
				conditionMap["response_header"] = []interface{}{respHeaderMap}
			} else {
				conditionMap["response_header"] = []interface{}{}
			}

			if condition.ResponseHeaderExists != nil {
				responseExistsMap := map[string]interface{}{}
				responseExistsMap["header"] = condition.ResponseHeaderExists.Header
				responseExistsMap["negation"] = condition.ResponseHeaderExists.Negation
				conditionMap["response_header_exists"] = []interface{}{responseExistsMap}
			} else {
				conditionMap["response_header_exists"] = []interface{}{}
			}

			if condition.HttpMethod != nil {
				httpMethodMap := map[string]interface{}{}
				httpMethodMap["http_method"] = condition.HttpMethod.HttpMethod
				httpMethodMap["negation"] = condition.HttpMethod.Negation
				conditionMap["http_method"] = []interface{}{httpMethodMap}
			} else {
				conditionMap["http_method"] = []interface{}{}
			}

			if condition.FileExtension != nil {
				if len(condition.FileExtension.FileExtension) > 0 {
					fileExtensionMap := map[string]interface{}{}
					fileExtensionMap["file_extension"] = condition.FileExtension.FileExtension
					fileExtensionMap["negation"] = condition.FileExtension.Negation
					conditionMap["file_extension"] = []interface{}{fileExtensionMap}
				} else {
					conditionMap["file_extension"] = []interface{}{}
				}
			}

			if condition.ContentType != nil {
				if len(condition.ContentType.ContentType) > 0 {
					contentTypeMap := map[string]interface{}{}
					contentTypeMap["content_type"] = condition.ContentType.ContentType
					contentTypeMap["negation"] = condition.ContentType.Negation
					conditionMap["content_type"] = []interface{}{contentTypeMap}
				} else {
					conditionMap["content_type"] = []interface{}{}
				}
			}

			if condition.Country != nil {
				if len(condition.Country.CountryCode) > 0 {
					countryMap := map[string]interface{}{}
					countryMap["country_code"] = condition.Country.CountryCode
					countryMap["negation"] = condition.Country.Negation
					conditionMap["country"] = []interface{}{countryMap}
				} else {
					conditionMap["country"] = []interface{}{}
				}
			}

			if condition.Organization != nil {
				orgMap := map[string]interface{}{}
				orgMap["organization"] = condition.Organization.Organization
				orgMap["negation"] = condition.Organization.Negation
				conditionMap["organization"] = []interface{}{orgMap}
			} else {
				conditionMap["organization"] = []interface{}{}
			}

			if condition.RequestRate != nil {

				requestRateMap := map[string]interface{}{}
				requestRateMap["negation"] = condition.RequestRate.Negation

				if len(*condition.RequestRate.Ips) > 0 {
					requestRateMap["ips"] = condition.RequestRate.Ips
				}
				if len(*condition.RequestRate.HttpMethods) > 0 {
					requestRateMap["http_methods"] = condition.RequestRate.HttpMethods
				}

				requestRateMap["path_pattern"] = condition.RequestRate.PathPattern
				requestRateMap["requests"] = condition.RequestRate.Requests
				requestRateMap["time"] = condition.RequestRate.Time

				if condition.RequestRate.UserDefinedTag != nil {
					requestRateMap["user_defined_tag"] = *condition.RequestRate.UserDefinedTag
				}

				conditionMap["request_rate"] = []interface{}{requestRateMap}
			} else {
				conditionMap["request_rate"] = []interface{}{}
			}

			if condition.OwnerTypes != nil {
				if len(*condition.OwnerTypes.OwnerTypes) > 0 {
					ownerMap := map[string]interface{}{}
					ownerMap["owner_types"] = condition.OwnerTypes.OwnerTypes
					ownerMap["negation"] = condition.OwnerTypes.Negation
					conditionMap["owner_types"] = []interface{}{ownerMap}
				} else {
					conditionMap["owner_types"] = []interface{}{}
				}
			}

			if len(condition.Tags.Tags) > 0 {
				tagMaps := map[string]interface{}{}
				tagMaps["tags"] = condition.Tags.Tags
				tagMaps["negation"] = condition.Tags.Negation
				conditionMap["tags"] = []interface{}{tagMaps}
			} else {
				conditionMap["tags"] = []interface{}{}
			}

			if condition.SessionRequestCount != nil {
				sessReqMap := map[string]interface{}{}
				sessReqMap["negation"] = condition.SessionRequestCount.Negation
				sessReqMap["request_count"] = condition.SessionRequestCount.RequestCount
				conditionMap["session_request_count"] = []interface{}{sessReqMap}
			} else {
				conditionMap["session_request_count"] = []interface{}{}
			}

			if len(condition.UserDefinedTags.Tags) > 0 {
				tagMaps := map[string]interface{}{}
				tagMaps["tags"] = condition.Tags.Tags
				tagMaps["negation"] = condition.Tags.Negation
				conditionMap["user_defined_tags"] = []interface{}{tagMaps}
			} else {
				conditionMap["user_defined_tags"] = []interface{}{}
			}

			conditionsList = append(conditionsList, conditionMap)
		}

		if len(conditionsList) > 0 {
			d.Set("conditions", conditionsList)
		}

	}

	log.Printf("[DEBUG] curr state: %+v", string(resp.Body))

	return nil
}

func resourceCustomRulesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Config).WaapClient
	parsed, _ := parseInput(d)
	domainID, _ := strconv.Atoi(d.Get("domain_id").(string))
	ruleId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Failed to load rule ID %s", err)
	}

	log.Printf("[DEBUG] Updating custom rule: ID: %+v | Domain ID: %+v", ruleId, domainID)

	updateRequest := waap.UpdateCustomRuleV1DomainsDomainIdCustomRulesRuleIdPatchJSONRequestBody{
		Action:      parsed.Action,
		Conditions:  parsed.Conditions,
		Description: parsed.Description,
		Enabled:     parsed.Enabled,
		Name:        parsed.Name,
	}

	resp, err := client.UpdateCustomRuleV1DomainsDomainIdCustomRulesRuleIdPatchWithResponse(ctx, domainID, ruleId, updateRequest)

	if err != nil {
		return diag.Errorf("Error while updating rule '%+v'", err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		return diag.Errorf("Failed to Update custom rule. Status code: %d with error: %s", resp.StatusCode(), string(resp.Body))
	}

	return resourceCustomRulesRead(ctx, d, m) // Refresh state
}

func resourceCustomRulesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Config).WaapClient

	ruleId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Failed to load rule ID %s", err)
	}

	domainID, _ := strconv.Atoi(d.Get("domain_id").(string))
	log.Printf("[DEBUG] Deleting custom rule: ID: %+v | Domain ID: %+v", ruleId, domainID)

	resp, err := client.DeleteCustomRuleV1DomainsDomainIdCustomRulesRuleIdDeleteWithResponse(ctx, domainID, ruleId)

	if err != nil {
		return diag.Errorf("Error while deleting rule '%+v'", err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		return diag.Errorf("Failed to Delete custom rule. Status code: %d with error: %s", resp.StatusCode(), resp.Body)
	}

	d.SetId("")
	return nil
}
