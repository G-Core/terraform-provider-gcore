// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster_image

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudGPUBaremetalClusterImageDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"image_id": schema.StringAttribute{
				Description: "Image ID",
				Required:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Required:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Required:    true,
			},
			"architecture": schema.StringAttribute{
				Description: "Image architecture type",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the image was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"id": schema.StringAttribute{
				Description: "Image ID",
				Computed:    true,
			},
			"min_disk": schema.Int64Attribute{
				Description: "Minimal boot volume required",
				Computed:    true,
			},
			"min_ram": schema.Int64Attribute{
				Description: "Minimal VM RAM required",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Image name",
				Computed:    true,
			},
			"os_distro": schema.StringAttribute{
				Description: "OS Distribution",
				Computed:    true,
			},
			"os_type": schema.StringAttribute{
				Description: "The operating system installed on the image",
				Computed:    true,
			},
			"os_version": schema.StringAttribute{
				Description: "OS version, i.e. 19.04 (for Ubuntu) or 9.4 for Debian",
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "Image size in bytes.",
				Computed:    true,
			},
			"ssh_key": schema.StringAttribute{
				Description: "Whether the image supports SSH key or not",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Image status",
				Computed:    true,
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the image was updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"visibility": schema.StringAttribute{
				Description: "Image visibility. Globally visible images are public",
				Computed:    true,
			},
			"tags": schema.ListNestedAttribute{
				Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClusterImageTagsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "Tag key. The maximum size for a key is 255 bytes.",
							Computed:    true,
						},
						"read_only": schema.BoolAttribute{
							Description: "If true, the tag is read-only and cannot be modified by the user",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "Tag value. The maximum size for a value is 1024 bytes.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudGPUBaremetalClusterImageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudGPUBaremetalClusterImageDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
