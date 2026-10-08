// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_policy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/waap"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*WaapDomainPolicyResource)(nil)
var _ resource.ResourceWithModifyPlan = (*WaapDomainPolicyResource)(nil)

func NewResource() resource.Resource {
	return &WaapDomainPolicyResource{}
}

// WaapDomainPolicyResource defines the resource implementation.
type WaapDomainPolicyResource struct {
	client *gcore.Client
}

func (r *WaapDomainPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_waap_domain_policy"
}

func (r *WaapDomainPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WaapDomainPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *WaapDomainPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dataBytes, err := data.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}
	res := new(http.Response)
	_, err = r.client.Waap.Domains.Policies.Update(
		ctx,
		data.PolicyID.ValueString(),
		waap.DomainPolicyUpdateParams{
			DomainID: data.DomainID.ValueInt64(),
		},
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
	data.ID = data.PolicyID

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WaapDomainPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *WaapDomainPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *WaapDomainPolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dataBytes, err := data.MarshalJSONForUpdate(*state)
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}
	res := new(http.Response)
	_, err = r.client.Waap.Domains.Policies.Update(
		ctx,
		data.PolicyID.ValueString(),
		waap.DomainPolicyUpdateParams{
			DomainID: data.DomainID.ValueInt64(),
		},
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
	data.ID = data.PolicyID

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WaapDomainPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *WaapDomainPolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// The API has no single-policy GET, so page through the domain's policies.
	// Drop this once GET /waap/v2/domains/{domain_id}/policies/{policy_id} exists:
	// the generator then emits Read and ImportState itself.
	iter := r.client.Waap.Domains.Policies.ListAutoPaging(
		ctx,
		data.DomainID.ValueInt64(),
		waap.DomainPolicyListParams{},
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	for iter.Next() {
		policy := iter.Current()
		if policy.ID != data.PolicyID.ValueString() {
			continue
		}
		data.Mode = types.BoolValue(policy.Mode)
		data.ID = data.PolicyID
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}
	var apiErr *gcore.Error
	if err := iter.Err(); err != nil && !(errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound) {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	resp.Diagnostics.AddWarning("Resource not found", "The resource was not found on the server and will be removed from state.")
	resp.State.RemoveResource(ctx)
}

func (r *WaapDomainPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *WaapDomainPolicyResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}
