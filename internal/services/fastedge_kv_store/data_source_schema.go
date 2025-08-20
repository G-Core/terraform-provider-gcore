// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_kv_store

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*FastedgeKvStoreDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"app_count": schema.Int64Attribute{
				Description: "The number of applications that use this store",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "A description of the store",
				Computed:    true,
			},
			"updated": schema.StringAttribute{
				Description: "Last update time",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"byod": schema.SingleNestedAttribute{
				Description: "BYOD (Bring Your Own Data) settings",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[FastedgeKvStoreByodDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"prefix": schema.StringAttribute{
						Description: "Key prefix",
						Computed:    true,
					},
					"url": schema.StringAttribute{
						Description: "URL to connect to",
						Computed:    true,
					},
				},
			},
			"stats": schema.SingleNestedAttribute{
				Description: "Store statistics",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[FastedgeKvStoreStatsDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"cf_count": schema.Int64Attribute{
						Description: "Total number of Cuckoo filter entries",
						Computed:    true,
					},
					"kv_count": schema.Int64Attribute{
						Description: "Total number of KV entries",
						Computed:    true,
					},
					"size": schema.Int64Attribute{
						Description: "Total store size in bytes",
						Computed:    true,
					},
					"zset_count": schema.Int64Attribute{
						Description: "Total number of sorted set entries",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *FastedgeKvStoreDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *FastedgeKvStoreDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
