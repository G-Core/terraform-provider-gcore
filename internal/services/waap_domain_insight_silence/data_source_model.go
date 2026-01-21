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

type WaapDomainInsightSilenceDataSourceModel struct {
	ID          types.String                                      `tfsdk:"id" path:"silence_id,computed"`
	SilenceID   types.String                                      `tfsdk:"silence_id" path:"silence_id,optional"`
	DomainID    types.Int64                                       `tfsdk:"domain_id" path:"domain_id,required"`
	Author      types.String                                      `tfsdk:"author" json:"author,computed"`
	Comment     types.String                                      `tfsdk:"comment" json:"comment,computed"`
	ExpireAt    timetypes.RFC3339                                 `tfsdk:"expire_at" json:"expire_at,computed" format:"date-time"`
	InsightType types.String                                      `tfsdk:"insight_type" json:"insight_type,computed"`
	Labels      customfield.Map[types.String]                     `tfsdk:"labels" json:"labels,computed"`
	FindOneBy   *WaapDomainInsightSilenceFindOneByDataSourceModel `tfsdk:"find_one_by"`
}

func (m *WaapDomainInsightSilenceDataSourceModel) toReadParams(_ context.Context) (params waap.DomainInsightSilenceGetParams, diags diag.Diagnostics) {
	params = waap.DomainInsightSilenceGetParams{
		DomainID: m.DomainID.ValueInt64(),
	}

	return
}

func (m *WaapDomainInsightSilenceDataSourceModel) toListParams(_ context.Context) (params waap.DomainInsightSilenceListParams, diags diag.Diagnostics) {
	mFindOneByID := []string{}
	if m.FindOneBy.ID != nil {
		for _, item := range *m.FindOneBy.ID {
			mFindOneByID = append(mFindOneByID, item.ValueString())
		}
	}
	mFindOneByInsightType := []string{}
	if m.FindOneBy.InsightType != nil {
		for _, item := range *m.FindOneBy.InsightType {
			mFindOneByInsightType = append(mFindOneByInsightType, item.ValueString())
		}
	}

	params = waap.DomainInsightSilenceListParams{
		DomainID:    m.DomainID.ValueInt64(),
		ID:          mFindOneByID,
		InsightType: mFindOneByInsightType,
	}

	if !m.FindOneBy.Author.IsNull() {
		params.Author = param.NewOpt(m.FindOneBy.Author.ValueString())
	}
	if !m.FindOneBy.Comment.IsNull() {
		params.Comment = param.NewOpt(m.FindOneBy.Comment.ValueString())
	}
	if !m.FindOneBy.Ordering.IsNull() {
		params.Ordering = waap.DomainInsightSilenceListParamsOrdering(m.FindOneBy.Ordering.ValueString())
	}

	return
}

type WaapDomainInsightSilenceFindOneByDataSourceModel struct {
	ID          *[]types.String `tfsdk:"id" query:"id,optional"`
	Author      types.String    `tfsdk:"author" query:"author,optional"`
	Comment     types.String    `tfsdk:"comment" query:"comment,optional"`
	InsightType *[]types.String `tfsdk:"insight_type" query:"insight_type,optional"`
	Ordering    types.String    `tfsdk:"ordering" query:"ordering,computed_optional"`
}
