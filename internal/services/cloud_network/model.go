// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudNetworkModel struct {
	ID                  types.String                   `tfsdk:"id" json:"id,computed"`
	ProjectID           types.Int64                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID            types.Int64                    `tfsdk:"region_id" path:"region_id,optional"`
	Type                types.String                   `tfsdk:"type" json:"type,computed_optional"`
	Name                types.String                   `tfsdk:"name" json:"name,required"`
	Tags                *map[string]types.String       `tfsdk:"tags" json:"tags,optional,no_refresh"`
	CreatedAt           timetypes.RFC3339              `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                   `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Default             types.Bool                     `tfsdk:"default" json:"default,computed"`
	External            types.Bool                     `tfsdk:"external" json:"external,computed"`
	Mtu                 types.Int64                    `tfsdk:"mtu" json:"mtu,computed"`
	PortSecurityEnabled types.Bool                     `tfsdk:"port_security_enabled" json:"port_security_enabled,computed"`
	Region              types.String                   `tfsdk:"region" json:"region,computed"`
	SegmentationID      types.Int64                    `tfsdk:"segmentation_id" json:"segmentation_id,computed"`
	Shared              types.Bool                     `tfsdk:"shared" json:"shared,computed"`
	UpdatedAt           timetypes.RFC3339              `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Subnets             customfield.List[types.String] `tfsdk:"subnets" json:"subnets,computed"`
}

func (m CloudNetworkModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudNetworkModel) MarshalJSONForUpdate(state CloudNetworkModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}
