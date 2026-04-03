// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_secret

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudSecretsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudSecretsItemsDataSourceModel] `json:"results,computed"`
}

type CloudSecretsDataSourceModel struct {
	ProjectID types.Int64                                                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID  types.Int64                                                    `tfsdk:"region_id" path:"region_id,optional"`
	Limit     types.Int64                                                    `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems  types.Int64                                                    `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudSecretsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudSecretsDataSourceModel) toListParams(_ context.Context) (params cloud.SecretListParams, diags diag.Diagnostics) {
	params = cloud.SecretListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}

	return
}

type CloudSecretsItemsDataSourceModel struct {
	ID           types.String                  `tfsdk:"id" json:"id,computed"`
	Name         types.String                  `tfsdk:"name" json:"name,computed"`
	SecretType   types.String                  `tfsdk:"secret_type" json:"secret_type,computed"`
	Status       types.String                  `tfsdk:"status" json:"status,computed"`
	Algorithm    types.String                  `tfsdk:"algorithm" json:"algorithm,computed"`
	BitLength    types.Int64                   `tfsdk:"bit_length" json:"bit_length,computed"`
	ContentTypes customfield.Map[types.String] `tfsdk:"content_types" json:"content_types,computed"`
	Created      timetypes.RFC3339             `tfsdk:"created" json:"created,computed" format:"date-time"`
	Expiration   timetypes.RFC3339             `tfsdk:"expiration" json:"expiration,computed" format:"date-time"`
	Mode         types.String                  `tfsdk:"mode" json:"mode,computed"`
}
