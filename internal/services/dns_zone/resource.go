// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/dns"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
	"github.com/stainless-sdks/gcore-terraform/internal/importpath"
	"github.com/stainless-sdks/gcore-terraform/internal/logging"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*DNSZoneResource)(nil)
var _ resource.ResourceWithModifyPlan = (*DNSZoneResource)(nil)
var _ resource.ResourceWithImportState = (*DNSZoneResource)(nil)

func NewResource() resource.Resource {
	return &DNSZoneResource{}
}

// DNSZoneResource defines the resource implementation.
type DNSZoneResource struct {
	client *gcore.Client
}

func (r *DNSZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zone"
}

func (r *DNSZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DNSZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DNSZoneModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save planned meta before any UnmarshalComputed calls — its OnlyNested
	// map behavior may overwrite optional map values with API-normalized JSON.
	plannedMeta := data.Meta

	dataBytes, err := data.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}
	res := new(http.Response)
	_, err = r.client.DNS.Zones.New(
		ctx,
		dns.ZoneNewParams{},
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
	data.Meta = plannedMeta

	// Toggle DNSSEC if the user explicitly enabled it
	if !data.DnssecEnabled.IsNull() && !data.DnssecEnabled.IsUnknown() && data.DnssecEnabled.ValueBool() {
		_, err = r.client.DNS.Zones.Dnssec.Update(
			ctx,
			data.Name.ValueString(),
			dns.ZoneDnssecUpdateParams{Enabled: param.NewOpt(true)},
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to enable DNSSEC", err.Error())
			return
		}
	}

	// The POST response returns only {id, warnings}. Do a follow-up GET to
	// populate all computed fields (serial, status, records, rrsets_amount, etc.).
	readRes := new(http.Response)
	_, err = r.client.DNS.Zones.Get(
		ctx,
		data.Name.ValueString(),
		option.WithResponseBodyInto(&readRes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to read zone after create", err.Error())
		return
	}
	readBytes, _ := io.ReadAll(readRes.Body)
	err = apijson.UnmarshalComputed(readBytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize read response", err.Error())
		return
	}
	data.Meta = plannedMeta
	data.Records = sortRecords(ctx, data.Records)

	data.ID = data.Name

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DNSZoneModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *DNSZoneModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save planned DNSSEC value — it's managed via a separate PATCH /dnssec
	// endpoint, not the zone PUT, so we strip it from the PUT body below.
	plannedDnssec := data.DnssecEnabled
	var err error

	// Only send the PUT if zone-level fields actually changed (skip when only
	// dnssec_enabled changed, since that uses a separate PATCH endpoint).
	zoneFieldsChanged := !data.Contact.Equal(state.Contact) ||
		!data.Expiry.Equal(state.Expiry) ||
		!data.NxTtl.Equal(state.NxTtl) ||
		!data.PrimaryServer.Equal(state.PrimaryServer) ||
		!data.Refresh.Equal(state.Refresh) ||
		!data.Retry.Equal(state.Retry) ||
		!data.Enabled.Equal(state.Enabled) ||
		!metaEqual(data.Meta, state.Meta)

	if zoneFieldsChanged {
		dataBytes, err := data.MarshalJSONForUpdate(*state)
		if err != nil {
			resp.Diagnostics.AddError("failed to serialize http request", err.Error())
			return
		}

		// Strip dnssec_enabled from the PUT body. MarshalForUpdate includes all
		// computed_optional fields for the PUT replace, but DNSSEC must only be
		// toggled via the dedicated PATCH /dnssec endpoint.
		var bodyMap map[string]json.RawMessage
		if err := json.Unmarshal(dataBytes, &bodyMap); err == nil {
			delete(bodyMap, "dnssec_enabled")
			dataBytes, _ = json.Marshal(bodyMap)
		}

		res := new(http.Response)
		_, err = r.client.DNS.Zones.Replace(
			ctx,
			data.Name.ValueString(),
			dns.ZoneReplaceParams{},
			option.WithRequestBody("application/json", dataBytes),
			option.WithResponseBodyInto(&res),
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to make http request", err.Error())
			return
		}
		bytes, _ := io.ReadAll(res.Body)
		plannedMeta := data.Meta
		err = apijson.UnmarshalComputed(bytes, &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
			return
		}
		data.Meta = plannedMeta

		// Restore planned value after UnmarshalComputed (PUT response may not include it).
		data.DnssecEnabled = plannedDnssec
	}

	// Toggle DNSSEC if it changed
	if !plannedDnssec.Equal(state.DnssecEnabled) && !plannedDnssec.IsNull() && !plannedDnssec.IsUnknown() {
		_, err = r.client.DNS.Zones.Dnssec.Update(
			ctx,
			data.Name.ValueString(),
			dns.ZoneDnssecUpdateParams{Enabled: param.NewOpt(data.DnssecEnabled.ValueBool())},
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to toggle DNSSEC", err.Error())
			return
		}
	}

	// The Replace API returns a partial response missing computed-only fields
	// (status, records, rrsets_amount, warnings).
	// Do a follow-up GET and selectively populate only the missing computed fields.
	readRes := new(http.Response)
	_, err = r.client.DNS.Zones.Get(
		ctx,
		data.Name.ValueString(),
		option.WithResponseBodyInto(&readRes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to read zone after update", err.Error())
		return
	}
	var readData DNSZoneModel
	readBytes, _ := io.ReadAll(readRes.Body)
	err = apijson.Unmarshal(readBytes, &readData)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize read response", err.Error())
		return
	}

	// Copy computed-only fields from the GET response that the Replace didn't return
	if data.DnssecEnabled.IsNull() || data.DnssecEnabled.IsUnknown() {
		data.DnssecEnabled = readData.DnssecEnabled
	}
	if data.Status.IsNull() || data.Status.IsUnknown() {
		data.Status = readData.Status
	}
	data.Serial = readData.Serial
	data.Warnings = readData.Warnings
	data.Records = sortRecords(ctx, readData.Records)
	data.RrsetsAmount = readData.RrsetsAmount
	data.ID = data.Name

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DNSZoneModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res := new(http.Response)
	_, err := r.client.DNS.Zones.Get(
		ctx,
		data.Name.ValueString(),
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
	data.Records = sortRecords(ctx, data.Records)
	data.ID = data.Name

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DNSZoneModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DNS.Zones.Delete(
		ctx,
		data.Name.ValueString(),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	data.ID = data.Name

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(DNSZoneModel)

	path := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<name>",
		&path,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Name = types.StringValue(path)

	res := new(http.Response)
	_, err := r.client.DNS.Zones.Get(
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
	err = apijson.Unmarshal(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	data.Records = sortRecords(ctx, data.Records)
	data.ID = data.Name

	// Clear meta after import to avoid drift from server-injected keys
	// (targeting, max_hosts, dnssec_status). The user's config value will
	// take precedence on the next plan/apply cycle.
	data.Meta = nil

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSZoneResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}

// sortRecords sorts the computed records list by name+type for deterministic ordering.
// The DNS API may return records in non-deterministic order which causes
// "inconsistent result after apply" errors when the order differs from prior state.
func sortRecords(ctx context.Context, records customfield.NestedObjectList[DNSZoneRecordsModel]) customfield.NestedObjectList[DNSZoneRecordsModel] {
	if records.IsNull() || records.IsUnknown() {
		return records
	}
	elements := records.ListValue.Elements()
	if len(elements) <= 1 {
		return records
	}
	sort.SliceStable(elements, func(i, j int) bool {
		oi := elements[i].(types.Object)
		oj := elements[j].(types.Object)
		nameI := oi.Attributes()["name"].(types.String).ValueString()
		nameJ := oj.Attributes()["name"].(types.String).ValueString()
		if nameI != nameJ {
			return nameI < nameJ
		}
		typeI := oi.Attributes()["type"].(types.String).ValueString()
		typeJ := oj.Attributes()["type"].(types.String).ValueString()
		return typeI < typeJ
	})
	result, diags := customfield.NewObjectListFromAttributes[DNSZoneRecordsModel](ctx, elements)
	if diags.HasError() {
		return records
	}
	return result
}

func metaEqual(a, b *map[string]customfield.MetaStringValue) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(*a) != len(*b) {
		return false
	}
	for k, v := range *a {
		if bv, ok := (*b)[k]; !ok || !v.Equal(bv) {
			return false
		}
	}
	return true
}
