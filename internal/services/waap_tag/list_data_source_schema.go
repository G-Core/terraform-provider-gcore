// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_tag

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapTagsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filter tags by their name. Supports '\\*' as a wildcard character.",
				Optional:    true,
			},
			"ordering": schema.StringAttribute{
				Description: "Determine the field to order results by\nAvailable values: \"name\", \"readable_name\", \"reserved\", \"-name\", \"-readable_name\", \"-reserved\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"name",
						"readable_name",
						"reserved",
						"-name",
						"-readable_name",
						"-reserved",
					),
				},
			},
			"readable_name": schema.StringAttribute{
				Description: "Filter tags by their readable name. Supports '\\*' as a wildcard character.",
				Optional:    true,
			},
			"reserved": schema.BoolAttribute{
				Description: "Filter to include only reserved tags.",
				Optional:    true,
			},
			"limit": schema.Int64Attribute{
				Description: "Number of items to return",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
			},
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
				CustomType:  customfield.NewNestedObjectListType[WaapTagsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Description: "A tag's human readable description",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of a tag that should be used in a WAAP rule condition",
							Computed:    true,
						},
						"readable_name": schema.StringAttribute{
							Description: "The display name of the tag",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *WaapTagsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *WaapTagsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
