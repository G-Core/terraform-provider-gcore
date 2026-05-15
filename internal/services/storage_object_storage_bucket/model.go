// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage_bucket

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageObjectStorageBucketModel struct {
	StorageID                           types.Int64                                                         `tfsdk:"storage_id" path:"storage_id,required"`
	Name                                types.String                                                        `tfsdk:"name" json:"name,required"`
	Cors                                *StorageObjectStorageBucketCorsModel                                `tfsdk:"cors" json:"cors,optional"`
	Policy                              *StorageObjectStorageBucketPolicyModel                              `tfsdk:"policy" json:"policy,optional"`
	StorageObjectStorageBucketLifecycle *StorageObjectStorageBucketStorageObjectStorageBucketLifecycleModel `tfsdk:"storage_object_storage_bucket_lifecycle" json:"lifecycle,optional"`
}

func (m StorageObjectStorageBucketModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StorageObjectStorageBucketModel) MarshalJSONForUpdate(state StorageObjectStorageBucketModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type StorageObjectStorageBucketCorsModel struct {
	AllowedOrigins *[]types.String `tfsdk:"allowed_origins" json:"allowed_origins,optional"`
}

type StorageObjectStorageBucketPolicyModel struct {
	IsPublic types.Bool `tfsdk:"is_public" json:"is_public,optional"`
}

type StorageObjectStorageBucketStorageObjectStorageBucketLifecycleModel struct {
	ExpirationDays types.Int64 `tfsdk:"expiration_days" json:"expiration_days,optional"`
}
