// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "WAAP domains enable Web Application and API Protection for monitoring and defending web applications against security threats.",
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
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
			"status": schema.StringAttribute{
				Description: "The different statuses a domain can have\nAvailable values: \"active\", \"bypass\", \"monitor\", \"locked\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"active",
						"bypass",
						"monitor",
						"locked",
					),
				},
			},
			"quotas": schema.MapNestedAttribute{
				Description: "Domain level quotas",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectMapType[WaapDomainQuotasDataSourceModel](ctx),
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

func (d *WaapDomainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapDomainDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
