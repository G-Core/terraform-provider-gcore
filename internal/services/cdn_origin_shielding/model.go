package cdn_origin_shielding

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CDNOriginShieldingModel is the Terraform state for the gcore_cdn_origin_shielding
// resource. Origin shielding is a settings-style resource on an existing CDN
// resource: it has no create/delete endpoints of its own, only get/replace on
// PUT /cdn/resources/{resource_id}/shielding_v2.
type CDNOriginShieldingModel struct {
	ResourceID   types.Int64 `tfsdk:"resource_id"`
	ShieldingPop types.Int64 `tfsdk:"shielding_pop"`
}
