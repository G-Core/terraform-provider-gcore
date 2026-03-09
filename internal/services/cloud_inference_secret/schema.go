// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_secret

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*CloudInferenceSecretResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Inference secrets store sensitive values such as AWS credentials used for SQS-based autoscaling triggers in deployments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "Secret name.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "Secret name.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.Int64Attribute{
				Description:   "Project ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"type": schema.StringAttribute{
				Description: "Secret type. Currently only `aws-iam` is supported.",
				Required:    true,
			},
			"data": schema.SingleNestedAttribute{
				Description: "Secret data.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"aws_access_key_id_wo": schema.StringAttribute{
						Description: "AWS IAM key ID.",
						Required:    true,
						WriteOnly:   true,
					},
					"aws_secret_access_key_wo": schema.StringAttribute{
						Description: "AWS IAM secret key.",
						Required:    true,
						WriteOnly:   true,
					},
				},
			},
			"data_wo_version": schema.Int64Attribute{
				Description: "The version of the data sensitive params - used to trigger updates of write-only params.",
				Required:    true,
			},
		},
	}
}

func (r *CloudInferenceSecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudInferenceSecretResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
