// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_insight_silence

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainInsightSilenceDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"silence_id": schema.StringAttribute{
				Description: "A generated unique identifier for the silence",
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
			"id": schema.StringAttribute{
				Description: "A generated unique identifier for the silence",
				Computed:    true,
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
		},
	}
}

func (d *WaapDomainInsightSilenceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapDomainInsightSilenceDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
