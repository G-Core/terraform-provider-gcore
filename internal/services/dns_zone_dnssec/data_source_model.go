// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_dnssec

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DNSZoneDnssecDataSourceModel struct {
	Name            types.String `tfsdk:"name" path:"name,required"`
	Algorithm       types.String `tfsdk:"algorithm" json:"algorithm,computed"`
	Digest          types.String `tfsdk:"digest" json:"digest,computed"`
	DigestAlgorithm types.String `tfsdk:"digest_algorithm" json:"digest_algorithm,computed"`
	DigestType      types.String `tfsdk:"digest_type" json:"digest_type,computed"`
	Ds              types.String `tfsdk:"ds" json:"ds,computed"`
	Flags           types.Int64  `tfsdk:"flags" json:"flags,computed"`
	KeyTag          types.Int64  `tfsdk:"key_tag" json:"key_tag,computed"`
	KeyType         types.String `tfsdk:"key_type" json:"key_type,computed"`
	PublicKey       types.String `tfsdk:"public_key" json:"public_key,computed"`
	Uuid            types.String `tfsdk:"uuid" json:"uuid,computed"`
}
