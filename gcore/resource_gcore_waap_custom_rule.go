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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWaapCustomRule() *schema.Resource {
	return &schema.Resource{
		Importer:      &schema.ResourceImporter{State: importWaapRule},
		CreateContext: resourceCustomRulesCreate,
		ReadContext:   resourceCustomRulesRead,
		UpdateContext: resourceCustomRulesUpdate,
		DeleteContext: resourceCustomRulesDelete,
		Description:   "Represent Custom Rules for a specific WAAP domain",
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The WAAP domain ID for which the Custom Rule is configured.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name assigned to the rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The description assigned to the rule.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the rule is enabled.",
			},
			"action": waapActionSchema,
			"conditions": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Description: "The conditions required for the WAAP engine to trigger the rule. " +
					"Rules may have between 1 and 5 conditions. All conditions must pass for the rule to trigger.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "IP address condition. This condition matches a single IP address.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"ip_address": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "A single IPv4 or IPv6 address",
									},
								},
							},
						},
						"ip_range": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "IP range condition. This condition matches a range of IP addresses.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"lower_bound": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The lower bound IPv4 or IPv6 address to match against.",
									},
									"upper_bound": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The upper bound IPv4 or IPv6 address to match against.",
									},
								},
							},
						},
						"url": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "URL condition. This condition matches a URL path.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The URL to match.",
									},
									"match_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "Contains",
										Description:  "The type of matching condition. Valid values are 'Exact', 'Contains', and 'Regex'. Default is 'Contains'.",
										ValidateFunc: validation.StringInSlice([]string{"Exact", "Contains", "Regex"}, false),
									},
								},
							},
						},
						"user_agent": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "User agent condition. This condition matches the user agent of the request.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"user_agent": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The user agent value to match.",
									},
									"match_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "Contains",
										Description:  "The type of matching condition. Valid values are 'Exact', 'Contains'. Default is 'Contains'.",
										ValidateFunc: validation.StringInSlice([]string{"Exact", "Contains"}, false),
									},
								},
							},
						},
						"header": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Request header condition. This condition matches a request header and its value.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"header": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The request header name.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The request header value.",
									},
									"match_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "Contains",
										Description:  "The type of matching condition. Valid values are 'Exact', 'Contains'. Default is 'Contains'.",
										ValidateFunc: validation.StringInSlice([]string{"Exact", "Contains"}, false),
									},
								},
							},
						},
						"header_exists": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Request header exists condition. This condition checks if a request header exists.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"header": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The request header name.",
									},
								},
							},
						},
						"response_header": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"header": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The request header name.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The request header value.",
									},
									"match_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "Contains",
										Description:  "The type of matching condition. Valid values are 'Exact', 'Contains'. Default is 'Contains'.",
										ValidateFunc: validation.StringInSlice([]string{"Exact", "Contains"}, false),
									}},
							},
						},
						"response_header_exists": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Response header exists condition. This condition checks if a response header exists.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"header": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The request header name.",
									},
								},
							},
						},
						"http_method": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "HTTP method condition. This condition matches the HTTP method of the request.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"http_method": {
										Type:     schema.TypeString,
										Required: true,
										Description: "The HTTP method to match against. " +
											"Valid values are 'CONNECT', 'DELETE', 'GET', 'HEAD', 'OPTIONS', 'PATCH', 'POST', 'PUT', and 'TRACE'.",
										ValidateFunc: validation.StringInSlice([]string{
											"CONNECT",
											"DELETE",
											"GET",
											"HEAD",
											"OPTIONS",
											"PATCH",
											"POST",
											"PUT",
											"TRACE",
										}, false),
									},
								},
							},
						},
						"file_extension": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "File extension condition. This condition matches the file extension of the request.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"file_extension": {
										Type:        schema.TypeSet,
										Required:    true,
										MinItems:    1,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The list of file extensions to match against.",
									},
								},
							},
						},
						"content_type": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Content type condition. This condition matches the content type of the request.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"content_type": {
										Type:        schema.TypeSet,
										Required:    true,
										MinItems:    1,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The list of content types to match against.",
									},
								},
							},
						},
						"country": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Country condition. This condition matches the country of the request based on the source IP address.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"country_code": {
										Type:        schema.TypeSet,
										Required:    true,
										MinItems:    1,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "A list of ISO 3166-1 alpha-2 formatted strings representing the countries to match against.",
									},
								},
							},
						},
						"organization": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Organization condition. This condition matches the organization of the request based on the source IP address.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"organization": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The organization to match against.",
									},
								},
							},
						},
						"request_rate": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Request rate condition. This condition matches the request rate.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ips": {
										Type:        schema.TypeSet,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "A list of source IPs that can trigger a request rate condition.",
									},
									"http_methods": {
										Type:     schema.TypeSet,
										Optional: true,
										Description: "Possible HTTP request methods that can trigger a request rate condition. " +
											"Valid values are 'CONNECT', 'DELETE', 'GET', 'HEAD', 'OPTIONS', 'PATCH', 'POST', 'PUT', and 'TRACE'.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"CONNECT",
												"DELETE",
												"GET",
												"HEAD",
												"OPTIONS",
												"PATCH",
												"POST",
												"PUT",
												"TRACE",
											}, false),
										},
									},
									"path_pattern": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "A regular expression matching the URL path of the incoming request.",
									},
									"requests": {
										Type:         schema.TypeInt,
										Required:     true,
										Description:  "The number of incoming requests over the given time that can trigger a request rate condition.",
										ValidateFunc: validation.IntAtLeast(20),
									},
									"time": {
										Type:         schema.TypeInt,
										Required:     true,
										Description:  "The number of seconds that the WAAP measures incoming requests over before triggering a request rate condition.",
										ValidateFunc: validation.IntAtLeast(1),
									},
									"user_defined_tag": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A user-defined tag that can be included in incoming requests and used to trigger a request rate condition.",
									},
								},
							},
						},
						"owner_types": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"owner_types": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Description: "Match the type of organization that owns the IP address making an incoming request. " +
											"Valid values are 'COMMERCIAL', 'EDUCATIONAL', 'GOVERNMENT', 'HOSTING_SERVICES', 'ISP', 'MOBILE_NETWORK', 'NETWORK', and 'RESERVED'.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"COMMERCIAL",
												"EDUCATIONAL",
												"GOVERNMENT",
												"HOSTING_SERVICES",
												"ISP",
												"MOBILE_NETWORK",
												"NETWORK",
												"RESERVED",
											}, false),
										},
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Tags condition. This condition matches the request tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"tags": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Description: "A list of tags to match against the request tags. " +
											"Tags can be obtained from the API endpoint /v1/tags or you can use the gcore_waap_tag data source.",
									},
								},
							},
						},
						"session_request_count": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Session request count condition. This condition matches the number of dynamic requests in the session.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"request_count": {
										Type:         schema.TypeInt,
										Required:     true,
										Description:  "The number of dynamic requests in the session.",
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"user_defined_tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"negation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not to apply a boolean NOT operation to the rule's condition.",
									},
									"tags": {
										Type:        schema.TypeSet,
										Required:    true,
										MinItems:    1,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "A list of user-defined tags to match against the request tags.",
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

func resourceCustomRulesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start WAAP Custom Rule creating")

	client := m.(*Config).WaapClient

	req := waap.CustomRule{
		Name:    d.Get("name").(string),
		Enabled: d.Get("enabled").(bool),
	}

	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		req.Description = &description
	}

	if v, ok := d.GetOk("action"); ok {
		if action := getWaapActionPayload(v); action != nil {
			req.Action = *action
		}
	}

	if v, ok := d.GetOk("conditions"); ok {
		req.Conditions = getConditionsPaylod(v)
	}

	resp, err := client.CreateCustomRuleV1DomainsDomainIdCustomRulesPostWithResponse(ctx, d.Get("domain_id").(int), req)
	if err != nil {
		return diag.Errorf("Failed to create Custom Rule: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return diag.Errorf("Failed to create Custom Rule. Status code: %d with error: %s", resp.StatusCode(), resp.Body)
	}

	d.SetId(strconv.Itoa(resp.JSON201.Id))

	log.Printf("[DEBUG] Finish WAAP Custom Rule creating (id=%s)\n", d.Id())

	return resourceCustomRulesRead(ctx, d, m)
}

func resourceCustomRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start WAAP Custom Rule reading (id=%s)", d.Id())

	client := m.(*Config).WaapClient

	ruleId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Failed to convert rule ID %s", err)
	}

	resp, err := client.GetCustomRuleV1DomainsDomainIdCustomRulesRuleIdGetWithResponse(ctx, d.Get("domain_id").(int), ruleId)
	if err != nil {
		return diag.Errorf("Failed to read Custom Rule: %w", err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		d.SetId("") // Resource not found, remove from state
		return diag.Diagnostics{
			{Severity: diag.Warning, Summary: fmt.Sprintf("Custom Rule (%s) was not found, removed from TF state", ruleId)},
		}
	}

	if resp.StatusCode() != http.StatusOK {
		return diag.Errorf("Failed to read Custom Rule. Status code: %d with error: %s", resp.StatusCode(), resp.Body)
	}

	d.Set("name", resp.JSON200.Name)
	d.Set("description", resp.JSON200.Description)
	d.Set("enabled", resp.JSON200.Enabled)
	d.Set("action", readWaapActionFromResponse(resp.JSON200.Action))

	err = d.Set("conditions", readConditionsFromResponse(resp.JSON200.Conditions))
	if err != nil {
		return diag.Errorf("Failed to save conditions to the state: %s", err)
	}

	log.Println("[DEBUG] Finish WAAP Custom Rule reading (id=%s)", ruleId)

	return nil
}

func resourceCustomRulesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start WAAP Custom Rule updating (id=%s)", d.Id())

	client := m.(*Config).WaapClient

	ruleId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Failed to convert rule ID %s", err)
	}

	name := d.Get("name").(string)
	enabled := d.Get("enabled").(bool)
	description := d.Get("description").(string)
	req := waap.UpdateCustomRule{
		Name:        &name,
		Enabled:     &enabled,
		Description: &description,
	}

	if d.HasChange("action") {
		if action := getWaapActionPayload(d.Get("action")); action != nil {
			req.Action = action
		}
	}

	if d.HasChange("conditions") {
		conditions := getConditionsPaylod(d.Get("conditions"))
		req.Conditions = &conditions
	}

	resp, err := client.UpdateCustomRuleV1DomainsDomainIdCustomRulesRuleIdPatchWithResponse(ctx, d.Get("domain_id").(int), ruleId, req)
	if err != nil {
		return diag.Errorf("Failed to update Custom Rule: %w", err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		return diag.Errorf("Failed to update Custom Rule. Status code: %d with error: %s", resp.StatusCode(), resp.Body)
	}

	log.Printf("[DEBUG] Finish WAAP Custom Rule updating (id=%s)", ruleId)
	return resourceCustomRulesRead(ctx, d, m)
}

func resourceCustomRulesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Start WAAP Custom Rule deleting (id=%s)\n", d.Id())

	client := m.(*Config).WaapClient

	ruleId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Failed to convert rule ID %s", err)
	}

	resp, err := client.DeleteCustomRuleV1DomainsDomainIdCustomRulesRuleIdDeleteWithResponse(ctx, d.Get("domain_id").(int), ruleId)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		return diag.Errorf("Failed to delete Custom Rule. Status code: %d with error: %s", resp.StatusCode(), resp.Body)
	}

	log.Printf("[DEBUG] Finish WAAP Custom Rule deleting (id=%s)\n", d.Id())
	d.SetId("")

	return nil
}

