// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_status

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudLoadBalancerStatusDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"loadbalancer_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.Int64Attribute{
				Required: true,
			},
			"region_id": schema.Int64Attribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Description: "UUID of the entity",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the load balancer",
				Computed:    true,
			},
			"operating_status": schema.StringAttribute{
				Description: "Operating status of the entity\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"DEGRADED",
						"DRAINING",
						"ERROR",
						"NO_MONITOR",
						"OFFLINE",
						"ONLINE",
					),
				},
			},
			"provisioning_status": schema.StringAttribute{
				Description: "Provisioning status of the entity\nAvailable values: \"ACTIVE\", \"DELETED\", \"ERROR\", \"PENDING_CREATE\", \"PENDING_DELETE\", \"PENDING_UPDATE\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"ACTIVE",
						"DELETED",
						"ERROR",
						"PENDING_CREATE",
						"PENDING_DELETE",
						"PENDING_UPDATE",
					),
				},
			},
			"listeners": schema.ListNestedAttribute{
				Description: "Listeners of the Load Balancer",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerStatusListenersDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "UUID of the entity",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the load balancer listener",
							Computed:    true,
						},
						"operating_status": schema.StringAttribute{
							Description: "Operating status of the entity\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"DEGRADED",
									"DRAINING",
									"ERROR",
									"NO_MONITOR",
									"OFFLINE",
									"ONLINE",
								),
							},
						},
						"pools": schema.ListNestedAttribute{
							Description: "Pools of the Listeners",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerStatusListenersPoolsDataSourceModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "UUID of the entity",
										Computed:    true,
									},
									"members": schema.ListNestedAttribute{
										Description: "Members (servers) of the pool",
										Computed:    true,
										CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerStatusListenersPoolsMembersDataSourceModel](ctx),
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Description: "UUID of the entity",
													Computed:    true,
												},
												"address": schema.StringAttribute{
													Description: "Address of the member (server)",
													Computed:    true,
												},
												"operating_status": schema.StringAttribute{
													Description: "Operating status of the entity\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
													Computed:    true,
													Validators: []validator.String{
														stringvalidator.OneOfCaseInsensitive(
															"DEGRADED",
															"DRAINING",
															"ERROR",
															"NO_MONITOR",
															"OFFLINE",
															"ONLINE",
														),
													},
												},
												"protocol_port": schema.Int64Attribute{
													Description: "Port of the member (server)",
													Computed:    true,
												},
												"provisioning_status": schema.StringAttribute{
													Description: "Provisioning status of the entity\nAvailable values: \"ACTIVE\", \"DELETED\", \"ERROR\", \"PENDING_CREATE\", \"PENDING_DELETE\", \"PENDING_UPDATE\".",
													Computed:    true,
													Validators: []validator.String{
														stringvalidator.OneOfCaseInsensitive(
															"ACTIVE",
															"DELETED",
															"ERROR",
															"PENDING_CREATE",
															"PENDING_DELETE",
															"PENDING_UPDATE",
														),
													},
												},
											},
										},
									},
									"name": schema.StringAttribute{
										Description: "Name of the load balancer pool",
										Computed:    true,
									},
									"operating_status": schema.StringAttribute{
										Description: "Operating status of the entity\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
										Computed:    true,
										Validators: []validator.String{
											stringvalidator.OneOfCaseInsensitive(
												"DEGRADED",
												"DRAINING",
												"ERROR",
												"NO_MONITOR",
												"OFFLINE",
												"ONLINE",
											),
										},
									},
									"provisioning_status": schema.StringAttribute{
										Description: "Provisioning status of the entity\nAvailable values: \"ACTIVE\", \"DELETED\", \"ERROR\", \"PENDING_CREATE\", \"PENDING_DELETE\", \"PENDING_UPDATE\".",
										Computed:    true,
										Validators: []validator.String{
											stringvalidator.OneOfCaseInsensitive(
												"ACTIVE",
												"DELETED",
												"ERROR",
												"PENDING_CREATE",
												"PENDING_DELETE",
												"PENDING_UPDATE",
											),
										},
									},
									"health_monitor": schema.SingleNestedAttribute{
										Description: "Health Monitor of the Pool",
										Computed:    true,
										CustomType:  customfield.NewNestedObjectType[CloudLoadBalancerStatusListenersPoolsHealthMonitorDataSourceModel](ctx),
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description: "UUID of the entity",
												Computed:    true,
											},
											"operating_status": schema.StringAttribute{
												Description: "Operating status of the entity\nAvailable values: \"DEGRADED\", \"DRAINING\", \"ERROR\", \"NO_MONITOR\", \"OFFLINE\", \"ONLINE\".",
												Computed:    true,
												Validators: []validator.String{
													stringvalidator.OneOfCaseInsensitive(
														"DEGRADED",
														"DRAINING",
														"ERROR",
														"NO_MONITOR",
														"OFFLINE",
														"ONLINE",
													),
												},
											},
											"provisioning_status": schema.StringAttribute{
												Description: "Provisioning status of the entity\nAvailable values: \"ACTIVE\", \"DELETED\", \"ERROR\", \"PENDING_CREATE\", \"PENDING_DELETE\", \"PENDING_UPDATE\".",
												Computed:    true,
												Validators: []validator.String{
													stringvalidator.OneOfCaseInsensitive(
														"ACTIVE",
														"DELETED",
														"ERROR",
														"PENDING_CREATE",
														"PENDING_DELETE",
														"PENDING_UPDATE",
													),
												},
											},
											"type": schema.StringAttribute{
												Description: "Type of the Health Monitor\nAvailable values: \"HTTP\", \"HTTPS\", \"K8S\", \"PING\", \"TCP\", \"TLS-HELLO\", \"UDP-CONNECT\".",
												Computed:    true,
												Validators: []validator.String{
													stringvalidator.OneOfCaseInsensitive(
														"HTTP",
														"HTTPS",
														"K8S",
														"PING",
														"TCP",
														"TLS-HELLO",
														"UDP-CONNECT",
													),
												},
											},
										},
									},
								},
							},
						},
						"provisioning_status": schema.StringAttribute{
							Description: "Provisioning status of the entity\nAvailable values: \"ACTIVE\", \"DELETED\", \"ERROR\", \"PENDING_CREATE\", \"PENDING_DELETE\", \"PENDING_UPDATE\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"ACTIVE",
									"DELETED",
									"ERROR",
									"PENDING_CREATE",
									"PENDING_DELETE",
									"PENDING_UPDATE",
								),
							},
						},
					},
				},
			},
			"tags": schema.ListNestedAttribute{
				Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerStatusTagsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "Tag key. The maximum size for a key is 255 bytes.",
							Computed:    true,
						},
						"read_only": schema.BoolAttribute{
							Description: "If true, the tag is read-only and cannot be modified by the user",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "Tag value. The maximum size for a value is 1024 bytes.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudLoadBalancerStatusDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudLoadBalancerStatusDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
