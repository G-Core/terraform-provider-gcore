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

type StorageSftpDataSourceModel struct {
	ID                  types.Int64                          `tfsdk:"id" path:"storage_id,computed"`
	StorageID           types.Int64                          `tfsdk:"storage_id" path:"storage_id,optional"`
	Address             types.String                         `tfsdk:"address" json:"address,computed"`
	CreatedAt           timetypes.RFC3339                    `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Expires             types.String                         `tfsdk:"expires" json:"expires,computed"`
	FullName            types.String                         `tfsdk:"full_name" json:"full_name,computed"`
	HasCustomConfigFile types.Bool                           `tfsdk:"has_custom_config_file" json:"has_custom_config_file,computed"`
	HasPassword         types.Bool                           `tfsdk:"has_password" json:"has_password,computed"`
	IsHTTPDisabled      types.Bool                           `tfsdk:"is_http_disabled" json:"is_http_disabled,computed"`
	LocationName        types.String                         `tfsdk:"location_name" json:"location_name,computed"`
	Name                types.String                         `tfsdk:"name" json:"name,computed"`
	Password            types.String                         `tfsdk:"password" json:"password,computed"`
	ProvisioningStatus  types.String                         `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	ServerAlias         types.String                         `tfsdk:"server_alias" json:"server_alias,computed"`
	SSHKeyIDs           customfield.List[types.Int64]        `tfsdk:"ssh_key_ids" json:"ssh_key_ids,computed"`
	FindOneBy           *StorageSftpFindOneByDataSourceModel `tfsdk:"find_one_by"`
}

func (m *StorageSftpDataSourceModel) toListParams(_ context.Context) (params storage.SftpStorageListParams, diags diag.Diagnostics) {
	params = storage.SftpStorageListParams{}

	if !m.FindOneBy.ID.IsNull() {
		params.ID = param.NewOpt(m.FindOneBy.ID.ValueString())
	}
	if !m.FindOneBy.LocationName.IsNull() {
		params.LocationName = param.NewOpt(m.FindOneBy.LocationName.ValueString())
	}
	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}
	if !m.FindOneBy.OrderBy.IsNull() {
		params.OrderBy = param.NewOpt(m.FindOneBy.OrderBy.ValueString())
	}
	if !m.FindOneBy.ProvisioningStatus.IsNull() {
		params.ProvisioningStatus = storage.SftpStorageListParamsProvisioningStatus(m.FindOneBy.ProvisioningStatus.ValueString())
	}
	if !m.FindOneBy.ShowDeleted.IsNull() {
		params.ShowDeleted = param.NewOpt(m.FindOneBy.ShowDeleted.ValueBool())
	}

	return
}

type StorageSftpFindOneByDataSourceModel struct {
	ID                 types.String `tfsdk:"id" query:"id,optional"`
	LocationName       types.String `tfsdk:"location_name" query:"location_name,optional"`
	Name               types.String `tfsdk:"name" query:"name,optional"`
	OrderBy            types.String `tfsdk:"order_by" query:"order_by,computed_optional"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status" query:"provisioning_status,optional"`
	ShowDeleted        types.Bool   `tfsdk:"show_deleted" query:"show_deleted,optional"`
}
