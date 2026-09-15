// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_sftp

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageSftpModel struct {
	ID           types.Int64  `tfsdk:"id" json:"id,computed"`
	LocationName types.String `tfsdk:"location_name" json:"location_name,required"`
	Name         types.String `tfsdk:"name" json:"name,required"`
	// no_refresh keeps API responses from writing the password to state. Create
	// uses the JSON name, while Update adds the password directly to the request.
	PasswordWo          types.String                  `tfsdk:"password_wo" json:"password,optional,no_refresh"`
	PasswordWoVersion   types.Int64                   `tfsdk:"password_wo_version"`
	Expires             types.String                  `tfsdk:"expires" json:"expires,optional"`
	HasCustomConfigFile types.Bool                    `tfsdk:"has_custom_config_file" json:"has_custom_config_file,computed_optional"`
	IsHTTPDisabled      types.Bool                    `tfsdk:"is_http_disabled" json:"is_http_disabled,computed_optional"`
	ServerAlias         types.String                  `tfsdk:"server_alias" json:"server_alias,computed_optional"`
	SSHKeyIDs           customfield.List[types.Int64] `tfsdk:"ssh_key_ids" json:"ssh_key_ids,computed_optional"`
	Address             types.String                  `tfsdk:"address" json:"address,computed"`
	CreatedAt           timetypes.RFC3339             `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	FullName            types.String                  `tfsdk:"full_name" json:"full_name,computed"`
	HasPassword         types.Bool                    `tfsdk:"has_password" json:"has_password,computed"`
	ProvisioningStatus  types.String                  `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
}

func (m StorageSftpModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StorageSftpModel) MarshalJSONForUpdate(state StorageSftpModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}
