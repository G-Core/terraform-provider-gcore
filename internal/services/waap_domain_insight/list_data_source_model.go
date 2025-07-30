// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_insight

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainInsightsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapDomainInsightsItemsDataSourceModel] `json:"results,computed"`
}

type WaapDomainInsightsDataSourceModel struct {
	DomainID    types.Int64                                                          `tfsdk:"domain_id" path:"domain_id,required"`
	Description types.String                                                         `tfsdk:"description" query:"description,optional"`
	Ordering    types.String                                                         `tfsdk:"ordering" query:"ordering,optional"`
	ID          *[]types.String                                                      `tfsdk:"id" query:"id,optional"`
	InsightType *[]types.String                                                      `tfsdk:"insight_type" query:"insight_type,optional"`
	Status      *[]types.String                                                      `tfsdk:"status" query:"status,optional"`
	Limit       types.Int64                                                          `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems    types.Int64                                                          `tfsdk:"max_items"`
	Items       customfield.NestedObjectList[WaapDomainInsightsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapDomainInsightsDataSourceModel) toListParams(_ context.Context) (params waap.DomainInsightListParams, diags diag.Diagnostics) {
	mID := []string{}
	for _, item := range *m.ID {
		mID = append(mID, item.ValueString())
	}
	mInsightType := []string{}
	for _, item := range *m.InsightType {
		mInsightType = append(mInsightType, item.ValueString())
	}
	mStatus := []waap.WaapInsightStatus{}
	for _, item := range *m.Status {
		mStatus = append(mStatus, waap.WaapInsightStatus(item.ValueString()))
	}

	params = waap.DomainInsightListParams{
		ID:          mID,
		InsightType: mInsightType,
		Status:      mStatus,
	}

	if !m.Description.IsNull() {
		params.Description = param.NewOpt(m.Description.ValueString())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = waap.WaapInsightSortBy(m.Ordering.ValueString())
	}

	return
}

type WaapDomainInsightsItemsDataSourceModel struct {
	ID               types.String                  `tfsdk:"id" json:"id,computed"`
	Description      types.String                  `tfsdk:"description" json:"description,computed"`
	FirstSeen        timetypes.RFC3339             `tfsdk:"first_seen" json:"first_seen,computed" format:"date-time"`
	InsightType      types.String                  `tfsdk:"insight_type" json:"insight_type,computed"`
	Labels           customfield.Map[types.String] `tfsdk:"labels" json:"labels,computed"`
	LastSeen         timetypes.RFC3339             `tfsdk:"last_seen" json:"last_seen,computed" format:"date-time"`
	LastStatusChange timetypes.RFC3339             `tfsdk:"last_status_change" json:"last_status_change,computed" format:"date-time"`
	Recommendation   types.String                  `tfsdk:"recommendation" json:"recommendation,computed"`
	Status           types.String                  `tfsdk:"status" json:"status,computed"`
}
