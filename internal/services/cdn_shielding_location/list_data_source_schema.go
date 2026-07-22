// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_shielding_location

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CDNShieldingLocationsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "CDN origin shielding locations are the available shield (precache) points of presence that can be referenced when configuring origin shielding for a CDN resource.",
		Attributes: map[string]schema.Attribute{
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CDNShieldingLocationsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Origin shielding location ID.",
							Computed:    true,
						},
						"city": schema.StringAttribute{
							Description: "City of origin shielding location.",
							Computed:    true,
						},
						"country": schema.StringAttribute{
							Description: "Country of origin shielding location.",
							Computed:    true,
						},
						"datacenter": schema.StringAttribute{
							Description: "Name of origin shielding location datacenter.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CDNShieldingLocationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CDNShieldingLocationsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
