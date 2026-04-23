// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_ssh_key

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/storage"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageSSHKeysResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[StorageSSHKeysItemsDataSourceModel] `json:"results,computed"`
}

type StorageSSHKeysDataSourceModel struct {
	Name     types.String                                                     `tfsdk:"name" query:"name,optional"`
	OrderBy  types.String                                                     `tfsdk:"order_by" query:"order_by,computed_optional"`
	MaxItems types.Int64                                                      `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[StorageSSHKeysItemsDataSourceModel] `tfsdk:"items"`
}

func (m *StorageSSHKeysDataSourceModel) toListParams(_ context.Context) (params storage.SSHKeyListParams, diags diag.Diagnostics) {
	params = storage.SSHKeyListParams{}

	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = param.NewOpt(m.OrderBy.ValueString())
	}

	return
}

type StorageSSHKeysItemsDataSourceModel struct {
	ID        types.Int64       `tfsdk:"id" json:"id,computed"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Name      types.String      `tfsdk:"name" json:"name,computed"`
	PublicKey types.String      `tfsdk:"public_key" json:"public_key,computed"`
}
