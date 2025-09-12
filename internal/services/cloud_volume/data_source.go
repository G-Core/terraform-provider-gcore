// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume

import (
	"context"
	"fmt"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/logging"
)

type CloudVolumeDataSource struct {
	client *gcore.Client
}

var _ datasource.DataSourceWithConfigure = (*CloudVolumeDataSource)(nil)

func NewCloudVolumeDataSource() datasource.DataSource {
	return &CloudVolumeDataSource{}
}

func (d *CloudVolumeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_volume"
}

func (d *CloudVolumeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *CloudVolumeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CloudVolumeDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params, diags := data.toReadParams(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	volume, err := d.client.Cloud.Volumes.Get(
		ctx,
		data.VolumeID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	// Use raw JSON from the response to unmarshal the "computed" fields into the data model
	err = apijson.UnmarshalComputed([]byte(volume.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
