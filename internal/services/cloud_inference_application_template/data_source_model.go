// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_application_template

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudInferenceApplicationTemplateDataSourceModel struct {
	ApplicationName types.String                                                                            `tfsdk:"application_name" path:"application_name,required"`
	CoverURL        types.String                                                                            `tfsdk:"cover_url" json:"cover_url,computed"`
	Description     types.String                                                                            `tfsdk:"description" json:"description,computed"`
	DisplayName     types.String                                                                            `tfsdk:"display_name" json:"display_name,computed"`
	Name            types.String                                                                            `tfsdk:"name" json:"name,computed"`
	Readme          types.String                                                                            `tfsdk:"readme" json:"readme,computed"`
	Tags            customfield.Map[types.String]                                                           `tfsdk:"tags" json:"tags,computed"`
	Components      customfield.NestedObjectMap[CloudInferenceApplicationTemplateComponentsDataSourceModel] `tfsdk:"components" json:"components,computed"`
}

type CloudInferenceApplicationTemplateComponentsDataSourceModel struct {
	Description     types.String                                                                                            `tfsdk:"description" json:"description,computed"`
	DisplayName     types.String                                                                                            `tfsdk:"display_name" json:"display_name,computed"`
	Exposable       types.Bool                                                                                              `tfsdk:"exposable" json:"exposable,computed"`
	LicenseURL      types.String                                                                                            `tfsdk:"license_url" json:"license_url,computed"`
	Parameters      customfield.NestedObjectMap[CloudInferenceApplicationTemplateComponentsParametersDataSourceModel]       `tfsdk:"parameters" json:"parameters,computed"`
	Readme          types.String                                                                                            `tfsdk:"readme" json:"readme,computed"`
	Required        types.Bool                                                                                              `tfsdk:"required" json:"required,computed"`
	SuitableFlavors customfield.NestedObjectList[CloudInferenceApplicationTemplateComponentsSuitableFlavorsDataSourceModel] `tfsdk:"suitable_flavors" json:"suitable_flavors,computed"`
}

type CloudInferenceApplicationTemplateComponentsParametersDataSourceModel struct {
	DefaultValue types.String                   `tfsdk:"default_value" json:"default_value,computed"`
	Description  types.String                   `tfsdk:"description" json:"description,computed"`
	DisplayName  types.String                   `tfsdk:"display_name" json:"display_name,computed"`
	EnumValues   customfield.List[types.String] `tfsdk:"enum_values" json:"enum_values,computed"`
	MaxValue     types.String                   `tfsdk:"max_value" json:"max_value,computed"`
	MinValue     types.String                   `tfsdk:"min_value" json:"min_value,computed"`
	Pattern      types.String                   `tfsdk:"pattern" json:"pattern,computed"`
	Required     types.Bool                     `tfsdk:"required" json:"required,computed"`
	Type         types.String                   `tfsdk:"type" json:"type,computed"`
}

type CloudInferenceApplicationTemplateComponentsSuitableFlavorsDataSourceModel struct {
	Name types.String `tfsdk:"name" json:"name,computed"`
}
