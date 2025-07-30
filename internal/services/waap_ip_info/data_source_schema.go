// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_ip_info

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapIPInfoDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description: "The IP address to check",
				Required:    true,
			},
			"risk_score": schema.StringAttribute{
				Description: "The risk score of the IP address\nAvailable values: \"NO_RISK\", \"LOW\", \"MEDIUM\", \"HIGH\", \"EXTREME\", \"NOT_ENOUGH_DATA\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"NO_RISK",
						"LOW",
						"MEDIUM",
						"HIGH",
						"EXTREME",
						"NOT_ENOUGH_DATA",
					),
				},
			},
			"tags": schema.ListAttribute{
				Description: "The tags associated with the IP address that affect the risk score",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"whois": schema.SingleNestedAttribute{
				Description: "The WHOIS information for the IP address",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[WaapIPInfoWhoisDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"abuse_mail": schema.StringAttribute{
						Description: "The abuse mail",
						Computed:    true,
					},
					"cidr": schema.Int64Attribute{
						Description: "The CIDR",
						Computed:    true,
					},
					"country": schema.StringAttribute{
						Description: "The country",
						Computed:    true,
					},
					"net_description": schema.StringAttribute{
						Description: "The network description",
						Computed:    true,
					},
					"net_name": schema.StringAttribute{
						Description: "The network name",
						Computed:    true,
					},
					"net_range": schema.StringAttribute{
						Description: "The network range",
						Computed:    true,
					},
					"net_type": schema.StringAttribute{
						Description: "The network type",
						Computed:    true,
					},
					"org_id": schema.StringAttribute{
						Description: "The organization ID",
						Computed:    true,
					},
					"org_name": schema.StringAttribute{
						Description: "The organization name",
						Computed:    true,
					},
					"owner_type": schema.StringAttribute{
						Description: "The owner type",
						Computed:    true,
					},
					"rir": schema.StringAttribute{
						Description: "The RIR",
						Computed:    true,
					},
					"state": schema.StringAttribute{
						Description: "The state",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *WaapIPInfoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapIPInfoDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
