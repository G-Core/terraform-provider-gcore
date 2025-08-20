// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_binary

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FastedgeBinaryDataSourceModel struct {
	ID         types.Int64  `tfsdk:"id" path:"id,required"`
	APIType    types.String `tfsdk:"api_type" json:"api_type,computed"`
	Checksum   types.String `tfsdk:"checksum" json:"checksum,computed"`
	Source     types.Int64  `tfsdk:"source" json:"source,computed"`
	Status     types.Int64  `tfsdk:"status" json:"status,computed"`
	UnrefSince types.String `tfsdk:"unref_since" json:"unref_since,computed"`
}
