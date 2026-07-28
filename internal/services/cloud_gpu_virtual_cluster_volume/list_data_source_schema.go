// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster_volume

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudGPUVirtualClusterVolumesDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Description: "Cluster unique identifier",
				Required:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Optional:    true,
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
				CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClusterVolumesItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Volume unique identifier",
							Computed:    true,
						},
						"bootable": schema.BoolAttribute{
							Description: "True if this is bootable volume",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Volume creation date and time",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"name": schema.StringAttribute{
							Description: "User defined name",
							Computed:    true,
						},
						"root_fs": schema.BoolAttribute{
							Description: "True if this volume contains root file system",
							Computed:    true,
						},
						"server_id": schema.StringAttribute{
							Description: "Server UUID",
							Computed:    true,
						},
						"size": schema.Int64Attribute{
							Description: "Volume size in GiB",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Current volume status\nAvailable values: \"attaching\", \"available\", \"awaiting-transfer\", \"backing-up\", \"creating\", \"deleting\", \"detaching\", \"downloading\", \"error\", \"error_backing-up\", \"error_deleting\", \"error_extending\", \"error_restoring\", \"extending\", \"in-use\", \"maintenance\", \"reserved\", \"restoring-backup\", \"retyping\", \"reverting\", \"uploading\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"attaching",
									"available",
									"awaiting-transfer",
									"backing-up",
									"creating",
									"deleting",
									"detaching",
									"downloading",
									"error",
									"error_backing-up",
									"error_deleting",
									"error_extending",
									"error_restoring",
									"extending",
									"in-use",
									"maintenance",
									"reserved",
									"restoring-backup",
									"retyping",
									"reverting",
									"uploading",
								),
							},
						},
						"tags": schema.ListNestedAttribute{
							Description: "User defined tags",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClusterVolumesTagsDataSourceModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Description: "Tag key. Maximum 255 characters. Cannot contain spaces, tabs, newlines, empty string or '=' character.",
										Computed:    true,
									},
									"read_only": schema.BoolAttribute{
										Description: "If true, the tag is read-only and cannot be modified by the user",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Tag value. Maximum 255 characters. Cannot contain spaces, tabs, newlines, empty string or '=' character.",
										Computed:    true,
									},
								},
							},
						},
						"type": schema.StringAttribute{
							Description: "Volume type",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudGPUVirtualClusterVolumesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudGPUVirtualClusterVolumesDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