func getConditionsPaylod(conditionsRaw any) []waap.CustomRuleConditionInput {
	conditions := conditionsRaw.([]interface{})
	result := []waap.CustomRuleConditionInput{}

	if len(conditions) == 0 || conditions[0] == nil {
		return result
	}

	for key, value := range conditions[0].(map[string]interface{}) {
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
				result = append(result, conditionRequest)
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
				result = append(result, conditionRequest)
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
				result = append(result, conditionRequest)
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
				result = append(result, conditionRequest)
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
				result = append(result, conditionRequest)
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
				result = append(result, conditionRequest)
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
				result = append(result, conditionRequest)
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
				result = append(result, conditionRequest)
			}

			if key == "http_method" {
				conditionRequest := waap.CustomRuleConditionInput{}
				httpMethodObj := item.(map[string]interface{})

				method := httpMethodObj["http_method"].(string)
				negation := httpMethodObj["negation"].(bool)

				httpMethod := waap.HttpMethodCondition{
					HttpMethod: waap.HTTPMethod(method),
					Negation:   &negation,
				}

				conditionRequest.HttpMethod = &httpMethod
				result = append(result, conditionRequest)
			}

			if key == "file_extension" {
				conditionRequest := waap.CustomRuleConditionInput{}
				obj := item.(map[string]interface{})

				negation := obj["negation"].(bool)
				condition := waap.FileExtensionCondition{
					FileExtension: convertSchemaSetToStringList(obj["file_extension"].(*schema.Set)),
					Negation:      &negation,
				}

				conditionRequest.FileExtension = &condition
				result = append(result, conditionRequest)
			}

			if key == "content_type" {
				conditionRequest := waap.CustomRuleConditionInput{}
				obj := item.(map[string]interface{})

				negation := obj["negation"].(bool)
				condition := waap.ContentTypeCondition{
					ContentType: convertSchemaSetToStringList(obj["content_type"].(*schema.Set)),
					Negation:    &negation,
				}

				conditionRequest.ContentType = &condition
				result = append(result, conditionRequest)
			}

			if key == "country" {
				conditionRequest := waap.CustomRuleConditionInput{}
				obj := item.(map[string]interface{})

				negation := obj["negation"].(bool)
				condition := waap.CountryCondition{
					CountryCode: convertSchemaSetToStringList(obj["country_code"].(*schema.Set)),
					Negation:    &negation,
				}

				conditionRequest.Country = &condition
				result = append(result, conditionRequest)
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
				result = append(result, conditionRequest)
			}

			if key == "request_rate" {
				conditionRequest := waap.CustomRuleConditionInput{}
				obj := item.(map[string]interface{})

				pattern := obj["path_pattern"].(string)
				requests := obj["requests"].(int)
				time := obj["time"].(int)

				ips := make([]waap.RequestRateCondition_Ips_Item, 0)
				for _, ip := range obj["ips"].(*schema.Set).List() {
					var ipAddress waap.RequestRateCondition_Ips_Item
					ipAddress.FromRequestRateConditionIps0(ip.(string))
					ips = append(ips, ipAddress)
				}

				methods := make([]waap.HTTPMethod, 0)
				for _, method := range obj["http_methods"].(*schema.Set).List() {
					methods = append(methods, waap.HTTPMethod(method.(string)))
				}

				condition := waap.RequestRateCondition{
					HttpMethods: &methods,
					PathPattern: pattern,
					Requests:    requests,
					Time:        time,
					Ips:         &ips,
				}

				if v, exists := obj["user_defined_tag"]; exists && v != "" {
					userDefinedTag := v.(string)
					condition.UserDefinedTag = &userDefinedTag
				}

				conditionRequest.RequestRate = &condition
				result = append(result, conditionRequest)
			}

			if key == "owner_types" {
				conditionRequest := waap.CustomRuleConditionInput{}
				obj := item.(map[string]interface{})

				negation := obj["negation"].(bool)
				ownerTypes := []waap.OwnerTypesConditionOwnerTypes{}
				for _, ownerType := range obj["owner_types"].(*schema.Set).List() {
					ownerTypes = append(ownerTypes, waap.OwnerTypesConditionOwnerTypes(ownerType.(string)))
				}

				condition := waap.OwnerTypesCondition{
					OwnerTypes: &ownerTypes,
					Negation:   &negation,
				}

				conditionRequest.OwnerTypes = &condition
				result = append(result, conditionRequest)
			}

			if key == "tags" {
				conditionRequest := waap.CustomRuleConditionInput{}
				obj := item.(map[string]interface{})

				negation := obj["negation"].(bool)
				condition := waap.TagsCondition{
					Tags:     convertSchemaSetToStringList(obj["tags"].(*schema.Set)),
					Negation: &negation,
				}

				conditionRequest.Tags = &condition
				result = append(result, conditionRequest)
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
				result = append(result, conditionRequest)
			}

			if key == "user_defined_tags" {
				conditionRequest := waap.CustomRuleConditionInput{}
				obj := item.(map[string]interface{})

				negation := obj["negation"].(bool)
				condition := waap.UserDefinedTagsCondition{
					Tags:     convertSchemaSetToStringList(obj["tags"].(*schema.Set)),
					Negation: &negation,
				}

				conditionRequest.UserDefinedTags = &condition
				result = append(result, conditionRequest)
			}
		}
	}

	return result
}

