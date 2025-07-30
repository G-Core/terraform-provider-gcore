// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_setting

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*WaapDomainSettingResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description:   "The domain ID",
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"api": schema.SingleNestedAttribute{
				Description: "Editable API settings of a domain",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"api_urls": schema.ListAttribute{
						Description: "The API URLs for a domain. If your domain has a common base URL for all API paths, it can be set here",
						Optional:    true,
						ElementType: types.StringType,
					},
					"is_api": schema.BoolAttribute{
						Description: "Indicates if the domain is an API domain. All requests to an API domain are treated as API requests. If this is set to true then the `api_urls` field is ignored.",
						Optional:    true,
					},
				},
			},
			"ddos": schema.SingleNestedAttribute{
				Description: "Editable DDoS settings for a domain.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"burst_threshold": schema.Int64Attribute{
						Description: "The burst threshold detects sudden rises in traffic. If it is met and the number of requests is at least five times the last 2-second interval, DDoS protection will activate. Default is 1000.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(30, 10000),
						},
					},
					"global_threshold": schema.Int64Attribute{
						Description: "The global threshold is responsible for identifying DDoS attacks with a slow rise in traffic. If the threshold is met and the current number of requests is at least double that of the previous 10-second window, DDoS protection will activate. Default is 5000.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(250, 50000),
						},
					},
				},
			},
		},
	}
}

func (r *WaapDomainSettingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *WaapDomainSettingResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
