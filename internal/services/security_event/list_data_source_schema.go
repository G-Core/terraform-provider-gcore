// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package security_event

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*SecurityEventsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"alert_type": schema.StringAttribute{
				Description: `Available values: "ddos_alert", "rtbh_alert".`,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("ddos_alert", "rtbh_alert"),
				},
			},
			"targeted_ip_addresses": schema.StringAttribute{
				Optional: true,
			},
			"limit": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 500),
				},
			},
			"ordering": schema.StringAttribute{
				Description: `Available values: "attack_start_time", "-attack_start_time", "attack_power_bps", "-attack_power_bps", "attack_power_pps", "-attack_power_pps", "number_of_ip_involved_in_attack", "-number_of_ip_involved_in_attack", "alert_type", "-alert_type".`,
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"attack_start_time",
						"-attack_start_time",
						"attack_power_bps",
						"-attack_power_bps",
						"attack_power_pps",
						"-attack_power_pps",
						"number_of_ip_involved_in_attack",
						"-number_of_ip_involved_in_attack",
						"alert_type",
						"-alert_type",
					),
				},
			},
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[SecurityEventsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"alert_type": schema.StringAttribute{
							Description: `Available values: "ddos_alert", "rtbh_alert".`,
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("ddos_alert", "rtbh_alert"),
							},
						},
						"attack_power_bps": schema.Float64Attribute{
							Computed: true,
						},
						"attack_power_pps": schema.Float64Attribute{
							Computed: true,
						},
						"attack_start_time": schema.StringAttribute{
							Computed:   true,
							CustomType: timetypes.RFC3339Type{},
						},
						"client_id": schema.Int64Attribute{
							Computed: true,
						},
						"notification_type": schema.StringAttribute{
							Computed: true,
						},
						"number_of_ip_involved_in_attack": schema.Int64Attribute{
							Computed: true,
						},
						"targeted_ip_addresses": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *SecurityEventsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *SecurityEventsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
