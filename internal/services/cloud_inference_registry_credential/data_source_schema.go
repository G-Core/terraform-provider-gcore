// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_registry_credential

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceRegistryCredentialDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Registry credentials store authentication details for private container registries used by inference deployments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Registry credential name.",
				Computed:    true,
			},
			"credential_name": schema.StringAttribute{
				Description: "Registry credential name.",
				Optional:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Registry credential name.",
				Computed:    true,
			},
			"registry_url": schema.StringAttribute{
				Description: "Registry URL.",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Registry username.",
				Computed:    true,
			},
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"limit": schema.Int64Attribute{
						Description: "Optional. Limit the number of returned items",
						Computed:    true,
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.AtMost(1000),
						},
					},
				},
			},
		},
	}
}

func (d *CloudInferenceRegistryCredentialDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudInferenceRegistryCredentialDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("credential_name"), path.MatchRoot("find_one_by")),
	}
}
