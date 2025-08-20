// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_dnssec

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*DNSZoneDnssecResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"enabled": schema.BoolAttribute{
				Optional: true,
			},
			"algorithm": schema.StringAttribute{
				Description: "Specifies the algorithm used for the key.",
				Computed:    true,
			},
			"digest": schema.StringAttribute{
				Description: "Represents the hashed value of the DS record.",
				Computed:    true,
			},
			"digest_algorithm": schema.StringAttribute{
				Description: "Specifies the algorithm used to generate the digest.",
				Computed:    true,
			},
			"digest_type": schema.StringAttribute{
				Description: "Specifies the type of the digest algorithm used.",
				Computed:    true,
			},
			"ds": schema.StringAttribute{
				Description: "Represents the complete DS record.",
				Computed:    true,
			},
			"flags": schema.Int64Attribute{
				Description: "Represents the flag for DNSSEC record.",
				Computed:    true,
			},
			"key_tag": schema.Int64Attribute{
				Description: "Represents the identifier of the DNSKEY record.",
				Computed:    true,
			},
			"key_type": schema.StringAttribute{
				Description: "Specifies the type of the key used in the algorithm.",
				Computed:    true,
			},
			"message": schema.StringAttribute{
				Computed: true,
			},
			"public_key": schema.StringAttribute{
				Description: "Represents the public key used in the DS record.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *DNSZoneDnssecResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *DNSZoneDnssecResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
