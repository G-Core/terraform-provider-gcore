// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudInferenceModelDataSourceModel struct {
	ModelID             types.String `tfsdk:"model_id" path:"model_id,required"`
	Category            types.String `tfsdk:"category" json:"category,computed"`
	DefaultFlavorName   types.String `tfsdk:"default_flavor_name" json:"default_flavor_name,computed"`
	Description         types.String `tfsdk:"description" json:"description,computed"`
	Developer           types.String `tfsdk:"developer" json:"developer,computed"`
	DocumentationPage   types.String `tfsdk:"documentation_page" json:"documentation_page,computed"`
	EulaURL             types.String `tfsdk:"eula_url" json:"eula_url,computed"`
	ExampleCurlRequest  types.String `tfsdk:"example_curl_request" json:"example_curl_request,computed"`
	HasEula             types.Bool   `tfsdk:"has_eula" json:"has_eula,computed"`
	ID                  types.String `tfsdk:"id" json:"id,computed"`
	ImageRegistryID     types.String `tfsdk:"image_registry_id" json:"image_registry_id,computed"`
	ImageURL            types.String `tfsdk:"image_url" json:"image_url,computed"`
	InferenceBackend    types.String `tfsdk:"inference_backend" json:"inference_backend,computed"`
	InferenceFrontend   types.String `tfsdk:"inference_frontend" json:"inference_frontend,computed"`
	Name                types.String `tfsdk:"name" json:"name,computed"`
	OpenAICompatibility types.String `tfsdk:"openai_compatibility" json:"openai_compatibility,computed"`
	Port                types.Int64  `tfsdk:"port" json:"port,computed"`
	Version             types.String `tfsdk:"version" json:"version,computed"`
}
