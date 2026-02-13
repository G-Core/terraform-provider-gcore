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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/logging"
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
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

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
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

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
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

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
