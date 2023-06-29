package gcore

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"

	gcdn "github.com/G-Core/gcorelabscdn-go/gcore"
	"github.com/G-Core/gcorelabscdn-go/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	optionsSchema = &schema.Schema{
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Computed:    true,
		Description: "Each option in CDN resource settings. Each option added to CDN resource settings should have the following mandatory request fields: enabled, value.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allowed_http_methods": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Computed:    true,
					Description: "The list of allowed HTTP methods.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Required:    true,
								Description: "Available methods: GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS.",
							},
						},
					},
				},
				"brotli_compression": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "The option allows to compress content with brotli on the CDN's end.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Required:    true,
								Description: "Specify the content-type for each type of content you wish to have compressed.",
							},
						},
					},
				},
				"browser_cache_settings": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Computed:    true,
					Description: "The cache expiration time for customers' browsers in seconds.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Use '0s' to disable caching.",
							},
						},
					},
				},
				"cache_http_headers": { // deprecated in favor of response_headers_hiding_policy
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Computed:    true,
					Description: "Legacy option. Use the response_headers_hiding_policy option instead.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Required:    true,
								Description: "List HTTP Headers that must be included in the response.",
							},
						},
					},
				},
				"cors": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "The option adds the Access-Control-Allow-Origin header to responses from CDN servers.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Required:    true,
								Description: "Specifies a value of the Access-Control-Allow-Origin header.",
							},
						},
					},
				},
				"country_acl": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Control access to the content for specified countries.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"policy_type": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Possible values: allow, deny.",
							},
							"excepted_values": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Required:    true,
								Description: "List of countries according to ISO-3166-1.",
							},
						},
					},
				},
				"disable_cache": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "When enabled the content caching is completely disabled.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"disable_proxy_force_ranges": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "It allows getting 206 responses regardless settings of an origin source.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"edge_cache_settings": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Computed:    true,
					Description: "The cache expiration time for CDN servers.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Caching time for a response with codes 200, 206, 301, 302. Responses with codes 4xx, 5xx will not be cached. Use '0s' disable to caching. Use custom_values field to specify a custom caching time for a response with specific codes.",
							},
							"default": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Content will be cached according to origin cache settings. The value applies for a response with codes 200, 201, 204, 206, 301, 302, 303, 304, 307, 308 if an origin server does not have caching HTTP headers. Responses with other codes will not be cached.",
							},
							"custom_values": {
								Type:     schema.TypeMap,
								Optional: true,
								Computed: true,
								DefaultFunc: func() (interface{}, error) {
									return map[string]interface{}{}, nil
								},
								Elem:        schema.TypeString,
								Description: "Caching time for a response with specific codes. These settings have a higher priority than the value field. Response code ('304', '404' for example). Use 'any' to specify caching time for all response codes. Caching time in seconds ('0s', '600s' for example). Use '0s' to disable caching for a specific response code.",
							},
						},
					},
				},
				"fetch_compressed": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "A CDN request and cache already compressed content. Your server should support compression. CDN servers won't ungzip your content even if a user's browser doesn't accept compression (nowadays almost all browsers support it).",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"follow_origin_redirect": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Redirection from origin. If the origin server returns a redirect, the option allows a CDN to pull the requested content from an origin server that was returned in the redirect.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"codes": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeInt},
								Required:    true,
								Description: "Specify the redirect status code that the origin server returns. Possible values: 301, 302, 303, 307, 308.",
							},
						},
					},
				},
				"force_return": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Allows to apply custom HTTP code to the CDN content.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"code": {
								Type:        schema.TypeInt,
								Required:    true,
								Description: "HTTP response status code. Available codes: 100 <= value <= 599. Reserved codes: 408, 444, 477, 494, 495, 496, 497, 499",
							},
							"body": {
								Type:        schema.TypeString,
								Optional:    true,
								Default:     "",
								Description: "Response text or URL if you're going to set up redirection. Max length = 100.",
							},
						},
					},
				},
				"forward_host_header": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "When a CDN requests content from an origin server the option allows forwarding the Host header used in the request made to a CDN.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"gzip_on": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "The option allows to compress content with gzip on the CDN`s end. CDN servers will request only uncompressed content from the origin.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"host_header": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Specify the Host header that CDN servers use when request content from an origin server. Your server must be able to process requests with the chosen header. If the option is in NULL state Host Header value is taken from the CNAME field.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
				"http3_enabled": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Use HTTP/3 protocol for content delivery when supported by the end users browser.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"ignore_cookie": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "By default, files pulled from an origin source with cookies are not cached in a CDN. Enable this option to cache such objects.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"ignore_query_string": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "This option determines how files with different query strings will be cached: either as one object (when this option is enabled) or as different objects (when this option is disabled).",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"image_stack": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "The option allows transforming JPG and PNG images (such as resizing or cropping) and automatically converting them to WebP or AVIF format. It is a paid option.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"avif_enabled": {
								Type:        schema.TypeBool,
								Optional:    true,
								Computed:    true,
								Description: "If enabled, JPG and PNG images automatically convert to AVIF format when supported by the end users browser.",
							},
							"webp_enabled": {
								Type:        schema.TypeBool,
								Optional:    true,
								Computed:    true,
								Description: "If enabled, JPG and PNG images automatically convert to WebP format when supported by the end users browser.",
							},
							"quality": {
								Type:        schema.TypeInt,
								Required:    true,
								Description: "Quality settings for JPG and PNG images. Specify a value from 1 to 100.",
							},
							"png_lossless": {
								Type:        schema.TypeBool,
								Optional:    true,
								Computed:    true,
								Description: "Represents compression without quality loss for PNG format.",
							},
						},
					},
				},
				"ip_address_acl": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Control access to the CDN Resource content for specified IP addresses.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"policy_type": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Possible values: allow, deny.",
							},
							"excepted_values": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Required:    true,
								Description: "List of IPs.",
							},
						},
					},
				},
				"limit_bandwidth": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "The option allows to control the download speed per connection.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"limit_type": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The way of controlling the download speed per each connection. Possible values are: static, dynamic.",
							},
							"speed": {
								Type:        schema.TypeInt,
								Optional:    true,
								Computed:    true,
								Description: "Maximum download speed per connection. Must be greater than 0.",
							},
							"buffer": {
								Type:        schema.TypeInt,
								Optional:    true,
								Computed:    true,
								Description: "Amount of downloaded data after which the user will be rate limited.",
							},
						},
					},
				},
				"proxy_cache_methods_set": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Allows caching for GET, HEAD and POST requests.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"query_params_blacklist": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Specify list of query strings. Files with those query strings will be cached as one object.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeSet,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Required: true,
							},
						},
					},
				},
				"query_params_whitelist": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Specify list of query strings. Files with those query strings will be cached as different objects.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeSet,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Required: true,
							},
						},
					},
				},
				"redirect_https_to_http": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "When enabled redirects HTTPS requests to HTTP.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"redirect_http_to_https": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Sets redirect from HTTP protocol to HTTPS for all resource requests.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"referrer_acl": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Control access to the CDN Resource content for specified domain names.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"policy_type": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Possible values: allow, deny.",
							},
							"excepted_values": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Required:    true,
								Description: "List of domain names.",
							},
						},
					},
				},
				"request_limiter": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "It allows to limit the amount of HTTP requests",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"rate": {
								Type:     schema.TypeInt,
								Required: true,
							},
							"burst": {
								Type:     schema.TypeInt,
								Required: true,
							},
							"rate_unit": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  "r/s",
							},
							"delay": {
								Type:     schema.TypeInt,
								Optional: true,
								Default:  0,
							},
						},
					},
				},
				"response_headers_hiding_policy": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Define HTTP headers (specified at an origin server) that a CDN server hides from the response.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"mode": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies a mode of hiding HTTP headers from the response. Possible values are: hide, show.",
							},
							"excepted": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Required:    true,
								Description: "List of HTTP headers.",
							},
						},
					},
				},
				"rewrite": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "The option allows changing or redirecting query paths.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"body": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The pattern for Rewrite.",
							},
							"flag": {
								Type:        schema.TypeString,
								Optional:    true,
								Default:     "break",
								Description: "Defines flag for the Rewrite option.",
							},
						},
					},
				},
				"secure_key": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "The option allows configuring an access with tokenized URLs.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"key": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The key generated on your side which will be used for URL signing.",
							},
							"type": {
								Type:        schema.TypeInt,
								Required:    true,
								Description: "Specify the type of URL Signing. It can be either 0 or 2.",
							},
						},
					},
				},
				"slice": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Files larger than 10 MB will be requested and cached in parts (no larger than 10 MB each part).",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"sni": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Computed:    true,
					Description: "Specify the SNI.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"sni_type": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Available values 'dynamic' or 'custom'.",
							},
							"custom_hostname": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Required to set custom hostname in case sni-type='custom'.",
							},
						},
					},
				},
				"stale": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Computed:    true,
					Description: "The list of errors which the option is applied for.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Required:    true,
								Description: "The list of errors which the option is applied for.",
							},
						},
					},
				},
				"static_headers": { // deprecated in favor of static_response_headers
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Legacy option. Use the static_response_headers option instead.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeMap,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Required: true,
							},
						},
					},
				},
				"static_request_headers": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Specify custom HTTP Headers for a CDN server to add to request.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeMap,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Required: true,
							},
						},
					},
				},
				"static_response_headers": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Specify custom HTTP Headers that a CDN server adds to a response.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeList,
								Required: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:        schema.TypeString,
											Required:    true,
											Description: "Header name.",
										},
										"value": {
											Type:        schema.TypeSet,
											Elem:        &schema.Schema{Type: schema.TypeString},
											Required:    true,
											Description: "Header value.",
										},
										"always": {
											Type:        schema.TypeBool,
											Optional:    true,
											Computed:    true,
											Description: "Specifies if the header will be added to a response from CDN regardless of response code.",
										},
									},
								},
							},
						},
					},
				},
				"tls_versions": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "The option specifies a list of allowed SSL/TLS protocol versions.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeSet,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Required: true,
							},
						},
					},
				},
				"use_default_le_chain": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "The option allows choosing a Let's Encrypt certificate chain.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"user_agent_acl": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Control access to the content for specified user-agent.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"policy_type": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Possible values: allow, deny.",
							},
							"excepted_values": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Required:    true,
								Description: "List of User-Agent.",
							},
						},
					},
				},
				"use_rsa_le_cert": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "The option allows choosing the RSA Let's Encrypt certificate type for the resource.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"webp": { // deprecated in favor of image_stack
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Legacy option. Use the image_stack option instead.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"jpg_quality": {
								Type:     schema.TypeInt,
								Required: true,
							},
							"png_quality": {
								Type:     schema.TypeInt,
								Required: true,
							},
							"png_lossless": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  false,
							},
						},
					},
				},
				"websockets": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "The option allows WebSockets connections to an origin server.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
							"value": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
			},
		},
	}
)

func resourceCDNResource() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"cname": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A CNAME that will be used to deliver content though a CDN. If you update this field new resource will be created.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom client description of the resource.",
			},
			"origin_group": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"origin_group",
					"origin",
				},
				Description: "ID of the Origins Group. Use one of your Origins Group or create a new one. You can use either 'origin' parameter or 'originGroup' in the resource definition.",
			},
			"origin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"origin_group",
					"origin",
				},
				Description: "A domain name or IP of your origin source. Specify a port if custom. You can use either 'origin' parameter or 'originGroup' in the resource definition.",
			},
			"origin_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "This option defines the protocol that will be used by CDN servers to request content from an origin source. If not specified, we will use HTTP to connect to an origin server. Possible values are: HTTPS, HTTP, MATCH.",
			},
			"secondary_hostnames": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return []string{}, nil
				},
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of additional CNAMEs.",
			},
			"ssl_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use HTTPS protocol for content delivery.",
			},
			"ssl_data": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"ssl_enabled"},
				Description:  "Specify the SSL Certificate ID which should be used for the CDN Resource.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The setting allows to enable or disable a CDN Resource",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of a CDN resource content availability. Possible values are: Active, Suspended, Processed.",
			},
			"options": optionsSchema,
		},
		CreateContext: resourceCDNResourceCreate,
		ReadContext:   resourceCDNResourceRead,
		UpdateContext: resourceCDNResourceUpdate,
		DeleteContext: resourceCDNResourceDelete,
		Description:   "Represent CDN resource",
	}
}

func resourceCDNResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN Resource creating")
	config := m.(*Config)
	client := config.CDNClient

	var req resources.CreateRequest
	req.Cname = d.Get("cname").(string)
	req.Description = d.Get("description").(string)
	req.Origin = d.Get("origin").(string)
	req.OriginGroup = d.Get("origin_group").(int)
	req.OriginProtocol = resources.Protocol(d.Get("origin_protocol").(string))
	req.SSlEnabled = d.Get("ssl_enabled").(bool)
	req.SSLData = d.Get("ssl_data").(int)

	req.Options = listToOptions(d.Get("options").([]interface{}))

	for _, hostname := range d.Get("secondary_hostnames").(*schema.Set).List() {
		req.SecondaryHostnames = append(req.SecondaryHostnames, hostname.(string))
	}

	result, err := client.Resources().Create(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	resourceCDNResourceRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN Resource creating (id=%d)\n", result.ID)
	return nil
}

func resourceCDNResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN Resource reading (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.Resources().Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("cname", result.Cname)
	d.Set("description", result.Description)
	d.Set("origin_group", result.OriginGroup)
	d.Set("origin_protocol", result.OriginProtocol)
	d.Set("secondary_hostnames", result.SecondaryHostnames)
	d.Set("ssl_enabled", result.SSlEnabled)
	d.Set("ssl_data", result.SSLData)
	d.Set("status", result.Status)
	d.Set("active", result.Active)
	if err := d.Set("options", optionsToList(result.Options)); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Resource reading")
	return nil
}

func resourceCDNResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN Resource updating (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req resources.UpdateRequest
	req.Active = d.Get("active").(bool)
	req.Description = d.Get("description").(string)
	req.OriginGroup = d.Get("origin_group").(int)
	req.SSlEnabled = d.Get("ssl_enabled").(bool)
	req.SSLData = d.Get("ssl_data").(int)
	req.OriginProtocol = resources.Protocol(d.Get("origin_protocol").(string))
	req.Options = listToOptions(d.Get("options").([]interface{}))
	for _, hostname := range d.Get("secondary_hostnames").(*schema.Set).List() {
		req.SecondaryHostnames = append(req.SecondaryHostnames, hostname.(string))
	}

	if _, err := client.Resources().Update(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Resource updating")
	return resourceCDNResourceRead(ctx, d, m)
}

func resourceCDNResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN Resource deleting (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Resources().Delete(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish CDN Resource deleting")
	return nil
}

func listToOptions(l []interface{}) *gcdn.Options {
	if len(l) == 0 {
		return nil
	}

	var opts gcdn.Options
	fields := l[0].(map[string]interface{})

	if opt, ok := getOptByName(fields, "allowed_http_methods"); ok {
		opts.AllowedHTTPMethods = &gcdn.AllowedHTTPMethods{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.AllowedHTTPMethods.Value = append(opts.AllowedHTTPMethods.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "brotli_compression"); ok {
		opts.BrotliCompression = &gcdn.BrotliCompression{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.BrotliCompression.Value = append(opts.BrotliCompression.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "browser_cache_settings"); ok {
		opts.BrowserCacheSettings = &gcdn.BrowserCacheSettings{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "cache_http_headers"); ok {
		opts.CacheHttpHeaders = &gcdn.CacheHttpHeaders{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.CacheHttpHeaders.Value = append(opts.CacheHttpHeaders.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "cors"); ok {
		opts.Cors = &gcdn.Cors{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.Cors.Value = append(opts.Cors.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "country_acl"); ok {
		opts.CountryACL = &gcdn.CountryACL{
			Enabled:    opt["enabled"].(bool),
			PolicyType: opt["policy_type"].(string),
		}
		for _, v := range opt["excepted_values"].(*schema.Set).List() {
			opts.CountryACL.ExceptedValues = append(opts.CountryACL.ExceptedValues, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "disable_cache"); ok {
		opts.DisableCache = &gcdn.DisableCache{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "disable_proxy_force_ranges"); ok {
		opts.DisableProxyForceRanges = &gcdn.DisableProxyForceRanges{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "edge_cache_settings"); ok {
		rawCustomVals := opt["custom_values"].(map[string]interface{})
		customVals := make(map[string]string, len(rawCustomVals))
		for key, value := range rawCustomVals {
			customVals[key] = value.(string)
		}

		opts.EdgeCacheSettings = &gcdn.EdgeCacheSettings{
			Enabled:      opt["enabled"].(bool),
			Value:        opt["value"].(string),
			CustomValues: customVals,
			Default:      opt["default"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "fetch_compressed"); ok {
		opts.FetchCompressed = &gcdn.FetchCompressed{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "follow_origin_redirect"); ok {
		opts.FollowOriginRedirect = &gcdn.FollowOriginRedirect{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["codes"].(*schema.Set).List() {
			opts.FollowOriginRedirect.Codes = append(opts.FollowOriginRedirect.Codes, v.(int))
		}
	}
	if opt, ok := getOptByName(fields, "force_return"); ok {
		opts.ForceReturn = &gcdn.ForceReturn{
			Enabled: opt["enabled"].(bool),
			Code:    opt["code"].(int),
			Body:    opt["body"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "forward_host_header"); ok {
		opts.ForwardHostHeader = &gcdn.ForwardHostHeader{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "gzip_on"); ok {
		opts.GzipOn = &gcdn.GzipOn{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "host_header"); ok {
		opts.HostHeader = &gcdn.HostHeader{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "http3_enabled"); ok {
		opts.HTTP3Enabled = &gcdn.HTTP3Enabled{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "ignore_cookie"); ok {
		opts.IgnoreCookie = &gcdn.IgnoreCookie{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "ignore_query_string"); ok {
		opts.IgnoreQueryString = &gcdn.IgnoreQueryString{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "image_stack"); ok {
		opts.ImageStack = &gcdn.ImageStack{
			Enabled: opt["enabled"].(bool),
			Quality: opt["quality"].(int),
		}
		if _, ok := opt["avif_enabled"]; ok {
			opts.ImageStack.AvifEnabled = opt["avif_enabled"].(bool)
		}
		if _, ok := opt["webp_enabled"]; ok {
			opts.ImageStack.WebpEnabled = opt["webp_enabled"].(bool)
		}
		if _, ok := opt["png_lossless"]; ok {
			opts.ImageStack.PngLossless = opt["png_lossless"].(bool)
		}
	}
	if opt, ok := getOptByName(fields, "ip_address_acl"); ok {
		opts.IPAddressACL = &gcdn.IPAddressACL{
			Enabled:    opt["enabled"].(bool),
			PolicyType: opt["policy_type"].(string),
		}
		for _, v := range opt["excepted_values"].(*schema.Set).List() {
			opts.IPAddressACL.ExceptedValues = append(opts.IPAddressACL.ExceptedValues, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "limit_bandwidth"); ok {
		opts.LimitBandwidth = &gcdn.LimitBandwidth{
			Enabled:   opt["enabled"].(bool),
			LimitType: opt["limit_type"].(string),
		}
		if _, ok := opt["speed"]; ok {
			opts.LimitBandwidth.Speed = opt["speed"].(int)
		}
		if _, ok := opt["buffer"]; ok {
			opts.LimitBandwidth.Buffer = opt["buffer"].(int)
		}
	}
	if opt, ok := getOptByName(fields, "proxy_cache_methods_set"); ok {
		opts.ProxyCacheMethodsSet = &gcdn.ProxyCacheMethodsSet{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "query_params_blacklist"); ok {
		opts.QueryParamsBlacklist = &gcdn.QueryParamsBlacklist{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.QueryParamsBlacklist.Value = append(opts.QueryParamsBlacklist.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "query_params_whitelist"); ok {
		opts.QueryParamsWhitelist = &gcdn.QueryParamsWhitelist{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.QueryParamsWhitelist.Value = append(opts.QueryParamsWhitelist.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "redirect_https_to_http"); ok {
		opts.RedirectHttpsToHttp = &gcdn.RedirectHttpsToHttp{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "redirect_http_to_https"); ok {
		opts.RedirectHttpToHttps = &gcdn.RedirectHttpToHttps{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "referrer_acl"); ok {
		opts.ReferrerACL = &gcdn.ReferrerACL{
			Enabled:    opt["enabled"].(bool),
			PolicyType: opt["policy_type"].(string),
		}
		for _, v := range opt["excepted_values"].(*schema.Set).List() {
			opts.ReferrerACL.ExceptedValues = append(opts.ReferrerACL.ExceptedValues, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "request_limiter"); ok {
		opts.RequestLimiter = &gcdn.RequestLimiter{
			Enabled:  opt["enabled"].(bool),
			Rate:     opt["rate"].(int),
			Burst:    opt["burst"].(int),
			RateUnit: opt["rate_unit"].(string),
			Delay:    opt["delay"].(int),
		}
	}
	if opt, ok := getOptByName(fields, "response_headers_hiding_policy"); ok {
		opts.ResponseHeadersHidingPolicy = &gcdn.ResponseHeadersHidingPolicy{
			Enabled: opt["enabled"].(bool),
			Mode:    opt["mode"].(string),
		}
		for _, v := range opt["excepted"].(*schema.Set).List() {
			opts.ResponseHeadersHidingPolicy.Excepted = append(opts.ResponseHeadersHidingPolicy.Excepted, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "rewrite"); ok {
		opts.Rewrite = &gcdn.Rewrite{
			Enabled: opt["enabled"].(bool),
			Body:    opt["body"].(string),
			Flag:    opt["flag"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "secure_key"); ok {
		opts.SecureKey = &gcdn.SecureKey{
			Enabled: opt["enabled"].(bool),
			Key:     opt["key"].(string),
			Type:    opt["type"].(int),
		}
	}
	if opt, ok := getOptByName(fields, "slice"); ok {
		opts.Slice = &gcdn.Slice{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "sni"); ok {
		opts.SNI = &gcdn.SNIOption{
			Enabled:        opt["enabled"].(bool),
			SNIType:        opt["sni_type"].(string),
			CustomHostname: opt["custom_hostname"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "stale"); ok {
		opts.Stale = &gcdn.Stale{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.Stale.Value = append(opts.Stale.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "static_headers"); ok {
		opts.StaticHeaders = &gcdn.StaticHeaders{
			Enabled: opt["enabled"].(bool),
			Value:   map[string]string{},
		}
		for k, v := range opt["value"].(map[string]interface{}) {
			opts.StaticHeaders.Value[k] = v.(string)
		}
	}
	if opt, ok := getOptByName(fields, "static_request_headers"); ok {
		opts.StaticRequestHeaders = &gcdn.StaticRequestHeaders{
			Enabled: opt["enabled"].(bool),
			Value:   map[string]string{},
		}
		for k, v := range opt["value"].(map[string]interface{}) {
			opts.StaticRequestHeaders.Value[k] = v.(string)
		}
	}
	if opt, ok := getOptByName(fields, "static_response_headers"); ok {
		opts.StaticResponseHeaders = &gcdn.StaticResponseHeaders{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].([]interface{}) {
			item_data := v.(map[string]interface{})
			item := &gcdn.StaticResponseHeadersItem{
				Name: item_data["name"].(string),
			}
			for _, v := range item_data["value"].(*schema.Set).List() {
				item.Value = append(item.Value, v.(string))
			}
			if _, ok := opt["always"]; ok {
				item.Always = item_data["always"].(bool)
			}
			opts.StaticResponseHeaders.Value = append(opts.StaticResponseHeaders.Value, *item)
		}
	}
	if opt, ok := getOptByName(fields, "tls_versions"); ok {
		opts.TLSVersions = &gcdn.TLSVersions{
			Enabled: opt["enabled"].(bool),
		}
		for _, v := range opt["value"].(*schema.Set).List() {
			opts.TLSVersions.Value = append(opts.TLSVersions.Value, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "use_default_le_chain"); ok {
		opts.UseDefaultLEChain = &gcdn.UseDefaultLEChain{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "user_agent_acl"); ok {
		opts.UserAgentACL = &gcdn.UserAgentACL{
			Enabled:    opt["enabled"].(bool),
			PolicyType: opt["policy_type"].(string),
		}
		for _, v := range opt["excepted_values"].(*schema.Set).List() {
			opts.UserAgentACL.ExceptedValues = append(opts.UserAgentACL.ExceptedValues, v.(string))
		}
	}
	if opt, ok := getOptByName(fields, "use_rsa_le_cert"); ok {
		opts.UseRSALECert = &gcdn.UseRSALECert{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "webp"); ok {
		opts.Webp = &gcdn.Webp{
			Enabled:     opt["enabled"].(bool),
			JPGQuality:  opt["jpg_quality"].(int),
			PNGQuality:  opt["png_quality"].(int),
			PNGLossless: opt["png_lossless"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "websockets"); ok {
		opts.WebSockets = &gcdn.WebSockets{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(bool),
		}
	}

	return &opts
}

func getOptByName(fields map[string]interface{}, name string) (map[string]interface{}, bool) {
	if _, ok := fields[name]; !ok {
		return nil, false
	}

	container, ok := fields[name].([]interface{})
	if !ok {
		return nil, false
	}

	if len(container) == 0 {
		return nil, false
	}

	opt, ok := container[0].(map[string]interface{})
	if !ok {
		return nil, false
	}

	return opt, true
}

func optionsToList(options *gcdn.Options) []interface{} {
	result := make(map[string][]interface{})
	if options.AllowedHTTPMethods != nil {
		m := structToMap(options.AllowedHTTPMethods)
		result["allowed_http_methods"] = []interface{}{m}
	}
	if options.BrotliCompression != nil {
		m := structToMap(options.BrotliCompression)
		result["brotli_compression"] = []interface{}{m}
	}
	if options.BrowserCacheSettings != nil {
		m := structToMap(options.BrowserCacheSettings)
		result["browser_cache_settings"] = []interface{}{m}
	}
	if options.CacheHttpHeaders != nil {
		m := structToMap(options.CacheHttpHeaders)
		result["cache_http_headers"] = []interface{}{m}
	}
	if options.Cors != nil {
		m := structToMap(options.Cors)
		result["cors"] = []interface{}{m}
	}
	if options.CountryACL != nil {
		m := structToMap(options.CountryACL)
		result["country_acl"] = []interface{}{m}
	}
	if options.DisableCache != nil {
		m := structToMap(options.DisableCache)
		result["disable_cache"] = []interface{}{m}
	}
	if options.DisableProxyForceRanges != nil {
		m := structToMap(options.DisableProxyForceRanges)
		result["disable_proxy_force_ranges"] = []interface{}{m}
	}
	if options.EdgeCacheSettings != nil {
		m := structToMap(options.EdgeCacheSettings)
		result["edge_cache_settings"] = []interface{}{m}
	}
	if options.FetchCompressed != nil {
		m := structToMap(options.FetchCompressed)
		result["fetch_compressed"] = []interface{}{m}
	}
	if options.FollowOriginRedirect != nil {
		m := structToMap(options.FollowOriginRedirect)
		result["follow_origin_redirect"] = []interface{}{m}
	}
	if options.ForceReturn != nil {
		m := structToMap(options.ForceReturn)
		result["force_return"] = []interface{}{m}
	}
	if options.ForwardHostHeader != nil {
		m := structToMap(options.ForwardHostHeader)
		result["forward_host_header"] = []interface{}{m}
	}
	if options.GzipOn != nil {
		m := structToMap(options.GzipOn)
		result["gzip_on"] = []interface{}{m}
	}
	if options.HostHeader != nil {
		m := structToMap(options.HostHeader)
		result["host_header"] = []interface{}{m}
	}
	if options.HTTP3Enabled != nil {
		m := structToMap(options.HTTP3Enabled)
		result["http3_enabled"] = []interface{}{m}
	}
	if options.IgnoreCookie != nil {
		m := structToMap(options.IgnoreCookie)
		result["ignore_cookie"] = []interface{}{m}
	}
	if options.IgnoreQueryString != nil {
		m := structToMap(options.IgnoreQueryString)
		result["ignore_query_string"] = []interface{}{m}
	}
	if options.ImageStack != nil {
		m := structToMap(options.ImageStack)
		result["image_stack"] = []interface{}{m}
	}
	if options.IPAddressACL != nil {
		m := structToMap(options.IPAddressACL)
		result["ip_address_acl"] = []interface{}{m}
	}
	if options.LimitBandwidth != nil {
		m := structToMap(options.LimitBandwidth)
		result["limit_bandwidth"] = []interface{}{m}
	}
	if options.ProxyCacheMethodsSet != nil {
		m := structToMap(options.ProxyCacheMethodsSet)
		result["proxy_cache_methods_set"] = []interface{}{m}
	}
	if options.QueryParamsBlacklist != nil {
		m := structToMap(options.QueryParamsBlacklist)
		result["query_params_blacklist"] = []interface{}{m}
	}
	if options.QueryParamsWhitelist != nil {
		m := structToMap(options.QueryParamsWhitelist)
		result["query_params_whitelist"] = []interface{}{m}
	}
	if options.RedirectHttpsToHttp != nil {
		m := structToMap(options.RedirectHttpsToHttp)
		result["redirect_https_to_http"] = []interface{}{m}
	}
	if options.RedirectHttpToHttps != nil {
		m := structToMap(options.RedirectHttpToHttps)
		result["redirect_http_to_https"] = []interface{}{m}
	}
	if options.ReferrerACL != nil {
		m := structToMap(options.ReferrerACL)
		result["referrer_acl"] = []interface{}{m}
	}
	if options.RequestLimiter != nil {
		m := structToMap(options.RequestLimiter)
		result["request_limiter"] = []interface{}{m}
	}
	if options.ResponseHeadersHidingPolicy != nil {
		m := structToMap(options.ResponseHeadersHidingPolicy)
		result["response_headers_hiding_policy"] = []interface{}{m}
	}
	if options.Rewrite != nil {
		m := structToMap(options.Rewrite)
		result["rewrite"] = []interface{}{m}
	}
	if options.SecureKey != nil {
		m := structToMap(options.SecureKey)
		result["secure_key"] = []interface{}{m}
	}
	if options.Slice != nil {
		m := structToMap(options.Slice)
		result["slice"] = []interface{}{m}
	}
	if options.SNI != nil {
		m := structToMap(options.SNI)
		result["sni"] = []interface{}{m}
	}
	if options.Stale != nil {
		m := structToMap(options.Stale)
		result["stale"] = []interface{}{m}
	}
	if options.StaticHeaders != nil {
		m := structToMap(options.StaticHeaders)
		result["static_headers"] = []interface{}{m}
	}
	if options.StaticRequestHeaders != nil {
		m := structToMap(options.StaticRequestHeaders)
		result["static_request_headers"] = []interface{}{m}
	}
	if options.StaticResponseHeaders != nil {
		m := structToMap(options.StaticResponseHeaders)
		items := []interface{}{}
		for _, v := range m["value"].([]gcdn.StaticResponseHeadersItem) {
			items = append(items, structToMap(v))
		}
		m["value"] = items
		result["static_response_headers"] = []interface{}{m}
	}
	if options.TLSVersions != nil {
		m := structToMap(options.TLSVersions)
		result["tls_versions"] = []interface{}{m}
	}
	if options.UseDefaultLEChain != nil {
		m := structToMap(options.UseDefaultLEChain)
		result["use_default_le_chain"] = []interface{}{m}
	}
	if options.UserAgentACL != nil {
		m := structToMap(options.UserAgentACL)
		result["user_agent_acl"] = []interface{}{m}
	}
	if options.UseRSALECert != nil {
		m := structToMap(options.UseRSALECert)
		result["use_rsa_le_cert"] = []interface{}{m}
	}
	if options.Webp != nil {
		m := structToMap(options.Webp)
		result["webp"] = []interface{}{m}
	}
	if options.WebSockets != nil {
		m := structToMap(options.WebSockets)
		result["websockets"] = []interface{}{m}
	}
	return []interface{}{result}
}

func structToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}
