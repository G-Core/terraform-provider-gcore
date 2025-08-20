// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_billing_reservation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudBillingReservationsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"activated_from": schema.StringAttribute{
				Description: "Lower bound, starting from what date the reservation was/will be activated",
				Optional:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"activated_to": schema.StringAttribute{
				Description: "High bound, before what date the reservation was/will be activated",
				Optional:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"created_from": schema.StringAttribute{
				Description: "Lower bound the filter, showing result(s) equal to or greater than date the reservation was created",
				Optional:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"created_to": schema.StringAttribute{
				Description: "High bound the filter, showing result(s) equal to or less date the reservation was created",
				Optional:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"deactivated_from": schema.StringAttribute{
				Description: "Lower bound, starting from what date the reservation was/will be deactivated",
				Optional:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"deactivated_to": schema.StringAttribute{
				Description: "High bound, before what date the reservation was/will be deactivated",
				Optional:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"metric_name": schema.StringAttribute{
				Description: "Name from billing features for specific resource",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region for reservation",
				Optional:    true,
			},
			"status": schema.ListAttribute{
				Description: "Field for fixed a status by reservation workflow",
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOfCaseInsensitive(
							"ACTIVATED",
							"APPROVED",
							"COPIED",
							"CREATED",
							"EXPIRED",
							"REJECTED",
							"RESERVED",
							"WAITING_FOR_PAYMENT",
						),
					),
				},
				ElementType: types.StringType,
			},
			"limit": schema.Int64Attribute{
				Description: "Limit of reservation list page",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"order_by": schema.StringAttribute{
				Description: "Order by field and direction.\nAvailable values: \"active_from.asc\", \"active_from.desc\", \"active_to.asc\", \"active_to.desc\", \"created_at.asc\", \"created_at.desc\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"active_from.asc",
						"active_from.desc",
						"active_to.asc",
						"active_to.desc",
						"created_at.asc",
						"created_at.desc",
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
				CustomType:  customfield.NewNestedObjectListType[CloudBillingReservationsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Reservation id",
							Computed:    true,
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
						"amount_prices": schema.SingleNestedAttribute{
							Description: "Reservation amount prices",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudBillingReservationsAmountPricesDataSourceModel](ctx),
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
						"is_expiration_message_visible": schema.BoolAttribute{
							Description: "Hide or show expiration message to customer.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Reservation name",
							Computed:    true,
						},
						"next_statuses": schema.ListAttribute{
							Description: "List of possible next reservation statuses",
							Computed:    true,
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
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
						"resources": schema.ListNestedAttribute{
							Description: "List of reservation resources",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudBillingReservationsResourcesDataSourceModel](ctx),
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
						"status": schema.StringAttribute{
							Description: "Reservation status",
							Computed:    true,
						},
						"user_status": schema.StringAttribute{
							Description: "User status",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudBillingReservationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudBillingReservationsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
