// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_sftp

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*StorageSftpResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "SFTP storages provide file transfer protocol access for securely uploading, downloading, and managing files over SSH.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "Unique identifier for the storage instance",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseNonNullStateForUnknown()},
			},
			"location_name": schema.StringAttribute{
				Description:   "Location code where the storage should be created",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "User-defined name for the storage instance",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"sftp_password": schema.StringAttribute{
				Description:   "SFTP password (8-63 chars). Required when `password_mode` is 'set'.\nMust be omitted when `password_mode` is 'auto' or 'none'.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"password_mode": schema.StringAttribute{
				Description: "Password handling mode for SFTP access:\n'auto': generate a random password (returned in the response)\n'set': use the password provided in `sftp_password`\n'none': no password (SSH-key-only access)\nAvailable values: \"auto\", \"set\", \"none\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"auto",
						"set",
						"none",
					),
				},
			},
			"expires": schema.StringAttribute{
				Description: `Duration when the storage should expire (e.g., "2 years 6 months"). Omit for no expiration.`,
				Optional:    true,
			},
			"has_custom_config_file": schema.BoolAttribute{
				Description: "Whether this storage should use a custom configuration file",
				Computed:    true,
				Optional:    true,
			},
			"is_http_disabled": schema.BoolAttribute{
				Description: "Whether HTTP access should be disabled (HTTPS only)",
				Computed:    true,
				Optional:    true,
			},
			"server_alias": schema.StringAttribute{
				Description: "Custom domain alias for accessing the storage. Omit for no alias.",
				Computed:    true,
				Optional:    true,
			},
			"ssh_key_ids": schema.ListAttribute{
				Description: "SSH key IDs to associate with this storage at creation time. If omitted, no keys are linked.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewListType[types.Int64](ctx),
				ElementType: types.Int64Type,
			},
			"address": schema.StringAttribute{
				Description: "Full hostname/address for accessing the storage endpoint",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "ISO 8601 timestamp when the storage was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"full_name": schema.StringAttribute{
				Description: "Read-only internal full name of the storage, composed as \"{`client_id`}-{name}\".\nUsed by the SFTP backend as the login username. Clients should use this value when connecting\nbut should continue to identify the storage by `name` in their own configuration.",
				Computed:    true,
			},
			"has_password": schema.BoolAttribute{
				Description: "Whether password authentication is configured for this storage",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description:   "SFTP password. Only returned when newly generated or set (create/patch). Omitted in GET/list responses.",
				Computed:      true,
				Sensitive:     true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"provisioning_status": schema.StringAttribute{
				Description: "Lifecycle status of the storage. Use this to check readiness before operations.\nAvailable values: \"creating\", \"active\", \"updating\", \"deleting\", \"deleted\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"creating",
						"active",
						"updating",
						"deleting",
						"deleted",
					),
				},
			},
		},
	}
}

func (r *StorageSftpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *StorageSftpResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
