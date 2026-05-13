// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_access_key

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/storage"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*StorageAccessKeyResource)(nil)
var _ resource.ResourceWithModifyPlan = (*StorageAccessKeyResource)(nil)
var _ resource.ResourceWithImportState = (*StorageAccessKeyResource)(nil)

func NewResource() resource.Resource {
	return &StorageAccessKeyResource{}
}

// StorageAccessKeyResource defines the resource implementation.
type StorageAccessKeyResource struct {
	client *gcore.Client
}

func (r *StorageAccessKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_access_key"
}

func (r *StorageAccessKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StorageAccessKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *StorageAccessKeyModel

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
	_, err = r.client.Storage.ObjectStorages.AccessKeys.New(
		ctx,
		data.StorageID.ValueInt64(),
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
	data.ID = data.AccessKey

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageAccessKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update is not supported for this resource
}

func (r *StorageAccessKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *StorageAccessKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res := new(http.Response)
	_, err := r.client.Storage.ObjectStorages.AccessKeys.Get(
		ctx,
		data.AccessKey.ValueString(),
		storage.ObjectStorageAccessKeyGetParams{
			StorageID: data.StorageID.ValueInt64(),
		},
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
	err = apijson.Unmarshal(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	data.ID = data.AccessKey

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageAccessKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *StorageAccessKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Storage.ObjectStorages.AccessKeys.Delete(
		ctx,
		data.AccessKey.ValueString(),
		storage.ObjectStorageAccessKeyDeleteParams{
			StorageID: data.StorageID.ValueInt64(),
		},
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	data.ID = data.AccessKey

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageAccessKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(StorageAccessKeyModel)

	path_storage_id := int64(0)
	path_access_key := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<storage_id>/<access_key>",
		&path_storage_id,
		&path_access_key,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.StorageID = types.Int64Value(path_storage_id)
	data.AccessKey = types.StringValue(path_access_key)

	res := new(http.Response)
	_, err := r.client.Storage.ObjectStorages.AccessKeys.Get(
		ctx,
		path_access_key,
		storage.ObjectStorageAccessKeyGetParams{
			StorageID: path_storage_id,
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
	data.ID = data.AccessKey

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageAccessKeyResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}
