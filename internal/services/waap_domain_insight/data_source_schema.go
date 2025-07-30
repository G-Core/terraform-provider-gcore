// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_insight

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainInsightDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"insight_id": schema.StringAttribute{
				Required: true,
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
			"id": schema.StringAttribute{
				Description: "A generated unique identifier for the insight",
				Computed:    true,
			},
			"insight_type": schema.StringAttribute{
				Description: "The type of the insight represented as a slug",
				Computed:    true,
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
			"labels": schema.MapAttribute{
				Description: "A hash table of label names and values that apply to the insight",
				Computed:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
		},
	}
}

func (d *WaapDomainInsightDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapDomainInsightDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
