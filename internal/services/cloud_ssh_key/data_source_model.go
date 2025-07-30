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
	ProjectID       types.Int64       `tfsdk:"project_id" path:"project_id,required"`
	SSHKeyID        types.String      `tfsdk:"ssh_key_id" path:"ssh_key_id,required"`
	CreatedAt       timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Fingerprint     types.String      `tfsdk:"fingerprint" json:"fingerprint,computed"`
	ID              types.String      `tfsdk:"id" json:"id,computed"`
	Name            types.String      `tfsdk:"name" json:"name,computed"`
	PublicKey       types.String      `tfsdk:"public_key" json:"public_key,computed"`
	SharedInProject types.Bool        `tfsdk:"shared_in_project" json:"shared_in_project,computed"`
	State           types.String      `tfsdk:"state" json:"state,computed"`
}

func (m *CloudSSHKeyDataSourceModel) toReadParams(_ context.Context) (params cloud.SSHKeyGetParams, diags diag.Diagnostics) {
	params = cloud.SSHKeyGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}
