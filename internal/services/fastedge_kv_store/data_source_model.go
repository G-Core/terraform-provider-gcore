// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_kv_store

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type FastedgeKvStoreDataSourceModel struct {
	ID       types.Int64                                                   `tfsdk:"id" path:"id,required"`
	AppCount types.Int64                                                   `tfsdk:"app_count" json:"app_count,computed"`
	Comment  types.String                                                  `tfsdk:"comment" json:"comment,computed"`
	Updated  timetypes.RFC3339                                             `tfsdk:"updated" json:"updated,computed" format:"date-time"`
	Byod     customfield.NestedObject[FastedgeKvStoreByodDataSourceModel]  `tfsdk:"byod" json:"byod,computed"`
	Stats    customfield.NestedObject[FastedgeKvStoreStatsDataSourceModel] `tfsdk:"stats" json:"stats,computed"`
}

type FastedgeKvStoreByodDataSourceModel struct {
	Prefix types.String `tfsdk:"prefix" json:"prefix,computed"`
	URL    types.String `tfsdk:"url" json:"url,computed"`
}

type FastedgeKvStoreStatsDataSourceModel struct {
	CfCount   types.Int64 `tfsdk:"cf_count" json:"cf_count,computed"`
	KvCount   types.Int64 `tfsdk:"kv_count" json:"kv_count,computed"`
	Size      types.Int64 `tfsdk:"size" json:"size,computed"`
	ZsetCount types.Int64 `tfsdk:"zset_count" json:"zset_count,computed"`
}
