// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_ssh_key

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/storage"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageSSHKeyDataSourceModel struct {
	ID        types.Int64                            `tfsdk:"id" path:"ssh_key_id,computed"`
	SSHKeyID  types.Int64                            `tfsdk:"ssh_key_id" path:"ssh_key_id,optional"`
	CreatedAt timetypes.RFC3339                      `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Name      types.String                           `tfsdk:"name" json:"name,computed"`
	PublicKey types.String                           `tfsdk:"public_key" json:"public_key,computed"`
	FindOneBy *StorageSSHKeyFindOneByDataSourceModel `tfsdk:"find_one_by"`
}

func (m *StorageSSHKeyDataSourceModel) toListParams(_ context.Context) (params storage.SSHKeyListParams, diags diag.Diagnostics) {
	params = storage.SSHKeyListParams{}

	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}
	if !m.FindOneBy.OrderBy.IsNull() {
		params.OrderBy = param.NewOpt(m.FindOneBy.OrderBy.ValueString())
	}

	return
}

type StorageSSHKeyFindOneByDataSourceModel struct {
	Name    types.String `tfsdk:"name" query:"name,optional"`
	OrderBy types.String `tfsdk:"order_by" query:"order_by,computed_optional"`
}
