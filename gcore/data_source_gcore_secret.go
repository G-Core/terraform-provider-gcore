package gcore

import (
	"context"
	"log"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/secret/v1/secrets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretRead,
		Description: "Represent secret",
		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"algorithm": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bit_length": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_types": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"expiration": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Datetime when the secret will expire. The format is 2025-12-28T19:14:44.180394",
				Computed:    true,
			},
			"created": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Datetime when the secret was created. The format is 2025-12-28T19:14:44.180394",
				Computed:    true,
			},
		},
	}
}

func dataSourceSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start secret reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	secretID := d.Id()
	log.Printf("[DEBUG] Secret id = %s", secretID)

	client, err := CreateClient(provider, d, secretPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	allSecrets, err := secrets.ListAll(client)
	if err != nil {
		return diag.Errorf("cannot get secrets. Error: %s", err.Error())
	}

	var found bool
	name := d.Get("name").(string)
	for _, secret := range allSecrets {
		if name == secret.Name {
			d.SetId(secret.ID)
			d.Set("name", name)
			d.Set("algorithm", secret.Algorithm)
			d.Set("bit_length", secret.BitLength)
			d.Set("mode", secret.Mode)
			d.Set("status", secret.Status)
			d.Set("expiration", secret.Expiration.Format(gcorecloud.RFC3339ZColon))
			d.Set("created", secret.CreatedAt.Format(gcorecloud.RFC3339ZColon))
			if err := d.Set("content_types", secret.ContentTypes); err != nil {
				return diag.FromErr(err)
			}
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("secret with name %s does not exit", name)
	}

	log.Println("[DEBUG] Finish secret reading")
	return diags
}
