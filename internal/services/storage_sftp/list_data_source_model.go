// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_sftp

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/storage"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageSftpsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[StorageSftpsItemsDataSourceModel] `json:"results,computed"`
}

type StorageSftpsDataSourceModel struct {
	ID                 types.String                                                   `tfsdk:"id" query:"id,optional"`
	LocationName       types.String                                                   `tfsdk:"location_name" query:"location_name,optional"`
	Name               types.String                                                   `tfsdk:"name" query:"name,optional"`
	ProvisioningStatus types.String                                                   `tfsdk:"provisioning_status" query:"provisioning_status,optional"`
	ShowDeleted        types.Bool                                                     `tfsdk:"show_deleted" query:"show_deleted,optional"`
	OrderBy            types.String                                                   `tfsdk:"order_by" query:"order_by,computed_optional"`
	MaxItems           types.Int64                                                    `tfsdk:"max_items"`
	Items              customfield.NestedObjectList[StorageSftpsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *StorageSftpsDataSourceModel) toListParams(_ context.Context) (params storage.SftpStorageListParams, diags diag.Diagnostics) {
	params = storage.SftpStorageListParams{}

	if !m.ID.IsNull() {
		params.ID = param.NewOpt(m.ID.ValueString())
	}
	if !m.LocationName.IsNull() {
		params.LocationName = param.NewOpt(m.LocationName.ValueString())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = param.NewOpt(m.OrderBy.ValueString())
	}
	if !m.ProvisioningStatus.IsNull() {
		params.ProvisioningStatus = storage.SftpStorageListParamsProvisioningStatus(m.ProvisioningStatus.ValueString())
	}
	if !m.ShowDeleted.IsNull() {
		params.ShowDeleted = param.NewOpt(m.ShowDeleted.ValueBool())
	}

	return
}

type StorageSftpsItemsDataSourceModel struct {
	ID                  types.Int64                   `tfsdk:"id" json:"id,computed"`
	Address             types.String                  `tfsdk:"address" json:"address,computed"`
	CreatedAt           timetypes.RFC3339             `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Expires             types.String                  `tfsdk:"expires" json:"expires,computed"`
	FullName            types.String                  `tfsdk:"full_name" json:"full_name,computed"`
	HasCustomConfigFile types.Bool                    `tfsdk:"has_custom_config_file" json:"has_custom_config_file,computed"`
	HasPassword         types.Bool                    `tfsdk:"has_password" json:"has_password,computed"`
	IsHTTPDisabled      types.Bool                    `tfsdk:"is_http_disabled" json:"is_http_disabled,computed"`
	LocationName        types.String                  `tfsdk:"location_name" json:"location_name,computed"`
	Name                types.String                  `tfsdk:"name" json:"name,computed"`
	ProvisioningStatus  types.String                  `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	ServerAlias         types.String                  `tfsdk:"server_alias" json:"server_alias,computed"`
	SSHKeyIDs           customfield.List[types.Int64] `tfsdk:"ssh_key_ids" json:"ssh_key_ids,computed"`
	Password            types.String                  `tfsdk:"password" json:"password,computed"`
}
