// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_origin_group

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CDNOriginGroupResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CDNOriginGroupResource)(nil)
var _ resource.ResourceWithImportState = (*CDNOriginGroupResource)(nil)

func NewResource() resource.Resource {
	return &CDNOriginGroupResource{}
}

// CDNOriginGroupResource defines the resource implementation.
type CDNOriginGroupResource struct {
	client *gcore.Client
}

func (r *CDNOriginGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_origin_group"
}

func (r *CDNOriginGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CDNOriginGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CDNOriginGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write-only fields are null in plan, get auth block from config
	var configAuth *CDNOriginGroupAuthModel
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("auth"), &configAuth)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Auth = configAuth

	dataBytes, err := data.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}
	res := new(http.Response)
	_, err = r.client.CDN.OriginGroups.New(
		ctx,
		cdn.OriginGroupNewParams{},
		option.WithRequestBody("application/json", dataBytes),
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CDNOriginGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CDNOriginGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *CDNOriginGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write-only fields are null in plan, get auth block from config
	var configAuth *CDNOriginGroupAuthModel
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("auth"), &configAuth)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Auth = configAuth

	dataBytes, err := data.MarshalJSONForUpdate(*state)
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}
	res := new(http.Response)
	_, err = r.client.CDN.OriginGroups.Update(
		ctx,
		data.ID.ValueInt64(),
		cdn.OriginGroupUpdateParams{},
		option.WithRequestBody("application/json", dataBytes),
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CDNOriginGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CDNOriginGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res := new(http.Response)
	_, err := r.client.CDN.OriginGroups.Get(
		ctx,
		data.ID.ValueInt64(),
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CDNOriginGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CDNOriginGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CDN.OriginGroups.Delete(
		ctx,
		data.ID.ValueInt64(),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CDNOriginGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CDNOriginGroupModel)

	path := int64(0)
	diags := importpath.ParseImportID(
		req.ID,
		"<origin_group_id>",
		&path,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.Int64Value(path)

	res := new(http.Response)
	_, err := r.client.CDN.OriginGroups.Get(
		ctx,
		path,
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)
	// Use Unmarshal (not UnmarshalComputed) for import to populate ALL fields from API response,
	// including required and optional fields like name, auth_type, sources.
	// UnmarshalComputed is for Read/Update where we preserve user-configured values.
	err = apijson.Unmarshal(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CDNOriginGroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip if this is a destroy or create operation
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var plan, state CDNOriginGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve purely computed fields from state to prevent drift
	// These fields are not user-configurable

	// has_related_resources: purely computed field - always preserve from state
	plan.HasRelatedResources = state.HasRelatedResources

	// proxy_next_upstream: computed+optional field
	// Only preserve from state if user didn't specify in config (IsUnknown)
	// This allows users to change the value while preventing drift when not specified
	if plan.ProxyNextUpstream.IsUnknown() {
		plan.ProxyNextUpstream = state.ProxyNextUpstream
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}
