// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_insight_silence

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainInsightSilenceDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "A generated unique identifier for the silence",
				Computed:    true,
			},
			"silence_id": schema.StringAttribute{
				Description: "A generated unique identifier for the silence",
				Optional:    true,
			},
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"author": schema.StringAttribute{
				Description: "The author of the silence",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "A comment explaining the reason for the silence",
				Computed:    true,
			},
			"expire_at": schema.StringAttribute{
				Description: "The date and time the silence expires in ISO 8601 format",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"insight_type": schema.StringAttribute{
				Description: "The slug of the insight type",
				Computed:    true,
			},
			"labels": schema.MapAttribute{
				Description: "A hash table of label names and values that apply to the insight silence",
				Computed:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.ListAttribute{
						Description: "The ID of the insight silence",
						Optional:    true,
						ElementType: types.StringType,
					},
					"author": schema.StringAttribute{
						Description: "The author of the insight silence",
						Optional:    true,
					},
					"comment": schema.StringAttribute{
						Description: "The comment of the insight silence",
						Optional:    true,
					},
					"insight_type": schema.ListAttribute{
						Description: "The type of the insight silence",
						Optional:    true,
						ElementType: types.StringType,
					},
					"ordering": schema.StringAttribute{
						Description: "Sort the response by given field.\nAvailable values: \"id\", \"-id\", \"insight_type\", \"-insight_type\", \"comment\", \"-comment\", \"author\", \"-author\", \"expire_at\", \"-expire_at\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"id",
								"-id",
								"insight_type",
								"-insight_type",
								"comment",
								"-comment",
								"author",
								"-author",
								"expire_at",
								"-expire_at",
							),
						},
					},
				},
			},
		},
	}
}

func (d *WaapDomainInsightSilenceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapDomainInsightSilenceDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("silence_id"), path.MatchRoot("find_one_by")),
	}
}
