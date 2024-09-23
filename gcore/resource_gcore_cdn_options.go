package gcore

import (
	"maps"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	commonOptions = map[string]*schema.Schema{
		"allowed_http_methods": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "Specify allowed HTTP methods.",
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
			Description: "Brotli compression option allows to compress content with brotli on the CDN's end. CDN servers will request only uncompressed content from the origin.",
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
			Description: "Specify the cache expiration time for customers' browsers in seconds.",
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
						Description: "Use '0s' to disable caching. The value applies for a response with codes 200, 201, 204, 206, 301, 302, 303, 304, 307, 308.",
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
			Description: "CORS header support option adds the Access-Control-Allow-Origin header to responses from CDN servers.",
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
						Description: "Specify a value of the Access-Control-Allow-Origin header. Possible values: '*', '$http_origin', 'example.com'.",
					},
					"always": {
						Type:        schema.TypeBool,
						Optional:    true,
						Computed:    true,
						Description: "Specify if the Access-Control-Allow-Origin header should be added to a response from CDN regardless of response code.",
					},
				},
			},
		},
		"country_acl": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Country access policy enables control access to content for specified countries.",
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
			Description: "Option enables browser caching. When enabled, content caching is completely disabled.",
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
			Description: "The option allows getting 206 responses regardless settings of an origin source. Enabled by default.",
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
						Description: "Specify caching time for the response with codes 200, 206, 301, 302. Responses with codes 4xx, 5xx will not be cached. Use '0s' to disable caching. Use custom_values field to specify a custom caching time for a response with specific codes.",
					},
					"default": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Content will be cached according to origin cache settings. The value applies for a response with codes 200, 201, 204, 206, 301, 302, 303, 304, 307, 308, if an origin server does not have caching HTTP headers. Responses with other codes will not be cached.",
					},
					"custom_values": {
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
						DefaultFunc: func() (interface{}, error) {
							return map[string]interface{}{}, nil
						},
						Elem:        schema.TypeString,
						Description: "Specify caching time in seconds ('0s', '600s' for example) for a response with specific response code ('304', '404' for example). Use 'any' to specify caching time for all response codes. Use '0s' to disable caching for a specific response code. These settings have a higher priority than the value field.",
					},
				},
			},
		},
		"fastedge": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Allows to configure FastEdge app to be called on different request/response phases.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"on_request_headers": {
						Type:        schema.TypeList,
						MaxItems:    1,
						Required:    true,
						Description: "Allows to configure FastEdge application that will be called to handle request headers as soon as CDN receives incoming HTTP request.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"enabled": {
									Type:        schema.TypeBool,
									Optional:    true,
									Default:     true,
									Description: "Determines if the FastEdge application should be called whenever HTTP request headers are received.",
								},
								"app_id": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "The ID of the application in FastEdge.",
								},
								"interrupt_on_error": {
									Type:        schema.TypeBool,
									Optional:    true,
									Default:     true,
									Description: "Determines if the request execution should be interrupted when an error occurs.",
								},
							},
						},
					},
				},
			},
		},
		"fetch_compressed": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Option allows to enable fetch compressed. CDN request and cache already compressed content. Your server should support compression. CDN servers will not ungzip your content even if a user's browser doesn't accept compression (nowadays almost all browsers support it).",
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
			Description: "Enable redirection from origin. If the origin server returns a redirect, the option allows the CDN to pull the requested content from the origin server that was returned in the redirect.",
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
			Description: "Allows to apply custom HTTP code to the CDN content. Specify HTTP-code you need and text or URL if you are going to set up redirect.",
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
			Description: "When a CDN requests content from an origin server, the option allows to forward the Host header used in the request made to a CDN.",
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
			Description: "GZip compression option allows to compress content with gzip on the CDN`s end. CDN servers will request only uncompressed content from the origin.",
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
			Description: "Option allows to set Host header that CDN servers use when request content from an origin server. Your server must be able to process requests with the chosen header. If the option is NULL, Host Header value is taken from the parent CDN resource's value.",
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
			Description: "Ignore query string option determines how files with different query strings will be cached: either as one object (option is enabled) or as different objects (option is disabled).",
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
			Description: "Image stack option allows transforming JPG and PNG images (such as resizing or cropping) and automatically converting them to WebP or AVIF format. It is a paid option.",
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
						Description: "Quality settings for JPG and PNG images. Specify a value from 1 to 100. The higher the value, the better the image quality and the larger the file size after conversion.",
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
			Description: "IP access policy option allows to control access to the CDN Resource content for specific IP addresses.",
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
						Description: "Specify list of IP address with a subnet mask.",
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
		"proxy_cache_key": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "The option allows to modify the cache key. If omitted, the default value is $request_uri. Warning: Enabling and changing this option can invalidate your current cache and affect the cache hit ratio. Furthermore, the \"Purge by pattern\" option will not work.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"value": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Key for caching. Should be a combination of the specified variables: $http_host, $request_uri, $scheme, $uri.",
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
		"proxy_connect_timeout": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "The time limit for establishing a connection with the origin.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"value": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Specify time in seconds ('1s', '30s' for example).",
					},
				},
			},
		},
		"proxy_read_timeout": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "The time limit for receiving a partial response from the origin. If no response is received within this time, the connection will be closed.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"value": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Specify time in seconds ('1s', '30s' for example).",
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
			Description: "When enabled, HTTPS requests are redirected to HTTP.",
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
			Description: "When enabled, HTTP requests are redirected to HTTPS.",
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
			Description: "Referrer access policy option allows to control access to the CDN Resource content for specified domain names.",
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
						Description: "Specify list of domain names or wildcard domains (without http:// or https://). For example, example.com or *.example.com.",
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
						Description: "Specify a mode of hiding HTTP headers from the response. Possible values are: hide, show.",
					},
					"excepted": {
						Type:        schema.TypeSet,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Required:    true,
						Description: "List of HTTP headers. The following required headers cannot be hidden from response: Connection, Content-Length, Content-Type, Date, Server.",
					},
				},
			},
		},
		"rewrite": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Rewrite option changes and redirects the requests from the CDN to the origin. It operates according to the Nginx configuration.",
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
						Description: "The pattern for Rewrite. At least one group should be specified. For Example: /rewrite_from/(.*) /rewrite_to/$1",
					},
					"flag": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "break",
						Description: "Define flag for the Rewrite option. Possible values: last, break, redirect, permanent.",
					},
				},
			},
		},
		"secure_key": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "The option allows configuring an access with tokenized URLs. It makes impossible to access content without a valid (unexpired) hash key. When enabled, you need to specify a key that you use to generate a token.",
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
						Description: "A key generated on your side that will be used for URL signing.",
					},
					"type": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "Specify the type of URL Signing. It can be either 0 or 2. Type 0 - includes end user's IP to secure token generation. Type 2 - excludes end user's IP from secure token generation.",
					},
				},
			},
		},
		"slice": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "When enabled, files larger than 10 MB are requested and cached in parts (no larger than 10 MB each). It reduces time to first byte. The origin must support HTTP Range requests.",
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
			Description: "Specify the SNI (Server Name Indication). SNI (Server Name Indication) is generally only required if your origin is using shared hosting or does not have a dedicated IP address. If the origin server presents multiple certificates, SNI allows the origin server to know which certificate to use for the connection. The option works only if originProtocol parameter is HTTPS or MATCH.",
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
						Description: "Specify SNI type. Possible values: dynamic, custom. dynamic - SNI hostname depends on the hostHeader and the forward_host_header options. custom - custom SNI hostname.",
					},
					"custom_hostname": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Custom SNI hostname. Required if sni_type is set to 'custom'.",
					},
				},
			},
		},
		"stale": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "The list of errors which Always Online option is applied for.",
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
						Description: "Possible values: error, http_403, http_404, http_429, http_500, http_502, http_503, http_504, invalid_header, timeout, updating.",
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
						Type:        schema.TypeMap,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Required:    true,
						Description: "Header name is restricted to 255 symbols and can contain latin letters (A-Z, a-z), numbers (0-9), dashes, and underscores. Header value is restricted to 512 symbols and can contain latin letters (a-z), numbers (0-9), spaces, underscores and symbols (-/.:). Space can be used only between words.",
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
		"user_agent_acl": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "User agents policy option allows to control access to the content for specified user-agent.",
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
						Description: "List of User-Agents. Use \"\" to allow/deny access when the User-Agent header is empty.",
					},
				},
			},
		},
		"waf": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Option allows to enable Basic WAF to protect you against the most common threats.",
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
		"websockets": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "WebSockets option allows WebSockets connections to an origin server.",
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
	}
)

