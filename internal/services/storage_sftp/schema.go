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
			"password_wo": schema.StringAttribute{
				Description: "SFTP password (8-63 chars). This is a write-only field — it is sent to the API " +
					"but never stored in state or plan artifacts. Set it together with " +
					"`password_wo_version`. Change the version to send a new password; the storage is " +
					"updated in place. Editing `password_wo` alone produces no plan. Remove both " +
					"arguments to clear password authentication in place. Omit both to create a " +
					"storage with SSH-key-only access.",
				Optional: true,
				// Sensitive does not redact source lines shown with diagnostics. Never
				// attach a validator or diagnostic to this attribute; the length check
				// lives in passwordPairValidator, which reports without a path
				Sensitive: true,
				WriteOnly: true,
			},
			"password_wo_version": schema.Int64Attribute{
				Description: "The version of `password_wo`. Terraform never sees the password value, so this " +
					"number is what tells it the password changed: change it (conventionally increment) " +
					"to re-send `password_wo` to the API in place — a rotation, or the first password on " +
					"a storage that has none or one Terraform did not set. Editing `password_wo` alone is " +
					"not detected. Remove it together with `password_wo` to clear password authentication " +
					"in place. Never causes a replacement.",
				Optional: true,
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
			// PlanHasPassword sets this only when the request decides the value.
			// Do not copy stale state when the version is still unknown.
			"has_password": schema.BoolAttribute{
				Description: "Whether password authentication is configured for this storage. Refreshed on " +
					"every read; if it becomes false while `password_wo` is configured, the next apply " +
					"re-applies the password.",
				Computed: true,
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
	return []resource.ConfigValidator{
		passwordPairValidator{},
	}
}
