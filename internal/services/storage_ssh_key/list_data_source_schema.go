// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_ssh_key

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*StorageSSHKeysDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "SSH keys enable secure access to SFTP storage by associating public keys with user accounts for authentication.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filter by name (partial match)",
				Optional:    true,
			},
			"order_by": schema.StringAttribute{
				Computed: true,
				Optional: true,
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
				CustomType:  customfield.NewNestedObjectListType[StorageSSHKeysItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Unique identifier for the SSH key",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "ISO 8601 timestamp when the SSH key was created",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"name": schema.StringAttribute{
							Description: "User-defined name for the SSH key",
							Computed:    true,
						},
						"public_key": schema.StringAttribute{
							Description: "The SSH public key content",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *StorageSSHKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *StorageSSHKeysDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
