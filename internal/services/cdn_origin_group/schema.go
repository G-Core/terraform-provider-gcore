// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_origin_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*CdnOriginGroupResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "Origin group ID.",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description: "Origin group name.",
				Required:    true,
			},
			"auth_type": schema.StringAttribute{
				Description: "Origin authentication type.\n\nPossible values:\n- **none** - Used for public origins.\n- **awsSignatureV4** - Used for S3 storage.",
				Optional:    true,
			},
			"path": schema.StringAttribute{
				Description: "Parameter is **deprecated**.",
				Optional:    true,
			},
			"use_next": schema.BoolAttribute{
				Description: "Defines whether to use the next origin from the origin group if origin responds with the cases specified in `proxy_next_upstream`.\nIf you enable it, you must specify cases in `proxy_next_upstream`.\n\nPossible values:\n- **true** - Option is enabled.\n- **false** - Option is disabled.",
				Optional:    true,
			},
			"auth": schema.SingleNestedAttribute{
				Description: "Credentials to access the private bucket.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"s3_access_key_id": schema.StringAttribute{
						Description: "Access key ID for the S3 account. \n\nRestrictions:\n- Latin letters (A-Z, a-z), numbers (0-9), colon, dash, and underscore.\n- From 3 to 512 characters.",
						Required:    true,
					},
					"s3_bucket_name": schema.StringAttribute{
						Description: "S3 bucket name. \n\nRestrictions:\n- Maximum 128 characters.",
						Required:    true,
					},
					"s3_secret_access_key": schema.StringAttribute{
						Description: "Secret access key for the S3 account. \n\nRestrictions:\n- Latin letters (A-Z, a-z), numbers (0-9), pluses, slashes, dashes, colons and underscores.\n- If \"`s3_type`\": amazon, length should be 40 characters.\n- If \"`s3_type`\": other, length should be from 16 to 255 characters.",
						Required:    true,
					},
					"s3_type": schema.StringAttribute{
						Description: "Storage type compatible with S3.\n\nPossible values:\n- **amazon** – AWS S3 storage.\n- **other** – Other (not AWS) S3 compatible storage.",
						Required:    true,
					},
					"s3_region": schema.StringAttribute{
						Description: "S3 storage region. \n\nThe parameter is required, if \"`s3_type`\": amazon.",
						Optional:    true,
					},
					"s3_storage_hostname": schema.StringAttribute{
						Description: "S3 storage hostname. \n\nThe parameter is required, if \"`s3_type`\": other.",
						Optional:    true,
					},
				},
			},
			"sources": schema.ListNestedAttribute{
				Description: "List of origin sources in the origin group.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"backup": schema.BoolAttribute{
							Description: "Defines whether the origin is a backup, meaning that it will not be used until one of active origins become unavailable.\n\nPossible values:\n- **true** - Origin is a backup.\n- **false** - Origin is not a backup.",
							Optional:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Enables or disables an origin source in the origin group.\n\nPossible values:\n- **true** - Origin is enabled and the CDN uses it to pull content.\n- **false** - Origin is disabled and the CDN does not use it to pull content.\n\nOrigin group must contain at least one enabled origin.",
							Optional:    true,
						},
						"source": schema.StringAttribute{
							Description: "IP address or domain name of the origin and the port, if custom port is used.",
							Optional:    true,
						},
					},
				},
			},
			"proxy_next_upstream": schema.ListAttribute{
				Description: "Defines cases when the request should be passed on to the next origin.\n\nPossible values:\n- **error** - an error occurred while establishing a connection with the origin, passing a request to it, or reading the response header\n- **timeout** - a timeout has occurred while establishing a connection with the origin, passing a request to it, or reading the response header\n- **`invalid_header`** - a origin returned an empty or invalid response\n- **`http_403`** - a origin returned a response with the code 403\n- **`http_404`** - a origin returned a response with the code 404\n- **`http_429`** - a origin returned a response with the code 429\n- **`http_500`** - a origin returned a response with the code 500\n- **`http_502`** - a origin returned a response with the code 502\n- **`http_503`** - a origin returned a response with the code 503\n- **`http_504`** - a origin returned a response with the code 504",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"has_related_resources": schema.BoolAttribute{
				Description: "Defines whether the origin group has related CDN resources.\n\nPossible values:\n- **true** - Origin group has related CDN resources.\n- **false** - Origin group does not have related CDN resources.",
				Computed:    true,
			},
		},
	}
}

func (r *CdnOriginGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CdnOriginGroupResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
