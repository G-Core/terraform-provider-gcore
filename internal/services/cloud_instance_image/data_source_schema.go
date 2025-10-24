// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance_image

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInstanceImageDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"image_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.Int64Attribute{
				Optional: true,
			},
			"region_id": schema.Int64Attribute{
				Optional: true,
			},
			"include_prices": schema.BoolAttribute{
				Description: "Show price",
				Optional:    true,
			},
			"architecture": schema.StringAttribute{
				Description: "An image architecture type: aarch64, `x86_64`\nAvailable values: \"aarch64\", \"x86_64\".",
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
			"hw_firmware_type": schema.StringAttribute{
				Description: "Specifies the type of firmware with which to boot the guest.\nAvailable values: \"bios\", \"uefi\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("bios", "uefi"),
				},
			},
			"hw_machine_type": schema.StringAttribute{
				Description: "A virtual chipset type.\nAvailable values: \"pc\", \"q35\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("pc", "q35"),
				},
			},
			"id": schema.StringAttribute{
				Description: "Image ID",
				Computed:    true,
			},
			"is_baremetal": schema.BoolAttribute{
				Description: "Set to true if the image will be used by bare metal servers. Defaults to false.",
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
			"region": schema.StringAttribute{
				Description: "Region name",
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
				CustomType:  customfield.NewNestedObjectListType[CloudInstanceImageTagsDataSourceModel](ctx),
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

func (d *CloudInstanceImageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudInstanceImageDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
