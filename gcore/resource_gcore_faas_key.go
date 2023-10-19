package gcore

import (
	"context"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFaaSKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFaaSKeyCreate,
		ReadContext:   resourceFaaSKeyRead,
		UpdateContext: resourceFaaSKeyUpdate,
		DeleteContext: resourceFaaSKeyDelete,
		Description:   "Represents FaaS API key",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, keyName, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(keyName)

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"functions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "API key secret",
			},
		},
	}
}

func resourceFaaSKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS key creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	keyName := d.Get("name").(string)
	log.Printf("[DEBUG] key = %s", keyName)

	client, err := CreateClient(provider, d, faasKeysPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := faas.CreateKeyOpts{
		Name:        keyName,
		Description: d.Get("description").(string),
	}

	expireRaw := d.Get("expire").(string)
	if len(expireRaw) > 0 {
		t, err := time.Parse(gcorecloud.RFC3339ZZ, expireRaw)
		if err != nil {
			return diag.FromErr(err)
		}
		expire := gcorecloud.JSONRFC3339ZZ{Time: t}
		opts.Expire = &expire
	}

	functionsRaw := d.Get("functions").([]any)
	if len(functionsRaw) > 0 {
		functions := make([]faas.KeysFunction, len(functionsRaw))
		for idx, fn := range functionsRaw {
			m := fn.(map[string]any)
			functions[idx] = faas.KeysFunction{
				Name:      m["name"].(string),
				Namespace: m["namespace"].(string),
			}
		}
		opts.Functions = functions
	}

	key, err := faas.CreateKey(client, opts)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(keyName)
	d.Set("secret", key.Secret)

	resourceFaaSKeyRead(ctx, d, m)

	log.Printf("[DEBUG] Finish FaaS key creating (%s)", keyName)
	return diags
}

func resourceFaaSKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS key reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	keyName := d.Id()
	log.Printf("[DEBUG] key = %s", keyName)

	client, err := CreateClient(provider, d, faasKeysPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	key, err := faas.GetKey(client, keyName).Extract()
	if err != nil {
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			log.Printf("[WARN] Removing key %s because resource doesn't exist anymore", d.Id())
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}
	d.Set("name", key.Name)
	d.Set("description", key.Description)
	d.Set("status", key.Status)
	if !key.Expire.IsZero() {
		d.Set("expire", key.Expire.Format(time.RFC3339))
	}
	d.Set("created_at", key.CreatedAt.Format(time.RFC3339))
	functions := make([]map[string]string, len(key.Functions))
	for idx, f := range key.Functions {
		functions[idx] = map[string]string{
			"name":      f.Name,
			"namespace": f.Namespace,
		}
	}
	if err := d.Set("functions", functions); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Finish FaaS key reading (%s)", keyName)
	return diags
}

func resourceFaaSKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS key updating")
	config := m.(*Config)
	provider := config.Provider
	keyName := d.Id()
	log.Printf("[DEBUG] key = %s", keyName)

	client, err := CreateClient(provider, d, faasKeysPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	var (
		needUpdate bool
		opts       faas.UpdateKeyOpts
	)

	if d.HasChange("description") {
		opts.Description = d.Get("description").(string)
		needUpdate = true
	}

	if d.HasChange("expire") {
		expireRaw := d.Get("expire").(string)
		if len(expireRaw) > 0 {
			t, err := time.Parse(gcorecloud.RFC3339ZZ, expireRaw)
			if err != nil {
				return diag.FromErr(err)
			}
			expire := gcorecloud.JSONRFC3339ZZ{Time: t}
			opts.Expire = &expire
		}
		needUpdate = true
	}

	if d.HasChange("functions") {
		functionsRaw := d.Get("functions").([]any)
		functions := make([]faas.KeysFunction, len(functionsRaw))
		for idx, fn := range functionsRaw {
			m := fn.(map[string]any)
			functions[idx] = faas.KeysFunction{
				Name:      m["name"].(string),
				Namespace: m["namespace"].(string),
			}
		}
		opts.Functions = functions
	}

	if needUpdate {
		_, err := faas.UpdateKey(client, keyName, opts)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] Finish FaaS key updating (%s)", keyName)
	return resourceFaaSKeyRead(ctx, d, m)
}

func resourceFaaSKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FaaS key deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	keyName := d.Id()
	log.Printf("[DEBUG] key = %s", keyName)

	client, err := CreateClient(provider, d, faasKeysPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	err = faas.DeleteKey(client, keyName)
	switch err.(type) {
	case nil, gcorecloud.ErrDefault404:
		d.SetId("")
		log.Printf("[DEBUG] Finish FaaS key deleting (%s)", keyName)
		return diags
	default:
		return diag.FromErr(err)
	}
}
