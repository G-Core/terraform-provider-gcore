// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudRegistryDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Required: true,
			},
			"region_id": schema.Int64Attribute{
				Required: true,
			},
			"registry_id": schema.Int64Attribute{
				Required: true,
			},
			"created_at": schema.StringAttribute{
				Description: "Registry creation date-time",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"id": schema.Int64Attribute{
				Description: "Registry ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Registry name",
				Computed:    true,
			},
			"repo_count": schema.Int64Attribute{
				Description: "Number of repositories in the registry",
				Computed:    true,
			},
			"storage_limit": schema.Int64Attribute{
				Description: "Registry storage limit, GiB",
				Computed:    true,
			},
			"storage_used": schema.Int64Attribute{
				Description: "Registry storage used, bytes",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Registry modification date-time",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"url": schema.StringAttribute{
				Description: "Registry url",
				Computed:    true,
			},
		},
	}
}

func (d *CloudRegistryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudRegistryDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
