package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/G-Core/gcorelabscdn-go/ruletemplates"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRuleTemplate() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Rule template name.",
			},
			"rule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to the file or folder for which the rule will be applied. The rule is applied if the requested URI matches the rule path. We add a leading forward slash to any rule path. Specify a path without a forward slash.",
			},
			"rule_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
				Description:  "Rule type. Possible values are: 0 - Regular expression. Must start with '^/' or '/'. 1 - Regular expression. Note that for this rule type we automatically add / to each rule pattern before your regular expression. This type is legacy, please use 0.",
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Rule execution order: from lowest (1) to highest. If requested URI matches multiple rules, the one higher in the order of the rules will be applied.",
			},
			"override_origin_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "MATCH"}, false),
				Description:  "Sets a protocol other than the one specified in the CDN resource settings to connect to the origin. If not specified, it will be inherited from the CDN resource settings. Possible values are: HTTPS, HTTP, MATCH.",
			},
			"options": ruleOptionsSchema,
		},
		CreateContext: resourceRuleTemplateCreate,
		ReadContext:   resourceRuleTemplateRead,
		UpdateContext: resourceRuleTemplateUpdate,
		DeleteContext: resourceRuleTemplateDelete,
		CustomizeDiff: validateCDNOptions,
		Description:   "Represent CDN rule template",
	}
}

func resourceRuleTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN Rule Template creating")
	config := m.(*Config)
	client := config.CDNClient

	var req ruletemplates.CreateRequest
	req.Name = d.Get("name").(string)
	req.Rule = d.Get("rule").(string)
	req.RuleType = ruletemplates.RuleType(d.Get("rule_type").(int))
	req.Weight = d.Get("weight").(int)

	if v, ok := d.GetOk("override_origin_protocol"); ok {
		overrideOriginProtocol := ruletemplates.Protocol(v.(string))
		req.OverrideOriginProtocol = &overrideOriginProtocol
	} else {
		req.OverrideOriginProtocol = nil
	}

	req.Options = listToOptions(d.Get("options").([]interface{}))

	result, err := client.RuleTemplates().Create(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	resourceRuleTemplateRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN Rule Template creating (id=%d)\n", result.ID)
	return nil
}

func resourceRuleTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ruleTemplateID := d.Id()
	log.Printf("[DEBUG] Start CDN Rule Template reading (id=%s)\n", ruleTemplateID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(ruleTemplateID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.RuleTemplates().Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", result.Name)
	d.Set("rule", result.Rule)
	d.Set("ruleType", result.RuleType)
	d.Set("weight", result.Weight)
	d.Set("overrideOriginProtocol", result.OverrideOriginProtocol)
	if err := d.Set("options", optionsToList(result.Options)); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Rule Template reading")
	return nil
}

func resourceRuleTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ruleTemplateID := d.Id()
	log.Printf("[DEBUG] Start CDN Rule Template updating (id=%s)\n", ruleTemplateID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(ruleTemplateID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req ruletemplates.UpdateRequest
	req.Name = d.Get("name").(string)
	req.Rule = d.Get("rule").(string)
	req.RuleType = ruletemplates.RuleType(d.Get("rule_type").(int))
	req.Weight = d.Get("weight").(int)

	if v, ok := d.GetOk("override_origin_protocol"); ok {
		overrideOriginProtocol := ruletemplates.Protocol(v.(string))
		req.OverrideOriginProtocol = &overrideOriginProtocol
	} else {
		req.OverrideOriginProtocol = nil
	}

	req.Options = listToOptions(d.Get("options").([]interface{}))

	if _, err := client.RuleTemplates().Update(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Rule Template updating")
	return resourceRuleTemplateRead(ctx, d, m)
}

func resourceRuleTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ruleTemplateID := d.Id()
	log.Printf("[DEBUG] Start CDN Rule Template deleting (id=%s)\n", ruleTemplateID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(ruleTemplateID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.RuleTemplates().Delete(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish CDN Rule Template deleting")
	return nil
}
