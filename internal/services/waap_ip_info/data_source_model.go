// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_ip_info

import (
	"context"

	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapIPInfoDataSourceModel struct {
	IP        types.String                                             `tfsdk:"ip" query:"ip,required"`
	RiskScore types.String                                             `tfsdk:"risk_score" json:"risk_score,computed"`
	Tags      customfield.List[types.String]                           `tfsdk:"tags" json:"tags,computed"`
	Whois     customfield.NestedObject[WaapIPInfoWhoisDataSourceModel] `tfsdk:"whois" json:"whois,computed"`
}

func (m *WaapIPInfoDataSourceModel) toReadParams(_ context.Context) (params waap.IPInfoGetParams, diags diag.Diagnostics) {
	params = waap.IPInfoGetParams{
		IP: m.IP.ValueString(),
	}

	return
}

type WaapIPInfoWhoisDataSourceModel struct {
	AbuseMail      types.String `tfsdk:"abuse_mail" json:"abuse_mail,computed"`
	Cidr           types.Int64  `tfsdk:"cidr" json:"cidr,computed"`
	Country        types.String `tfsdk:"country" json:"country,computed"`
	NetDescription types.String `tfsdk:"net_description" json:"net_description,computed"`
	NetName        types.String `tfsdk:"net_name" json:"net_name,computed"`
	NetRange       types.String `tfsdk:"net_range" json:"net_range,computed"`
	NetType        types.String `tfsdk:"net_type" json:"net_type,computed"`
	OrgID          types.String `tfsdk:"org_id" json:"org_id,computed"`
	OrgName        types.String `tfsdk:"org_name" json:"org_name,computed"`
	OwnerType      types.String `tfsdk:"owner_type" json:"owner_type,computed"`
	Rir            types.String `tfsdk:"rir" json:"rir,computed"`
	State          types.String `tfsdk:"state" json:"state,computed"`
}
