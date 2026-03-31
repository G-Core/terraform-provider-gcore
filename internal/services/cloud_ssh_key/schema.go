// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_ssh_key

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*CloudSSHKeyResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "SSH key pairs provide secure authentication to cloud instances, supporting both generated and imported public keys.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "SSH key ID",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.Int64Attribute{
				Description:   "Project ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
			},
			"name": schema.StringAttribute{
				Description:   "SSH key name",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"public_key": schema.StringAttribute{
				Description:   "The public part of an SSH key is the shareable portion of an SSH key pair. It can be safely sent to servers or services to grant access. It does not contain sensitive information. You must provide your own public key (usually found in a file like `id_ed25519.pub` or `id_rsa.pub`). Generate your SSH keypair locally using `ssh-keygen` before providing it here.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"shared_in_project": schema.BoolAttribute{
				Description: "SSH key is shared with all users in the project",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
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

func (r *CloudSSHKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudSSHKeyResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
