// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_binary

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FastedgeBinaryModel struct {
	ID       types.Int64  `tfsdk:"id" json:"id,computed"`
	Filename types.String `tfsdk:"filename"`
	Checksum types.String `tfsdk:"checksum" json:"checksum,computed"`
}
