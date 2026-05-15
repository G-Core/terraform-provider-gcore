// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage_bucket

import (
	"context"

	"github.com/G-Core/gcore-go/storage"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageObjectStorageBucketsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[StorageObjectStorageBucketsItemsDataSourceModel] `json:"results,computed"`
}

type StorageObjectStorageBucketsDataSourceModel struct {
	StorageID types.Int64                                                                   `tfsdk:"storage_id" path:"storage_id,required"`
	MaxItems  types.Int64                                                                   `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[StorageObjectStorageBucketsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *StorageObjectStorageBucketsDataSourceModel) toListParams(_ context.Context) (params storage.ObjectStorageBucketListParams, diags diag.Diagnostics) {
	params = storage.ObjectStorageBucketListParams{}

	return
}

type StorageObjectStorageBucketsItemsDataSourceModel struct {
	Cors                                customfield.NestedObject[StorageObjectStorageBucketsCorsDataSourceModel]                                `tfsdk:"cors" json:"cors,computed"`
	StorageObjectStorageBucketLifecycle customfield.NestedObject[StorageObjectStorageBucketsStorageObjectStorageBucketLifecycleDataSourceModel] `tfsdk:"storage_object_storage_bucket_lifecycle" json:"lifecycle,computed"`
	Name                                types.String                                                                                            `tfsdk:"name" json:"name,computed"`
	Policy                              customfield.NestedObject[StorageObjectStorageBucketsPolicyDataSourceModel]                              `tfsdk:"policy" json:"policy,computed"`
	StorageID                           types.Int64                                                                                             `tfsdk:"storage_id" json:"storage_id,computed"`
}

type StorageObjectStorageBucketsCorsDataSourceModel struct {
	AllowedOrigins customfield.List[types.String] `tfsdk:"allowed_origins" json:"allowed_origins,computed"`
}

type StorageObjectStorageBucketsStorageObjectStorageBucketLifecycleDataSourceModel struct {
	ExpirationDays types.Int64 `tfsdk:"expiration_days" json:"expiration_days,computed"`
}

type StorageObjectStorageBucketsPolicyDataSourceModel struct {
	IsPublic types.Bool `tfsdk:"is_public" json:"is_public,computed"`
}
