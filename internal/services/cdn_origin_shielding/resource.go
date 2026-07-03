// Custom code: origin shielding has no create/delete endpoints in the API, only
// get + replace (PUT /cdn/resources/{resource_id}/shielding_v2). Stainless therefore
// generates a data source only. This hand-written resource maps the Terraform CRUD
// lifecycle onto get/replace, mirroring the legacy provider's gcore_cdn_origin_shielding:
//   - create/update -> PUT { shielding_pop: <id> }
//   - delete        -> PUT { shielding_pop: null }  (disables shielding)

package cdn_origin_shielding

import (
	"context"
	"fmt"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CDNOriginShieldingResource)(nil)
var _ resource.ResourceWithImportState = (*CDNOriginShieldingResource)(nil)

func NewResource() resource.Resource {
	return &CDNOriginShieldingResource{}
}

// CDNOriginShieldingResource defines the resource implementation.
type CDNOriginShieldingResource struct {
	client *gcore.Client
}

func (r *CDNOriginShieldingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_origin_shielding"
}

func (r *CDNOriginShieldingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// apply changes origin shielding for a CDN resource. Pass a set option to enable
// shielding at a location, or param.Null to disable it.
func (r *CDNOriginShieldingResource) apply(ctx context.Context, resourceID int64, pop param.Opt[int64]) error {
	_, err := r.client.CDN.CDNResources.Shield.Replace(
		ctx,
		resourceID,
		cdn.CDNResourceShieldReplaceParams{
			OriginShielding: cdn.OriginShieldingParam{ShieldingPop: pop},
		},
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	return err
}

func (r *CDNOriginShieldingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CDNOriginShieldingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.apply(ctx, data.ResourceID.ValueInt64(), param.NewOpt(data.ShieldingPop.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CDNOriginShieldingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CDNOriginShieldingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.apply(ctx, data.ResourceID.ValueInt64(), param.NewOpt(data.ShieldingPop.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CDNOriginShieldingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CDNOriginShieldingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CDN.CDNResources.Shield.Get(
		ctx,
		data.ResourceID.ValueInt64(),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	// A null (0) shielding_pop means origin shielding is disabled for this CDN
	// resource, so the managed resource no longer exists.
	if result.ShieldingPop == 0 {
		resp.Diagnostics.AddWarning("Origin shielding disabled", "Origin shielding is no longer enabled for this CDN resource; removing from state.")
		resp.State.RemoveResource(ctx)
		return
	}

	data.ShieldingPop = types.Int64Value(result.ShieldingPop)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CDNOriginShieldingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CDNOriginShieldingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Disabling shielding is a PUT with a null shielding_pop.
	err := r.apply(ctx, data.ResourceID.ValueInt64(), param.Null[int64]())
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
}

func (r *CDNOriginShieldingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CDNOriginShieldingModel)

	resourceID := int64(0)
	diags := importpath.ParseImportID(
		req.ID,
		"<resource_id>",
		&resourceID,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ResourceID = types.Int64Value(resourceID)

	result, err := r.client.CDN.CDNResources.Shield.Get(
		ctx,
		resourceID,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	data.ShieldingPop = types.Int64Value(result.ShieldingPop)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
