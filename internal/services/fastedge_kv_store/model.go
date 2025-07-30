// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_kv_store

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type FastedgeKvStoreModel struct {
	ID       types.Int64                                         `tfsdk:"id" json:"id,computed"`
	Comment  types.String                                        `tfsdk:"comment" json:"comment,optional"`
	Byod     *FastedgeKvStoreByodModel                           `tfsdk:"byod" json:"byod,optional"`
	AppCount types.Int64                                         `tfsdk:"app_count" json:"app_count,computed"`
	Updated  timetypes.RFC3339                                   `tfsdk:"updated" json:"updated,computed" format:"date-time"`
	Stats    customfield.NestedObject[FastedgeKvStoreStatsModel] `tfsdk:"stats" json:"stats,computed"`
}

func (m FastedgeKvStoreModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m FastedgeKvStoreModel) MarshalJSONForUpdate(state FastedgeKvStoreModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type FastedgeKvStoreByodModel struct {
	Prefix types.String `tfsdk:"prefix" json:"prefix,required"`
	URL    types.String `tfsdk:"url" json:"url,required"`
}

type FastedgeKvStoreStatsModel struct {
	CfCount   types.Int64 `tfsdk:"cf_count" json:"cf_count,computed"`
	KvCount   types.Int64 `tfsdk:"kv_count" json:"kv_count,computed"`
	Size      types.Int64 `tfsdk:"size" json:"size,computed"`
	ZsetCount types.Int64 `tfsdk:"zset_count" json:"zset_count,computed"`
}
