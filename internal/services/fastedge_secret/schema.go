// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_secret

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*FastedgeSecretResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "FastEdge secrets store sensitive values such as API keys and tokens that can be referenced by FastEdge applications.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "The unique identifier of the secret.",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown(), int64planmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "The unique name of the secret.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"comment": schema.StringAttribute{
				Description:   "A description or comment about the secret.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"secret_slots": schema.SetNestedAttribute{
				Description: "A list of secret slots associated with this secret.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectSetType[FastedgeSecretSecretSlotsModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"slot": schema.Int64Attribute{
							Description: "Unix timestamp (seconds since epoch) indicating when this secret version becomes active. Use for time-based secret rotation.",
							Required:    true,
						},
						"checksum": schema.StringAttribute{
							Description: "SHA-256 hash of the decrypted value for integrity verification (auto-generated)",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "The plaintext secret value. Will be encrypted with AES-256-GCM before storage.",
							Optional:    true,
						},
					},
				},
				PlanModifiers: []planmodifier.Set{setplanmodifier.RequiresReplaceIfConfigured()},
			},
			"app_count": schema.Int64Attribute{
				Description: "The number of applications that use this secret.",
				Computed:    true,
			},
		},
	}
}

func (r *FastedgeSecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *FastedgeSecretResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
