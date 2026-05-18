// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageObjectStorageModel struct {
	ID                 types.Int64                                                       `tfsdk:"id" json:"id,computed"`
	LocationName       types.String                                                      `tfsdk:"location_name" json:"location_name,required"`
	Name               types.String                                                      `tfsdk:"name" json:"name,required"`
	Address            types.String                                                      `tfsdk:"address" json:"address,computed"`
	CreatedAt          timetypes.RFC3339                                                 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	FullName           types.String                                                      `tfsdk:"full_name" json:"full_name,computed"`
	ProvisioningStatus types.String                                                      `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	AccessKeys         customfield.NestedObjectList[StorageObjectStorageAccessKeysModel] `tfsdk:"access_keys" json:"access_keys,computed,no_refresh"`
}

func (m StorageObjectStorageModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StorageObjectStorageModel) MarshalJSONForUpdate(state StorageObjectStorageModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type StorageObjectStorageAccessKeysModel struct {
	AccessKey types.String `tfsdk:"access_key" json:"access_key,computed"`
	SecretKey types.String `tfsdk:"secret_key" json:"secret_key,computed"`
}
