// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_secret

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/fastedge"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*FastedgeSecretResource)(nil)
var _ resource.ResourceWithModifyPlan = (*FastedgeSecretResource)(nil)
var _ resource.ResourceWithImportState = (*FastedgeSecretResource)(nil)

func NewResource() resource.Resource {
	return &FastedgeSecretResource{}
}

// FastedgeSecretResource defines the resource implementation.
type FastedgeSecretResource struct {
	client *gcore.Client
}

func (r *FastedgeSecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fastedge_secret"
}

func (r *FastedgeSecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// clearSecretValues removes secret values from the model before saving to state.
// The API never returns secret values (they are write-only), so we must not persist
// the values from the plan/config into state. This prevents sensitive data leakage
// into the state file. Only the checksum (returned by the API) is kept for integrity
// verification.
func clearSecretValues(ctx context.Context, data *FastedgeSecretModel) error {
	if data.SecretSlots.IsNull() || data.SecretSlots.IsUnknown() {
		return nil
	}

	slots, diags := data.SecretSlots.AsStructSliceT(ctx)
	if diags.HasError() {
		return fmt.Errorf("failed to read secret slots: %v", diags)
	}

	for i := range slots {
		slots[i].Value = types.StringNull()
	}

	newSlots, diags := customfield.NewObjectSet[FastedgeSecretSecretSlotsModel](ctx, slots)
	if diags.HasError() {
		return fmt.Errorf("failed to rebuild secret slots: %v", diags)
	}
	data.SecretSlots = newSlots
	return nil
}

func (r *FastedgeSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FastedgeSecretModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read secret values from config since they are stripped from the plan
	// by ModifyPlan to prevent state leakage.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("secret_slots"), &data.SecretSlots)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataBytes, err := data.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}
	res := new(http.Response)
	_, err = r.client.Fastedge.Secrets.New(
		ctx,
		fastedge.SecretNewParams{},
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

	if err := clearSecretValues(ctx, data); err != nil {
		resp.Diagnostics.AddError("failed to clear secret values from state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FastedgeSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *FastedgeSecretModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *FastedgeSecretModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve ID since the PATCH response doesn't include it
	id := data.ID

	// Check if secret_slots actually changed by comparing plan vs state.
	// ModifyPlan copies state.SecretSlots when nothing changed, so if they're
	// equal, we know secrets didn't change and shouldn't be in the PATCH body.
	slotsChanged := !data.SecretSlots.Equal(state.SecretSlots)

	if slotsChanged {
		// Read secret values from config since they are stripped from the plan.
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("secret_slots"), &data.SecretSlots)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	dataBytes, err := data.MarshalJSONForUpdate(*state)
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}
	res := new(http.Response)
	_, err = r.client.Fastedge.Secrets.Update(
		ctx,
		data.ID.ValueInt64(),
		fastedge.SecretUpdateParams{},
		option.WithRequestBody("application/json", dataBytes),
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)

	// Parse API response to get updated checksums.
	var apiData FastedgeSecretModel
	err = apijson.UnmarshalComputed(bytes, &apiData)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	// Use plan values for required/optional fields (name, comment) since
	// UnmarshalComputed skips them. Only use apiData for computed fields.
	data.ID = id
	// data.Name and data.Comment already have the correct plan values.
	data.AppCount = apiData.AppCount

	// Build API checksum map.
	apiChecksumBySlot := map[int64]string{}
	if !apiData.SecretSlots.IsNull() && !apiData.SecretSlots.IsUnknown() {
		apiSlots, d := apiData.SecretSlots.AsStructSliceT(ctx)
		resp.Diagnostics.Append(d...)
		for _, s := range apiSlots {
			if !s.Checksum.IsNull() && !s.Checksum.IsUnknown() {
				apiChecksumBySlot[s.Slot.ValueInt64()] = s.Checksum.ValueString()
			}
		}
	}

	// Handle secret_slots in the result:
	if !slotsChanged {
		// Secret slots didn't change — preserve exact state to match the plan.
		data.SecretSlots = state.SecretSlots
	} else if len(apiChecksumBySlot) > 0 {
		// API returned new checksums — build result with updated checksums.
		configSlots, d := data.SecretSlots.AsStructSliceT(ctx)
		resp.Diagnostics.Append(d...)
		updatedSlots := make([]FastedgeSecretSecretSlotsModel, 0, len(configSlots))
		for _, cs := range configSlots {
			slotID := cs.Slot.ValueInt64()
			checksum := types.StringNull()
			if apiCS, ok := apiChecksumBySlot[slotID]; ok {
				checksum = types.StringValue(apiCS)
			}
			updatedSlots = append(updatedSlots, FastedgeSecretSecretSlotsModel{
				Slot:     cs.Slot,
				Checksum: checksum,
				Value:    types.StringNull(),
			})
		}
		newSlots, d := customfield.NewObjectSet[FastedgeSecretSecretSlotsModel](ctx, updatedSlots)
		resp.Diagnostics.Append(d...)
		data.SecretSlots = newSlots
	} else {
		// Fallback: use API response directly.
		data.SecretSlots = apiData.SecretSlots
		if err := clearSecretValues(ctx, data); err != nil {
			resp.Diagnostics.AddError("failed to clear secret values from state", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FastedgeSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FastedgeSecretModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve ID since the GET response doesn't include it
	id := data.ID

	res := new(http.Response)
	_, err := r.client.Fastedge.Secrets.Get(
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

	// Restore ID since it's not in the GET response body
	data.ID = id

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FastedgeSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FastedgeSecretModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res := new(http.Response)
	err := r.client.Fastedge.Secrets.Delete(
		ctx,
		data.ID.ValueInt64(),
		fastedge.SecretDeleteParams{
			Force: param.NewOpt(true), // Force delete even if secret is in use
		},
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)

	// Handle 409 Conflict - secret is referenced and cannot be deleted
	if res != nil && res.StatusCode == http.StatusConflict {
		resp.Diagnostics.AddWarning(
			"Secret is referenced",
			fmt.Sprintf("Secret (%d) is referenced by one or more applications and cannot be deleted. The resource will remain in state.", data.ID.ValueInt64()),
		)
		return // Don't clear state - secret still exists
	}

	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	// Successfully deleted - state will be cleared automatically
}

func (r *FastedgeSecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(FastedgeSecretModel)

	importID := int64(0)
	diags := importpath.ParseImportID(
		req.ID,
		"<id>",
		&importID,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.Int64Value(importID)

	res := new(http.Response)
	_, err := r.client.Fastedge.Secrets.Get(
		ctx,
		importID,
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

	// Restore ID since the GET response doesn't include it
	data.ID = types.Int64Value(importID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func secretChecksum(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (r *FastedgeSecretResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Destroy — nothing to do.
	if req.Plan.Raw.IsNull() {
		return
	}

	// Create — no state to compare against; just clear values so the
	// post-apply state (which also has null values) is consistent.
	if req.State.Raw.IsNull() {
		var plan FastedgeSecretModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if err := clearSecretValues(ctx, &plan); err != nil {
			resp.Diagnostics.AddError("failed to clear secret values from plan", err.Error())
			return
		}
		resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
		return
	}

	// Update — compare per-slot checksums to detect real changes.
	var config FastedgeSecretModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FastedgeSecretModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If config has no secret_slots, let Terraform handle it normally.
	if config.SecretSlots.IsNull() || config.SecretSlots.IsUnknown() {
		return
	}

	configSlots, diags := config.SecretSlots.AsStructSliceT(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch true slot checksums directly from the API. We cannot rely on
	// req.State because SetNestedAttribute's semantic equality can swap
	// element mappings when matching the refreshed Read response against
	// prior state elements.
	apiChecksumBySlot := map[int64]string{}
	if !state.ID.IsNull() && !state.ID.IsUnknown() && r.client != nil {
		apiRes := new(http.Response)
		_, err := r.client.Fastedge.Secrets.Get(
			ctx,
			state.ID.ValueInt64(),
			option.WithResponseBodyInto(&apiRes),
		)
		if err == nil && apiRes.StatusCode == http.StatusOK {
			apiBytes, _ := io.ReadAll(apiRes.Body)
			var apiData FastedgeSecretModel
			if unmarshalErr := apijson.UnmarshalComputed(apiBytes, &apiData); unmarshalErr == nil {
				if !apiData.SecretSlots.IsNull() && !apiData.SecretSlots.IsUnknown() {
					apiSlots, d := apiData.SecretSlots.AsStructSliceT(ctx)
					resp.Diagnostics.Append(d...)
					for _, s := range apiSlots {
						if !s.Checksum.IsNull() && !s.Checksum.IsUnknown() {
							apiChecksumBySlot[s.Slot.ValueInt64()] = s.Checksum.ValueString()
						}
					}
				}
			}
		}
	}

	// Determine if any slots actually changed.
	changed := false
	resolved := make([]FastedgeSecretSecretSlotsModel, 0, len(configSlots))
	for _, cs := range configSlots {
		slotID := cs.Slot.ValueInt64()
		apiChecksum, inAPI := apiChecksumBySlot[slotID]

		var checksum types.String
		if !cs.Value.IsNull() && !cs.Value.IsUnknown() && inAPI {
			if secretChecksum(cs.Value.ValueString()) == apiChecksum {
				checksum = types.StringValue(apiChecksum) // unchanged
			} else {
				checksum = types.StringUnknown() // changed → triggers update
				changed = true
			}
		} else if cs.Value.IsNull() && inAPI {
			checksum = types.StringValue(apiChecksum)
		} else {
			checksum = types.StringUnknown()
			changed = true
		}

		resolved = append(resolved, FastedgeSecretSecretSlotsModel{
			Slot:     cs.Slot,
			Checksum: checksum,
			Value:    types.StringNull(),
		})
	}

	// Check if slots were removed.
	if len(apiChecksumBySlot) != len(configSlots) {
		changed = true
	}

	// When nothing changed, copy state's secret_slots wholesale to preserve
	// exact Set element identity and avoid spurious diffs.
	if !changed {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("secret_slots"), state.SecretSlots)...)
		return
	}

	newSlots, diags := customfield.NewObjectSet[FastedgeSecretSecretSlotsModel](ctx, resolved)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("secret_slots"), newSlots)...)
}
