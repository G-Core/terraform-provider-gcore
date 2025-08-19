// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_insight_silence

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainInsightSilencesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapDomainInsightSilencesItemsDataSourceModel] `json:"results,computed"`
}

type WaapDomainInsightSilencesDataSourceModel struct {
	DomainID    types.Int64                                                                 `tfsdk:"domain_id" path:"domain_id,required"`
	Author      types.String                                                                `tfsdk:"author" query:"author,optional"`
	Comment     types.String                                                                `tfsdk:"comment" query:"comment,optional"`
	ID          *[]types.String                                                             `tfsdk:"id" query:"id,optional"`
	InsightType *[]types.String                                                             `tfsdk:"insight_type" query:"insight_type,optional"`
	Limit       types.Int64                                                                 `tfsdk:"limit" query:"limit,computed_optional"`
	Ordering    types.String                                                                `tfsdk:"ordering" query:"ordering,computed_optional"`
	MaxItems    types.Int64                                                                 `tfsdk:"max_items"`
	Items       customfield.NestedObjectList[WaapDomainInsightSilencesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapDomainInsightSilencesDataSourceModel) toListParams(_ context.Context) (params waap.DomainInsightSilenceListParams, diags diag.Diagnostics) {
	mID := []string{}
	if m.ID != nil {
		for _, item := range *m.ID {
			mID = append(mID, item.ValueString())
		}
	}
	mInsightType := []string{}
	if m.InsightType != nil {
		for _, item := range *m.InsightType {
			mInsightType = append(mInsightType, item.ValueString())
		}
	}

	params = waap.DomainInsightSilenceListParams{
		ID:          mID,
		InsightType: mInsightType,
	}

	if !m.Author.IsNull() {
		params.Author = param.NewOpt(m.Author.ValueString())
	}
	if !m.Comment.IsNull() {
		params.Comment = param.NewOpt(m.Comment.ValueString())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = waap.DomainInsightSilenceListParamsOrdering(m.Ordering.ValueString())
	}

	return
}

type WaapDomainInsightSilencesItemsDataSourceModel struct {
	ID          types.String                  `tfsdk:"id" json:"id,computed"`
	Author      types.String                  `tfsdk:"author" json:"author,computed"`
	Comment     types.String                  `tfsdk:"comment" json:"comment,computed"`
	ExpireAt    timetypes.RFC3339             `tfsdk:"expire_at" json:"expire_at,computed" format:"date-time"`
	InsightType types.String                  `tfsdk:"insight_type" json:"insight_type,computed"`
	Labels      customfield.Map[types.String] `tfsdk:"labels" json:"labels,computed"`
}
