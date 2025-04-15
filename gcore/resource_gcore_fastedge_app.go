package gcore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

func resourceFastEdgeApp() *schema.Resource {
	return &schema.Resource{
		Description: "FastEdge application.",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Application name.",
				Type:        schema.TypeString,
				Computed:    true, // if app name is not provided, it will be generated
				Optional:    true,
			},
			"binary": {
				Description: "WebAssembly binary id.",
				Type:        schema.TypeInt,
				Computed:    true, // if template is specified, binary id is returned
				Optional:    true,
				ExactlyOneOf: []string{
					"binary",
					"template",
				},
			},
			"template": {
				Description: "Application template id.",
				Type:        schema.TypeInt,
				Optional:    true,
				ExactlyOneOf: []string{
					"binary",
					"template",
				},
			},
			"status": {
				Description: "Status code. Possible values are: enabled, disabled, suspended.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: func(v any, k string) ([]string, []error) {
					status := strings.ToLower(v.(string))
					if status != "enabled" && status != "disabled" {
						return nil, []error{errors.New("only 'enabled' or 'disabled' status can be set by the user")}
					}
					return nil, nil
				},
			},
			"env": {
				Description: "Environment variables.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        schema.TypeString,
			},
			"secrets": {
				Description: "Secret variables.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        schema.TypeInt,
			},
			"rsp_headers": {
				Description: "Response headers.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        schema.TypeString,
			},
			"debug": {
				Description: "Logging enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"comment": {
				Description: "Application comment.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
		CreateContext: resourceFastEdgeAppCreate,
		ReadContext:   resourceFastEdgeAppRead,
		UpdateContext: resourceFastEdgeAppUpdate,
		DeleteContext: resourceFastEdgeAppDelete,
	}
}

func resourceFastEdgeAppCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge app creation")
	config := m.(*Config)
	client := config.FastEdgeClient

	status := statusToInt(d.Get("status").(string))
	app := sdk.App{
		Name:       fieldValue[string](d, "name"),
		Binary:     fieldValueInt64(d, "binary"),
		Template:   fieldValueInt64(d, "template"),
		Debug:      fieldValue[bool](d, "debug"),
		Env:        fieldValueStringMap(d, "env"),
		RspHeaders: fieldValueStringMap(d, "rsp_headers"),
		Secrets:    fieldValueSecretMap(d, "secrets"),
		Comment:    fieldValue[string](d, "comment"),
		Status:     &status,
	}
	rsp, err := client.AddAppWithResponse(ctx, app)
	if err != nil {
		return diag.Errorf("calling AddApp API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		return diag.Errorf("calling AddApp API: %s", extractErrorMessage(rsp.Body))
	}

	d.SetId(strconv.FormatInt(rsp.JSON200.Id, 10))
	d.Set("name", rsp.JSON200.Name)
	d.Set("binary", rsp.JSON200.Binary)
	d.Set("status", statusToString(rsp.JSON200.Status)) // return status may differ from the one set by the user

	log.Printf("[DEBUG] Finish FastEdge app creation (id=%d)\n", rsp.JSON200.Id)
	return nil
}

func resourceFastEdgeAppRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge app read")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("converting id to number: %v", err)
	}

	rsp, err := client.GetAppWithResponse(ctx, id)
	if err != nil {
		return diag.Errorf("calling GetApp API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("[FastEdge app (%d) was not found, removed from TF state", id),
				},
			}
		}
		return diag.Errorf("calling GetApp API: %s", extractErrorMessage(rsp.Body))
	}

	app := rsp.JSON200
	setField(d, "name", app.Name)
	setField(d, "binary", app.Binary)
	setField(d, "template", app.Template)
	setField(d, "debug", app.Debug)
	setField(d, "env", app.Env)
	setField(d, "rsp_headers", app.RspHeaders)
	setField(d, "comment", app.Comment)
	d.Set("status", statusToString(*app.Status))
	if app.Secrets != nil {
		secrets := make(map[string]any, len(*app.Secrets))
		for k, v := range *app.Secrets {
			secrets[k] = v.Id
		}
		d.Set("secrets", secrets)
	}

	log.Println("[DEBUG] Finish FastEdge app read")
	return nil
}

func resourceFastEdgeAppUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge app update")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("converting id to number: %v", err)
	}

	status := statusToInt(d.Get("status").(string))
	app := sdk.App{
		Name:       fieldValue[string](d, "name"),
		Binary:     fieldValueInt64(d, "binary"),
		Template:   fieldValueInt64(d, "template"),
		Debug:      fieldValue[bool](d, "debug"),
		Env:        fieldValueStringMap(d, "env"),
		RspHeaders: fieldValueStringMap(d, "rsp_headers"),
		Secrets:    fieldValueSecretMap(d, "secrets"),
		Comment:    fieldValue[string](d, "comment"),
		Status:     &status,
	}
	rsp, err := client.UpdateAppWithResponse(ctx, id, sdk.UpdateAppJSONRequestBody{App: app})
	if err != nil {
		return diag.Errorf("calling UpdateApp API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("[FastEdge app (%d) was not found, removed from TF state", id),
				},
			}
		}
		return diag.Errorf("calling UpdateApp API: %s", extractErrorMessage(rsp.Body))
	}

	d.Set("name", rsp.JSON200.Name)                     // name may change as a result of the update
	d.Set("status", statusToString(rsp.JSON200.Status)) // return status may differ from the one set by the user

	log.Println("[DEBUG] Finish FastEdge app update")
	return nil
}

func resourceFastEdgeAppDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Println("[DEBUG] Start FastEdge app deletion")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("converting id to number: %v", err)
	}

	rsp, err := client.DelAppWithResponse(ctx, id)
	if err != nil {
		return diag.Errorf("calling DelApp API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() == http.StatusConflict {
			diags = diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("App (%d) is referenced so cannot be deleted, but removed from TF state", id),
				},
			}
		} else {
			return diag.Errorf("calling DelApp API: %s", extractErrorMessage(rsp.Body))
		}
	}

	d.SetId("")
	log.Println("[DEBUG] Finish FastEdge app deletion")
	return diags
}

func fieldValue[T any](d *schema.ResourceData, name string) *T {
	v := d.Get(name)
	if v == nil {
		return nil
	}
	val, ok := v.(T)
	if !ok {
		return nil
	}
	return &val
}

func fieldValueInt64(d *schema.ResourceData, name string) *int64 {
	v := d.Get(name)
	if v == nil {
		return nil
	}
	intVal, ok := v.(int)
	if !ok || intVal == 0 {
		return nil
	}
	val := int64(intVal)
	return &val
}

func fieldValueStringMap(d *schema.ResourceData, name string) *map[string]string {
	v := d.Get(name)
	if v == nil {
		return nil
	}
	tmpVal, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	val := convertStringMap(tmpVal)
	return &val
}

func fieldValueSecretMap(d *schema.ResourceData, name string) *map[string]sdk.AppSecretShort {
	v := d.Get(name)
	if v == nil {
		return nil
	}
	tmpVal, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	val := make(map[string]sdk.AppSecretShort, len(tmpVal))
	for k, v := range tmpVal {
		val[k] = sdk.AppSecretShort{
			Id: int64(v.(int)),
		}
	}
	return &val
}

func convertStringMap(in map[string]any) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v.(string)
	}
	return out
}

func setField[T any](d *schema.ResourceData, name string, value *T) {
	if value != nil {
		d.Set(name, *value)
	}
}

func statusToInt(status string) int {
	switch status {
	case "enabled":
		return 1
	case "disabled":
		return 2
	default:
		return -1 // will fail on server side
	}
}

func statusToString(status int) string {
	switch status {
	case 1:
		return "enabled"
	case 2:
		return "disabled"
	case 5:
		return "suspended"
	default:
		return ""
	}
}
