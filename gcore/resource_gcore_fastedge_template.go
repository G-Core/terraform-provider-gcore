package gcore

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

var paramTypes = []string{"string", "number", "date", "time", "secret"}

func resourceFastEdgeTemplate() *schema.Resource {
	return &schema.Resource{
		Description: "FastEdge application template.",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Template name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"binary": {
				Description: "WebAssembly binary id.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"short_descr": {
				Description: "Short description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"long_descr": {
				Description: "Instruction how to configure the template.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"param": {
				Description: "Template parameter.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Parameter name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"type": {
							Description:  "Parameter type. Possible values are: " + strings.Join(paramTypes, ", ") + ".",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(paramTypes, false),
						},
						"mandatory": {
							Description: "Is parameter mandatory, true/false (false by default).",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"descr": {
							Description: "Parameter description.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
		CreateContext: resourceFastEdgeTemplateCreate,
		ReadContext:   resourceFastEdgeTemplateRead,
		UpdateContext: resourceFastEdgeTemplateUpdate,
		DeleteContext: resourceFastEdgeTemplateDelete,
	}
}

func resourceFastEdgeTemplateCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge template creation")
	config := m.(*Config)
	client := config.FastEdgeClient

	template := sdk.Template{
		Name:       *fieldValue[string](d, "name"),
		BinaryId:   *fieldValueInt64(d, "binary"),
		ShortDescr: fieldValue[string](d, "short_descr"),
		LongDescr:  fieldValue[string](d, "long_descr"),
		Params:     procParams(d),
	}

	rsp, err := client.AddTemplateWithResponse(ctx, template)
	if err != nil {
		return diag.Errorf("calling AddTemplate API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		return diag.Errorf("calling AddTemplate API: %s", extractErrorMessage(rsp.Body))
	}

	d.SetId(strconv.FormatInt(rsp.JSON200.Id, 10))

	log.Printf("[DEBUG] Finish FastEdge template creation (id=%d)\n", rsp.JSON200.Id)
	return nil
}

func resourceFastEdgeTemplateRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge template read")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("converting id to number: %v", err)
	}

	rsp, err := client.GetTemplateWithResponse(ctx, id)
	if err != nil {
		return diag.Errorf("calling GetTemplate API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("[FastEdge template (%d) was not found, removed from TF state", id),
				},
			}
		}
		return diag.Errorf("calling GetTemplate API: %s", extractErrorMessage(rsp.Body))
	}

	template := rsp.JSON200
	setField(d, "name", &template.Name)
	setField(d, "binary", &template.BinaryId)
	setField(d, "short_descr", template.ShortDescr)
	setField(d, "long_descr", template.LongDescr)
	if template.Params != nil {
		params := make([]any, len(template.Params))
		for i, param := range template.Params {
			params[i] = map[string]any{
				"name":      param.Name,
				"type":      param.DataType,
				"mandatory": param.Mandatory,
				"descr":     param.Descr,
			}
		}
		d.Set("param", params)
	} else {
		d.Set("param", nil)
	}

	log.Println("[DEBUG] Finish FastEdge template read")
	return nil
}

func resourceFastEdgeTemplateUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge template update")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("converting id to number: %v", err)
	}

	template := sdk.Template{
		Name:       *fieldValue[string](d, "name"),
		BinaryId:   *fieldValueInt64(d, "binary"),
		ShortDescr: fieldValue[string](d, "short_descr"),
		LongDescr:  fieldValue[string](d, "long_descr"),
		Params:     procParams(d),
	}

	rsp, err := client.UpdateTemplateWithResponse(ctx, id, template)
	if err != nil {
		return diag.Errorf("calling UpdateTemplate API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("[FastEdge template (%d) was not found, removed from TF state", id),
				},
			}
		}
		return diag.Errorf("calling UpdateTemplate API: %s", extractErrorMessage(rsp.Body))
	}

	log.Printf("[DEBUG] Finish FastEdge template update (id=%d)\n", rsp.JSON200.Id)
	return nil
}

func resourceFastEdgeTemplateDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Println("[DEBUG] Start FastEdge template deletion")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("converting id to number: %v", err)
	}

	rsp, err := client.DelTemplateWithResponse(ctx, id)
	if err != nil {
		return diag.Errorf("calling DelTemplate API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		return diag.Errorf("calling DelTemplate API: %s", extractErrorMessage(rsp.Body))
	}

	d.SetId("")
	log.Println("[DEBUG] Finish FastEdge template deletion")
	return diags
}

func procParams(d *schema.ResourceData) []sdk.TemplateParam {
	if v, ok := d.Get("param").(*schema.Set); ok {
		list := v.List()
		res := make([]sdk.TemplateParam, len(list))
		for i, v := range list {
			p := v.(map[string]any)
			var descr *string
			if tmp, ok := p["descr"].(string); ok {
				descr = &tmp
			}
			res[i] = sdk.TemplateParam{
				Name:      p["name"].(string),
				DataType:  sdk.TemplateParamDataType(p["type"].(string)),
				Mandatory: p["mandatory"].(bool),
				Descr:     descr,
			}
		}
		return res
	}
	return []sdk.TemplateParam{}
}
