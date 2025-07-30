// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_setting

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type WaapDomainSettingModel struct {
	DomainID types.Int64                 `tfsdk:"domain_id" path:"domain_id,required"`
	API      *WaapDomainSettingAPIModel  `tfsdk:"api" json:"api,optional"`
	DDOS     *WaapDomainSettingDDOSModel `tfsdk:"ddos" json:"ddos,optional"`
}

func (m WaapDomainSettingModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m WaapDomainSettingModel) MarshalJSONForUpdate(state WaapDomainSettingModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type WaapDomainSettingAPIModel struct {
	APIURLs *[]types.String `tfsdk:"api_urls" json:"api_urls,optional"`
	IsAPI   types.Bool      `tfsdk:"is_api" json:"is_api,optional"`
}

type WaapDomainSettingDDOSModel struct {
	BurstThreshold  types.Int64 `tfsdk:"burst_threshold" json:"burst_threshold,optional"`
	GlobalThreshold types.Int64 `tfsdk:"global_threshold" json:"global_threshold,optional"`
}
