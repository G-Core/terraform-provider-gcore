package gcore

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	waap "github.com/G-Core/gcore-waap-sdk-go"
)

func resourceWaapCustomPageSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapCustomPageSetCreate,
		ReadContext:   resourceWaapCustomPageSetRead,
		UpdateContext: resourceWaapCustomPageSetUpdate,
		DeleteContext: resourceWaapCustomPageSetDelete,
		Description:   "Represent WAAP custom page set",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the custom page set.",
			},
			"domains": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of domain IDs associated with this custom page set.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"block": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logo": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A base64 encoded image of the logo to present",
						},
						"header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the header of the custom page",
						},
						"title": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the title of the custom page",
						},
						"text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the body of the custom page",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates whether the custom custom page is active or inactive",
						},
					},
				},
			},
			"block_csrf": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logo": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A base64 encoded image of the logo to present",
						},
						"header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the header of the custom page",
						},
						"title": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the title of the custom page",
						},
						"text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the body of the custom page",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates whether the custom custom page is active or inactive",
						},
					},
				},
			},
			"captcha": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logo": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A base64 encoded image of the logo to present",
						},
						"header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the header of the custom page",
						},
						"title": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the title of the custom page",
						},
						"text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the body of the custom page",
						},
						"error": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Error message",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates whether the custom custom page is active or inactive",
						},
					},
				},
			},
			"cookie_disabled": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the header of the custom page",
						},
						"text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the body of the custom page",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates whether the custom custom page is active or inactive",
						},
					},
				},
			},
			"handshake": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logo": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A base64 encoded image of the logo to present",
						},
						"header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the header of the custom page",
						},
						"title": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the title of the custom page",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates whether the custom custom page is active or inactive",
						},
					},
				},
			},
			"javascript_disabled": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the header of the custom page",
						},
						"text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to display in the body of the custom page",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates whether the custom custom page is active or inactive",
						},
					},
				},
			},
		},
	}
}

func resourceWaapCustomPageSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	createReq := waap.CreateCustomPageSetV1CustomPageSetsPostJSONRequestBody{
		Name: d.Get("name").(string),
	}

	if domains, ok := d.GetOk("domains"); ok {
		domainsList := expandIntList(domains.([]interface{}))
		createReq.Domains = &domainsList
	}

	if blockConfigs, ok := d.GetOk("block"); ok && len(blockConfigs.([]interface{})) > 0 {
		blockConfig := blockConfigs.([]interface{})[0].(map[string]interface{})
		block := waap.BlockPageData{}

		if v, ok := blockConfig["logo"].(string); ok && v != "" {
			block.Logo = &v
		}
		if v, ok := blockConfig["header"].(string); ok && v != "" {
			block.Header = &v
		}
		if v, ok := blockConfig["title"].(string); ok && v != "" {
			block.Title = &v
		}
		if v, ok := blockConfig["text"].(string); ok && v != "" {
			block.Text = &v
		}
		if v, ok := blockConfig["enabled"].(bool); ok {
			block.Enabled = v
		}

		createReq.Block = &block
	}

	if blockCsrfConfigs, ok := d.GetOk("block_csrf"); ok && len(blockCsrfConfigs.([]interface{})) > 0 {
		blockCsrfConfig := blockCsrfConfigs.([]interface{})[0].(map[string]interface{})
		blockCsrf := waap.BlockCsrfPageData{}

		if v, ok := blockCsrfConfig["logo"].(string); ok && v != "" {
			blockCsrf.Logo = &v
		}
		if v, ok := blockCsrfConfig["header"].(string); ok && v != "" {
			blockCsrf.Header = &v
		}
		if v, ok := blockCsrfConfig["title"].(string); ok && v != "" {
			blockCsrf.Title = &v
		}
		if v, ok := blockCsrfConfig["text"].(string); ok && v != "" {
			blockCsrf.Text = &v
		}
		if v, ok := blockCsrfConfig["enabled"].(bool); ok {
			blockCsrf.Enabled = v
		}

		createReq.BlockCsrf = &blockCsrf
	}

	if captchaConfigs, ok := d.GetOk("captcha"); ok && len(captchaConfigs.([]interface{})) > 0 {
		captchaConfig := captchaConfigs.([]interface{})[0].(map[string]interface{})
		captcha := waap.CaptchaPageData{}

		if v, ok := captchaConfig["logo"].(string); ok && v != "" {
			captcha.Logo = &v
		}
		if v, ok := captchaConfig["header"].(string); ok && v != "" {
			captcha.Header = &v
		}
		if v, ok := captchaConfig["title"].(string); ok && v != "" {
			captcha.Title = &v
		}
		if v, ok := captchaConfig["text"].(string); ok && v != "" {
			captcha.Text = &v
		}
		if v, ok := captchaConfig["error"].(string); ok && v != "" {
			captcha.Error = &v
		}
		if v, ok := captchaConfig["enabled"].(bool); ok {
			captcha.Enabled = v
		}

		createReq.Captcha = &captcha
	}

	if cookieDisabledConfigs, ok := d.GetOk("cookie_disabled"); ok && len(cookieDisabledConfigs.([]interface{})) > 0 {
		cookieDisabledConfig := cookieDisabledConfigs.([]interface{})[0].(map[string]interface{})
		cookieDisabled := waap.CookieDisabledPageData{}

		if v, ok := cookieDisabledConfig["header"].(string); ok && v != "" {
			cookieDisabled.Header = &v
		}
		if v, ok := cookieDisabledConfig["text"].(string); ok && v != "" {
			cookieDisabled.Text = &v
		}
		if v, ok := cookieDisabledConfig["enabled"].(bool); ok {
			cookieDisabled.Enabled = v
		}

		createReq.CookieDisabled = &cookieDisabled
	}

	if handshakeConfigs, ok := d.GetOk("handshake"); ok && len(handshakeConfigs.([]interface{})) > 0 {
		handshakeConfig := handshakeConfigs.([]interface{})[0].(map[string]interface{})
		handshake := waap.HandshakePageData{}

		if v, ok := handshakeConfig["logo"].(string); ok && v != "" {
			handshake.Logo = &v
		}
		if v, ok := handshakeConfig["header"].(string); ok && v != "" {
			handshake.Header = &v
		}
		if v, ok := handshakeConfig["title"].(string); ok && v != "" {
			handshake.Title = &v
		}
		if v, ok := handshakeConfig["enabled"].(bool); ok {
			handshake.Enabled = v
		}

		createReq.Handshake = &handshake
	}

	if jsDisabledConfigs, ok := d.GetOk("javascript_disabled"); ok && len(jsDisabledConfigs.([]interface{})) > 0 {
		jsDisabledConfig := jsDisabledConfigs.([]interface{})[0].(map[string]interface{})
		jsDisabled := waap.JavascriptDisabledPageData{}

		if v, ok := jsDisabledConfig["header"].(string); ok && v != "" {
			jsDisabled.Header = &v
		}
		if v, ok := jsDisabledConfig["text"].(string); ok && v != "" {
			jsDisabled.Text = &v
		}
		if v, ok := jsDisabledConfig["enabled"].(bool); ok {
			jsDisabled.Enabled = v
		}

		createReq.JavascriptDisabled = &jsDisabled
	}

	resp, err := client.CreateCustomPageSetV1CustomPageSetsPostWithResponse(ctx, createReq)
	if err != nil {
		return diag.Errorf("failed to create custom page set: %s", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return diag.Errorf("failed to create custom page set: status code %d with error: %s", resp.StatusCode(), resp.Body)
	}

	if resp.JSON201 == nil {
		return diag.Errorf("failed to create custom page set: empty response")
	}

	d.SetId(fmt.Sprintf("%d", resp.JSON201.Id))

	return resourceWaapCustomPageSetRead(ctx, d, m)
}

func resourceWaapCustomPageSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Config).WaapClient

	setID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to convert custom page set ID: %s", err)
	}

	resp, err := client.GetCustomPageSetV1CustomPageSetsSetIdGetWithResponse(ctx, setID)
	if err != nil {
		return diag.Errorf("failed to get custom page set: %s", err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return diag.Diagnostics{
			{Severity: diag.Warning, Summary: fmt.Sprintf("Custom Page Set (%s) was not found, removed from TF state", setID)},
		}
	}

	if resp.StatusCode() != http.StatusOK {
		return diag.Errorf("failed to read custom page set: status code %d with error: %s", resp.StatusCode(), resp.Body)
	}

	if resp.JSON200 == nil {
		return diag.Errorf("failed to read custom page set: empty response")
	}

	d.Set("name", resp.JSON200.Name)

	if resp.JSON200.Domains != nil {
		d.Set("domains", resp.JSON200.Domains)
	}

	if resp.JSON200.Block != nil {
		blockMap := map[string]interface{}{}

		if resp.JSON200.Block.Logo != nil {
			blockMap["logo"] = *resp.JSON200.Block.Logo
		}
		if resp.JSON200.Block.Header != nil {
			blockMap["header"] = *resp.JSON200.Block.Header
		}
		if resp.JSON200.Block.Title != nil {
			blockMap["title"] = *resp.JSON200.Block.Title
		}
		if resp.JSON200.Block.Text != nil {
			blockMap["text"] = *resp.JSON200.Block.Text
		}
		blockMap["enabled"] = resp.JSON200.Block.Enabled

		if len(blockMap) > 0 {
			d.Set("block", []interface{}{blockMap})
		}
	}

	if resp.JSON200.BlockCsrf != nil {
		blockCsrfMap := map[string]interface{}{}

		if resp.JSON200.BlockCsrf.Logo != nil {
			blockCsrfMap["logo"] = *resp.JSON200.BlockCsrf.Logo
		}
		if resp.JSON200.BlockCsrf.Header != nil {
			blockCsrfMap["header"] = *resp.JSON200.BlockCsrf.Header
		}
		if resp.JSON200.BlockCsrf.Title != nil {
			blockCsrfMap["title"] = *resp.JSON200.BlockCsrf.Title
		}
		if resp.JSON200.BlockCsrf.Text != nil {
			blockCsrfMap["text"] = *resp.JSON200.BlockCsrf.Text
		}
		blockCsrfMap["enabled"] = resp.JSON200.BlockCsrf.Enabled

		if len(blockCsrfMap) > 0 {
			d.Set("block_csrf", []interface{}{blockCsrfMap})
		}
	}

	if resp.JSON200.Captcha != nil {
		captchaMap := map[string]interface{}{}

		if resp.JSON200.Captcha.Logo != nil {
			captchaMap["logo"] = *resp.JSON200.Captcha.Logo
		}
		if resp.JSON200.Captcha.Header != nil {
			captchaMap["header"] = *resp.JSON200.Captcha.Header
		}
		if resp.JSON200.Captcha.Title != nil {
			captchaMap["title"] = *resp.JSON200.Captcha.Title
		}
		if resp.JSON200.Captcha.Text != nil {
			captchaMap["text"] = *resp.JSON200.Captcha.Text
		}
		if resp.JSON200.Captcha.Error != nil {
			captchaMap["error"] = *resp.JSON200.Captcha.Error
		}
		captchaMap["enabled"] = resp.JSON200.Captcha.Enabled

		if len(captchaMap) > 0 {
			d.Set("captcha", []interface{}{captchaMap})
		}
	}

	if resp.JSON200.CookieDisabled != nil {
		cookieDisabledMap := map[string]interface{}{}

		if resp.JSON200.CookieDisabled.Header != nil {
			cookieDisabledMap["header"] = *resp.JSON200.CookieDisabled.Header
		}
		if resp.JSON200.CookieDisabled.Text != nil {
			cookieDisabledMap["text"] = *resp.JSON200.CookieDisabled.Text
		}
		cookieDisabledMap["enabled"] = resp.JSON200.CookieDisabled.Enabled

		if len(cookieDisabledMap) > 0 {
			d.Set("cookie_disabled", []interface{}{cookieDisabledMap})
		}
	}

	if resp.JSON200.Handshake != nil {
		handshakeMap := map[string]interface{}{}

		if resp.JSON200.Handshake.Logo != nil {
			handshakeMap["logo"] = *resp.JSON200.Handshake.Logo
		}
		if resp.JSON200.Handshake.Header != nil {
			handshakeMap["header"] = *resp.JSON200.Handshake.Header
		}
		if resp.JSON200.Handshake.Title != nil {
			handshakeMap["title"] = *resp.JSON200.Handshake.Title
		}
		handshakeMap["enabled"] = resp.JSON200.Handshake.Enabled

		if len(handshakeMap) > 0 {
			d.Set("handshake", []interface{}{handshakeMap})
		}
	}

	if resp.JSON200.JavascriptDisabled != nil {
		jsDisabledMap := map[string]interface{}{}

		if resp.JSON200.JavascriptDisabled.Header != nil {
			jsDisabledMap["header"] = *resp.JSON200.JavascriptDisabled.Header
		}
		if resp.JSON200.JavascriptDisabled.Text != nil {
			jsDisabledMap["text"] = *resp.JSON200.JavascriptDisabled.Text
		}
		jsDisabledMap["enabled"] = resp.JSON200.JavascriptDisabled.Enabled

		if len(jsDisabledMap) > 0 {
			d.Set("javascript_disabled", []interface{}{jsDisabledMap})
		}
	}
	return nil
}

func resourceWaapCustomPageSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient

	setID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to convert custom page set ID: %s", err)
	}

	updateReq := waap.UpdateCustomPageSetV1CustomPageSetsSetIdPatchJSONRequestBody{}

	name := d.Get("name").(string)
	updateReq.Name = &name

	if domains, ok := d.GetOk("domains"); ok {
		domainsList := expandIntList(domains.([]interface{}))
		updateReq.Domains = &domainsList
	} else {
		emptyList := []int{}
		updateReq.Domains = &emptyList
	}

	if blockConfigs, ok := d.GetOk("block"); ok && len(blockConfigs.([]interface{})) > 0 {
		blockConfig := blockConfigs.([]interface{})[0].(map[string]interface{})
		block := waap.BlockPageData{}

		if v, ok := blockConfig["logo"].(string); ok && v != "" {
			block.Logo = &v
		}
		if v, ok := blockConfig["header"].(string); ok && v != "" {
			block.Header = &v
		}
		if v, ok := blockConfig["title"].(string); ok && v != "" {
			block.Title = &v
		}
		if v, ok := blockConfig["text"].(string); ok && v != "" {
			block.Text = &v
		}
		if v, ok := blockConfig["enabled"].(bool); ok {
			block.Enabled = v
		}

		updateReq.Block = &block
	}

	if blockCsrfConfigs, ok := d.GetOk("block_csrf"); ok && len(blockCsrfConfigs.([]interface{})) > 0 {
		blockCsrfConfig := blockCsrfConfigs.([]interface{})[0].(map[string]interface{})
		blockCsrf := waap.BlockCsrfPageData{}

		if v, ok := blockCsrfConfig["logo"].(string); ok && v != "" {
			blockCsrf.Logo = &v
		}
		if v, ok := blockCsrfConfig["header"].(string); ok && v != "" {
			blockCsrf.Header = &v
		}
		if v, ok := blockCsrfConfig["title"].(string); ok && v != "" {
			blockCsrf.Title = &v
		}
		if v, ok := blockCsrfConfig["text"].(string); ok && v != "" {
			blockCsrf.Text = &v
		}
		if v, ok := blockCsrfConfig["enabled"].(bool); ok {
			blockCsrf.Enabled = v
		}

		updateReq.BlockCsrf = &blockCsrf
	}

	if captchaConfigs, ok := d.GetOk("captcha"); ok && len(captchaConfigs.([]interface{})) > 0 {
		captchaConfig := captchaConfigs.([]interface{})[0].(map[string]interface{})
		captcha := waap.CaptchaPageData{}

		if v, ok := captchaConfig["logo"].(string); ok && v != "" {
			captcha.Logo = &v
		}
		if v, ok := captchaConfig["header"].(string); ok && v != "" {
			captcha.Header = &v
		}
		if v, ok := captchaConfig["title"].(string); ok && v != "" {
			captcha.Title = &v
		}
		if v, ok := captchaConfig["text"].(string); ok && v != "" {
			captcha.Text = &v
		}
		if v, ok := captchaConfig["error"].(string); ok && v != "" {
			captcha.Error = &v
		}
		if v, ok := captchaConfig["enabled"].(bool); ok {
			captcha.Enabled = v
		}

		updateReq.Captcha = &captcha
	}

	if cookieDisabledConfigs, ok := d.GetOk("cookie_disabled"); ok && len(cookieDisabledConfigs.([]interface{})) > 0 {
		cookieDisabledConfig := cookieDisabledConfigs.([]interface{})[0].(map[string]interface{})
		cookieDisabled := waap.CookieDisabledPageData{}

		if v, ok := cookieDisabledConfig["header"].(string); ok && v != "" {
			cookieDisabled.Header = &v
		}
		if v, ok := cookieDisabledConfig["text"].(string); ok && v != "" {
			cookieDisabled.Text = &v
		}
		if v, ok := cookieDisabledConfig["enabled"].(bool); ok {
			cookieDisabled.Enabled = v
		}

		updateReq.CookieDisabled = &cookieDisabled
	}

	if handshakeConfigs, ok := d.GetOk("handshake"); ok && len(handshakeConfigs.([]interface{})) > 0 {
		handshakeConfig := handshakeConfigs.([]interface{})[0].(map[string]interface{})
		handshake := waap.HandshakePageData{}

		if v, ok := handshakeConfig["logo"].(string); ok && v != "" {
			handshake.Logo = &v
		}
		if v, ok := handshakeConfig["header"].(string); ok && v != "" {
			handshake.Header = &v
		}
		if v, ok := handshakeConfig["title"].(string); ok && v != "" {
			handshake.Title = &v
		}
		if v, ok := handshakeConfig["enabled"].(bool); ok {
			handshake.Enabled = v
		}

		updateReq.Handshake = &handshake
	}

	if jsDisabledConfigs, ok := d.GetOk("javascript_disabled"); ok && len(jsDisabledConfigs.([]interface{})) > 0 {
		jsDisabledConfig := jsDisabledConfigs.([]interface{})[0].(map[string]interface{})
		jsDisabled := waap.JavascriptDisabledPageData{}

		if v, ok := jsDisabledConfig["header"].(string); ok && v != "" {
			jsDisabled.Header = &v
		}
		if v, ok := jsDisabledConfig["text"].(string); ok && v != "" {
			jsDisabled.Text = &v
		}
		if v, ok := jsDisabledConfig["enabled"].(bool); ok {
			jsDisabled.Enabled = v
		}

		updateReq.JavascriptDisabled = &jsDisabled
	}

	resp, err := client.UpdateCustomPageSetV1CustomPageSetsSetIdPatchWithResponse(ctx, setID, updateReq)
	if err != nil {
		return diag.Errorf("failed to update custom page set: %s", err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		return diag.Errorf("failed to update custom page set: status code %d with error: %s", resp.StatusCode(), resp.Body)
	}

	return resourceWaapCustomPageSetRead(ctx, d, m)
}

func resourceWaapCustomPageSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient
	setID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to convert custom page set ID: %s", err)
	}

	resp, err := client.DeleteCustomPageSetV1CustomPageSetsSetIdDeleteWithResponse(ctx, setID)
	if err != nil {
		return diag.Errorf("failed to delete custom page set: %s", err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		return diag.Errorf("failed to delete custom page set: status code %d with error: %s", resp.StatusCode(), resp.Body)
	}

	d.SetId("")
	return nil
}

// Helper function to convert []interface{} to []int
func expandIntList(list []interface{}) []int {
	result := make([]int, len(list))
	for i, v := range list {
		result[i] = v.(int)
	}
	return result
}
