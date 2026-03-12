// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FastedgeAppModel struct {
	ID            types.Int64                                          `tfsdk:"id" json:"id,computed"`
	Template      types.Int64                                          `tfsdk:"template" json:"template,optional"`
	Binary        types.Int64                                          `tfsdk:"binary" json:"binary,computed_optional"`
	Comment       types.String                                         `tfsdk:"comment" json:"comment,computed_optional"`
	Debug         types.Bool                                           `tfsdk:"debug" json:"debug,computed_optional"`
	Log           types.String                                         `tfsdk:"log" json:"log,computed_optional"`
	Name          types.String                                         `tfsdk:"name" json:"name,computed_optional"`
	Status        types.Int64                                          `tfsdk:"status" json:"status,computed_optional"`
	Env           customfield.Map[types.String]                        `tfsdk:"env" json:"env,computed_optional"`
	RspHeaders    customfield.Map[types.String]                        `tfsdk:"rsp_headers" json:"rsp_headers,computed_optional"`
	Secrets       customfield.NestedObjectMap[FastedgeAppSecretsModel] `tfsdk:"secrets" json:"secrets,computed_optional"`
	Stores        customfield.NestedObjectMap[FastedgeAppStoresModel]  `tfsdk:"stores" json:"stores,computed_optional"`
	APIType       types.String                                         `tfsdk:"api_type" json:"api_type,computed"`
	DebugUntil    timetypes.RFC3339                                    `tfsdk:"debug_until" json:"debug_until,computed" format:"date-time"`
	Plan          types.String                                         `tfsdk:"plan" json:"plan,computed"`
	PlanID        types.Int64                                          `tfsdk:"plan_id" json:"plan_id,computed"`
	TemplateName  types.String                                         `tfsdk:"template_name" json:"template_name,computed"`
	UpgradeableTo types.Int64                                          `tfsdk:"upgradeable_to" json:"upgradeable_to,computed"`
	URL           types.String                                         `tfsdk:"url" json:"url,computed"`
	Networks      customfield.List[types.String]                       `tfsdk:"networks" json:"networks,computed"`
}

func (m FastedgeAppModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m FastedgeAppModel) MarshalJSONForUpdate(state FastedgeAppModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type FastedgeAppSecretsModel struct {
	ID      types.Int64  `tfsdk:"id" json:"id,required"`
	Comment types.String `tfsdk:"comment" json:"comment,computed"`
	Name    types.String `tfsdk:"name" json:"name,computed"`
}

type FastedgeAppStoresModel struct {
	ID      types.Int64  `tfsdk:"id" json:"id,required"`
	Comment types.String `tfsdk:"comment" json:"comment,computed"`
	Name    types.String `tfsdk:"name" json:"name,computed"`
}
