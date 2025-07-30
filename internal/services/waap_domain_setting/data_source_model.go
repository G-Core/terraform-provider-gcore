// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_setting

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainSettingDataSourceModel struct {
	DomainID types.Int64                                                    `tfsdk:"domain_id" path:"domain_id,required"`
	API      customfield.NestedObject[WaapDomainSettingAPIDataSourceModel]  `tfsdk:"api" json:"api,computed"`
	DDOS     customfield.NestedObject[WaapDomainSettingDDOSDataSourceModel] `tfsdk:"ddos" json:"ddos,computed"`
}

type WaapDomainSettingAPIDataSourceModel struct {
	APIURLs customfield.List[types.String] `tfsdk:"api_urls" json:"api_urls,computed"`
	IsAPI   types.Bool                     `tfsdk:"is_api" json:"is_api,computed"`
}

type WaapDomainSettingDDOSDataSourceModel struct {
	BurstThreshold     types.Int64 `tfsdk:"burst_threshold" json:"burst_threshold,computed"`
	GlobalThreshold    types.Int64 `tfsdk:"global_threshold" json:"global_threshold,computed"`
	SubSecondThreshold types.Int64 `tfsdk:"sub_second_threshold" json:"sub_second_threshold,computed"`
}
