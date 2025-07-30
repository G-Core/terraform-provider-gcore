// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_insight

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainInsightsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the insight. Supports '\\*' as a wildcard.",
				Optional:    true,
			},
			"ordering": schema.StringAttribute{
				Description: "Sort the response by given field.\nAvailable values: \"id\", \"-id\", \"insight_type\", \"-insight_type\", \"first_seen\", \"-first_seen\", \"last_seen\", \"-last_seen\", \"last_status_change\", \"-last_status_change\", \"status\", \"-status\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"id",
						"-id",
						"insight_type",
						"-insight_type",
						"first_seen",
						"-first_seen",
						"last_seen",
						"-last_seen",
						"last_status_change",
						"-last_status_change",
						"status",
						"-status",
					),
				},
			},
			"id": schema.ListAttribute{
				Description: "The ID of the insight",
				Optional:    true,
				ElementType: types.StringType,
			},
			"insight_type": schema.ListAttribute{
				Description: "The type of the insight",
				Optional:    true,
				ElementType: types.StringType,
			},
			"status": schema.ListAttribute{
				Description: "The status of the insight",
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOfCaseInsensitive(
							"OPEN",
							"ACKED",
							"CLOSED",
						),
					),
				},
				ElementType: types.StringType,
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
				CustomType:  customfield.NewNestedObjectListType[WaapDomainInsightsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "A generated unique identifier for the insight",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the insight",
							Computed:    true,
						},
						"first_seen": schema.StringAttribute{
							Description: "The date and time the insight was first seen in ISO 8601 format",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"insight_type": schema.StringAttribute{
							Description: "The type of the insight represented as a slug",
							Computed:    true,
						},
						"labels": schema.MapAttribute{
							Description: "A hash table of label names and values that apply to the insight",
							Computed:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"last_seen": schema.StringAttribute{
							Description: "The date and time the insight was last seen in ISO 8601 format",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"last_status_change": schema.StringAttribute{
							Description: "The date and time the insight was last seen in ISO 8601 format",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"recommendation": schema.StringAttribute{
							Description: "The recommended action to perform to resolve the insight",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The different statuses an insight can have\nAvailable values: \"OPEN\", \"ACKED\", \"CLOSED\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"OPEN",
									"ACKED",
									"CLOSED",
								),
							},
						},
					},
				},
			},
		},
	}
}

func (d *WaapDomainInsightsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *WaapDomainInsightsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