var (
	resourceOptions = map[string]*schema.Schema{
		"http3_enabled": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Use HTTP/3 protocol for content delivery if supported by the end users browser.",
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
		"tls_versions": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "The option specifies a list of allowed SSL/TLS protocol versions. The list cannot be empty. By default, the option is disabled (all protocols versions are allowed). ",
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
						Description: "Possible values (case sensitive): SSLv3, TLSv1, TLSv1.1, TLSv1.2, TLSv1.3.",
					},
				},
			},
		},
		"use_default_le_chain": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "The option allows choosing a Let's Encrypt certificate chain. The specified chain will be used during the next Let's Encrypt certificate issue or renewal. ",
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
	}
)

var (
	resourceOptionsSchema = &schema.Schema{
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Computed:    true,
		Description: "Each option in CDN resource settings. Each option added to CDN resource settings should have the following mandatory request fields: enabled, value.",
		Elem: &schema.Resource{
			Schema: resourceOptions,
		},
	}
)

var (
	ruleOptionsSchema = &schema.Schema{
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Computed:    true,
		Description: "Each option in CDN rule settings. Each option added to CDN rule settings should have the following mandatory request fields: enabled, value.",
		Elem: &schema.Resource{
			Schema: commonOptions,
		},
	}
)

func init() {
	maps.Copy(resourceOptions, commonOptions)
}
