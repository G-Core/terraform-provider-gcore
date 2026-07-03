// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_origin_shielding

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNOriginShieldingDataSourceModel struct {
	ResourceID   types.Int64 `tfsdk:"resource_id" path:"resource_id,required"`
	ShieldingPop types.Int64 `tfsdk:"shielding_pop" json:"shielding_pop,computed"`
}
