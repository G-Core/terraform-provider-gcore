// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_api_key

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceAPIKeysDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"limit": schema.Int64Attribute{
				Description: "Optional. Limit the number of returned items",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtMost(1000),
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
				CustomType:  customfield.NewNestedObjectListType[CloudInferenceAPIKeysItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Description: "Timestamp when the API Key was created.",
							Computed:    true,
						},
						"deployment_names": schema.ListAttribute{
							Description: "List of inference deployment names to which this API Key has been attached.",
							Computed:    true,
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
						},
						"description": schema.StringAttribute{
							Description: "Description of the API Key.",
							Computed:    true,
						},
						"expires_at": schema.StringAttribute{
							Description: "Timestamp when the API Key will expire.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "API Key name.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudInferenceAPIKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudInferenceAPIKeysDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
