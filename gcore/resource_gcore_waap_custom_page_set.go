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

	name, domains, block, blockCsrf, captcha, cookieDisabled, handshake, javascriptDisabled :=
		prepareWaapCustomPageSetPayload(d)

	createReq := waap.CreateCustomPageSetV1CustomPageSetsPostJSONRequestBody{
		Name:               name,
		Domains:            domains,
		Block:              block,
		BlockCsrf:          blockCsrf,
		Captcha:            captcha,
		CookieDisabled:     cookieDisabled,
		Handshake:          handshake,
		JavascriptDisabled: javascriptDisabled,
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

	domains := []interface{}{}
	if resp.JSON200.Domains != nil {
		d.Set("domains", resp.JSON200.Domains)
	} else {
		d.Set("domains", domains)
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

		d.Set("block", []interface{}{blockMap})
	} else {
		d.Set("block", []interface{}{})
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

		d.Set("block_csrf", []interface{}{blockCsrfMap})
	} else {
		d.Set("block_csrf", []interface{}{})
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

		d.Set("captcha", []interface{}{captchaMap})
	} else {
		d.Set("captcha", []interface{}{})
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

		d.Set("cookie_disabled", []interface{}{cookieDisabledMap})
	} else {
		d.Set("cookie_disabled", []interface{}{})
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

		d.Set("handshake", []interface{}{handshakeMap})
	} else {
		d.Set("handshake", []interface{}{})
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

		d.Set("javascript_disabled", []interface{}{jsDisabledMap})
	} else {
		d.Set("javascript_disabled", []interface{}{})
	}

	return nil
}
func resourceWaapCustomPageSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).WaapClient

	setID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to convert custom page set ID: %s", err)
	}

	name, domains, block, blockCsrf, captcha, cookieDisabled, handshake, javascriptDisabled :=
		prepareWaapCustomPageSetPayload(d)

	updateReq := waap.UpdateCustomPageSetV1CustomPageSetsSetIdPatchJSONRequestBody{
		Name:               &name,
		Domains:            domains,
		Block:              block,
		BlockCsrf:          blockCsrf,
		Captcha:            captcha,
		CookieDisabled:     cookieDisabled,
		Handshake:          handshake,
		JavascriptDisabled: javascriptDisabled,
	}

	if domains == nil {
		emptyList := []int{}
		updateReq.Domains = &emptyList
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

func prepareWaapCustomPageSetPayload(d *schema.ResourceData) (
	string, // name
	*[]int, // domains
	*waap.BlockPageData,
	*waap.BlockCsrfPageData,
	*waap.CaptchaPageData,
	*waap.CookieDisabledPageData,
	*waap.HandshakePageData,
	*waap.JavascriptDisabledPageData,
) {

	name := d.Get("name").(string)

	var domains *[]int
	if domainsRaw, ok := d.GetOk("domains"); ok {
		domainsList := expandIntList(domainsRaw.([]interface{}))
		domains = &domainsList
	}

	var block *waap.BlockPageData
	if blockConfigs, ok := d.GetOk("block"); ok && len(blockConfigs.([]interface{})) > 0 {
		blockConfig := blockConfigs.([]interface{})[0].(map[string]interface{})
		blockData := waap.BlockPageData{}

		if v, ok := blockConfig["logo"].(string); ok && v != "" {
			blockData.Logo = &v
		}
		if v, ok := blockConfig["header"].(string); ok && v != "" {
			blockData.Header = &v
		}
		if v, ok := blockConfig["title"].(string); ok && v != "" {
			blockData.Title = &v
		}
		if v, ok := blockConfig["text"].(string); ok && v != "" {
			blockData.Text = &v
		}
		if v, ok := blockConfig["enabled"].(bool); ok {
			blockData.Enabled = v
		}

		block = &blockData
	}

	var blockCsrf *waap.BlockCsrfPageData
	if blockCsrfConfigs, ok := d.GetOk("block_csrf"); ok && len(blockCsrfConfigs.([]interface{})) > 0 {
		blockCsrfConfig := blockCsrfConfigs.([]interface{})[0].(map[string]interface{})
		blockCsrfData := waap.BlockCsrfPageData{}

		if v, ok := blockCsrfConfig["logo"].(string); ok && v != "" {
			blockCsrfData.Logo = &v
		}
		if v, ok := blockCsrfConfig["header"].(string); ok && v != "" {
			blockCsrfData.Header = &v
		}
		if v, ok := blockCsrfConfig["title"].(string); ok && v != "" {
			blockCsrfData.Title = &v
		}
		if v, ok := blockCsrfConfig["text"].(string); ok && v != "" {
			blockCsrfData.Text = &v
		}
		if v, ok := blockCsrfConfig["enabled"].(bool); ok {
			blockCsrfData.Enabled = v
		}

		blockCsrf = &blockCsrfData
	}

	var captcha *waap.CaptchaPageData
	if captchaConfigs, ok := d.GetOk("captcha"); ok && len(captchaConfigs.([]interface{})) > 0 {
		captchaConfig := captchaConfigs.([]interface{})[0].(map[string]interface{})
		captchaData := waap.CaptchaPageData{}

		if v, ok := captchaConfig["logo"].(string); ok && v != "" {
			captchaData.Logo = &v
		}
		if v, ok := captchaConfig["header"].(string); ok && v != "" {
			captchaData.Header = &v
		}
		if v, ok := captchaConfig["title"].(string); ok && v != "" {
			captchaData.Title = &v
		}
		if v, ok := captchaConfig["text"].(string); ok && v != "" {
			captchaData.Text = &v
		}
		if v, ok := captchaConfig["error"].(string); ok && v != "" {
			captchaData.Error = &v
		}
		if v, ok := captchaConfig["enabled"].(bool); ok {
			captchaData.Enabled = v
		}

		captcha = &captchaData
	}

	var cookieDisabled *waap.CookieDisabledPageData
	if cookieDisabledConfigs, ok := d.GetOk("cookie_disabled"); ok && len(cookieDisabledConfigs.([]interface{})) > 0 {
		cookieDisabledConfig := cookieDisabledConfigs.([]interface{})[0].(map[string]interface{})
		cookieDisabledData := waap.CookieDisabledPageData{}

		if v, ok := cookieDisabledConfig["header"].(string); ok && v != "" {
			cookieDisabledData.Header = &v
		}
		if v, ok := cookieDisabledConfig["text"].(string); ok && v != "" {
			cookieDisabledData.Text = &v
		}
		if v, ok := cookieDisabledConfig["enabled"].(bool); ok {
			cookieDisabledData.Enabled = v
		}

		cookieDisabled = &cookieDisabledData
	}

	var handshake *waap.HandshakePageData
	if handshakeConfigs, ok := d.GetOk("handshake"); ok && len(handshakeConfigs.([]interface{})) > 0 {
		handshakeConfig := handshakeConfigs.([]interface{})[0].(map[string]interface{})
		handshakeData := waap.HandshakePageData{}

		if v, ok := handshakeConfig["logo"].(string); ok && v != "" {
			handshakeData.Logo = &v
		}
		if v, ok := handshakeConfig["header"].(string); ok && v != "" {
			handshakeData.Header = &v
		}
		if v, ok := handshakeConfig["title"].(string); ok && v != "" {
			handshakeData.Title = &v
		}
		if v, ok := handshakeConfig["enabled"].(bool); ok {
			handshakeData.Enabled = v
		}

		handshake = &handshakeData
	}

	var javascriptDisabled *waap.JavascriptDisabledPageData
	if jsDisabledConfigs, ok := d.GetOk("javascript_disabled"); ok && len(jsDisabledConfigs.([]interface{})) > 0 {
		jsDisabledConfig := jsDisabledConfigs.([]interface{})[0].(map[string]interface{})
		jsDisabledData := waap.JavascriptDisabledPageData{}

		if v, ok := jsDisabledConfig["header"].(string); ok && v != "" {
			jsDisabledData.Header = &v
		}
		if v, ok := jsDisabledConfig["text"].(string); ok && v != "" {
			jsDisabledData.Text = &v
		}
		if v, ok := jsDisabledConfig["enabled"].(bool); ok {
			jsDisabledData.Enabled = v
		}

		javascriptDisabled = &jsDisabledData
	}

	return name, domains, block, blockCsrf, captcha, cookieDisabled, handshake, javascriptDisabled
}
