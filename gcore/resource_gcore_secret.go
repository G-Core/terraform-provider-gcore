package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/secret/v1/secrets"
	secretsV2 "github.com/G-Core/gcorelabscloud-go/gcore/secret/v2/secrets"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SecretDeleting int = 1200
const SecretCreatingTimeout int = 1200
const secretPoint = "secrets"

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecretCreate,
		ReadContext:   resourceSecretRead,
		DeleteContext: resourceSecretDelete,
		Description:   "Represent secret",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, secretID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(secretID)

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
			"private_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SSL private key in PEM format",
			},
			"certificate_chain": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SSL certificate chain of intermediates and root certificates in PEM format",
			},
			"certificate": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SSL certificate in PEM format",
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
				Description: "Datetime when the secret will expire. The format is 2025-12-28T19:14:44",
				Optional:    true,
				Computed:    true,
				StateFunc: func(val interface{}) string {
					expTime, _ := time.Parse(gcorecloud.RFC3339NoZ, val.(string))
					return expTime.Format(gcorecloud.RFC3339NoZ)
				},
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					rawTime := i.(string)
					_, err := time.Parse(gcorecloud.RFC3339NoZ, rawTime)
					if err != nil {
						return diag.FromErr(err)
					}
					return nil
				},
			},
			"created": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Datetime when the secret was created. The format is 2025-12-28T19:14:44.180394",
				Computed:    true,
			},
		},
	}
}

func resourceSecretCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Secret creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, secretPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := secretsV2.CreateOpts{
		Name: d.Get("name").(string),
		Payload: secretsV2.PayloadOpts{
			CertificateChain: d.Get("certificate_chain").(string),
			Certificate:      d.Get("certificate").(string),
			PrivateKey:       d.Get("private_key").(string),
		},
	}
	if rawTime := d.Get("expiration").(string); rawTime != "" {
		expiration, err := time.Parse(gcorecloud.RFC3339NoZ, rawTime)
		if err != nil {
			return diag.FromErr(err)
		}
		opts.Expiration = &expiration
	}

	results, err := secretsV2.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)

	clientV1, err := CreateClient(provider, d, secretPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	secretID, err := tasks.WaitTaskAndReturnResult(clientV1, taskID, true, SecretCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(clientV1, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		Secret, err := secrets.ExtractSecretIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve Secret ID from task info: %w", err)
		}
		return Secret, nil
	},
	)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Secret id (%s)", secretID)

	d.SetId(secretID.(string))

	resourceSecretRead(ctx, d, m)

	log.Printf("[DEBUG] Finish Secret creating (%s)", secretID)
	return diags
}

func resourceSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	secret, err := secrets.Get(client, secretID).Extract()
	if err != nil {
		return diag.Errorf("cannot get secret with ID: %s. Error: %s", secretID, err.Error())
	}
	d.Set("name", secret.Name)
	d.Set("algorithm", secret.Algorithm)
	d.Set("bit_length", secret.BitLength)
	d.Set("mode", secret.Mode)
	d.Set("status", secret.Status)
	d.Set("expiration", secret.Expiration.Format(gcorecloud.RFC3339NoZ))
	d.Set("created", secret.CreatedAt.Format(gcorecloud.RFC3339MilliNoZ))
	if err := d.Set("content_types", secret.ContentTypes); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish secret reading")
	return diags
}

func resourceSecretDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start secret deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	secretID := d.Id()
	log.Printf("[DEBUG] Secret id = %s", secretID)

	client, err := CreateClient(provider, d, secretPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	results, err := secrets.Delete(client, secretID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, SecretDeleting, func(task tasks.TaskID) (interface{}, error) {
		_, err := secrets.Get(client, secretID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete secret with ID: %s", secretID)
		}
		return nil, nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Finish of secret deleting")
	return diags
}
