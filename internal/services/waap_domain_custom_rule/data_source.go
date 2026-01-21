// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_custom_rule

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

type WaapDomainCustomRuleDataSource struct {
	client *gcore.Client
}

var _ datasource.DataSourceWithConfigure = (*WaapDomainCustomRuleDataSource)(nil)

func NewWaapDomainCustomRuleDataSource() datasource.DataSource {
	return &WaapDomainCustomRuleDataSource{}
}

func (d *WaapDomainCustomRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_waap_domain_custom_rule"
}

func (d *WaapDomainCustomRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WaapDomainCustomRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *WaapDomainCustomRuleDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.FindOneBy != nil {
		params, diags := data.toListParams(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		env := WaapDomainCustomRulesResultsListDataSourceEnvelope{}
		page, err := d.client.Waap.Domains.CustomRules.List(
			ctx,
			data.DomainID.ValueInt64(),
			params,
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to make http request", err.Error())
			return
		}

		bytes := []byte(page.RawJSON())
		err = apijson.UnmarshalComputed(bytes, &env)
		if err != nil {
			resp.Diagnostics.AddError("failed to unmarshal http request", err.Error())
			return
		}

		if count := len(env.Results.Elements()); count != 1 {
			resp.Diagnostics.AddError("failed to find exactly one result", fmt.Sprint(count)+" found")
			return
		}
		ts, diags := env.Results.AsStructSliceT(ctx)
		resp.Diagnostics.Append(diags...)
		data.RuleID = ts[0].ID
	}

	params, diags := data.toReadParams(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res := new(http.Response)
	_, err := d.client.Waap.Domains.CustomRules.Get(
		ctx,
		data.RuleID.ValueInt64(),
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
	data.ID = data.RuleID

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
