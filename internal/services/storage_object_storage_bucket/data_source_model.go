// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage_bucket

import (
	"context"

	"github.com/G-Core/gcore-go/storage"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageObjectStorageBucketDataSourceModel struct {
	Name                                types.String                                                                                           `tfsdk:"name" path:"name,required"`
	StorageID                           types.Int64                                                                                            `tfsdk:"storage_id" path:"storage_id,required"`
	Cors                                customfield.NestedObject[StorageObjectStorageBucketCorsDataSourceModel]                                `tfsdk:"cors" json:"cors,computed"`
	Policy                              customfield.NestedObject[StorageObjectStorageBucketPolicyDataSourceModel]                              `tfsdk:"policy" json:"policy,computed"`
	StorageObjectStorageBucketLifecycle customfield.NestedObject[StorageObjectStorageBucketStorageObjectStorageBucketLifecycleDataSourceModel] `tfsdk:"storage_object_storage_bucket_lifecycle" json:"lifecycle,computed"`
}

func (m *StorageObjectStorageBucketDataSourceModel) toReadParams(_ context.Context) (params storage.ObjectStorageBucketGetParams, diags diag.Diagnostics) {
	params = storage.ObjectStorageBucketGetParams{
		StorageID: m.StorageID.ValueInt64(),
	}

	return
}

type StorageObjectStorageBucketCorsDataSourceModel struct {
	AllowedOrigins customfield.List[types.String] `tfsdk:"allowed_origins" json:"allowed_origins,computed"`
}

type StorageObjectStorageBucketPolicyDataSourceModel struct {
	IsPublic types.Bool `tfsdk:"is_public" json:"is_public,computed"`
}

type StorageObjectStorageBucketStorageObjectStorageBucketLifecycleDataSourceModel struct {
	ExpirationDays types.Int64 `tfsdk:"expiration_days" json:"expiration_days,computed"`
}
