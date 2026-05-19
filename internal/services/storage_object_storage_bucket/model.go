// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage_bucket

import (
	"encoding/json"

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

// MarshalJSONForUpdate produces the PATCH body for the bucket. It wraps
// apijson.MarshalForPatch with Storage v4 sentinel translation: when the user
// transitions a nested config (cors, lifecycle) from set to null, the v4 API
// requires the explicit "clear" sentinel form rather than an absent field.
//
//	cors      → {"allowed_origins": []}   removes CORS configuration
//	lifecycle → {"expiration_days": 0}    removes the lifecycle rule
//
// Without this, MarshalForPatch omits the null field and the API silently
// preserves the existing value, causing state divergence (see GCLOUD2-25203).
func (m StorageObjectStorageBucketModel) MarshalJSONForUpdate(state StorageObjectStorageBucketModel) (data []byte, err error) {
	data, err = apijson.MarshalForPatch(m, state)
	if err != nil {
		return nil, err
	}

	corsCleared := state.Cors != nil && m.Cors == nil
	lifecycleCleared := state.StorageObjectStorageBucketLifecycle != nil && m.StorageObjectStorageBucketLifecycle == nil
	if !corsCleared && !lifecycleCleared {
		return data, nil
	}

	var patch map[string]any
	if err := json.Unmarshal(data, &patch); err != nil {
		return nil, err
	}
	if corsCleared {
		patch["cors"] = map[string]any{"allowed_origins": []string{}}
	}
	if lifecycleCleared {
		patch["lifecycle"] = map[string]any{"expiration_days": 0}
	}
	return json.Marshal(patch)
}

type StorageObjectStorageBucketCorsModel struct {
	AllowedOrigins *[]types.String `tfsdk:"allowed_origins" json:"allowed_origins,computed_optional"`
}

type StorageObjectStorageBucketPolicyModel struct {
	IsPublic types.Bool `tfsdk:"is_public" json:"is_public,computed_optional"`
}

type StorageObjectStorageBucketStorageObjectStorageBucketLifecycleModel struct {
	ExpirationDays types.Int64 `tfsdk:"expiration_days" json:"expiration_days,computed_optional"`
}
