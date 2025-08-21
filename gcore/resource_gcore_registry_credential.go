package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/inference/v3/credentials"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRegistryCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRegistryCredentialCreate,
		ReadContext:   resourceRegistryCredentialRead,
		UpdateContext: resourceRegistryCredentialUpdate,
		DeleteContext: resourceRegistryCredentialDelete,
		Description:   "Represent inference registry credential",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, credsName, err := ImportStringParserWithNoRegion(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.SetId(credsName)

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
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"registry_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRegistryCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Inference deployment creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	credsName := d.Get("name").(string)
	opts := credentials.CreateRegistryCredentialOpts{
		Name:        credsName,
		Username:    d.Get("username").(string),
		Password:    d.Get("password").(string),
		RegistryURL: d.Get("registry_url").(string),
	}

	cr, err := credentials.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Registry credential name (%s)", credsName)
	d.SetId(cr.Name)

	resourceInferenceDeploymentRead(ctx, d, m)

	log.Printf("[DEBUG] Finish registry credential creating (%s)", credsName)
	return diags
}

func resourceRegistryCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start registry credential reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	credsName := d.Id()
	log.Printf("[DEBUG] registry credential name = %s", credsName)

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	creds, err := credentials.Get(client, credsName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", creds.Name)
	d.Set("username", creds.Username)
	d.Set("registry_url", creds.RegistryURL)
	d.Set("project_id", creds.ProjectID)

	log.Println("[DEBUG] Finish registry credential reading")
	return diags
}

func resourceRegistryCredentialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start registry credential updating")
	config := m.(*Config)
	provider := config.Provider
	credsName := d.Id()
	log.Printf("[DEBUG] registry credential = %s", credsName)

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := credentials.UpdateRegistryCredentialOpts{
		Username:    d.Get("username").(string),
		Password:    d.Get("password").(string),
		RegistryURL: d.Get("registry_url").(string),
	}

	_, err = credentials.Update(client, credsName, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Finish of registry credential updating")
	return resourceRegistryCredentialRead(ctx, d, m)
}

func resourceRegistryCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start registry credential deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	credsName := d.Id()
	log.Printf("[DEBUG] registry credential = %s", credsName)

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := credentials.Delete(client, credsName).ExtractErr(); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Finish of registry credential deleting")
	return diags
}
