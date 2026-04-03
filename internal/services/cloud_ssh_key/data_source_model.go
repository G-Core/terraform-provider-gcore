// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_ssh_key

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudSSHKeyDataSourceModel struct {
	ID              types.String                         `tfsdk:"id" path:"ssh_key_id,computed"`
	SSHKeyID        types.String                         `tfsdk:"ssh_key_id" path:"ssh_key_id,optional"`
	ProjectID       types.Int64                          `tfsdk:"project_id" path:"project_id,optional"`
	CreatedAt       timetypes.RFC3339                    `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Fingerprint     types.String                         `tfsdk:"fingerprint" json:"fingerprint,computed"`
	Name            types.String                         `tfsdk:"name" json:"name,computed"`
	PublicKey       types.String                         `tfsdk:"public_key" json:"public_key,computed"`
	SharedInProject types.Bool                           `tfsdk:"shared_in_project" json:"shared_in_project,computed"`
	State           types.String                         `tfsdk:"state" json:"state,computed"`
	FindOneBy       *CloudSSHKeyFindOneByDataSourceModel `tfsdk:"find_one_by"`
}

func (m *CloudSSHKeyDataSourceModel) toReadParams(_ context.Context) (params cloud.SSHKeyGetParams, diags diag.Diagnostics) {
	params = cloud.SSHKeyGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

func (m *CloudSSHKeyDataSourceModel) toListParams(_ context.Context) (params cloud.SSHKeyListParams, diags diag.Diagnostics) {
	params = cloud.SSHKeyListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.FindOneBy.Limit.IsNull() {
		params.Limit = param.NewOpt(m.FindOneBy.Limit.ValueInt64())
	}
	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}
	if !m.FindOneBy.OrderBy.IsNull() {
		params.OrderBy = cloud.SSHKeyListParamsOrderBy(m.FindOneBy.OrderBy.ValueString())
	}

	return
}

type CloudSSHKeyFindOneByDataSourceModel struct {
	Limit   types.Int64  `tfsdk:"limit" query:"limit,computed_optional"`
	Name    types.String `tfsdk:"name" query:"name,optional"`
	OrderBy types.String `tfsdk:"order_by" query:"order_by,computed_optional"`
}
