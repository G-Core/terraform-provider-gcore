// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_ssh_key

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudSSHKeysResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudSSHKeysItemsDataSourceModel] `json:"results,computed"`
}

type CloudSSHKeysDataSourceModel struct {
	ProjectID types.Int64                                                    `tfsdk:"project_id" path:"project_id,optional"`
	Name      types.String                                                   `tfsdk:"name" query:"name,optional"`
	OrderBy   types.String                                                   `tfsdk:"order_by" query:"order_by,computed_optional"`
	MaxItems  types.Int64                                                    `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudSSHKeysItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudSSHKeysDataSourceModel) toListParams(_ context.Context) (params cloud.SSHKeyListParams, diags diag.Diagnostics) {
	params = cloud.SSHKeyListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = cloud.SSHKeyListParamsOrderBy(m.OrderBy.ValueString())
	}

	return
}

type CloudSSHKeysItemsDataSourceModel struct {
	ID              types.String      `tfsdk:"id" json:"id,computed"`
	CreatedAt       timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Fingerprint     types.String      `tfsdk:"fingerprint" json:"fingerprint,computed"`
	Name            types.String      `tfsdk:"name" json:"name,computed"`
	ProjectID       types.Int64       `tfsdk:"project_id" json:"project_id,computed"`
	PublicKey       types.String      `tfsdk:"public_key" json:"public_key,computed"`
	SharedInProject types.Bool        `tfsdk:"shared_in_project" json:"shared_in_project,computed"`
	State           types.String      `tfsdk:"state" json:"state,computed"`
}
