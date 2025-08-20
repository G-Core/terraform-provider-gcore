// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_insight_silence

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type WaapDomainInsightSilenceModel struct {
	ID          types.String             `tfsdk:"id" json:"id,computed"`
	DomainID    types.Int64              `tfsdk:"domain_id" path:"domain_id,required"`
	InsightType types.String             `tfsdk:"insight_type" json:"insight_type,required"`
	Author      types.String             `tfsdk:"author" json:"author,required"`
	Comment     types.String             `tfsdk:"comment" json:"comment,required"`
	Labels      *map[string]types.String `tfsdk:"labels" json:"labels,required"`
	ExpireAt    timetypes.RFC3339        `tfsdk:"expire_at" json:"expire_at,optional" format:"date-time"`
}

func (m WaapDomainInsightSilenceModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m WaapDomainInsightSilenceModel) MarshalJSONForUpdate(state WaapDomainInsightSilenceModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}
