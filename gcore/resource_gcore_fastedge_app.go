package gcore

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

func resourceFastEdgeApp() *schema.Resource {
	return &schema.Resource{
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
			},
			"template": {
				Description: "Application template id.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"status": {
				Description: "Status code. Possible values are: enabled, disabled, suspended.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: func(v interface{}, k string) ([]string, []error) {
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
		Description:   "FastEdge application.",
	}
}

func resourceFastEdgeAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge app creation")
	config := m.(*Config)
	client := config.FastEdgeClient

	status := statusToInt(d.Get("status").(string))
	app := sdk.App{
		Name:       fieldValue[string](d, "name"),
		Binary:     fieldValueInt64(d, "binary"),
		Template:   fieldValueInt64(d, "template"),
		Debug:      fieldValue[bool](d, "debug"),
		Env:        fieldValueMap(d, "env"),
		RspHeaders: fieldValueMap(d, "rsp_headers"),
		Comment:    fieldValue[string](d, "comment"),
		Status:     &status,
	}
	rsp, err := client.AddAppWithResponse(ctx, app)
	if err != nil {
		return diag.FromErr(err)
	}
	if !statusOK(rsp.StatusCode()) {
		return diag.FromErr(errors.New(extractErrorMessage(rsp.Body)))
	}
	d.SetId(strconv.FormatInt(rsp.JSON200.Id, 10))
	d.Set("name", rsp.JSON200.Name)
	d.Set("binary", rsp.JSON200.Binary)

	log.Printf("[DEBUG] Finish FastEdge app creation (id=%d)\n", rsp.JSON200.Id)
	return nil
}

func resourceFastEdgeAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge app read")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	rsp, err := client.GetAppWithResponse(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() == http.StatusNotFound {
			log.Printf("[WARN] FastEdge app (%d) was not found, removing from TF state", id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(errors.New(extractErrorMessage(rsp.Body)))
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

	log.Println("[DEBUG] Finish FastEdge app read")
	return nil
}

func resourceFastEdgeAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge app update")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var statusCode *int
	status := fieldValueIfChanged[string](d, "status")
	if status != nil {
		statusCode = new(int)
		*statusCode = statusToInt(*status)
	}
	app := sdk.App{
		Name:       fieldValueIfChanged[string](d, "name"),
		Binary:     fieldValueIfChangedInt64(d, "binary"),
		Template:   fieldValueIfChangedInt64(d, "template"),
		Debug:      fieldValueIfChanged[bool](d, "debug"),
		Env:        fieldValueIfChangedMap(d, "env"),
		RspHeaders: fieldValueIfChangedMap(d, "rsp_headers"),
		Comment:    fieldValueIfChanged[string](d, "comment"),
		Status:     statusCode,
	}
	rsp, err := client.PatchAppWithResponse(ctx, id, app)
	if err != nil {
		return diag.FromErr(err)
	}
	if !statusOK(rsp.StatusCode()) {
		return diag.FromErr(errors.New(extractErrorMessage(rsp.Body)))
	}
	d.Set("name", rsp.JSON200.Name) // name may change as a result of the update

	log.Println("[DEBUG] Finish FastEdge app update")
	return nil
}

func resourceFastEdgeAppDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge app deletion")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	rsp, err := client.DelAppWithResponse(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() != http.StatusConflict {
			return diag.FromErr(errors.New(extractErrorMessage(rsp.Body)))
		}
		return diag.FromErr(errors.New("FastEdge app is referenced, cannot delete"))
	}

	d.SetId("")
	log.Println("[DEBUG] Finish FastEdge app deletion")
	return nil
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

func fieldValueMap(d *schema.ResourceData, name string) *map[string]string {
	v := d.Get(name)
	if v == nil {
		return nil
	}
	tmpVal, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	val := convertMap(tmpVal)
	return &val
}

func convertMap(in map[string]interface{}) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v.(string)
	}
	return out
}

func fieldValueIfChanged[T any](d *schema.ResourceData, name string) *T {
	oldVal, newVal := d.GetChange(name)
	if cmp.Equal(newVal, oldVal) {
		return nil
	}
	val, ok := newVal.(T)
	if !ok {
		return nil
	}
	return &val
}

func fieldValueIfChangedInt64(d *schema.ResourceData, name string) *int64 {
	oldVal, newVal := d.GetChange(name)
	if cmp.Equal(newVal, oldVal) {
		return nil
	}
	intVal, ok := newVal.(int)
	if !ok {
		return nil
	}
	val := int64(intVal)
	return &val
}

func fieldValueIfChangedMap(d *schema.ResourceData, name string) *map[string]string {
	oldVal, newVal := d.GetChange(name)
	if cmp.Equal(newVal, oldVal) {
		return nil
	}
	tmpVal, ok := newVal.(map[string]interface{})
	if !ok {
		return nil
	}
	val := convertMap(tmpVal)
	return &val
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
