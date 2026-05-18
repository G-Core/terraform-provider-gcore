// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/storage"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageObjectStorageDataSourceModel struct {
	ID                 types.Int64                                   `tfsdk:"id" path:"storage_id,computed"`
	StorageID          types.Int64                                   `tfsdk:"storage_id" path:"storage_id,optional"`
	Address            types.String                                  `tfsdk:"address" json:"address,computed"`
	CreatedAt          timetypes.RFC3339                             `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	FullName           types.String                                  `tfsdk:"full_name" json:"full_name,computed"`
	LocationName       types.String                                  `tfsdk:"location_name" json:"location_name,computed"`
	Name               types.String                                  `tfsdk:"name" json:"name,computed"`
	ProvisioningStatus types.String                                  `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	FindOneBy          *StorageObjectStorageFindOneByDataSourceModel `tfsdk:"find_one_by"`
}

func (m *StorageObjectStorageDataSourceModel) toListParams(_ context.Context) (params storage.ObjectStorageListParams, diags diag.Diagnostics) {
	params = storage.ObjectStorageListParams{}

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
		params.ProvisioningStatus = storage.ObjectStorageListParamsProvisioningStatus(m.FindOneBy.ProvisioningStatus.ValueString())
	}
	if !m.FindOneBy.ShowDeleted.IsNull() {
		params.ShowDeleted = param.NewOpt(m.FindOneBy.ShowDeleted.ValueBool())
	}

	return
}

type StorageObjectStorageFindOneByDataSourceModel struct {
	ID                 types.String `tfsdk:"id" query:"id,optional"`
	LocationName       types.String `tfsdk:"location_name" query:"location_name,optional"`
	Name               types.String `tfsdk:"name" query:"name,optional"`
	OrderBy            types.String `tfsdk:"order_by" query:"order_by,computed_optional"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status" query:"provisioning_status,optional"`
	ShowDeleted        types.Bool   `tfsdk:"show_deleted" query:"show_deleted,optional"`
}
