// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_image

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudBaremetalImagesDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Bare metal images are operating system images used to boot bare metal servers, filterable by name, visibility, OS distribution, and architecture.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Optional:    true,
			},
			"architecture": schema.StringAttribute{
				Description: "Filter by image architecture.\nAvailable values: \"aarch64\", \"x86_64\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("aarch64", "x86_64"),
				},
			},
			"include_prices": schema.BoolAttribute{
				Description: "Show price.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Filter by image name (case-insensitive substring match)",
				Optional:    true,
			},
			"os_distro": schema.StringAttribute{
				Description: "Filter by OS distribution (case-insensitive). E.g. `ubuntu`, `centos`, `debian`",
				Optional:    true,
			},
			"os_version": schema.StringAttribute{
				Description: "Filter by OS version (case-insensitive). E.g. `22.04`",
				Optional:    true,
			},
			"private": schema.StringAttribute{
				Description: "Any value to show private images",
				Optional:    true,
			},
			"tag_key_value": schema.StringAttribute{
				Description: "Optional. Filter by tag key-value pairs.",
				Optional:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "Image visibility. Globally visible images are public\nAvailable values: \"private\", \"public\", \"shared\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"private",
						"public",
						"shared",
					),
				},
			},
			"tag_key": schema.ListAttribute{
				Description: "Optional. Filter by tag keys. ?`tag_key`=key1&`tag_key`=key2",
				Optional:    true,
				ElementType: types.StringType,
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
				CustomType:  customfield.NewNestedObjectListType[CloudBaremetalImagesItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Image ID",
							Computed:    true,
						},
						"architecture": schema.StringAttribute{
							Description: "An image architecture type: aarch64, `x86_64`.\nAvailable values: \"aarch64\", \"x86_64\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("aarch64", "x86_64"),
							},
						},
						"created_at": schema.StringAttribute{
							Description: "Datetime when the image was created",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"creator_task_id": schema.StringAttribute{
							Description: "Task that created this entity",
							Computed:    true,
						},
						"currency_code": schema.StringAttribute{
							Description: "Currency code. Shown if the `include_prices` query parameter if set to true",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Image description",
							Computed:    true,
						},
						"disk_format": schema.StringAttribute{
							Description: "Disk format",
							Computed:    true,
						},
						"display_order": schema.Int64Attribute{
							Computed: true,
						},
						"gpu_driver": schema.StringAttribute{
							Description: "Name of the GPU driver vendor",
							Computed:    true,
						},
						"gpu_driver_type": schema.StringAttribute{
							Description: "Type of the GPU driver",
							Computed:    true,
						},
						"gpu_driver_version": schema.StringAttribute{
							Description: "Version of the installed GPU driver",
							Computed:    true,
						},
						"hw_firmware_type": schema.StringAttribute{
							Description: "Specifies the type of firmware with which to boot the guest.\nAvailable values: \"bios\", \"uefi\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("bios", "uefi"),
							},
						},
						"hw_machine_type": schema.StringAttribute{
							Description: "A virtual chipset type.\nAvailable values: \"i440\", \"q35\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("i440", "q35"),
							},
						},
						"is_baremetal": schema.BoolAttribute{
							Description: "For bare metal servers this value is always set to true",
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
							Description: "Image display name",
							Computed:    true,
						},
						"os_distro": schema.StringAttribute{
							Description: "OS Distribution, i.e. Debian, CentOS, Ubuntu, CoreOS etc.",
							Computed:    true,
						},
						"os_type": schema.StringAttribute{
							Description: "The operating system installed on the image.\nAvailable values: \"linux\", \"windows\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("linux", "windows"),
							},
						},
						"os_version": schema.StringAttribute{
							Description: "OS version, i.e. 19.04 (for Ubuntu) or 9.4 for Debian",
							Computed:    true,
						},
						"price_per_hour": schema.Float64Attribute{
							Description: "Price per hour. Shown if the `include_prices` query parameter if set to true",
							Computed:    true,
						},
						"price_per_month": schema.Float64Attribute{
							Description: "Price per month. Shown if the `include_prices` query parameter if set to true",
							Computed:    true,
						},
						"price_status": schema.StringAttribute{
							Description: "Price status for the UI\nAvailable values: \"error\", \"hide\", \"show\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"error",
									"hide",
									"show",
								),
							},
						},
						"project_id": schema.Int64Attribute{
							Description: "Project ID",
							Computed:    true,
						},
						"region": schema.StringAttribute{
							Description: "Region name",
							Computed:    true,
						},
						"region_id": schema.Int64Attribute{
							Description: "Region ID",
							Computed:    true,
						},
						"size": schema.Int64Attribute{
							Description: "Image size in bytes",
							Computed:    true,
						},
						"ssh_key": schema.StringAttribute{
							Description: "Whether the image supports SSH key or not\nAvailable values: \"allow\", \"deny\", \"required\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"allow",
									"deny",
									"required",
								),
							},
						},
						"status": schema.StringAttribute{
							Description: "Image status, i.e. active",
							Computed:    true,
						},
						"tags_v2": schema.ListNestedAttribute{
							Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudBaremetalImagesTagsV2DataSourceModel](ctx),
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
					},
				},
			},
		},
	}
}

func (d *CloudBaremetalImagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudBaremetalImagesDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
