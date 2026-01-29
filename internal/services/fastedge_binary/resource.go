// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_binary

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/importpath"
	"github.com/stainless-sdks/gcore-terraform/internal/logging"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*FastedgeBinaryResource)(nil)
var _ resource.ResourceWithModifyPlan = (*FastedgeBinaryResource)(nil)
var _ resource.ResourceWithImportState = (*FastedgeBinaryResource)(nil)

func NewResource() resource.Resource {
	return &FastedgeBinaryResource{}
}

// FastedgeBinaryResource defines the resource implementation.
type FastedgeBinaryResource struct {
	client *gcore.Client
}

func (r *FastedgeBinaryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fastedge_binary"
}

func (r *FastedgeBinaryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FastedgeBinaryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FastedgeBinaryModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	//dataBytes, err := data.MarshalJSON()
	//if err != nil {
	//	resp.Diagnostics.AddError("failed to serialize http request", err.Error())
	//	return
	//}
	res := new(http.Response)
	//_, err = r.client.Fastedge.Binaries.New(
	//	ctx,
	//	data.Body.ValueString(),
	//	option.WithRequestBody("application/json", dataBytes),
	//	option.WithResponseBodyInto(&res),
	//	option.WithMiddleware(logging.Middleware(ctx)),
	//)
	//if err != nil {
	//	resp.Diagnostics.AddError("failed to make http request", err.Error())
	//	return
	//}
	bytes, _ := io.ReadAll(res.Body)
	err := apijson.UnmarshalComputed(bytes, data.Body)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FastedgeBinaryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update is not supported for this resource
}

func (r *FastedgeBinaryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FastedgeBinaryModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res := new(http.Response)
	_, err := r.client.Fastedge.Binaries.Get(
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
	err = apijson.Unmarshal(bytes, data.Body)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FastedgeBinaryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FastedgeBinaryModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Fastedge.Binaries.Delete(
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

func (r *FastedgeBinaryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(FastedgeBinaryModel)

	path := int64(0)
	diags := importpath.ParseImportID(
		req.ID,
		"<id>",
		&path,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.Int64Value(path)

	res := new(http.Response)
	_, err := r.client.Fastedge.Binaries.Get(
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
	err = apijson.Unmarshal(bytes, data.Body)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FastedgeBinaryResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}
