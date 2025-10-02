// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_floating_ip

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudFloatingIPModel struct {
	ID                types.String                   `tfsdk:"id" json:"id,computed"`
	ProjectID         types.Int64                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID          types.Int64                    `tfsdk:"region_id" path:"region_id,optional"`
	FixedIPAddress    types.String                   `tfsdk:"fixed_ip_address" json:"fixed_ip_address,optional"`
	PortID            types.String                   `tfsdk:"port_id" json:"port_id,optional"`
	Tags              *map[string]types.String       `tfsdk:"tags" json:"tags,optional,no_refresh"`
	CreatedAt         timetypes.RFC3339              `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID     types.String                   `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	FloatingIPAddress types.String                   `tfsdk:"floating_ip_address" json:"floating_ip_address,computed"`
	Region            types.String                   `tfsdk:"region" json:"region,computed"`
	RouterID          types.String                   `tfsdk:"router_id" json:"router_id,computed"`
	Status            types.String                   `tfsdk:"status" json:"status,computed"`
	TaskID            types.String                   `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt         timetypes.RFC3339              `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Tasks             customfield.List[types.String] `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
}

func (m CloudFloatingIPModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudFloatingIPModel) MarshalJSONForUpdate(state CloudFloatingIPModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}
