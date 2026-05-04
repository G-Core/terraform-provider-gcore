// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_router

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudNetworkRouterModel struct {
	ID                  types.String                                                         `tfsdk:"id" json:"id,computed"`
	ProjectID           types.Int64                                                          `tfsdk:"project_id" path:"project_id,optional"`
	RegionID            types.Int64                                                          `tfsdk:"region_id" path:"region_id,optional"`
	Interfaces          customfield.NestedObjectList[CloudNetworkRouterInterfacesModel]      `tfsdk:"interfaces" json:"interfaces,computed_optional"`
	Name                types.String                                                         `tfsdk:"name" json:"name,required"`
	ExternalGatewayInfo customfield.NestedObject[CloudNetworkRouterExternalGatewayInfoModel] `tfsdk:"external_gateway_info" json:"external_gateway_info,computed_optional"`
	Routes              customfield.NestedObjectList[CloudNetworkRouterRoutesModel]          `tfsdk:"routes" json:"routes,computed_optional"`
	CreatedAt           timetypes.RFC3339                                                    `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                         `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Distributed         types.Bool                                                           `tfsdk:"distributed" json:"distributed,computed"`
	Region              types.String                                                         `tfsdk:"region" json:"region,computed"`
	Status              types.String                                                         `tfsdk:"status" json:"status,computed"`
}

func (m CloudNetworkRouterModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudNetworkRouterModel) MarshalJSONForUpdate(state CloudNetworkRouterModel) (data []byte, err error) {
	// Interfaces are managed via attach/detach API calls, exclude from PATCH
	mCopy := m
	mCopy.Interfaces = state.Interfaces

	// The v1 router PATCH endpoint requires the FULL external_gateway_info
	// block whenever any subfield changes — network_id is server-required
	// even when the user only toggles enable_snat. MarshalForPatch normally
	// emits per-subfield diffs, which would send {"enable_snat":false} alone
	// and the API rejects with `400 ... 'network_id': ['Field required']`.
	// Force full-object emission by zeroing the state's gateway when it
	// differs from the plan, so the encoder sees plan=non-null, state=null
	// and emits every subfield.
	stateCopy := state
	if !m.ExternalGatewayInfo.IsNull() && !m.ExternalGatewayInfo.IsUnknown() &&
		!m.ExternalGatewayInfo.Equal(state.ExternalGatewayInfo) {
		stateCopy.ExternalGatewayInfo = customfield.NestedObject[CloudNetworkRouterExternalGatewayInfoModel]{}
	}

	return apijson.MarshalForPatch(mCopy, stateCopy)
}

// NormalizeExternalGateway reconciles the external_gateway_info field with the
// shape Terraform expects, after unmarshalling an API response. Two
// adjustments:
//
//  1. The SDK's Router.ExternalGatewayInfo is a value type, so a router with no
//     gateway round-trips as an empty object ({"network_id":"","enable_snat":false,...})
//     rather than null. Detect that empty shape and force the field to null so
//     plan == state for gateway-less routers.
//  2. The API GET payload omits the type field. The platform only exposes
//     "manual" as a stable post-create value (the "default" form is a
//     Create-time convenience that resolves to a concrete network_id), so we
//     treat any returned gateway with a non-empty network_id as manual.
func (m *CloudNetworkRouterModel) NormalizeExternalGateway(ctx context.Context) {
	if m.ExternalGatewayInfo.IsNull() || m.ExternalGatewayInfo.IsUnknown() {
		return
	}
	gw, diags := m.ExternalGatewayInfo.Value(ctx)
	if diags.HasError() || gw == nil {
		return
	}
	hasNetworkID := !gw.NetworkID.IsNull() && !gw.NetworkID.IsUnknown() && gw.NetworkID.ValueString() != ""
	if !hasNetworkID {
		// Empty gateway object from the SDK's value-type struct → null in state.
		m.ExternalGatewayInfo = customfield.NullObject[CloudNetworkRouterExternalGatewayInfoModel](ctx)
		return
	}
	if !gw.Type.IsNull() && !gw.Type.IsUnknown() && gw.Type.ValueString() != "" {
		return
	}
	gw.Type = types.StringValue("manual")
	updated, diags := customfield.NewObject(ctx, gw)
	if diags.HasError() {
		return
	}
	m.ExternalGatewayInfo = updated
}

type CloudNetworkRouterInterfacesModel struct {
	SubnetID types.String `tfsdk:"subnet_id" json:"subnet_id,required,no_refresh"`
	Type     types.String `tfsdk:"type" json:"type,computed_optional,no_refresh"`
}

type CloudNetworkRouterExternalGatewayInfoModel struct {
	NetworkID  types.String `tfsdk:"network_id" json:"network_id,computed_optional"`
	EnableSnat types.Bool   `tfsdk:"enable_snat" json:"enable_snat,computed_optional"`
	Type       types.String `tfsdk:"type" json:"type,computed_optional"`
}

type CloudNetworkRouterRoutesModel struct {
	Destination types.String `tfsdk:"destination" json:"destination,required"`
	Nexthop     types.String `tfsdk:"nexthop" json:"nexthop,required"`
}
