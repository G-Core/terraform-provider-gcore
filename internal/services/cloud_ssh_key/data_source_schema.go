// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_ssh_key

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudSSHKeyDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Required:    true,
			},
			"ssh_key_id": schema.StringAttribute{
				Description: "SSH key ID",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "SSH key creation time",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"fingerprint": schema.StringAttribute{
				Description: "Fingerprint",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "SSH key ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "SSH key name",
				Computed:    true,
			},
			"public_key": schema.StringAttribute{
				Description: "The public part of an SSH key is the shareable portion of an SSH key pair. It can be safely sent to servers or services to grant access. It does not contain sensitive information.",
				Computed:    true,
			},
			"shared_in_project": schema.BoolAttribute{
				Description: "SSH key will be visible to all users in the project",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "SSH key state\nAvailable values: \"ACTIVE\", \"DELETING\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("ACTIVE", "DELETING"),
				},
			},
		},
	}
}

func (d *CloudSSHKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudSSHKeyDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
