// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_application_deployment

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudInferenceApplicationDeploymentModel struct {
	DeploymentName          types.String                                                                `tfsdk:"deployment_name" path:"deployment_name,optional"`
	ProjectID               types.Int64                                                                 `tfsdk:"project_id" path:"project_id,optional"`
	ApplicationName         types.String                                                                `tfsdk:"application_name" json:"application_name,required"`
	Name                    types.String                                                                `tfsdk:"name" json:"name,required"`
	Regions                 *[]types.Int64                                                              `tfsdk:"regions" json:"regions,required"`
	ComponentsConfiguration *map[string]CloudInferenceApplicationDeploymentComponentsConfigurationModel `tfsdk:"components_configuration" json:"components_configuration,required"`
	APIKeys                 *[]types.String                                                             `tfsdk:"api_keys" json:"api_keys,optional"`
	Tasks                   customfield.List[types.String]                                              `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
	Status                  customfield.NestedObject[CloudInferenceApplicationDeploymentStatusModel]    `tfsdk:"status" json:"status,computed"`
}

func (m CloudInferenceApplicationDeploymentModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudInferenceApplicationDeploymentModel) MarshalJSONForUpdate(state CloudInferenceApplicationDeploymentModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type CloudInferenceApplicationDeploymentComponentsConfigurationModel struct {
	Exposed            types.Bool                                                                                    `tfsdk:"exposed" json:"exposed,required"`
	Flavor             types.String                                                                                  `tfsdk:"flavor" json:"flavor,required"`
	Scale              *CloudInferenceApplicationDeploymentComponentsConfigurationScaleModel                         `tfsdk:"scale" json:"scale,required"`
	ParameterOverrides *map[string]CloudInferenceApplicationDeploymentComponentsConfigurationParameterOverridesModel `tfsdk:"parameter_overrides" json:"parameter_overrides,optional"`
}

type CloudInferenceApplicationDeploymentComponentsConfigurationScaleModel struct {
	Max types.Int64 `tfsdk:"max" json:"max,required"`
	Min types.Int64 `tfsdk:"min" json:"min,required"`
}

type CloudInferenceApplicationDeploymentComponentsConfigurationParameterOverridesModel struct {
	Value types.String `tfsdk:"value" json:"value,required"`
}

type CloudInferenceApplicationDeploymentStatusModel struct {
	ComponentInferences customfield.NestedObjectMap[CloudInferenceApplicationDeploymentStatusComponentInferencesModel] `tfsdk:"component_inferences" json:"component_inferences,computed"`
	ConsolidatedStatus  types.String                                                                                   `tfsdk:"consolidated_status" json:"consolidated_status,computed"`
	ExposeAddresses     customfield.NestedObjectMap[CloudInferenceApplicationDeploymentStatusExposeAddressesModel]     `tfsdk:"expose_addresses" json:"expose_addresses,computed"`
	Regions             customfield.NestedObjectList[CloudInferenceApplicationDeploymentStatusRegionsModel]            `tfsdk:"regions" json:"regions,computed"`
}

type CloudInferenceApplicationDeploymentStatusComponentInferencesModel struct {
	Flavor types.String `tfsdk:"flavor" json:"flavor,computed"`
	Name   types.String `tfsdk:"name" json:"name,computed"`
}

type CloudInferenceApplicationDeploymentStatusExposeAddressesModel struct {
	Address types.String `tfsdk:"address" json:"address,computed"`
}

type CloudInferenceApplicationDeploymentStatusRegionsModel struct {
	Components customfield.NestedObjectMap[CloudInferenceApplicationDeploymentStatusRegionsComponentsModel] `tfsdk:"components" json:"components,computed"`
	RegionID   types.Int64                                                                                  `tfsdk:"region_id" json:"region_id,computed"`
	Status     types.String                                                                                 `tfsdk:"status" json:"status,computed"`
}

type CloudInferenceApplicationDeploymentStatusRegionsComponentsModel struct {
	Error  types.String `tfsdk:"error" json:"error,computed"`
	Status types.String `tfsdk:"status" json:"status,computed"`
}
