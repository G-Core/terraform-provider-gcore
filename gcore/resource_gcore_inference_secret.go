package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/inference/v3/secrets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInferenceSecrets() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInferenceSecretCreate,
		ReadContext:   resourceInferenceSecretRead,
		UpdateContext: resourceInferenceSecretUpdate,
		DeleteContext: resourceInferenceSecretDelete,
		Description:   "Represent inference secret. Specify this secret if you are using an AWS SQS-based trigger for inference deployment.",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, inferenceName, err := ImportStringParserWithNoRegion(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("name", inferenceName)
				d.SetId(inferenceName)

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
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_aws_access_key_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"data_aws_secret_access_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceInferenceSecretCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Inference secret creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	secretName := d.Get("name").(string)
	opts := secrets.CreateInferenceSecretOpts{
		Name: secretName,
		Type: "aws-iam",
	}

	opts.Data.AWSSecretKeyID = d.Get("data_aws_access_key_id").(string)
	opts.Data.AWSSecretAccessKey = d.Get("data_aws_secret_access_key").(string)

	secret, err := secrets.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(secret.Name)
	d.Set("data_aws_access_key_id", secret.Data.AWSSecretKeyID)
	d.Set("data_aws_secret_access_key", secret.Data.AWSSecretAccessKey)
	resourceInferenceDeploymentRead(ctx, d, m)

	log.Printf("[DEBUG] Finish inference secret creating (%s)", secretName)
	return diags
}

func resourceInferenceSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Inference secret reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	secretName := d.Id()
	secret, err := secrets.Get(client, secretName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", secret.Name)
	d.Set("type", secret.Type)

	log.Printf("[DEBUG] Finish inference secret reading (%s)", secretName)
	return diags
}

func resourceInferenceSecretUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Inference secret updating")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	secretName := d.Id()
	opts := secrets.UpdateInferenceSecretOpts{
		Type: "aws-iam",
	}

	opts.Data.AWSSecretKeyID = d.Get("data_aws_access_key_id").(string)
	opts.Data.AWSSecretAccessKey = d.Get("data_aws_secret_access_key").(string)

	if d.HasChange("data_aws_access_key_id") || d.HasChange("data_aws_secret_access_key") {
		secret, err := secrets.Update(client, secretName, opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("data_aws_access_key_id", secret.Data.AWSSecretKeyID)
		d.Set("data_aws_secret_access_key", secret.Data.AWSSecretAccessKey)
	}

	log.Printf("[DEBUG] Finish inference secret updating (%s)", secretName)
	return resourceInferenceSecretRead(ctx, d, m)
}

func resourceInferenceSecretDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start inference secret deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	secretName := d.Id()
	log.Printf("[DEBUG] inference secret = %s", secretName)

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := secrets.Delete(client, secretName).ExtractErr(); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Finish of inference secret deleting (%s)", secretName)
	return diags
}
