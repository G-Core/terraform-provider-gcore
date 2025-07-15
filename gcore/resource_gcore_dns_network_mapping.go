package gcore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"net"

	dnssdk "github.com/G-Core/gcore-dns-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	DNSNetworkMappingResource = "gcore_dns_network_mapping"

	DNSNetworkMappingSchemaName    = "name"
	DNSNetworkMappingSchemaMapping = "mapping"
	DNSNetworkMappingSchemaTags    = "tags"
	DNSNetworkMappingSchemaCidr4   = "cidr4"
	DNSNetworkMappingSchemaCidr6   = "cidr6"
)

func resourceDNSNetworkMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSNetworkMappingCreate,
		ReadContext:   resourceDNSNetworkMappingRead,
		UpdateContext: resourceDNSNetworkMappingUpdate,
		DeleteContext: resourceDNSNetworkMappingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			DNSNetworkMappingSchemaName: {
				Type:     schema.TypeString,
				Required: true,
			},
			DNSNetworkMappingSchemaMapping: {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						DNSNetworkMappingSchemaTags: {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						DNSNetworkMappingSchemaCidr4: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsCIDR,
							},
						},
						DNSNetworkMappingSchemaCidr6: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsCIDR,
							},
						},
					},
				},
			},
		},
	}
}

func resourceDNSNetworkMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).DNSClient

	name := d.Get(DNSNetworkMappingSchemaName).(string)
	mapping := d.Get(DNSNetworkMappingSchemaMapping).([]interface{})

	sdkMapping := dnssdk.NetworkMappingRequest{
		Name:    name,
		Mapping: expandMappingEntries(mapping),
	}

	id, err := client.CreateNetworkMapping(ctx, sdkMapping)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(id))

	return resourceDNSNetworkMappingRead(ctx, d, m)
}

func resourceDNSNetworkMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).DNSClient
	id, err := strconv.ParseUint(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not parse network mapping ID: %w", err))
	}

	mapping, err := client.GetNetworkMapping(ctx, id)
	if err != nil {
		var apiErr *dnssdk.APIError
		if errors.As(err, apiErr) && apiErr.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set(DNSNetworkMappingSchemaName, mapping.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(DNSNetworkMappingSchemaMapping, flattenMappingEntries(mapping.Mapping)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDNSNetworkMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).DNSClient
	id, err := strconv.ParseUint(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not parse network mapping ID: %w", err))
	}

	if d.HasChange(DNSNetworkMappingSchemaName) || d.HasChange(DNSNetworkMappingSchemaMapping) {
		name := d.Get(DNSNetworkMappingSchemaName).(string)
		mapping := d.Get(DNSNetworkMappingSchemaMapping).([]interface{})

		sdkMapping := dnssdk.NetworkMappingRequest{
			Name:    name,
			Mapping: expandMappingEntries(mapping),
		}

		err := client.UpdateNetworkMapping(ctx, id, sdkMapping)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDNSNetworkMappingRead(ctx, d, m)
}

func resourceDNSNetworkMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).DNSClient
	id, err := strconv.ParseUint(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not parse network mapping ID: %w", err))
	}

	err = client.DeleteNetworkMapping(ctx, id)
	if err != nil {
		var apiErr *dnssdk.APIError
		if errors.As(err, apiErr) && apiErr.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}

func expandMappingEntries(mappings []interface{}) []dnssdk.MappingEntry {
	if len(mappings) == 0 {
		return nil
	}

	sdkMappings := make([]dnssdk.MappingEntry, len(mappings))

	for i, m := range mappings {
		mappingData := m.(map[string]interface{})

		var entry dnssdk.MappingEntry

		if tags, ok := mappingData[DNSNetworkMappingSchemaTags].([]interface{}); ok {
			entry.Tags = expandStringList(tags)
		}
		if cidr4, ok := mappingData[DNSNetworkMappingSchemaCidr4].([]interface{}); ok {
			entry.CIDR4 = expandIPNetList(cidr4)
		}
		if cidr6, ok := mappingData[DNSNetworkMappingSchemaCidr6].([]interface{}); ok {
			entry.CIDR6 = expandIPNetList(cidr6)
		}
		sdkMappings[i] = entry
	}
	return sdkMappings
}

func expandStringList(list []interface{}) []string {
	ss := make([]string, len(list))
	for i, v := range list {
		ss[i] = v.(string)
	}
	return ss
}

func expandIPNetList(list []interface{}) []dnssdk.IPNet {
	nets := make([]dnssdk.IPNet, len(list))
	for i, v := range list {
		_, ipnet, _ := net.ParseCIDR(v.(string))
		if ipnet != nil {
			nets[i] = dnssdk.IPNet{IPNet: *ipnet}
		}
	}
	return nets
}

func flattenMappingEntries(entries []dnssdk.MappingEntry) []interface{} {
	if entries == nil {
		return make([]interface{}, 0)
	}

	mappings := make([]interface{}, len(entries))
	for i, entry := range entries {
		m := make(map[string]interface{})
		m[DNSNetworkMappingSchemaTags] = flattenStringList(entry.Tags)
		m[DNSNetworkMappingSchemaCidr4] = flattenIPNetList(entry.CIDR4)
		m[DNSNetworkMappingSchemaCidr6] = flattenIPNetList(entry.CIDR6)
		mappings[i] = m
	}
	return mappings
}

func flattenStringList(list []string) []interface{} {
	vs := make([]interface{}, len(list))
	for i, v := range list {
		vs[i] = v
	}
	return vs
}

func flattenIPNetList(list []dnssdk.IPNet) []interface{} {
	vs := make([]interface{}, len(list))
	for i, v := range list {
		vs[i] = v.String()
	}
	return vs
}
