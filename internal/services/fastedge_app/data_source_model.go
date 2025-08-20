// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type FastedgeAppDataSourceModel struct {
	ID           types.Int64                                                    `tfsdk:"id" path:"id,required"`
	APIType      types.String                                                   `tfsdk:"api_type" json:"api_type,computed"`
	Binary       types.Int64                                                    `tfsdk:"binary" json:"binary,computed"`
	Comment      types.String                                                   `tfsdk:"comment" json:"comment,computed"`
	Debug        types.Bool                                                     `tfsdk:"debug" json:"debug,computed"`
	DebugUntil   timetypes.RFC3339                                              `tfsdk:"debug_until" json:"debug_until,computed" format:"date-time"`
	Log          types.String                                                   `tfsdk:"log" json:"log,computed"`
	Name         types.String                                                   `tfsdk:"name" json:"name,computed"`
	Plan         types.String                                                   `tfsdk:"plan" json:"plan,computed"`
	PlanID       types.Int64                                                    `tfsdk:"plan_id" json:"plan_id,computed"`
	Status       types.Int64                                                    `tfsdk:"status" json:"status,computed"`
	Template     types.Int64                                                    `tfsdk:"template" json:"template,computed"`
	TemplateName types.String                                                   `tfsdk:"template_name" json:"template_name,computed"`
	URL          types.String                                                   `tfsdk:"url" json:"url,computed"`
	Env          customfield.Map[types.String]                                  `tfsdk:"env" json:"env,computed"`
	Networks     customfield.List[types.String]                                 `tfsdk:"networks" json:"networks,computed"`
	RspHeaders   customfield.Map[types.String]                                  `tfsdk:"rsp_headers" json:"rsp_headers,computed"`
	Stores       customfield.Map[types.Int64]                                   `tfsdk:"stores" json:"stores,computed"`
	Secrets      customfield.NestedObjectMap[FastedgeAppSecretsDataSourceModel] `tfsdk:"secrets" json:"secrets,computed"`
}

type FastedgeAppSecretsDataSourceModel struct {
	ID      types.Int64  `tfsdk:"id" json:"id,computed"`
	Comment types.String `tfsdk:"comment" json:"comment,computed"`
	Name    types.String `tfsdk:"name" json:"name,computed"`
}
