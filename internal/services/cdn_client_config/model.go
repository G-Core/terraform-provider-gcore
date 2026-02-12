// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_client_config

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CDNClientConfigModel struct {
	ID                        types.Int64                                           `tfsdk:"id" json:"id,computed"`
	UtilizationLevel          types.Int64                                           `tfsdk:"utilization_level" json:"utilization_level,optional"`
	AutoSuspendEnabled        types.Bool                                            `tfsdk:"auto_suspend_enabled" json:"auto_suspend_enabled,computed"`
	CDNResourcesRulesMaxCount types.Int64                                           `tfsdk:"cdn_resources_rules_max_count" json:"cdn_resources_rules_max_count,computed"`
	Cname                     types.String                                          `tfsdk:"cname" json:"cname,computed"`
	Created                   types.String                                          `tfsdk:"created" json:"created,computed"`
	Updated                   types.String                                          `tfsdk:"updated" json:"updated,computed"`
	UseBalancer               types.Bool                                            `tfsdk:"use_balancer" json:"use_balancer,computed"`
	Service                   customfield.NestedObject[CDNClientConfigServiceModel] `tfsdk:"service" json:"service,computed"`
}

func (m CDNClientConfigModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CDNClientConfigModel) MarshalJSONForUpdate(state CDNClientConfigModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CDNClientConfigServiceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Status  types.String `tfsdk:"status" json:"status,computed"`
	Updated types.String `tfsdk:"updated" json:"updated,computed"`
}
