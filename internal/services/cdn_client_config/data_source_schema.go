// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_client_config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CDNClientConfigDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"auto_suspend_enabled": schema.BoolAttribute{
				Description: "Defines whether resources will be deactivated automatically by inactivity.\n\nPossible values:\n- **true** - Resources will be deactivated.\n- **false** - Resources will not be deactivated.",
				Computed:    true,
			},
			"cdn_resources_rules_max_count": schema.Int64Attribute{
				Description: "Limit on the number of rules for each CDN resource.",
				Computed:    true,
			},
			"cname": schema.StringAttribute{
				Description: "Domain zone to which a CNAME record of your CDN resources should be pointed.",
				Computed:    true,
			},
			"created": schema.StringAttribute{
				Description: "Date of the first synchronization with the Platform (ISO 8601/RFC 3339 format, UTC.)",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "Account ID.",
				Computed:    true,
			},
			"updated": schema.StringAttribute{
				Description: "Date of the last update of information about CDN service (ISO 8601/RFC 3339 format, UTC.)",
				Computed:    true,
			},
			"use_balancer": schema.BoolAttribute{
				Description: "Defines whether custom balancing is used for content delivery.\n\nPossible values:\n- **true** - Custom balancing is used for content delivery.\n- **false** - Custom balancing is not used for content delivery.",
				Computed:    true,
			},
			"utilization_level": schema.Int64Attribute{
				Description: "CDN traffic usage limit in gigabytes.\n\nWhen the limit is reached, we will send an email notification.",
				Computed:    true,
			},
			"service": schema.SingleNestedAttribute{
				Description: "Information about the CDN service status.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CDNClientConfigServiceDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Defines whether the CDN service is activated.\n\nPossible values:\n- **true** - Service is activated.\n- **false** - Service is not activated.",
						Computed:    true,
					},
					"status": schema.StringAttribute{
						Description: "CDN service status.\n\nPossible values:\n- **new** - CDN service is not activated.\n- **trial** - Free trial is in progress.\n- **trialend** - Free trial has ended and CDN service is stopped. All CDN resources are suspended.\n- **activating** - CDN service is being activated. It can take up to 15 minutes.\n- **active** - CDN service is active.\n- **paused** - CDN service is stopped. All CDN resources are suspended.\n- **deleted** - CDN service is stopped. All CDN resources are deleted.",
						Computed:    true,
					},
					"updated": schema.StringAttribute{
						Description: "Date of the last CDN service status update (ISO 8601/RFC 3339 format, UTC.)",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *CDNClientConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CDNClientConfigDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
