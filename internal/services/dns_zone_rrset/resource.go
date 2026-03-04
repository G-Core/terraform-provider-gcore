// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_rrset

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/dns"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*DNSZoneRrsetResource)(nil)
var _ resource.ResourceWithModifyPlan = (*DNSZoneRrsetResource)(nil)
var _ resource.ResourceWithImportState = (*DNSZoneRrsetResource)(nil)

func NewResource() resource.Resource {
	return &DNSZoneRrsetResource{}
}

// DNSZoneRrsetResource defines the resource implementation.
type DNSZoneRrsetResource struct {
	client *gcore.Client
}

func (r *DNSZoneRrsetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zone_rrset"
}

func (r *DNSZoneRrsetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DNSZoneRrsetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DNSZoneRrsetModel

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
	_, err = r.client.DNS.Zones.Rrsets.New(
		ctx,
		data.RrsetType.ValueString(),
		dns.ZoneRrsetNewParams{
			ZoneName:  data.ZoneName.ValueString(),
			RrsetName: data.RrsetName.ValueString(),
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
	// Save planned meta values before UnmarshalComputed — its OnlyNested map
	// behavior may overwrite optional map values with API-normalized JSON
	// (e.g. the API returns [50,30] for [50.0,30.0]).
	plannedData := *data
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	preservePlannedMeta(data, &plannedData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSZoneRrsetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DNSZoneRrsetModel
	var state *DNSZoneRrsetModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
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
	_, err = r.client.DNS.Zones.Rrsets.Replace(
		ctx,
		data.RrsetType.ValueString(),
		dns.ZoneRrsetReplaceParams{
			ZoneName:  data.ZoneName.ValueString(),
			RrsetName: data.RrsetName.ValueString(),
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
	// Save planned meta values before UnmarshalComputed (same reason as Create)
	plannedData := *data
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	preservePlannedMeta(data, &plannedData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSZoneRrsetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DNSZoneRrsetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res := new(http.Response)
	_, err := r.client.DNS.Zones.Rrsets.Get(
		ctx,
		data.RrsetType.ValueString(),
		dns.ZoneRrsetGetParams{
			ZoneName:  data.ZoneName.ValueString(),
			RrsetName: data.RrsetName.ValueString(),
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
	// Save resource_records before UnmarshalComputed. The OnlyNested behavior
	// processes arrays positionally, which cross-contaminates computed fields
	// (id) and optional map values (meta) when the API returns records in a
	// different order than state. Deep-copy because the pointer-to-slice means
	// UnmarshalComputed modifies through the same pointer.
	savedRecords := copyResourceRecords(data.ResourceRecords)
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	data.ResourceRecords = savedRecords

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSZoneRrsetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DNSZoneRrsetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res := new(http.Response)
	_, err := r.client.DNS.Zones.Rrsets.Delete(
		ctx,
		data.RrsetType.ValueString(),
		dns.ZoneRrsetDeleteParams{
			ZoneName:  data.ZoneName.ValueString(),
			RrsetName: data.RrsetName.ValueString(),
		},
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	// If resource already deleted externally, treat as success
	if res != nil && res.StatusCode == 404 {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSZoneRrsetResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}

func (r *DNSZoneRrsetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: zone_name/rrset_name/rrset_type
	// Example: maxima.lt/www.maxima.lt/A
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected format: zone_name/rrset_name/rrset_type, got: %s", req.ID),
		)
		return
	}

	zoneName := parts[0]
	rrsetName := parts[1]
	rrsetType := parts[2]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_name"), types.StringValue(zoneName))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("rrset_name"), types.StringValue(rrsetName))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("rrset_type"), types.StringValue(rrsetType))...)
}

// preservePlannedMeta restores planned meta values after UnmarshalComputed.
// The apijson decoder's OnlyNested behavior for maps overwrites individual
// values via json.Unmarshaler, which can cause the API's normalized JSON
// (e.g. [50,30] for [50.0,30.0]) to replace the planned values, triggering
// Terraform's "inconsistent result after apply" error.
// copyResourceRecords deep-copies the resource_records slice so that
// UnmarshalComputed's in-place modifications don't affect the saved copy.
func copyResourceRecords(src *[]*DNSZoneRrsetResourceRecordsModel) *[]*DNSZoneRrsetResourceRecordsModel {
	if src == nil {
		return nil
	}
	dst := make([]*DNSZoneRrsetResourceRecordsModel, len(*src))
	for i, rec := range *src {
		if rec == nil {
			continue
		}
		cp := *rec
		// Deep-copy meta map
		if rec.Meta != nil {
			m := make(map[string]customfield.MetaStringValue, len(*rec.Meta))
			for k, v := range *rec.Meta {
				m[k] = v
			}
			cp.Meta = &m
		}
		// Deep-copy content slice
		if rec.Content != nil {
			c := make([]jsontypes.Normalized, len(*rec.Content))
			copy(c, *rec.Content)
			cp.Content = &c
		}
		dst[i] = &cp
	}
	return &dst
}

func preservePlannedMeta(data, plan *DNSZoneRrsetModel) {
	data.Meta = plan.Meta
	if data.ResourceRecords != nil && plan.ResourceRecords != nil {
		for i := range *data.ResourceRecords {
			if i < len(*plan.ResourceRecords) {
				(*data.ResourceRecords)[i].Meta = (*plan.ResourceRecords)[i].Meta
			}
		}
	}
}
