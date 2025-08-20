// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_billing_reservation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudBillingReservationDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"reservation_id": schema.Int64Attribute{
				Description: "ID of the reservation",
				Required:    true,
			},
			"active_from": schema.StringAttribute{
				Description: "Reservation active from date",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"active_to": schema.StringAttribute{
				Description: "Reservation active to date",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"activity_period": schema.StringAttribute{
				Description: "Name of the billing period, e.g month",
				Computed:    true,
			},
			"activity_period_length": schema.Int64Attribute{
				Description: "Length of the full reservation period by `activity_period`",
				Computed:    true,
			},
			"billing_plan_id": schema.Int64Attribute{
				Description: "Billing plan id",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Reservation creation date",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"error": schema.StringAttribute{
				Description: "Error message if any occured during reservation",
				Computed:    true,
			},
			"eta": schema.StringAttribute{
				Description: "ETA delivery if bare metal out of stock. Value None means that bare metal in stock.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"id": schema.Int64Attribute{
				Description: "Reservation id",
				Computed:    true,
			},
			"is_expiration_message_visible": schema.BoolAttribute{
				Description: "Hide or show expiration message to customer.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Reservation name",
				Computed:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region id",
				Computed:    true,
			},
			"region_name": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"remind_expiration_message": schema.StringAttribute{
				Description: "The date when show expiration date to customer",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"status": schema.StringAttribute{
				Description: "Reservation status",
				Computed:    true,
			},
			"user_status": schema.StringAttribute{
				Description: "User status",
				Computed:    true,
			},
			"next_statuses": schema.ListAttribute{
				Description: "List of possible next reservation statuses",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"amount_prices": schema.SingleNestedAttribute{
				Description: "Reservation amount prices",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudBillingReservationAmountPricesDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"commit_price_per_month": schema.StringAttribute{
						Description: "Commit price of the item charged per month",
						Computed:    true,
					},
					"commit_price_per_unit": schema.StringAttribute{
						Description: "Commit price of the item charged per hour",
						Computed:    true,
					},
					"commit_price_total": schema.StringAttribute{
						Description: "Commit price of the item charged for all period reservation",
						Computed:    true,
					},
					"currency_code": schema.StringAttribute{
						Description: "Currency code (3 letter code per ISO 4217)",
						Computed:    true,
					},
					"overcommit_price_per_month": schema.StringAttribute{
						Description: "Overcommit price of the item charged per month",
						Computed:    true,
					},
					"overcommit_price_per_unit": schema.StringAttribute{
						Description: "Overcommit price of the item charged per hour",
						Computed:    true,
					},
					"overcommit_price_total": schema.StringAttribute{
						Description: "Overcommit price of the item charged for all period reservation",
						Computed:    true,
					},
				},
			},
			"resources": schema.ListNestedAttribute{
				Description: "List of reservation resources",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudBillingReservationResourcesDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"activity_period": schema.StringAttribute{
							Description: "Name of the billing period, e.g month",
							Computed:    true,
						},
						"activity_period_length": schema.Int64Attribute{
							Description: "Length of the full reservation period by `activity_period`",
							Computed:    true,
						},
						"billing_plan_item_id": schema.Int64Attribute{
							Description: "Billing plan item id",
							Computed:    true,
						},
						"commit_price_per_month": schema.StringAttribute{
							Description: "Commit price of the item charged per month",
							Computed:    true,
						},
						"commit_price_per_unit": schema.StringAttribute{
							Description: "Commit price of the item charged per hour",
							Computed:    true,
						},
						"commit_price_total": schema.StringAttribute{
							Description: "Commit price of the item charged for all period reservation",
							Computed:    true,
						},
						"overcommit_billing_plan_item_id": schema.Int64Attribute{
							Description: "Overcommit billing plan item id",
							Computed:    true,
						},
						"overcommit_price_per_month": schema.StringAttribute{
							Description: "Overcommit price of the item charged per month",
							Computed:    true,
						},
						"overcommit_price_per_unit": schema.StringAttribute{
							Description: "Overcommit price of the item charged per hour",
							Computed:    true,
						},
						"overcommit_price_total": schema.StringAttribute{
							Description: "Overcommit price of the item charged for all period reservation",
							Computed:    true,
						},
						"resource_count": schema.Int64Attribute{
							Description: "Number of reserved resource items",
							Computed:    true,
						},
						"resource_name": schema.StringAttribute{
							Description: "Resource name",
							Computed:    true,
						},
						"resource_type": schema.StringAttribute{
							Description: "Resource type\nAvailable values: \"flavor\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("flavor"),
							},
						},
						"unit_name": schema.StringAttribute{
							Description: "Billing unit name",
							Computed:    true,
						},
						"unit_size_month": schema.StringAttribute{
							Description: "Minimal billing size, for example it is 744 hours per 1 month.",
							Computed:    true,
						},
						"unit_size_total": schema.StringAttribute{
							Description: "Unit size month multiplied by count of resources in the reservation",
							Computed:    true,
						},
						"cpu": schema.StringAttribute{
							Description: "Baremetal CPU description",
							Computed:    true,
						},
						"disk": schema.StringAttribute{
							Description: "Baremetal disk description",
							Computed:    true,
						},
						"ram": schema.StringAttribute{
							Description: "Baremetal RAM description",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudBillingReservationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudBillingReservationDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
