// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_secret

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudSecretResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Secrets store sensitive data such as TLS certificates and private keys in encrypted form within a cloud region.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.Int64Attribute{
				Description:   "Project ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"region_id": schema.Int64Attribute{
				Description:   "Region ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "Secret name",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"payload": schema.SingleNestedAttribute{
				Description: "Secret payload.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"certificate_wo": schema.StringAttribute{
						Description: "SSL certificate in PEM format. " +
							"This is a write-only field — it will be sent to the API but never stored in state.",
						Required:  true,
						WriteOnly: true,
					},
					"certificate_chain_wo": schema.StringAttribute{
						Description: "SSL certificate chain of intermediates and root certificates in PEM format. " +
							"This is a write-only field — it will be sent to the API but never stored in state.",
						Required:  true,
						WriteOnly: true,
					},
					"private_key_wo": schema.StringAttribute{
						Description: "SSL private key in PEM format. " +
							"This is a write-only field — it will be sent to the API but never stored in state.",
						Required:  true,
						WriteOnly: true,
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
			},
			"payload_wo_version": schema.Int64Attribute{
				Description: "The version of the payload sensitive params — increment this value to force " +
					"Terraform to re-create the secret with updated payload values.",
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"expiration": schema.StringAttribute{
				Description:   "Datetime when the secret will expire. Defaults to None",
				Optional:      true,
				CustomType:    timetypes.RFC3339Type{},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"algorithm": schema.StringAttribute{
				Description: "Metadata provided by a user or system for informational purposes. Defaults to None",
				Computed:    true,
			},
			"bit_length": schema.Int64Attribute{
				Description: "Metadata provided by a user or system for informational purposes. Value must be greater than zero. Defaults to None",
				Computed:    true,
			},
			"created": schema.StringAttribute{
				Description: "Datetime when the secret was created. The format is 2020-01-01T12:00:00+00:00",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"mode": schema.StringAttribute{
				Description: "Metadata provided by a user or system for informational purposes. Defaults to None",
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
			"content_types": schema.MapAttribute{
				Description: "Describes the content-types that can be used to retrieve the payload. The content-type used with symmetric secrets is application/octet-stream",
				Computed:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
		},
	}
}

func (r *CloudSecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudSecretResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