func readConditionsFromResponse(conditions []waap.CustomRuleConditionOutput) []interface{} {
	conditionMap := map[string]interface{}{
		"ip":                     []interface{}{},
		"ip_range":               []interface{}{},
		"url":                    []interface{}{},
		"user_agent":             []interface{}{},
		"header":                 []interface{}{},
		"header_exists":          []interface{}{},
		"response_header":        []interface{}{},
		"response_header_exists": []interface{}{},
		"http_method":            []interface{}{},
		"file_extension":         []interface{}{},
		"content_type":           []interface{}{},
		"country":                []interface{}{},
		"organization":           []interface{}{},
		"request_rate":           []interface{}{},
		"owner_types":            []interface{}{},
		"tags":                   []interface{}{},
		"session_request_count":  []interface{}{},
		"user_defined_tags":      []interface{}{},
	}

	for _, condition := range conditions {
		if condition.Ip != nil {
			ipMap := map[string]interface{}{}
			ipMap["ip_address"] = marshalStructToJSONString(condition.Ip.IpAddress)
			ipMap["negation"] = condition.Ip.Negation
			conditionMap["ip"] = append(conditionMap["ip"].([]interface{}), ipMap)
		}

		if condition.IpRange != nil {
			ipRangeMap := map[string]interface{}{}
			ipRangeMap["lower_bound"] = marshalStructToJSONString(condition.IpRange.LowerBound)
			ipRangeMap["upper_bound"] = marshalStructToJSONString(condition.IpRange.UpperBound)
			ipRangeMap["negation"] = condition.IpRange.Negation
			conditionMap["ip_range"] = append(conditionMap["ip_range"].([]interface{}), ipRangeMap)
		}

		if condition.Url != nil {
			urlMap := map[string]interface{}{}
			urlMap["url"] = condition.Url.Url
			urlMap["negation"] = condition.Url.Negation
			urlMap["match_type"] = condition.Url.MatchType
			conditionMap["url"] = append(conditionMap["url"].([]interface{}), urlMap)
		}

		if condition.UserAgent != nil {
			userAgentMap := map[string]interface{}{}
			userAgentMap["user_agent"] = condition.UserAgent.UserAgent
			userAgentMap["match_type"] = condition.UserAgent.MatchType
			userAgentMap["negation"] = condition.UserAgent.Negation
			conditionMap["user_agent"] = append(conditionMap["user_agent"].([]interface{}), userAgentMap)
		}

		if condition.Header != nil {
			headerMap := map[string]interface{}{}
			headerMap["header"] = condition.Header.Header
			headerMap["value"] = condition.Header.Value
			headerMap["negation"] = condition.Header.Negation
			headerMap["match_type"] = condition.Header.MatchType
			conditionMap["header"] = append(conditionMap["header"].([]interface{}), headerMap)
		}

		if condition.HeaderExists != nil {
			headerExistsMap := map[string]interface{}{}
			headerExistsMap["header"] = condition.HeaderExists.Header
			headerExistsMap["negation"] = condition.HeaderExists.Negation
			conditionMap["header_exists"] = append(conditionMap["header_exists"].([]interface{}), headerExistsMap)
		}

		if condition.ResponseHeader != nil {
			respHeaderMap := map[string]interface{}{}
			respHeaderMap["header"] = condition.ResponseHeader.Header
			respHeaderMap["value"] = condition.ResponseHeader.Value
			respHeaderMap["negation"] = condition.ResponseHeader.Negation
			respHeaderMap["match_type"] = condition.ResponseHeader.MatchType
			conditionMap["response_header"] = append(conditionMap["response_header"].([]interface{}), respHeaderMap)
		}

		if condition.ResponseHeaderExists != nil {
			responseExistsMap := map[string]interface{}{}
			responseExistsMap["header"] = condition.ResponseHeaderExists.Header
			responseExistsMap["negation"] = condition.ResponseHeaderExists.Negation
			conditionMap["response_header_exists"] = append(conditionMap["response_header_exists"].([]interface{}), responseExistsMap)
		}

		if condition.HttpMethod != nil {
			httpMethodMap := map[string]interface{}{}
			httpMethodMap["http_method"] = condition.HttpMethod.HttpMethod
			httpMethodMap["negation"] = condition.HttpMethod.Negation
			conditionMap["http_method"] = append(conditionMap["http_method"].([]interface{}), httpMethodMap)
		}

		if condition.FileExtension != nil && len(condition.FileExtension.FileExtension) > 0 {
			fileExtensionMap := map[string]interface{}{}
			fileExtensionMap["file_extension"] = condition.FileExtension.FileExtension
			fileExtensionMap["negation"] = condition.FileExtension.Negation
			conditionMap["file_extension"] = append(conditionMap["file_extension"].([]interface{}), fileExtensionMap)
		}

		if condition.ContentType != nil && len(condition.ContentType.ContentType) > 0 {
			contentTypeMap := map[string]interface{}{}
			contentTypeMap["content_type"] = condition.ContentType.ContentType
			contentTypeMap["negation"] = condition.ContentType.Negation
			conditionMap["content_type"] = append(conditionMap["content_type"].([]interface{}), contentTypeMap)
		}

		if condition.Country != nil && len(condition.Country.CountryCode) > 0 {
			countryMap := map[string]interface{}{}
			countryMap["country_code"] = condition.Country.CountryCode
			countryMap["negation"] = condition.Country.Negation
			conditionMap["country"] = append(conditionMap["country"].([]interface{}), countryMap)
		}

		if condition.Organization != nil {
			orgMap := map[string]interface{}{}
			orgMap["organization"] = condition.Organization.Organization
			orgMap["negation"] = condition.Organization.Negation
			conditionMap["organization"] = append(conditionMap["organization"].([]interface{}), orgMap)
		}

		if condition.RequestRate != nil {
			requestRateMap := map[string]interface{}{}

			if condition.RequestRate.Ips != nil && len(*condition.RequestRate.Ips) > 0 {
				var marshaledIps []string
				for _, ip := range *condition.RequestRate.Ips {
					marshaledIps = append(marshaledIps, marshalStructToJSONString(ip))
				}
				requestRateMap["ips"] = marshaledIps
			}

			if condition.RequestRate.HttpMethods != nil && len(*condition.RequestRate.HttpMethods) > 0 {
				var httpMethods []string
				for _, httpMethod := range *condition.RequestRate.HttpMethods {
					httpMethods = append(httpMethods, string(httpMethod))
				}
				requestRateMap["http_methods"] = httpMethods
			}

			requestRateMap["path_pattern"] = condition.RequestRate.PathPattern
			requestRateMap["requests"] = condition.RequestRate.Requests
			requestRateMap["time"] = condition.RequestRate.Time

			if condition.RequestRate.UserDefinedTag != nil {
				requestRateMap["user_defined_tag"] = *condition.RequestRate.UserDefinedTag
			}

			conditionMap["request_rate"] = append(conditionMap["request_rate"].([]interface{}), requestRateMap)
		}

		if condition.OwnerTypes != nil && len(*condition.OwnerTypes.OwnerTypes) > 0 {
			ownerMap := map[string]interface{}{}
			var ownerTypes []string
			for _, ownerType := range *condition.OwnerTypes.OwnerTypes {
				ownerTypes = append(ownerTypes, string(ownerType))
			}
			ownerMap["owner_types"] = ownerTypes
			ownerMap["negation"] = condition.OwnerTypes.Negation
			conditionMap["owner_types"] = append(conditionMap["owner_types"].([]interface{}), ownerMap)
		}

		if condition.Tags != nil && len(condition.Tags.Tags) > 0 {
			tagMaps := map[string]interface{}{}
			tagMaps["tags"] = condition.Tags.Tags
			tagMaps["negation"] = condition.Tags.Negation
			conditionMap["tags"] = append(conditionMap["tags"].([]interface{}), tagMaps)
		}

		if condition.SessionRequestCount != nil {
			sessReqMap := map[string]interface{}{}
			sessReqMap["negation"] = condition.SessionRequestCount.Negation
			sessReqMap["request_count"] = condition.SessionRequestCount.RequestCount
			conditionMap["session_request_count"] = append(conditionMap["session_request_count"].([]interface{}), sessReqMap)
		}

		if condition.UserDefinedTags != nil && len(condition.UserDefinedTags.Tags) > 0 {
			tagMaps := map[string]interface{}{}
			tagMaps["tags"] = condition.UserDefinedTags.Tags
			tagMaps["negation"] = condition.UserDefinedTags.Negation
			conditionMap["user_defined_tags"] = append(conditionMap["user_defined_tags"].([]interface{}), tagMaps)
		}
	}

	return []interface{}{conditionMap}
}
