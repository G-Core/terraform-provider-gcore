// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_registry_credential

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*CloudInferenceRegistryCredentialResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Registry credentials store authentication details for private container registries used by inference deployments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "Registry credential name.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "Registry credential name.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.Int64Attribute{
				Description:   "Project ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
			},
			"password_wo": schema.StringAttribute{
				Description: "Registry password.",
				Required:    true,
				WriteOnly:   true,
			},
			"registry_url": schema.StringAttribute{
				Description: "Registry URL.",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Registry username.",
				Required:    true,
			},
			"password_wo_version": schema.Int64Attribute{
				Description: "Registry credential password write-only version. Used to trigger updates of the " +
					"write-only password field.",
				Required: true,
			},
		},
	}
}

func (r *CloudInferenceRegistryCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudInferenceRegistryCredentialResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
