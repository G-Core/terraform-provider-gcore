// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/shared/constant"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudInstanceResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudInstanceResource)(nil)
var _ resource.ResourceWithImportState = (*CloudInstanceResource)(nil)

func NewResource() resource.Resource {
	return &CloudInstanceResource{}
}

// CloudInstanceResource defines the resource implementation.
type CloudInstanceResource struct {
	client *gcore.Client
}

func (r *CloudInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance"
}

func (r *CloudInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*gcore.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"unexpected resource configure type",
			fmt.Sprintf("Expected *gcore.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *CloudInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudInstanceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.InstanceNewParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo"), &data.Password)...)

	dataBytes, err := data.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}

	instance, err := r.client.Cloud.Instances.NewAndPoll(
		ctx,
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	err = apijson.UnmarshalComputed([]byte(instance.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	// Extract fields from the API response that aren't mapped via json tags
	// API returns: {"flavor": {"flavor_id": "..."}, "project_id": N, "region_id": N, ...}
	var rawResponse map[string]interface{}
	if err := json.Unmarshal([]byte(instance.RawJSON()), &rawResponse); err == nil {
		if flavorObj, ok := rawResponse["flavor"].(map[string]interface{}); ok {
			if flavorID, ok := flavorObj["flavor_id"].(string); ok {
				data.Flavor = types.StringValue(flavorID)
			}
		}
		// Note: volumes are user-provided (existing volumes only), no extraction needed
	}

	// Extract interface port_id and ip_address from the interfaces list
	if data.Interfaces != nil && len(*data.Interfaces) > 0 {
		listParams := cloud.InstanceInterfaceListParams{}
		if !data.ProjectID.IsNull() {
			listParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			listParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		interfaces, err := r.client.Cloud.Instances.Interfaces.List(
			ctx,
			data.ID.ValueString(),
			listParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err == nil && interfaces != nil {
			mergeInterfaceComputedFields(data.Interfaces, interfaces.Results, false)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudInstanceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudInstanceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := data.ID.ValueString()

	// Track if we've handled updates via specialized endpoints that require a state refresh
	stateHasChanged := false

	// Preserve the user's desired vm_state before any operations that might overwrite data
	// This is needed because UnmarshalComputed will overwrite data.VmState with API response
	// which may contain transient states like "resized" during operations
	plannedVmState := data.VmState

	// Handle flavor changes using the /changeflavor endpoint
	// The changeflavor API is eventually consistent - after task completion,
	// the GET may return stale data. We poll until flavor matches expected value.
	if !data.Flavor.Equal(state.Flavor) {
		resizeParams := cloud.InstanceResizeParams{
			FlavorID: data.Flavor.ValueString(),
		}

		if !data.ProjectID.IsNull() {
			resizeParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}

		if !data.RegionID.IsNull() {
			resizeParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		// Call Resize (changeflavor endpoint) - NOT ResizeAndPoll
		// ResizeAndPoll immediately calls GET after task completion, which returns stale data
		taskList, err := r.client.Cloud.Instances.Resize(
			ctx,
			instanceID,
			resizeParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to resize instance (change flavor)", err.Error())
			return
		}

		if len(taskList.Tasks) == 0 {
			resp.Diagnostics.AddError("resize returned no tasks", "expected at least one task from changeflavor")
			return
		}

		// Poll task until completion
		taskID := taskList.Tasks[0]
		_, err = r.client.Cloud.Tasks.Poll(ctx, taskID, option.WithMiddleware(logging.Middleware(ctx)))
		if err != nil {
			resp.Diagnostics.AddError("failed to poll resize task", err.Error())
			return
		}

		// Poll instance until status == ACTIVE and flavor matches expected value
		// This handles eventual consistency where GET returns stale data immediately after task completion
		expectedFlavor := data.Flavor.ValueString()
		getParams := cloud.InstanceGetParams{}
		if !data.ProjectID.IsNull() {
			getParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			getParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		// Poll with timeout for flavor propagation
		maxAttempts := 60 // 60 attempts * 2 seconds = 2 minutes max
		for attempt := 0; attempt < maxAttempts; attempt++ {
			instance, err := r.client.Cloud.Instances.Get(
				ctx,
				instanceID,
				getParams,
				option.WithMiddleware(logging.Middleware(ctx)),
			)
			if err != nil {
				resp.Diagnostics.AddError("failed to get instance after resize", err.Error())
				return
			}

			// Check if instance is ACTIVE and flavor matches
			if instance.Status == cloud.InstanceStatusActive &&
				(instance.Flavor.FlavorID == expectedFlavor || instance.Flavor.FlavorName == expectedFlavor) {
				// Flavor has propagated, update data with the response
				err = apijson.UnmarshalComputed([]byte(instance.RawJSON()), &data)
				if err != nil {
					resp.Diagnostics.AddError("failed to deserialize instance after resize", err.Error())
					return
				}
				tflog.Info(ctx, "Flavor change completed", map[string]interface{}{
					"instance_id":     instanceID,
					"expected_flavor": expectedFlavor,
					"actual_flavor":   instance.Flavor.FlavorID,
					"attempts":        attempt + 1,
				})
				break
			}

			// If we've exhausted attempts, return error
			if attempt == maxAttempts-1 {
				resp.Diagnostics.AddError(
					"timeout waiting for flavor change",
					fmt.Sprintf("instance flavor did not update to %s within timeout (current: %s, status: %s)",
						expectedFlavor, instance.Flavor.FlavorName, instance.Status),
				)
				return
			}

			// Wait before next poll
			time.Sleep(2 * time.Second)
		}

		stateHasChanged = true
	}

	// Note: Volume resizing is handled via gcore_cloud_volume resource, not through instance
	// Volume attach/detach is handled below after vm_state changes

	// Handle vm_state changes using the /action endpoint
	// Use plannedVmState (the user's desired value) instead of data.VmState
	// because data may have been overwritten with API response containing transient states
	if !plannedVmState.IsNull() && !plannedVmState.IsUnknown() && !plannedVmState.Equal(state.VmState) {
		desiredState := plannedVmState.ValueString()

		actionParams := cloud.InstanceActionParams{}

		if !data.ProjectID.IsNull() {
			actionParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}

		if !data.RegionID.IsNull() {
			actionParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		switch desiredState {
		case "active":
			actionParams.OfStartActionInstanceSerializer = &cloud.InstanceActionParamsBodyStartActionInstanceSerializer{
				Action: constant.ValueOf[constant.Start](),
			}
		case "stopped":
			actionParams.OfBasicActionInstanceSerializer = &cloud.InstanceActionParamsBodyBasicActionInstanceSerializer{
				Action: "stop",
			}
		default:
			resp.Diagnostics.AddError(
				"invalid vm_state",
				fmt.Sprintf("vm_state must be 'active' or 'stopped', got: %s", desiredState),
			)
			return
		}

		instance, err := r.client.Cloud.Instances.ActionAndPoll(
			ctx,
			instanceID,
			actionParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to change vm_state", err.Error())
			return
		}

		// Update state with new instance data after action
		err = apijson.UnmarshalComputed([]byte(instance.RawJSON()), &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize instance after vm_state change", err.Error())
			return
		}
		stateHasChanged = true
	}

	// Handle volume changes using attach/detach endpoints
	if data.Volumes != nil && state.Volumes != nil {
		// Build maps of volume_id -> volume for comparison
		stateVols := make(map[string]*CloudInstanceVolumesModel)
		planVols := make(map[string]*CloudInstanceVolumesModel)

		for _, vol := range *state.Volumes {
			if vol != nil && !vol.VolumeID.IsNull() {
				stateVols[vol.VolumeID.ValueString()] = vol
			}
		}

		for _, vol := range *data.Volumes {
			if vol != nil && !vol.VolumeID.IsNull() {
				planVols[vol.VolumeID.ValueString()] = vol
			}
		}

		// Detach removed volumes (volumes in state but not in plan)
		for volID := range stateVols {
			if _, exists := planVols[volID]; !exists {
				detachParams := cloud.VolumeDetachFromInstanceParams{
					InstanceID: instanceID,
				}
				if !data.ProjectID.IsNull() {
					detachParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
				}
				if !data.RegionID.IsNull() {
					detachParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
				}

				err := r.client.Cloud.Volumes.DetachFromInstanceAndPoll(
					ctx,
					volID,
					detachParams,
					option.WithMiddleware(logging.Middleware(ctx)),
				)
				if err != nil {
					resp.Diagnostics.AddError(
						"failed to detach volume",
						fmt.Sprintf("Volume %s: %s", volID, err.Error()),
					)
					return
				}
				stateHasChanged = true
			}
		}

		// Attach new volumes (volumes in plan but not in state)
		for volID, vol := range planVols {
			if _, exists := stateVols[volID]; !exists {
				attachParams := cloud.VolumeAttachToInstanceParams{
					InstanceID: instanceID,
				}
				if !data.ProjectID.IsNull() {
					attachParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
				}
				if !data.RegionID.IsNull() {
					attachParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
				}
				// Include attachment_tag if specified
				if vol != nil && !vol.AttachmentTag.IsNull() && !vol.AttachmentTag.IsUnknown() {
					attachParams.AttachmentTag = param.NewOpt(vol.AttachmentTag.ValueString())
				}

				err := r.client.Cloud.Volumes.AttachToInstanceAndPoll(
					ctx,
					volID,
					attachParams,
					option.WithMiddleware(logging.Middleware(ctx)),
				)
				if err != nil {
					resp.Diagnostics.AddError(
						"failed to attach volume",
						fmt.Sprintf("Volume %s: %s", volID, err.Error()),
					)
					return
				}
				stateHasChanged = true
			}
		}
	}

	// Handle interface changes using attach/detach endpoints
	if data.Interfaces != nil && state.Interfaces != nil {
		planCount := len(*data.Interfaces)
		stateCount := len(*state.Interfaces)

		// Detach removed interfaces (from the end)
		if stateCount > planCount {
			for i := stateCount - 1; i >= planCount; i-- {
				stateIface := (*state.Interfaces)[i]
				if stateIface.PortID.IsNull() || stateIface.IPAddress.IsNull() {
					resp.Diagnostics.AddError(
						"cannot detach interface",
						fmt.Sprintf("Interface at index %d is missing port_id or ip_address required for detach", i),
					)
					return
				}

				detachParams := cloud.InstanceInterfaceDetachParams{
					IPAddress: stateIface.IPAddress.ValueString(),
					PortID:    stateIface.PortID.ValueString(),
				}
				if !data.ProjectID.IsNull() {
					detachParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
				}
				if !data.RegionID.IsNull() {
					detachParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
				}

				_, err := r.client.Cloud.Instances.Interfaces.DetachAndPoll(
					ctx,
					instanceID,
					detachParams,
					option.WithMiddleware(logging.Middleware(ctx)),
				)
				if err != nil {
					resp.Diagnostics.AddError(
						"failed to detach interface",
						fmt.Sprintf("Interface at index %d: %s", i, err.Error()),
					)
					return
				}
				stateHasChanged = true
			}
		}

		// Attach new interfaces
		if planCount > stateCount {
			for i := stateCount; i < planCount; i++ {
				planIface := (*data.Interfaces)[i]
				ifaceType := planIface.Type.ValueString()

				attachParams := cloud.InstanceInterfaceAttachParams{}
				if !data.ProjectID.IsNull() {
					attachParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
				}
				if !data.RegionID.IsNull() {
					attachParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
				}

				switch ifaceType {
				case "external":
					attachParams.OfNewInterfaceExternalExtendSchemaWithDDOS = &cloud.InstanceInterfaceAttachParamsBodyNewInterfaceExternalExtendSchemaWithDDOS{
						Type: param.NewOpt("external"),
					}
					if !planIface.IPFamily.IsNull() {
						attachParams.OfNewInterfaceExternalExtendSchemaWithDDOS.IPFamily = planIface.IPFamily.ValueString()
					}
				case "subnet":
					if planIface.SubnetID.IsNull() {
						resp.Diagnostics.AddError(
							"missing subnet_id",
							fmt.Sprintf("Interface at index %d with type 'subnet' requires subnet_id", i),
						)
						return
					}
					attachParams.OfNewInterfaceSpecificSubnetSchema = &cloud.InstanceInterfaceAttachParamsBodyNewInterfaceSpecificSubnetSchema{
						Type:     param.NewOpt("subnet"),
						SubnetID: planIface.SubnetID.ValueString(),
					}
				case "any_subnet":
					if planIface.NetworkID.IsNull() {
						resp.Diagnostics.AddError(
							"missing network_id",
							fmt.Sprintf("Interface at index %d with type 'any_subnet' requires network_id", i),
						)
						return
					}
					attachParams.OfNewInterfaceAnySubnetSchema = &cloud.InstanceInterfaceAttachParamsBodyNewInterfaceAnySubnetSchema{
						Type:      param.NewOpt("any_subnet"),
						NetworkID: planIface.NetworkID.ValueString(),
					}
					if !planIface.IPFamily.IsNull() {
						attachParams.OfNewInterfaceAnySubnetSchema.IPFamily = planIface.IPFamily.ValueString()
					}
				case "reserved_fixed_ip":
					if planIface.PortID.IsNull() {
						resp.Diagnostics.AddError(
							"missing port_id",
							fmt.Sprintf("Interface at index %d with type 'reserved_fixed_ip' requires port_id", i),
						)
						return
					}
					attachParams.OfNewInterfaceReservedFixedIPSchema = &cloud.InstanceInterfaceAttachParamsBodyNewInterfaceReservedFixedIPSchema{
						Type:   param.NewOpt("reserved_fixed_ip"),
						PortID: planIface.PortID.ValueString(),
					}
				default:
					resp.Diagnostics.AddError(
						"invalid interface type",
						fmt.Sprintf("Interface at index %d has invalid type: %s", i, ifaceType),
					)
					return
				}

				_, err := r.client.Cloud.Instances.Interfaces.AttachAndPoll(
					ctx,
					instanceID,
					attachParams,
					option.WithMiddleware(logging.Middleware(ctx)),
				)
				if err != nil {
					resp.Diagnostics.AddError(
						"failed to attach interface",
						fmt.Sprintf("Interface at index %d: %s", i, err.Error()),
					)
					return
				}
				stateHasChanged = true
			}
		}
	}

	// Build a set of ports that already have floating IPs attached
	// This is needed to avoid creating duplicate FIPs when source="new" is used
	// (GCLOUD2-21138: FIP ID is not stored for source="new", so we must check the API)
	portsWithFIP := make(map[string]bool)
	if data.Interfaces != nil {
		listParams := cloud.InstanceInterfaceListParams{}
		if !data.ProjectID.IsNull() {
			listParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			listParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		interfacesResp, err := r.client.Cloud.Instances.Interfaces.List(
			ctx,
			instanceID,
			listParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err == nil && interfacesResp != nil {
			for _, iface := range interfacesResp.Results {
				if len(iface.FloatingipDetails) > 0 {
					portsWithFIP[iface.PortID] = true
				}
			}
		}
	}

	// Handle floating IP changes on existing interfaces using FloatingIPs.Assign/Unassign
	if data.Interfaces != nil && state.Interfaces != nil {
		minCount := len(*data.Interfaces)
		if len(*state.Interfaces) < minCount {
			minCount = len(*state.Interfaces)
		}

		for i := 0; i < minCount; i++ {
			planIface := (*data.Interfaces)[i]
			stateIface := (*state.Interfaces)[i]

			// Get floating IP IDs from state and plan
			var stateFloatingID, planFloatingID string
			var planSource string
			if stateIface.FloatingIP != nil && !stateIface.FloatingIP.ExistingFloatingID.IsNull() {
				stateFloatingID = stateIface.FloatingIP.ExistingFloatingID.ValueString()
			}
			if planIface.FloatingIP != nil {
				if !planIface.FloatingIP.ExistingFloatingID.IsNull() {
					planFloatingID = planIface.FloatingIP.ExistingFloatingID.ValueString()
				}
				if !planIface.FloatingIP.Source.IsNull() {
					planSource = planIface.FloatingIP.Source.ValueString()
				}
			}

			// Get port_id from state (needed for Assign)
			portID := ""
			if !stateIface.PortID.IsNull() {
				portID = stateIface.PortID.ValueString()
			}

			// Check if we need to create a NEW floating IP
			// This happens when source="new", there's no existing FIP in state, AND the port
			// doesn't already have a FIP attached (GCLOUD2-21138: prevents duplicate FIP creation on updates)
			if planSource == "new" && stateFloatingID == "" && portID != "" {
				// Skip if port already has a FIP (GCLOUD2-21138: prevents duplicate FIP on update)
				if portsWithFIP[portID] {
					tflog.Debug(ctx, "Skipping FIP creation - port already has FIP attached", map[string]interface{}{
						"port_id":   portID,
						"interface": i,
					})
					continue
				}
				// Create a new floating IP with port_id to assign during creation
				fipParams := cloud.FloatingIPNewParams{
					PortID: param.NewOpt(portID),
				}
				if !data.ProjectID.IsNull() {
					fipParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
				}
				if !data.RegionID.IsNull() {
					fipParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
				}

				newFIP, err := r.client.Cloud.FloatingIPs.NewAndPoll(
					ctx,
					fipParams,
					option.WithMiddleware(logging.Middleware(ctx)),
				)
				if err != nil {
					resp.Diagnostics.AddError(
						"failed to create floating IP",
						fmt.Sprintf("Interface at index %d: %s", i, err.Error()),
					)
					return
				}

				// FIP created and assigned to port successfully.
				// Keep source="new" to match user config - don't change to "existing"
				// Don't populate existing_floating_id - it's only for user-specified existing FIPs
				// The FIP is assigned to the port and will be visible via the interface's IP address
				tflog.Info(ctx, "Created and assigned new floating IP", map[string]interface{}{
					"floating_ip_id": newFIP.ID,
					"port_id":        portID,
					"interface":      i,
				})

				stateHasChanged = true
				continue
			}

			// Handle removal of FIP that was created with source="new"
			// In this case, stateFloatingID is empty (we don't store FIP ID for source="new")
			// but we need to look up the FIP from the port and unassign it
			var stateSource string
			if stateIface.FloatingIP != nil && !stateIface.FloatingIP.Source.IsNull() {
				stateSource = stateIface.FloatingIP.Source.ValueString()
			}
			if stateSource == "new" && planIface.FloatingIP == nil && portID != "" {
				// User had source="new" and now removed the floating_ip block
				// We need to find the FIP attached to this port and unassign it
				tflog.Info(ctx, "Removing floating IP that was created with source=new", map[string]interface{}{
					"port_id":   portID,
					"interface": i,
				})

				// List floating IPs to find the one attached to this port
				listParams := cloud.FloatingIPListParams{}
				if !data.ProjectID.IsNull() {
					listParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
				}
				if !data.RegionID.IsNull() {
					listParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
				}

				// Find the FIP attached to our port using auto-paging
				var fipToUnassign string
				pager := r.client.Cloud.FloatingIPs.ListAutoPaging(ctx, listParams, option.WithMiddleware(logging.Middleware(ctx)))
				for pager.Next() {
					fip := pager.Current()
					if fip.PortID == portID {
						fipToUnassign = fip.ID
						break
					}
				}
				if err := pager.Err(); err != nil {
					resp.Diagnostics.AddError(
						"failed to list floating IPs",
						fmt.Sprintf("Interface at index %d: %s", i, err.Error()),
					)
					return
				}

				if fipToUnassign != "" {
					updateParams := cloud.FloatingIPUpdateParams{
						PortID: param.Null[string](),
					}
					if !data.ProjectID.IsNull() {
						updateParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
					}
					if !data.RegionID.IsNull() {
						updateParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
					}

					_, err := r.client.Cloud.FloatingIPs.UpdateAndPoll(
						ctx,
						fipToUnassign,
						updateParams,
						option.WithMiddleware(logging.Middleware(ctx)),
					)
					if err != nil {
						resp.Diagnostics.AddError(
							"failed to unassign floating IP",
							fmt.Sprintf("Interface at index %d, floating IP %s: %s", i, fipToUnassign, err.Error()),
						)
						return
					}
					tflog.Info(ctx, "Unassigned floating IP created with source=new", map[string]interface{}{
						"floating_ip_id": fipToUnassign,
						"port_id":        portID,
						"interface":      i,
					})
					stateHasChanged = true
				}
				continue
			}

			// Skip if no change in existing FIP assignment
			if stateFloatingID == planFloatingID {
				continue
			}

			// Unassign old floating IP if one was attached
			if stateFloatingID != "" {
				updateParams := cloud.FloatingIPUpdateParams{
					PortID: param.Null[string](),
				}
				if !data.ProjectID.IsNull() {
					updateParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
				}
				if !data.RegionID.IsNull() {
					updateParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
				}

				_, err := r.client.Cloud.FloatingIPs.UpdateAndPoll(
					ctx,
					stateFloatingID,
					updateParams,
					option.WithMiddleware(logging.Middleware(ctx)),
				)
				if err != nil {
					resp.Diagnostics.AddError(
						"failed to unassign floating IP",
						fmt.Sprintf("Interface at index %d, floating IP %s: %s", i, stateFloatingID, err.Error()),
					)
					return
				}
				stateHasChanged = true
			}

			// Assign new floating IP if one is specified
			if planFloatingID != "" && portID != "" {
				updateParams := cloud.FloatingIPUpdateParams{
					PortID: param.Opt[string]{Value: portID},
				}
				if !data.ProjectID.IsNull() {
					updateParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
				}
				if !data.RegionID.IsNull() {
					updateParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
				}

				_, err := r.client.Cloud.FloatingIPs.UpdateAndPoll(
					ctx,
					planFloatingID,
					updateParams,
					option.WithMiddleware(logging.Middleware(ctx)),
				)
				if err != nil {
					resp.Diagnostics.AddError(
						"failed to assign floating IP",
						fmt.Sprintf("Interface at index %d, floating IP %s: %s", i, planFloatingID, err.Error()),
					)
					return
				}
				stateHasChanged = true
			}
		}
	}

	// Handle security group changes using assign/unassign endpoints
	if data.SecurityGroups != nil || state.SecurityGroups != nil {
		// Get security group IDs from state and plan
		stateIDs := make(map[string]bool)
		planIDs := make(map[string]bool)

		if state.SecurityGroups != nil {
			for _, sg := range *state.SecurityGroups {
				if !sg.ID.IsNull() && !sg.ID.IsUnknown() {
					stateIDs[sg.ID.ValueString()] = true
				}
			}
		}

		if data.SecurityGroups != nil {
			for _, sg := range *data.SecurityGroups {
				if !sg.ID.IsNull() && !sg.ID.IsUnknown() {
					planIDs[sg.ID.ValueString()] = true
				}
			}
		}

		// Find added and removed security groups
		var addedIDs []string
		var removedIDs []string

		for id := range planIDs {
			if !stateIDs[id] {
				addedIDs = append(addedIDs, id)
			}
		}

		for id := range stateIDs {
			if !planIDs[id] {
				removedIDs = append(removedIDs, id)
			}
		}

		// If there are changes, fetch security groups to build ID->Name map
		if len(addedIDs) > 0 || len(removedIDs) > 0 {
			sgListParams := cloud.SecurityGroupListParams{}
			if !data.ProjectID.IsNull() {
				sgListParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
			}
			if !data.RegionID.IsNull() {
				sgListParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
			}

			sgList, err := r.client.Cloud.SecurityGroups.List(
				ctx,
				sgListParams,
				option.WithMiddleware(logging.Middleware(ctx)),
			)
			if err != nil {
				resp.Diagnostics.AddError("failed to list security groups for ID to name mapping", err.Error())
				return
			}

			// Build ID to Name map
			idToName := make(map[string]string)
			for _, sg := range sgList.Results {
				idToName[sg.ID] = sg.Name
			}

			// Unassign removed security groups
			for _, id := range removedIDs {
				name, ok := idToName[id]
				if !ok {
					resp.Diagnostics.AddError(
						"security group not found",
						fmt.Sprintf("Could not find security group with ID %s", id),
					)
					return
				}

				unassignParams := cloud.InstanceUnassignSecurityGroupParams{
					Name: param.NewOpt(name),
				}
				if !data.ProjectID.IsNull() {
					unassignParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
				}
				if !data.RegionID.IsNull() {
					unassignParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
				}

				err := r.client.Cloud.Instances.UnassignSecurityGroup(
					ctx,
					instanceID,
					unassignParams,
					option.WithMiddleware(logging.Middleware(ctx)),
				)
				if err != nil {
					resp.Diagnostics.AddError(
						"failed to unassign security group",
						fmt.Sprintf("Security group %s (ID: %s): %s", name, id, err.Error()),
					)
					return
				}
				stateHasChanged = true
			}

			// Assign new security groups
			for _, id := range addedIDs {
				name, ok := idToName[id]
				if !ok {
					resp.Diagnostics.AddError(
						"security group not found",
						fmt.Sprintf("Could not find security group with ID %s", id),
					)
					return
				}

				assignParams := cloud.InstanceAssignSecurityGroupParams{
					Name: param.NewOpt(name),
				}
				if !data.ProjectID.IsNull() {
					assignParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
				}
				if !data.RegionID.IsNull() {
					assignParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
				}

				err := r.client.Cloud.Instances.AssignSecurityGroup(
					ctx,
					instanceID,
					assignParams,
					option.WithMiddleware(logging.Middleware(ctx)),
				)
				if err != nil {
					resp.Diagnostics.AddError(
						"failed to assign security group",
						fmt.Sprintf("Security group %s (ID: %s): %s", name, id, err.Error()),
					)
					return
				}
				stateHasChanged = true
			}
		}
	}

	// Handle servergroup/placement group changes using AddToPlacementGroup/RemoveFromPlacementGroup endpoints
	stateServergroupID := ""
	planServergroupID := ""
	if !state.ServergroupID.IsNull() && !state.ServergroupID.IsUnknown() {
		stateServergroupID = state.ServergroupID.ValueString()
	}
	if !data.ServergroupID.IsNull() && !data.ServergroupID.IsUnknown() {
		planServergroupID = data.ServergroupID.ValueString()
	}

	if stateServergroupID != planServergroupID {
		// First, remove from old servergroup if any
		if stateServergroupID != "" {
			removeParams := cloud.InstanceRemoveFromPlacementGroupParams{}
			if !data.ProjectID.IsNull() {
				removeParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
			}
			if !data.RegionID.IsNull() {
				removeParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
			}

			_, err := r.client.Cloud.Instances.RemoveFromPlacementGroupAndPoll(
				ctx,
				instanceID,
				removeParams,
				option.WithMiddleware(logging.Middleware(ctx)),
			)
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to remove instance from placement group",
					fmt.Sprintf("Placement group %s: %s", stateServergroupID, err.Error()),
				)
				return
			}
			stateHasChanged = true
		}

		// Then, add to new servergroup if any
		if planServergroupID != "" {
			addParams := cloud.InstanceAddToPlacementGroupParams{
				ServergroupID: planServergroupID,
			}
			if !data.ProjectID.IsNull() {
				addParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
			}
			if !data.RegionID.IsNull() {
				addParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
			}

			_, err := r.client.Cloud.Instances.AddToPlacementGroupAndPoll(
				ctx,
				instanceID,
				addParams,
				option.WithMiddleware(logging.Middleware(ctx)),
			)
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to add instance to placement group",
					fmt.Sprintf("Placement group %s: %s", planServergroupID, err.Error()),
				)
				return
			}
			stateHasChanged = true
		}
	}

	// Check if there are simple field changes that need the standard PATCH endpoint.
	// This runs AFTER specialized endpoints, allowing combined updates like:
	//   name = "new-name" (PATCH) + flavor = "g1-standard-2" (specialized /changeflavor)
	// The PATCH endpoint handles: name, tags (and potentially other simple fields in the future).
	nameChanged := !data.Name.Equal(state.Name)
	tagsChanged := !data.Tags.Equal(state.Tags)

	if nameChanged || tagsChanged {
		params := cloud.InstanceUpdateParams{}

		if !data.ProjectID.IsNull() {
			params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}

		if !data.RegionID.IsNull() {
			params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		dataBytes, err := data.MarshalJSONForUpdate(*state)
		if err != nil {
			resp.Diagnostics.AddError("failed to serialize http request", err.Error())
			return
		}

		_, err = r.client.Cloud.Instances.Update(
			ctx,
			instanceID,
			params,
			option.WithRequestBody("application/json", dataBytes),
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to update instance", err.Error())
			return
		}
		stateHasChanged = true
	}

	// If any changes were made (specialized endpoints or PATCH), refresh state from API
	if stateHasChanged {
		// Re-read instance to get latest state after all operations
		readParams := cloud.InstanceGetParams{}
		if !data.ProjectID.IsNull() {
			readParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			readParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		instance, err := r.client.Cloud.Instances.Get(
			ctx,
			instanceID,
			readParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to read instance after update", err.Error())
			return
		}
		err = apijson.UnmarshalComputed([]byte(instance.RawJSON()), &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize instance after update", err.Error())
			return
		}

		// Extract flavor_id from the flavor object in the API response
		// API returns: {"flavor": {"flavor_id": "...", ...}}
		// We need to extract just the flavor_id string
		var rawResponse map[string]interface{}
		if err := json.Unmarshal([]byte(instance.RawJSON()), &rawResponse); err == nil {
			if flavorObj, ok := rawResponse["flavor"].(map[string]interface{}); ok {
				if flavorID, ok := flavorObj["flavor_id"].(string); ok {
					data.Flavor = types.StringValue(flavorID)
				}
			}
		}

		// Extract interface port_id and ip_address from the interfaces list
		if data.Interfaces != nil && len(*data.Interfaces) > 0 {
			listParams := cloud.InstanceInterfaceListParams{}
			if !data.ProjectID.IsNull() {
				listParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
			}
			if !data.RegionID.IsNull() {
				listParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
			}

			interfaces, err := r.client.Cloud.Instances.Interfaces.List(
				ctx,
				instanceID,
				listParams,
				option.WithMiddleware(logging.Middleware(ctx)),
			)
			if err == nil && interfaces != nil {
				mergeInterfaceComputedFields(data.Interfaces, interfaces.Results, false)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudInstanceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.InstanceGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	// Preserve interfaces and volumes before UnmarshalComputed overwrites them.
	// This pattern is unique to cloud_instance because:
	// - GET /instances doesn't return full interface details (no floating_ip config, etc.)
	// - GET /instances doesn't return volume boot_index
	// - We'd lose user-specified config on every refresh without this
	// Other resources don't need this because their APIs return complete data.
	priorInterfaces := data.Interfaces
	priorVolumes := data.Volumes

	instance, err := r.client.Cloud.Instances.Get(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		// Check if it's a 404 error
		if apiErr, ok := err.(*gcore.Error); ok && apiErr.StatusCode == 404 {
			resp.Diagnostics.AddWarning("Resource not found", "The resource was not found on the server and will be removed from state.")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	err = apijson.UnmarshalComputed([]byte(instance.RawJSON()), &data)

	// Restore interfaces and volumes from prior state after UnmarshalComputed
	data.Interfaces = priorInterfaces
	data.Volumes = priorVolumes
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	// Extract fields from the API response that aren't mapped via json tags
	// API returns: {"flavor": {"flavor_id": "..."}, "project_id": N, "region_id": N, ...}
	var rawResponse map[string]interface{}
	if err := json.Unmarshal([]byte(instance.RawJSON()), &rawResponse); err == nil {
		if flavorObj, ok := rawResponse["flavor"].(map[string]interface{}); ok {
			if flavorID, ok := flavorObj["flavor_id"].(string); ok {
				data.Flavor = types.StringValue(flavorID)
			}
		}
		// Note: volumes are user-provided (existing volumes only), preserved from prior state above
	}

	// Extract interface port_id, ip_address, and floating_ip from the interfaces list
	if data.Interfaces != nil && len(*data.Interfaces) > 0 {
		listParams := cloud.InstanceInterfaceListParams{}
		if !data.ProjectID.IsNull() {
			listParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			listParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		interfaces, err := r.client.Cloud.Instances.Interfaces.List(
			ctx,
			data.ID.ValueString(),
			listParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err == nil && interfaces != nil {
			mergeInterfaceComputedFields(data.Interfaces, interfaces.Results, true)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudInstanceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.InstanceDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.Instances.DeleteAndPoll(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudInstanceModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_instance_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<instance_id>",
		&path_project_id,
		&path_region_id,
		&path_instance_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.ID = types.StringValue(path_instance_id)

	instance, err := r.client.Cloud.Instances.Get(
		ctx,
		path_instance_id,
		cloud.InstanceGetParams{
			ProjectID: param.NewOpt(path_project_id),
			RegionID:  param.NewOpt(path_region_id),
		},
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	rawJSON := instance.RawJSON()

	err = apijson.UnmarshalComputed([]byte(rawJSON), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	// Extract computed values from the API response
	var rawResponse map[string]interface{}
	if err := json.Unmarshal([]byte(rawJSON), &rawResponse); err == nil {
		// Extract flavor_id from the flavor object
		if flavorObj, ok := rawResponse["flavor"].(map[string]interface{}); ok {
			if flavorID, ok := flavorObj["flavor_id"].(string); ok {
				data.Flavor = types.StringValue(flavorID)
			}
		}

		// Extract name from API response
		if name, ok := rawResponse["name"].(string); ok {
			data.Name = types.StringValue(name)
		}

		// Create volumes array from API response for import
		// Only extract volume_id and boot_index (all other fields are managed via gcore_cloud_volume)
		// API returns: {"volumes": [{"id": "...", "delete_on_termination": bool}]}
		// We fetch each volume from the Volume API to check its "bootable" field
		if volumesArr, ok := rawResponse["volumes"].([]interface{}); ok {
			volumes := make([]*CloudInstanceVolumesModel, 0, len(volumesArr))
			for _, vol := range volumesArr {
				volMap, ok := vol.(map[string]interface{})
				if !ok {
					continue
				}
				volModel := &CloudInstanceVolumesModel{}
				volumeID, ok := volMap["id"].(string)
				if !ok {
					continue
				}
				volModel.VolumeID = types.StringValue(volumeID)

				// Note: We intentionally do NOT set boot_index during import.
				// This matches the old provider behavior where boot_index is only
				// stored if the user explicitly specifies it in their config.
				// This prevents drift when user imports without specifying boot_index.
				// boot_index stays null - user must specify it if they want to manage it
				volumes = append(volumes, volModel)
			}
			data.Volumes = &volumes
		}

		// Extract tags from API response
		// API returns: "tags": [{"key": "k1", "value": "v1", "read_only": bool}, ...]
		// Resource model uses: map[string]string
		// Note: Skip read_only tags as they are system-generated (e.g., image_id, os_distro)
		if tagsArr, ok := rawResponse["tags"].([]interface{}); ok && len(tagsArr) > 0 {
			tags := make(map[string]types.String, len(tagsArr))
			for _, tag := range tagsArr {
				tagMap, ok := tag.(map[string]interface{})
				if !ok {
					continue
				}
				// Skip read-only tags (system-generated)
				if readOnly, ok := tagMap["read_only"].(bool); ok && readOnly {
					continue
				}
				key, keyOk := tagMap["key"].(string)
				value, valueOk := tagMap["value"].(string)
				if keyOk && valueOk {
					tags[key] = types.StringValue(value)
				}
			}
			if len(tags) > 0 {
				data.Tags = customfield.NewMapMust[types.String](ctx, tags)
			}
		}
	}

	// Fetch interfaces from the API to populate port_id and ip_address
	// This is needed for import to have complete interface state
	listParams := cloud.InstanceInterfaceListParams{
		ProjectID: param.NewOpt(path_project_id),
		RegionID:  param.NewOpt(path_region_id),
	}

	interfacesResp, err := r.client.Cloud.Instances.Interfaces.List(
		ctx,
		path_instance_id,
		listParams,
		option.WithMiddleware(logging.Middleware(ctx)),
	)

	// Populate interfaces if API call succeeded
	if err == nil && interfacesResp != nil && len(interfacesResp.Results) > 0 {
		interfaces := make([]*CloudInstanceInterfacesModel, 0, len(interfacesResp.Results))
		for _, iface := range interfacesResp.Results {
			ifaceModel := &CloudInstanceInterfacesModel{
				PortID: types.StringValue(iface.PortID),
			}

			// Get IP address from ip_assignments
			if len(iface.IPAssignments) > 0 {
				ifaceModel.IPAddress = types.StringValue(iface.IPAssignments[0].IPAddress)
			}

			// Determine interface type based on network_details.external flag
			// External networks have this flag set to true
			if iface.NetworkDetails.External {
				ifaceModel.Type = types.StringValue("external")
			} else {
				// Private network - set as subnet type
				ifaceModel.Type = types.StringValue("subnet")
				ifaceModel.NetworkID = types.StringValue(iface.NetworkID)
				// Try to get subnet_id from ip_assignments
				if len(iface.IPAssignments) > 0 {
					ifaceModel.SubnetID = types.StringValue(iface.IPAssignments[0].SubnetID)
				}
			}

			// Extract interface_name if available
			if iface.InterfaceName != "" {
				ifaceModel.InterfaceName = types.StringValue(iface.InterfaceName)
			}

			// Extract floating IP if attached
			if len(iface.FloatingipDetails) > 0 {
				// Use the first floating IP (most common case)
				floatingIP := iface.FloatingipDetails[0]
				ifaceModel.FloatingIP = &CloudInstanceInterfacesFloatingIPModel{
					Source:             types.StringValue("existing"),
					ExistingFloatingID: types.StringValue(floatingIP.ID),
				}
			}

			interfaces = append(interfaces, ifaceModel)
		}
		data.Interfaces = &interfaces
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudInstanceResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}

// mergeInterfaceComputedFields matches Terraform interfaces to API interfaces and copies
// computed fields (port_id, ip_address). Uses port_id for stable matching when available,
// falls back to index for new interfaces. Set updateFloatingIP=true to also sync floating IP state.
func mergeInterfaceComputedFields(tfInterfaces *[]*CloudInstanceInterfacesModel, apiInterfaces []cloud.NetworkInterface, updateFloatingIP bool) {
	if tfInterfaces == nil || len(*tfInterfaces) == 0 || len(apiInterfaces) == 0 {
		return
	}

	// Build map of port_id -> index for efficient lookup
	apiIfaceByPort := make(map[string]int)
	for idx, iface := range apiInterfaces {
		apiIfaceByPort[iface.PortID] = idx
	}

	for i := range *tfInterfaces {
		tfIface := (*tfInterfaces)[i]
		apiIdx := -1

		// Match by port_id when available
		if !tfIface.PortID.IsNull() && !tfIface.PortID.IsUnknown() {
			if idx, ok := apiIfaceByPort[tfIface.PortID.ValueString()]; ok {
				apiIdx = idx
			}
		}

		// Fall back to index for new interfaces without port_id yet
		if apiIdx < 0 && i < len(apiInterfaces) {
			apiIdx = i
		}

		if apiIdx < 0 || apiIdx >= len(apiInterfaces) {
			continue
		}

		apiIface := apiInterfaces[apiIdx]
		(*tfInterfaces)[i].PortID = types.StringValue(apiIface.PortID)

		if len(apiIface.IPAssignments) > 0 {
			(*tfInterfaces)[i].IPAddress = types.StringValue(apiIface.IPAssignments[0].IPAddress)
		}

		if updateFloatingIP {
			if len(apiIface.FloatingipDetails) > 0 {
				floatingIP := apiIface.FloatingipDetails[0]
				currentSource := "existing"
				if (*tfInterfaces)[i].FloatingIP != nil && !(*tfInterfaces)[i].FloatingIP.Source.IsNull() {
					currentSource = (*tfInterfaces)[i].FloatingIP.Source.ValueString()
				}
				if currentSource == "new" {
					(*tfInterfaces)[i].FloatingIP = &CloudInstanceInterfacesFloatingIPModel{
						Source:             types.StringValue("new"),
						ExistingFloatingID: types.StringNull(),
					}
				} else {
					(*tfInterfaces)[i].FloatingIP = &CloudInstanceInterfacesFloatingIPModel{
						Source:             types.StringValue("existing"),
						ExistingFloatingID: types.StringValue(floatingIP.ID),
					}
				}
			} else {
				(*tfInterfaces)[i].FloatingIP = nil
			}
		}
	}
}
