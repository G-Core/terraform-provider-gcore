// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_model

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudInferenceModelsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudInferenceModelsItemsDataSourceModel] `json:"results,computed"`
}

type CloudInferenceModelsDataSourceModel struct {
	Limit    types.Int64                                                            `tfsdk:"limit" query:"limit,computed_optional"`
	OrderBy  types.String                                                           `tfsdk:"order_by" query:"order_by,computed_optional"`
	MaxItems types.Int64                                                            `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[CloudInferenceModelsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudInferenceModelsDataSourceModel) toListParams(_ context.Context) (params cloud.InferenceModelListParams, diags diag.Diagnostics) {
	params = cloud.InferenceModelListParams{}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = cloud.InferenceModelListParamsOrderBy(m.OrderBy.ValueString())
	}

	return
}

type CloudInferenceModelsItemsDataSourceModel struct {
	ID                  types.String `tfsdk:"id" json:"id,computed"`
	Category            types.String `tfsdk:"category" json:"category,computed"`
	DefaultFlavorName   types.String `tfsdk:"default_flavor_name" json:"default_flavor_name,computed"`
	Description         types.String `tfsdk:"description" json:"description,computed"`
	Developer           types.String `tfsdk:"developer" json:"developer,computed"`
	DocumentationPage   types.String `tfsdk:"documentation_page" json:"documentation_page,computed"`
	EulaURL             types.String `tfsdk:"eula_url" json:"eula_url,computed"`
	ExampleCurlRequest  types.String `tfsdk:"example_curl_request" json:"example_curl_request,computed"`
	HasEula             types.Bool   `tfsdk:"has_eula" json:"has_eula,computed"`
	ImageRegistryID     types.String `tfsdk:"image_registry_id" json:"image_registry_id,computed"`
	ImageURL            types.String `tfsdk:"image_url" json:"image_url,computed"`
	InferenceBackend    types.String `tfsdk:"inference_backend" json:"inference_backend,computed"`
	InferenceFrontend   types.String `tfsdk:"inference_frontend" json:"inference_frontend,computed"`
	ModelID             types.String `tfsdk:"model_id" json:"model_id,computed"`
	Name                types.String `tfsdk:"name" json:"name,computed"`
	OpenAICompatibility types.String `tfsdk:"openai_compatibility" json:"openai_compatibility,computed"`
	Port                types.Int64  `tfsdk:"port" json:"port,computed"`
	Version             types.String `tfsdk:"version" json:"version,computed"`
}
