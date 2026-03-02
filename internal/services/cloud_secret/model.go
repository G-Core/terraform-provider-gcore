// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_secret

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudSecretModel struct {
	ID               types.String                   `tfsdk:"id" json:"id,computed"`
	ProjectID        types.Int64                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID         types.Int64                    `tfsdk:"region_id" path:"region_id,optional"`
	Name             types.String                   `tfsdk:"name" json:"name,required"`
	Payload          *CloudSecretPayloadModel       `tfsdk:"payload" json:"payload,required,no_refresh"`
	PayloadWoVersion types.Int64                    `tfsdk:"payload_wo_version"`
	Expiration       timetypes.RFC3339              `tfsdk:"expiration" json:"expiration,optional" format:"date-time"`
	Algorithm        types.String                   `tfsdk:"algorithm" json:"algorithm,computed"`
	BitLength        types.Int64                    `tfsdk:"bit_length" json:"bit_length,computed"`
	Created          timetypes.RFC3339              `tfsdk:"created" json:"created,computed" format:"date-time"`
	Mode             types.String                   `tfsdk:"mode" json:"mode,computed"`
	SecretType       types.String                   `tfsdk:"secret_type" json:"secret_type,computed"`
	Status           types.String                   `tfsdk:"status" json:"status,computed"`
	ContentTypes     customfield.Map[types.String]  `tfsdk:"content_types" json:"content_types,computed"`
	Tasks            customfield.List[types.String] `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
}

func (m CloudSecretModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudSecretModel) MarshalJSONForUpdate(state CloudSecretModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type CloudSecretPayloadModel struct {
	Certificate      types.String `tfsdk:"certificate_wo" json:"certificate,required"`
	CertificateChain types.String `tfsdk:"certificate_chain_wo" json:"certificate_chain,required"`
	PrivateKey       types.String `tfsdk:"private_key_wo" json:"private_key,required"`
}
