// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_rrset

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*DNSZoneRrsetDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rrset_name": schema.StringAttribute{
				Required: true,
			},
			"rrset_type": schema.StringAttribute{
				Required: true,
			},
			"zone_name": schema.StringAttribute{
				Required: true,
			},
			"filter_set_id": schema.Int64Attribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"ttl": schema.Int64Attribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Description: "RRSet type\nAvailable values: \"A\", \"AAAA\", \"NS\", \"CNAME\", \"MX\", \"TXT\", \"SRV\", \"SOA\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"A",
						"AAAA",
						"NS",
						"CNAME",
						"MX",
						"TXT",
						"SRV",
						"SOA",
					),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp marshals/unmarshals date and time as timestamp in json",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"warning": schema.StringAttribute{
				Description: "Warning about some possible side effects without strictly disallowing operations on rrset\nreadonly\nDeprecated: use Warnings instead",
				Computed:    true,
			},
			"meta": schema.MapAttribute{
				Description: "Meta information for rrset. Map with string key and any valid json as value, with valid keys\n1. `failover` (object, beta feature, might be changed in the future) can have fields\n1.1. `protocol` (string, required, HTTP, TCP, UDP, ICMP)\n1.2. `port` (int, required, 1-65535)\n1.3. `frequency` (int, required, in seconds 10-3600)\n1.4. `timeout` (int, required, in seconds 1-10),\n1.5. `method` (string, only for protocol=HTTP)\n1.6. `command` (string, bytes to be sent only for protocol=TCP/UDP)\n1.7. `url` (string, only for protocol=HTTP)\n1.8. `tls` (bool, only for protocol=HTTP)\n1.9. `regexp` (string regex to match, only for non-ICMP)\n1.10. `http_status_code` (int, only for protocol=HTTP)\n1.11. `host` (string, only for protocol=HTTP)\n2. `geodns_link` (string) - name of the geodns link to use, if previously set, must re-send when updating or\nCDN integration will be removed for this RRSet",
				Computed:    true,
				CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
				ElementType: jsontypes.NormalizedType{},
			},
			"pickers": schema.ListNestedAttribute{
				Description: "Set of pickers",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[DNSZoneRrsetPickersDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "Filter type\nAvailable values: \"geodns\", \"asn\", \"country\", \"continent\", \"region\", \"ip\", \"geodistance\", \"weighted_shuffle\", \"default\", \"first_n\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"geodns",
									"asn",
									"country",
									"continent",
									"region",
									"ip",
									"geodistance",
									"weighted_shuffle",
									"default",
									"first_n",
								),
							},
						},
						"limit": schema.Int64Attribute{
							Description: "Limits the number of records returned by the filter\nCan be a positive value for a specific limit. Use zero or leave it blank to indicate no limits.",
							Computed:    true,
						},
						"strict": schema.BoolAttribute{
							Description: "if strict=false, then the filter will return all records if no records match the filter",
							Computed:    true,
						},
					},
				},
			},
			"resource_records": schema.ListNestedAttribute{
				Description: "List of resource record from rrset",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[DNSZoneRrsetResourceRecordsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"content": schema.ListAttribute{
							Description: "Content of resource record\nThe exact length of the array depends on the type of rrset,\neach individual record parameter must be a separate element of the array. For example\n+ SRV-record: `[100, 1, 5061, \"example.com\"]`\n+ CNAME-record: `[ \"the.target.domain\" ]`\n+ A-record: `[ \"1.2.3.4\", \"5.6.7.8\" ]`\n+ AAAA-record: `[ \"2001:db8::1\", \"2001:db8::2\" ]`\n+ MX-record: `[ \"mail1.example.com\", \"mail2.example.com\" ]`\n+ SVCB/HTTPS-record: `[ 1, \".\", [\"alpn\", \"h3\", \"h2\"], [ \"port\", 1443 ], [ \"ipv4hint\", \"10.0.0.1\" ], [ \"ech\", \"AEn+DQBFKwAgACABWIHUGj4u+PIggYXcR5JF0gYk3dCRioBW8uJq9H4mKAAIAAEAAQABAANAEnB1YmxpYy50bHMtZWNoLmRldgAA\" ] ]`",
							Computed:    true,
							CustomType:  customfield.NewListType[jsontypes.Normalized](ctx),
							ElementType: jsontypes.NormalizedType{},
						},
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"enabled": schema.BoolAttribute{
							Computed: true,
						},
						"meta": schema.MapAttribute{
							Description: "Meta information for record\nMap with string key and any valid json as value, with valid keys\n1. `asn` (array of int)\n2. `continents` (array of string)\n3. `countries` (array of string)\n4. `latlong` (array of float64, latitude and longitude)\n5. `fallback` (bool)\n6. `backup` (bool)\n7. `notes` (string)\n8. `weight` (float)\n9. `ip` (string)\nSome keys are reserved for balancing, @see https://api.gcore.com/dns/v2/info/meta\nThis meta will be used to decide which resource record should pass\nthrough filters from the filter set",
							Computed:    true,
							CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
							ElementType: jsontypes.NormalizedType{},
						},
					},
				},
			},
			"warnings": schema.ListNestedAttribute{
				Description: "Warning about some possible side effects without strictly disallowing operations on rrset\nreadonly",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[DNSZoneRrsetWarningsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed: true,
						},
						"message": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *DNSZoneRrsetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *DNSZoneRrsetDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
