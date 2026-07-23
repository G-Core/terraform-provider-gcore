// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share_access_rule

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudFileShareAccessRulesDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "File share access rules control which IP addresses can mount a file share and their permissions (read-only or read-write).",
		Attributes: map[string]schema.Attribute{
			"file_share_id": schema.StringAttribute{
				Description: "File Share ID",
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
				CustomType:  customfield.NewNestedObjectListType[CloudFileShareAccessRulesItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Access Rule ID",
							Computed:    true,
						},
						"access_level": schema.StringAttribute{
							Description: "Access mode\nAvailable values: \"ro\", \"rw\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("ro", "rw"),
							},
						},
						"access_to": schema.StringAttribute{
							Description: "Source IP or network",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "Access Rule state\nAvailable values: \"active\", \"applying\", \"denying\", \"error\", \"new\", \"queued_to_apply\", \"queued_to_deny\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"active",
									"applying",
									"denying",
									"error",
									"new",
									"queued_to_apply",
									"queued_to_deny",
								),
							},
						},
					},
				},
			},
		},
	}
}

func (d *CloudFileShareAccessRulesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudFileShareAccessRulesDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
