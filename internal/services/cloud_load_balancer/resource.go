// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/custom"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudLoadBalancerResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudLoadBalancerResource)(nil)
var _ resource.ResourceWithImportState = (*CloudLoadBalancerResource)(nil)

func NewResource() resource.Resource {
	return &CloudLoadBalancerResource{}
}

// CloudLoadBalancerResource defines the resource implementation.
type CloudLoadBalancerResource struct {
	client *gcore.Client
}

func (r *CloudLoadBalancerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_load_balancer"
}

func (r *CloudLoadBalancerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudLoadBalancerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudLoadBalancerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dataBytes, err := data.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}

	params := cloud.LoadBalancerNewParams{}
	if err := params.UnmarshalJSON(dataBytes); err != nil {
		resp.Diagnostics.AddError("failed to deserialize into params", err.Error())
		return
	}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err = r.client.Cloud.LoadBalancers.NewAndPoll(
		ctx,
		params,
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
		data.Tags = tags
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudLoadBalancerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudLoadBalancerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read the adoption marker before the flavor short-circuit: every
	// state-writing path below, including resize, must verify adopted values
	// and close the adoption window, or a resize right after import would
	// adopt silently and leave the marker set forever.
	adopting, diags := importAdoptionPending(ctx, req.Private)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle flavor changes with ResizeAndPoll - if only flavor changed, skip regular update
	flavorChanged := !data.Flavor.Equal(state.Flavor) && !data.Flavor.IsNull()

	if flavorChanged {
		resizeParams := cloud.LoadBalancerResizeParams{
			Flavor: data.Flavor.ValueString(),
		}

		if !data.ProjectID.IsNull() {
			resizeParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}

		if !data.RegionID.IsNull() {
			resizeParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		res := new(http.Response)
		_, err := r.client.Cloud.LoadBalancers.ResizeAndPoll(
			ctx,
			data.ID.ValueString(),
			resizeParams,
			option.WithResponseBodyInto(&res),
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to resize load balancer", err.Error())
			return
		}
		bytes, _ := io.ReadAll(res.Body)
		err = apijson.UnmarshalComputed(bytes, &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize resize response", err.Error())
			return
		}
		if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
			data.Tags = tags
		}

		// After resize, set state and return (don't call regular Update).
		// Verify before writing, clear only after writing: a mismatch must
		// leave both state and the marker untouched so the next apply replays
		// the same adoption decision, and a marker cleared before success
		// would turn a transient failure into a destroy+recreate on retry.
		// ResizeAndPoll finishes with a GET, so bytes carries the evidence.
		if adopting {
			resp.Diagnostics.Append(r.verifyImportAdoption(ctx, data, state, bytes)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		if adopting && !resp.Diagnostics.HasError() {
			resp.Diagnostics.Append(clearImportAdoptionPending(ctx, resp.Private)...)
		}
		return
	}

	params := cloud.LoadBalancerUpdateParams{}

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

	// vip_network_id and vip_subnet_id are create-only and null in
	// state after an import; adopting the configured values into state must not
	// send them to the PATCH endpoint, which does not accept them. Outside that
	// adoption window the plan modifiers force replacement instead, so nothing
	// is stripped and a stray field fails loudly against the API.
	dataBytes, err = stripAdoptedCreateOnlyFields(dataBytes, *state, adopting)
	if err != nil {
		resp.Diagnostics.AddError("failed to filter update request body", err.Error())
		return
	}

	dataStr := strings.TrimSpace(string(dataBytes))

	// If no fields have changed, skip the update and just refresh from API
	if dataStr == "{}" || dataStr == "null" || len(dataBytes) == 0 {
		// No changes to send - just read current state
		res := new(http.Response)
		_, err := r.client.Cloud.LoadBalancers.Get(
			ctx,
			data.ID.ValueString(),
			cloud.LoadBalancerGetParams{
				ProjectID: params.ProjectID,
				RegionID:  params.RegionID,
			},
			option.WithResponseBodyInto(&res),
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to read load balancer", err.Error())
			return
		}
		bytes, _ := io.ReadAll(res.Body)
		err = apijson.UnmarshalComputed(bytes, &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize response", err.Error())
			return
		}
		if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
			data.Tags = tags
		}
		// Verify before writing, clear only after writing: a mismatch must
		// leave both state and the marker untouched so the next apply replays
		// the same adoption decision, and a marker cleared before success
		// would turn a transient failure into a destroy+recreate on retry.
		if adopting {
			resp.Diagnostics.Append(r.verifyImportAdoption(ctx, data, state, bytes)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		if adopting && !resp.Diagnostics.HasError() {
			resp.Diagnostics.Append(clearImportAdoptionPending(ctx, resp.Private)...)
		}
		return
	}

	// Call Update - note: API returns empty vrrp_ips despite docs claiming otherwise
	_, err = r.client.Cloud.LoadBalancers.Update(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to update load balancer", err.Error())
		return
	}

	// IMPORTANT: The PATCH endpoint returns vrrp_ips: [] (empty array) even though docs say otherwise.
	// We must do an explicit GET to retrieve all computed fields including vrrp_ips.
	// This matches the old provider behavior: resourceLoadBalancerV2Read(ctx, d, m)
	res := new(http.Response)
	_, err = r.client.Cloud.LoadBalancers.Get(
		ctx,
		data.ID.ValueString(),
		cloud.LoadBalancerGetParams{
			ProjectID: params.ProjectID,
			RegionID:  params.RegionID,
		},
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
		data.Tags = tags
	}

	// Verify before writing, clear only after writing: a mismatch must leave
	// both state and the marker untouched so the next apply replays the same
	// adoption decision, and a marker cleared before success would turn a
	// transient failure into a destroy+recreate on retry.
	if adopting {
		resp.Diagnostics.Append(r.verifyImportAdoption(ctx, data, state, bytes)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if adopting && !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(clearImportAdoptionPending(ctx, resp.Private)...)
	}
}

// verifyImportAdoption checks every create-only value this apply adopts into
// state against the load balancer the API just described, using the response
// body already in hand plus at most a couple of subnet lookups. It must run on
// every state-writing path of Update while the adoption marker is set, BEFORE
// resp.State.Set and before the marker is cleared: a mismatch has to abort the
// apply with the marker intact, so a retry replays the adoption decision
// instead of degrading into a destroy+recreate of live infrastructure.
func (r *CloudLoadBalancerResource) verifyImportAdoption(ctx context.Context, plan, state *CloudLoadBalancerModel, lbBytes []byte) diag.Diagnostics {
	var diags diag.Diagnostics

	evidence, err := parseAdoptionEvidence(lbBytes)
	if err != nil {
		diags.AddError("failed to inspect load balancer for import adoption", err.Error())
		return diags
	}

	subnetNetwork := func(ctx context.Context, subnetID string) (string, error) {
		params := cloud.NetworkSubnetGetParams{}
		if !plan.ProjectID.IsNull() {
			params.ProjectID = param.NewOpt(plan.ProjectID.ValueInt64())
		}
		if !plan.RegionID.IsNull() {
			params.RegionID = param.NewOpt(plan.RegionID.ValueInt64())
		}
		subnet, err := r.client.Cloud.Networks.Subnets.Get(ctx, subnetID, params, option.WithMiddleware(logging.Middleware(ctx)))
		if err != nil {
			return "", err
		}
		return subnet.NetworkID, nil
	}

	findings, err := verifyAdoptedCreateOnlyFields(ctx, *plan, *state, evidence, subnetNetwork)
	if err != nil {
		diags.AddError("failed to verify adopted load balancer attributes", err.Error())
		return diags
	}

	for _, finding := range findings {
		switch finding.verdict {
		case adoptionMismatch:
			diags.AddError(
				fmt.Sprintf("configured %s does not match the imported load balancer", finding.attribute),
				finding.detail+"\n\nThe load balancer update API does not accept this attribute, so the "+
					"value cannot be applied in place, and adopting it into state would record "+
					"infrastructure that does not exist. Either set the attribute to the value the "+
					"load balancer really has (or remove it), or recreate the load balancer with: "+
					"terraform apply -replace=<resource address>",
			)
		case adoptionUnverifiable:
			diags.AddWarning(
				fmt.Sprintf("adopted %s without verification", finding.attribute),
				finding.detail+" The configured value was recorded in state as-is; if it does not "+
					"match the real load balancer, state will not reflect reality.",
			)
		}
	}

	return diags
}

func (r *CloudLoadBalancerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudLoadBalancerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.LoadBalancerGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.LoadBalancers.Get(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if res != nil && res.StatusCode == 404 {
		resp.Diagnostics.AddWarning("Resource not found", "The resource was not found on the server and will be removed from state.")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
		data.Tags = tags
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudLoadBalancerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.LoadBalancerDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.LoadBalancers.DeleteAndPoll(
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

func (r *CloudLoadBalancerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudLoadBalancerModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_load_balancer_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<load_balancer_id>",
		&path_project_id,
		&path_region_id,
		&path_load_balancer_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.ID = types.StringValue(path_load_balancer_id)

	res := new(http.Response)
	_, err := r.client.Cloud.LoadBalancers.Get(
		ctx,
		path_load_balancer_id,
		cloud.LoadBalancerGetParams{
			ProjectID: param.NewOpt(path_project_id),
			RegionID:  param.NewOpt(path_region_id),
		},
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.Unmarshal(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
		data.Tags = tags
	}

	// Populate flavor from the nested flavor object in API response.
	// The API returns flavor as {"flavor_id": "...", "flavor_name": "lb1-1-2", ...}
	// but Terraform expects just the flavor name string.
	// Without this, terraform plan after import would show flavor changing from null
	// to the config value, triggering an unnecessary resize operation.
	// See: GCLOUD2-20778
	if data.Flavor.IsNull() || data.Flavor.IsUnknown() {
		var rawResponse struct {
			Flavor *struct {
				FlavorName string `json:"flavor_name"`
			} `json:"flavor"`
		}
		if err := json.Unmarshal(bytes, &rawResponse); err == nil && rawResponse.Flavor != nil && rawResponse.Flavor.FlavorName != "" {
			data.Flavor = types.StringValue(rawResponse.Flavor.FlavorName)
		}
	}

	// vip_network_id and vip_subnet_id are create-only and tagged
	// no_refresh, so the GET above cannot populate them and they stay null in
	// state. Mark the adoption window so the plan modifiers on those attributes
	// know this null came from an import rather than from a load balancer that
	// was created without them - the first apply adopts the configured values
	// into state and clears the marker.
	resp.Diagnostics.Append(markImportAdoptionPending(ctx, resp.Private)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {
	// Deliberately empty. Clearing the import-adoption marker here looks
	// attractive but cannot work: a plan whose only effect would be the clear
	// is a no-op plan, and Terraform never persists PlannedPrivate for no-op
	// plans (plain `terraform plan` writes no state at all, and a no-op apply
	// short-circuits before writing resource state). Staleness of the marker
	// is instead made harmless by verifying every adoption against the API in
	// Update before anything reaches state.
}
