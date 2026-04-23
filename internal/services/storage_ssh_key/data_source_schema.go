// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_ssh_key

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

var _ datasource.DataSourceWithConfigValidators = (*StorageSSHKeyDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "SSH keys enable secure access to SFTP storage by associating public keys with user accounts for authentication.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"ssh_key_id": schema.Int64Attribute{
				Optional: true,
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
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Filter by name (partial match)",
						Optional:    true,
					},
					"order_by": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
				},
			},
		},
	}
}

func (d *StorageSSHKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *StorageSSHKeyDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("ssh_key_id"), path.MatchRoot("find_one_by")),
	}
}
