package gcore

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

var AllowDeletionOfSecretsWithSlots = true

func resourceFastEdgeSecret() *schema.Resource {
	return &schema.Resource{
		Description: "Represents FastEdge secret",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secret name.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Secret description.",
			},
			"slot": {
				Description: "Secret slot.",
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Set: func(val any) int {
					// for existing resources old secret value (from state) is empty while new
					// is not therefore default hash differs even for unchanged value, so
					// exclude 'value' from hash calcuilation
					var buf bytes.Buffer
					slot := val.(map[string]any)
					buf.WriteString(strconv.Itoa(slot["id"].(int)))
					buf.WriteString(slot["checksum"].(string))
					return schema.HashString(buf.String())
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Secret slot id, often used as 'effective from' timestamp.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"value": {
							Description: "Secret value.",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
						},
						"checksum": {
							Description: "Secret value checksum.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
		CreateContext: resourceFastEdgeSecretCreate,
		ReadContext:   resourceFastEdgeSecretRead,
		UpdateContext: resourceFastEdgeSecretUpdate,
		DeleteContext: resourceFastEdgeSecretDelete,
		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, meta any) error {
			if diff.GetRawState().IsNull() { // adding new resource - not need for checksum
				return nil
			}
			if v, ok := diff.Get("slot").(*schema.Set); ok {
				res := schema.NewSet(v.F, nil)
				list := v.List()
				for _, v := range list {
					p := v.(map[string]any)
					value := p["value"].(string)
					checksum := secretChecksum(value)
					res.Add(map[string]any{
						"id":       p["id"],
						"checksum": checksum,
						"value":    value,
					})
				}
				diff.SetNew("slot", res)
			}
			return nil
		},
	}
}

func resourceFastEdgeSecretCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge template creation")
	config := m.(*Config)
	client := config.FastEdgeClient

	secret := sdk.Secret{
		Name:        fieldValue[string](d, "name"),
		Comment:     fieldValue[string](d, "comment"),
		SecretSlots: getSlots(d),
	}

	rsp, err := client.AddSecretWithResponse(ctx, sdk.AddSecretJSONRequestBody{Secret: secret})
	if err != nil {
		return diag.Errorf("calling AddSecret API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		return diag.Errorf("calling AddSecret API: %s", extractErrorMessage(rsp.Body))
	}

	d.SetId(strconv.FormatInt(*rsp.JSON200.Id, 10))
	d.Set("slot", parseSecretSlots(rsp.JSON200.SecretSlots))

	log.Printf("[DEBUG] Finish FastEdge secret creation (id=%d)\n", rsp.JSON200.Id)
	return nil
}

func resourceFastEdgeSecretRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge secret read")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("converting id to number: %v", err)
	}
	rsp, err := client.GetSecretWithResponse(ctx, id)
	if err != nil {
		return diag.Errorf("calling GetSecret API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		return diag.Errorf("calling GetSecret API: %s", extractErrorMessage(rsp.Body))
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("[FastEdge secret (%d) was not found, removed from TF state", id),
				},
			}
		}
		return diag.Errorf("calling GetSecret API: %s", extractErrorMessage(rsp.Body))
	}

	secret := rsp.JSON200
	setField(d, "name", &secret.Name)
	setField(d, "comment", &secret.Comment)
	d.Set("slot", parseSecretSlots(rsp.JSON200.SecretSlots))

	log.Println("[DEBUG] Finish FastEdge secret read")
	return nil
}

func resourceFastEdgeSecretUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge secret update")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("converting id to number: %v", err)
	}
	secret := sdk.Secret{
		Name:        fieldValue[string](d, "name"),
		Comment:     fieldValue[string](d, "comment"),
		SecretSlots: getSlots(d),
	}

	rsp, err := client.UpdateSecretWithResponse(ctx, id, sdk.UpdateSecretJSONRequestBody{Secret: secret})
	if err != nil {
		return diag.Errorf("calling UpdateSecret API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("[FastEdge secret (%d) was not found, removed from TF state", id),
				},
			}
		}
		return diag.Errorf("calling UpdateSecret API: %s", extractErrorMessage(rsp.Body))
	}

	d.Set("slot", parseSecretSlots(rsp.JSON200.SecretSlots))

	log.Println("[DEBUG] Finish FastEdge secret update")
	return nil
}

func resourceFastEdgeSecretDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Println("[DEBUG] Start FastEdge secret deletion")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("converting id to number: %v", err)
	}

	rsp, err := client.DeleteSecretWithResponse(ctx, id, &sdk.DeleteSecretParams{Force: &AllowDeletionOfSecretsWithSlots})
	if err != nil {
		return diag.Errorf("calling DeleteSecret API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() == http.StatusConflict {
			diags = diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Secret (%d) is referenced so cannot be deleted, but removed from TF state", id),
				},
			}
		} else {
			return diag.Errorf("calling DeleteSecret API: %s", extractErrorMessage(rsp.Body))
		}
	}

	d.SetId("")
	log.Println("[DEBUG] Finish FastEdge secret deletion")
	return diags
}

func getSlots(d *schema.ResourceData) *[]sdk.SecretSlot {
	if v, ok := d.Get("slot").(*schema.Set); ok {
		list := v.List()
		slots := make([]sdk.SecretSlot, len(list))
		for i, v := range list {
			p := v.(map[string]any)
			value := p["value"].(string)
			slots[i] = sdk.SecretSlot{
				Slot:  int64(p["id"].(int)),
				Value: &value,
			}
		}
		return &slots
	}
	return nil
}

func parseSecretSlots(slots *[]sdk.SecretSlot) []map[string]any {
	if slots != nil {
		s := make([]map[string]any, len(*slots))
		for i, v := range *slots {
			s[i] = map[string]any{
				"id":       v.Slot,
				"checksum": *v.Checksum,
			}
		}
		return s
	}
	return nil
}

func secretChecksum(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
