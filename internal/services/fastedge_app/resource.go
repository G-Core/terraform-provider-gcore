// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/fastedge"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*FastedgeAppResource)(nil)
var _ resource.ResourceWithModifyPlan = (*FastedgeAppResource)(nil)

func NewResource() resource.Resource {
	return &FastedgeAppResource{}
}

// FastedgeAppResource defines the resource implementation.
type FastedgeAppResource struct {
	client *gcore.Client
}

func (r *FastedgeAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fastedge_app"
}

func (r *FastedgeAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FastedgeAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FastedgeAppModel

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
	_, err = r.client.Fastedge.Apps.New(
		ctx,
		fastedge.AppNewParams{},
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

func (r *FastedgeAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update is not supported for this resource
}

func (r *FastedgeAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *FastedgeAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *FastedgeAppResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"This resource cannot be destroyed from Terraform. If you create this resource, it will be "+
				"present in the API until manually deleted.",
		)
	}
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will remove the resource from the Terraform state "+
				"but will not change it in the API. If you would like to destroy or reset this resource "+
				"in the API, refer to the documentation for how to do it manually.",
		)
	}
}
