// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_application_deployment

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudInferenceApplicationDeploymentDataSourceModel struct {
	DeploymentName          types.String                                                                                           `tfsdk:"deployment_name" path:"deployment_name,required"`
	ProjectID               types.Int64                                                                                            `tfsdk:"project_id" path:"project_id,required"`
	ApplicationName         types.String                                                                                           `tfsdk:"application_name" json:"application_name,computed"`
	Name                    types.String                                                                                           `tfsdk:"name" json:"name,computed"`
	APIKeys                 customfield.List[types.String]                                                                         `tfsdk:"api_keys" json:"api_keys,computed"`
	Regions                 customfield.List[types.Int64]                                                                          `tfsdk:"regions" json:"regions,computed"`
	ComponentsConfiguration customfield.NestedObjectMap[CloudInferenceApplicationDeploymentComponentsConfigurationDataSourceModel] `tfsdk:"components_configuration" json:"components_configuration,computed"`
	Status                  customfield.NestedObject[CloudInferenceApplicationDeploymentStatusDataSourceModel]                     `tfsdk:"status" json:"status,computed"`
}

func (m *CloudInferenceApplicationDeploymentDataSourceModel) toReadParams(_ context.Context) (params cloud.InferenceApplicationDeploymentGetParams, diags diag.Diagnostics) {
	params = cloud.InferenceApplicationDeploymentGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

type CloudInferenceApplicationDeploymentComponentsConfigurationDataSourceModel struct {
	Exposed            types.Bool                                                                                                               `tfsdk:"exposed" json:"exposed,computed"`
	Flavor             types.String                                                                                                             `tfsdk:"flavor" json:"flavor,computed"`
	ParameterOverrides customfield.NestedObjectMap[CloudInferenceApplicationDeploymentComponentsConfigurationParameterOverridesDataSourceModel] `tfsdk:"parameter_overrides" json:"parameter_overrides,computed"`
	Scale              customfield.NestedObject[CloudInferenceApplicationDeploymentComponentsConfigurationScaleDataSourceModel]                 `tfsdk:"scale" json:"scale,computed"`
}

type CloudInferenceApplicationDeploymentComponentsConfigurationParameterOverridesDataSourceModel struct {
	Value types.String `tfsdk:"value" json:"value,computed"`
}

type CloudInferenceApplicationDeploymentComponentsConfigurationScaleDataSourceModel struct {
	Max types.Int64 `tfsdk:"max" json:"max,computed"`
	Min types.Int64 `tfsdk:"min" json:"min,computed"`
}

type CloudInferenceApplicationDeploymentStatusDataSourceModel struct {
	ComponentInferences customfield.NestedObjectMap[CloudInferenceApplicationDeploymentStatusComponentInferencesDataSourceModel] `tfsdk:"component_inferences" json:"component_inferences,computed"`
	ConsolidatedStatus  types.String                                                                                             `tfsdk:"consolidated_status" json:"consolidated_status,computed"`
	ExposeAddresses     customfield.NestedObjectMap[CloudInferenceApplicationDeploymentStatusExposeAddressesDataSourceModel]     `tfsdk:"expose_addresses" json:"expose_addresses,computed"`
	Regions             customfield.NestedObjectList[CloudInferenceApplicationDeploymentStatusRegionsDataSourceModel]            `tfsdk:"regions" json:"regions,computed"`
}

type CloudInferenceApplicationDeploymentStatusComponentInferencesDataSourceModel struct {
	Flavor types.String `tfsdk:"flavor" json:"flavor,computed"`
	Name   types.String `tfsdk:"name" json:"name,computed"`
}

type CloudInferenceApplicationDeploymentStatusExposeAddressesDataSourceModel struct {
	Address types.String `tfsdk:"address" json:"address,computed"`
}

type CloudInferenceApplicationDeploymentStatusRegionsDataSourceModel struct {
	Components customfield.NestedObjectMap[CloudInferenceApplicationDeploymentStatusRegionsComponentsDataSourceModel] `tfsdk:"components" json:"components,computed"`
	RegionID   types.Int64                                                                                            `tfsdk:"region_id" json:"region_id,computed"`
	Status     types.String                                                                                           `tfsdk:"status" json:"status,computed"`
}

type CloudInferenceApplicationDeploymentStatusRegionsComponentsDataSourceModel struct {
	Error  types.String `tfsdk:"error" json:"error,computed"`
	Status types.String `tfsdk:"status" json:"status,computed"`
}
