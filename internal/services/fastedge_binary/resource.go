// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_binary

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-framework/path"
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

	filename := data.Filename.ValueString()

	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		resp.Diagnostics.AddError("failed to open file", fmt.Sprintf("Cannot open file %s: %s", filename, err.Error()))
		return
	}
	defer file.Close()

	// Save expected checksum before unmarshal overwrites it
	expectedChecksum := data.Checksum.ValueString()

	// Upload the binary
	result, err := r.client.Fastedge.Binaries.New(
		ctx,
		file,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to upload binary", err.Error())
		return
	}

	// Verify checksum matches what we calculated
	if result.Checksum != expectedChecksum {
		resp.Diagnostics.AddError(
			"checksum mismatch",
			fmt.Sprintf("Uploaded binary checksum (%s) does not match expected (%s). The file may have been corrupted during upload.", result.Checksum, expectedChecksum),
		)
		return
	}

	// The Create response (BinaryShort) has the real ID but lacks the source field.
	// The Get response includes source but returns id=0 (API quirk).
	// Read the full binary to populate all computed fields.
	realID := result.ID
	fullResult, err := r.client.Fastedge.Binaries.Get(
		ctx,
		realID,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to read binary after upload", err.Error())
		return
	}
	err = apijson.UnmarshalComputed([]byte(fullResult.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	// ID is not in json tags (API returns id=0 in Get), set from Create response
	data.ID = types.Int64Value(realID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FastedgeBinaryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// The only update scenario is a filename path change without content change
	// (e.g., after import, or renaming the local file). No API call needed.
	var data *FastedgeBinaryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FastedgeBinaryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FastedgeBinaryModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.Fastedge.Binaries.Get(
		ctx,
		data.ID.ValueInt64(),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		if apiErr, ok := err.(*gcore.Error); ok && apiErr.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning("Resource not found", "The resource was not found on the server and will be removed from state.")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	// ID is not in json tags (API returns id=0 in Get), preserve from state
	id := data.ID
	err = apijson.UnmarshalComputed([]byte(result.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	data.ID = id

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
		if apiErr, ok := err.(*gcore.Error); ok && apiErr.StatusCode == http.StatusConflict {
			resp.Diagnostics.AddWarning(
				"Binary in use",
				fmt.Sprintf("FastEdge binary (%d) is referenced by an application and cannot be deleted. Remove the application first, then delete the binary.", data.ID.ValueInt64()),
			)
			return
		}
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
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

	result, err := r.client.Fastedge.Binaries.Get(
		ctx,
		path,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	err = apijson.Unmarshal([]byte(result.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	// ID is not in json tags (API returns id=0 in Get), set from import path
	data.ID = types.Int64Value(path)

	// Filename is unknown after import - user must provide it
	data.Filename = types.StringNull()

	resp.Diagnostics.AddWarning(
		"Filename required after import",
		"The 'filename' attribute is not known after import. Please update your configuration with the correct local file path. A plan will show changes until the filename is set.",
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FastedgeBinaryResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip if destroying
	if req.Plan.Raw.IsNull() {
		return
	}

	// Skip if creating (no prior state)
	if req.State.Raw.IsNull() {
		var plan *FastedgeBinaryModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Skip if filename is unknown
		if plan.Filename.IsUnknown() || plan.Filename.IsNull() {
			return
		}

		// Calculate checksum for new resource
		checksum, err := fileChecksum(plan.Filename.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to read file",
				fmt.Sprintf("Cannot read file %s to calculate checksum: %s", plan.Filename.ValueString(), err.Error()),
			)
			return
		}
		plan.Checksum = types.StringValue(checksum)
		resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
		return
	}

	var plan *FastedgeBinaryModel
	var state *FastedgeBinaryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Skip if filename is unknown (during import)
	if plan.Filename.IsUnknown() || plan.Filename.IsNull() {
		return
	}

	filename := plan.Filename.ValueString()

	// Calculate MD5 checksum of the local file
	checksum, err := fileChecksum(filename)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read file",
			fmt.Sprintf("Cannot read file %s to calculate checksum: %s", filename, err.Error()),
		)
		return
	}

	// Set the checksum in the plan
	plan.Checksum = types.StringValue(checksum)

	// If checksum differs from state, trigger replacement
	if !state.Checksum.IsNull() && state.Checksum.ValueString() != checksum {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("checksum"))
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

// fileChecksum calculates the MD5 checksum of a file
func fileChecksum(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
