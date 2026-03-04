// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_secret

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

var _ datasource.DataSourceWithConfigValidators = (*CloudSecretsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
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
				CustomType:  customfield.NewNestedObjectListType[CloudSecretsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Secret uuid",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Secret name",
							Computed:    true,
						},
						"secret_type": schema.StringAttribute{
							Description: "Secret type, base64 encoded. symmetric - Used for storing byte arrays such as keys suitable for symmetric encryption; public - Used for storing the public key of an asymmetric keypair; private - Used for storing the private key of an asymmetric keypair; passphrase - Used for storing plain text passphrases; certificate - Used for storing cryptographic certificates such as X.509 certificates; opaque - Used for backwards compatibility with previous versions of the API\nAvailable values: \"certificate\", \"opaque\", \"passphrase\", \"private\", \"public\", \"symmetric\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"certificate",
									"opaque",
									"passphrase",
									"private",
									"public",
									"symmetric",
								),
							},
						},
						"status": schema.StringAttribute{
							Description: "Status",
							Computed:    true,
						},
						"algorithm": schema.StringAttribute{
							Description: "Metadata provided by a user or system for informational purposes. Defaults to None",
							Computed:    true,
						},
						"bit_length": schema.Int64Attribute{
							Description: "Metadata provided by a user or system for informational purposes. Value must be greater than zero. Defaults to None",
							Computed:    true,
						},
						"content_types": schema.MapAttribute{
							Description: "Describes the content-types that can be used to retrieve the payload. The content-type used with symmetric secrets is application/octet-stream",
							Computed:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"created": schema.StringAttribute{
							Description: "Datetime when the secret was created. The format is 2020-01-01T12:00:00+00:00",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"expiration": schema.StringAttribute{
							Description: "Datetime when the secret will expire. The format is 2020-01-01T12:00:00+00:00. Defaults to None",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"mode": schema.StringAttribute{
							Description: "Metadata provided by a user or system for informational purposes. Defaults to None",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudSecretsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudSecretsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
