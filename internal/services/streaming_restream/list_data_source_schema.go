// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_restream

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*StreamingRestreamsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
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
				CustomType:  customfield.NewNestedObjectListType[StreamingRestreamsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"active": schema.BoolAttribute{
							Description: "Enables/Disables restream. Has two possible values:\n\n* **true** — restream is enabled and can be started\n* **false** — restream is disabled.\n\n  \nDefault is true",
							Computed:    true,
						},
						"client_user_id": schema.Int64Attribute{
							Description: "Custom field where you can specify user ID in your system",
							Computed:    true,
						},
						"live": schema.BoolAttribute{
							Description: "Indicates that the stream is being published. Has two possible values:\n\n* **true** — stream is being published\n* **false** — stream isn't published",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Restream name",
							Computed:    true,
						},
						"stream_id": schema.Int64Attribute{
							Description: "ID of the stream to restream",
							Computed:    true,
						},
						"uri": schema.StringAttribute{
							Description: "A URL to push the stream to",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *StreamingRestreamsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *StreamingRestreamsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
