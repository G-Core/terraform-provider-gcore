// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_secret

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/logging"
)

type CloudSecretDataSource struct {
	client *gcore.Client
}

var _ datasource.DataSourceWithConfigure = (*CloudSecretDataSource)(nil)

func NewCloudSecretDataSource() datasource.DataSource {
	return &CloudSecretDataSource{}
}

func (d *CloudSecretDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_secret"
}

func (d *CloudSecretDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudSecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CloudSecretDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params, diags := data.toReadParams(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res := new(http.Response)
	_, err := d.client.Cloud.Secrets.Get(
		ctx,
		data.SecretID.ValueString(),
		params,
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
