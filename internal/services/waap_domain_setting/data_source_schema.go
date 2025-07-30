// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_setting

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapDomainSettingDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.Int64Attribute{
				Description: "The domain ID",
				Required:    true,
			},
			"api": schema.SingleNestedAttribute{
				Description: "API settings of a domain",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[WaapDomainSettingAPIDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"api_urls": schema.ListAttribute{
						Description: "The API URLs for a domain. If your domain has a common base URL for all API paths, it can be set here",
						Computed:    true,
						CustomType:  customfield.NewListType[types.String](ctx),
						ElementType: types.StringType,
					},
					"is_api": schema.BoolAttribute{
						Description: "Indicates if the domain is an API domain. All requests to an API domain are treated as API requests. If this is set to true then the `api_urls` field is ignored.",
						Computed:    true,
					},
				},
			},
			"ddos": schema.SingleNestedAttribute{
				Description: "DDoS settings for a domain.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[WaapDomainSettingDDOSDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"burst_threshold": schema.Int64Attribute{
						Description: "The burst threshold detects sudden rises in traffic. If it is met and the number of requests is at least five times the last 2-second interval, DDoS protection will activate. Default is 1000.",
						Computed:    true,
						Validators: []validator.Int64{
							int64validator.Between(30, 10000),
						},
					},
					"global_threshold": schema.Int64Attribute{
						Description: "The global threshold is responsible for identifying DDoS attacks with a slow rise in traffic. If the threshold is met and the current number of requests is at least double that of the previous 10-second window, DDoS protection will activate. Default is 5000.",
						Computed:    true,
						Validators: []validator.Int64{
							int64validator.Between(250, 50000),
						},
					},
					"sub_second_threshold": schema.Int64Attribute{
						Description: "The sub-second threshold protects WAAP servers against attacks from traffic bursts. When this threshold is reached, the DDoS mode will activate on the affected WAAP server, not the whole WAAP cluster. Default is 50.",
						Computed:    true,
						Validators: []validator.Int64{
							int64validator.Between(25, 1000),
						},
					},
				},
			},
		},
	}
}

func (d *WaapDomainSettingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *WaapDomainSettingDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
