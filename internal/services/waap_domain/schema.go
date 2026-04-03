// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*WaapDomainResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "WAAP domains enable Web Application and API Protection for monitoring and defending web applications against security threats.",
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description:   "The domain ID",
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"status": schema.StringAttribute{
				Description: "The current status of the domain\nAvailable values: \"active\", \"monitor\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("active", "monitor"),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "The date and time the domain was created in ISO 8601 format",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"custom_page_set": schema.Int64Attribute{
				Description: "The ID of the custom page set",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "The domain ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The domain name",
				Computed:    true,
			},
			"quotas": schema.MapNestedAttribute{
				Description: "Domain level quotas",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectMapType[WaapDomainQuotasModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allowed": schema.Int64Attribute{
							Description: "The maximum allowed number of this resource",
							Computed:    true,
						},
						"current": schema.Int64Attribute{
							Description: "The current number of this resource",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *WaapDomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *WaapDomainResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
