package highwinds

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/pkg/debug"
	"github.com/openwurl/wurlwind/pkg/utilities"
	"github.com/openwurl/wurlwind/striketracker/models"
)

// resourceConfiguration manages configuration scopes (except for default/root scope)
func resourceConfiguration() *schema.Resource {
	// Scope
	scopeSchema := &schema.Schema{
		Type:        schema.TypeMap,
		Required:    true,
		Description: "Fields concerning the identity of this scope configuration",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"platform": {
					Type:        schema.TypeString,
					Description: "The CDN platform this scope is utilizing",
					Default:     "CDS",
					Optional:    true,
				},
				"path": {
					Type:        schema.TypeString,
					Description: "The URI path of this scope configuration",
					Default:     "/",
					Optional:    true,
				},
				"name": {
					Type:        schema.TypeString,
					Description: "The name of this scope configuration",
					Required:    true,
				},
			},
		},
	}

	// Origin Pull host
	originHostSchema := &schema.Schema{
		Description: "The origins this scope will use to populate its cache",
		Type:        schema.TypeSet,
		MaxItems:    1,
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"primary": {
					Type:        schema.TypeInt,
					Description: "The ID of the primary origin to use",
					Optional:    true,
				},
				"secondary": {
					Type:        schema.TypeInt,
					Description: "The ID of the secondary origin to use",
					Optional:    true,
				},
				"path": {
					Type:        schema.TypeString,
					Description: "The URI path to use at this origin",
					Default:     "/",
					Optional:    true,
				},
			},
		},
	}

	// Stale Cache Extension
	originPullCacheExtensionSchema := &schema.Schema{
		Description: "The stale cache extension settings for this scope",
		Type:        schema.TypeSet,
		MaxItems:    1,
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:        schema.TypeBool,
					Description: "Whether or not stale cache extension is enabled",
					Optional:    true,
				},
				"expired_cache_extension": {
					Type:        schema.TypeInt,
					Description: "The TTL extension for expired cache",
					Optional:    true,
				},
				"origin_unreachable_cache_extension": {
					Type:        schema.TypeInt,
					Description: "The TTL extension for origin fetch failures",
					Optional:    true,
				},
			},
		},
	}

	// Cache Policy / OriginPullPolicy - weighted
	originPullPolicySchema := &schema.Schema{
		Description: "Cache control policies to apply to an origin",
		Type:        schema.TypeSet,
		MinItems:    1,
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"weight": {
					Type:        schema.TypeInt,
					Description: "Indicates the position in the ordered stack of this CDN cache rule",
					Required:    true,
				},
				"enabled": {
					Type:        schema.TypeBool,
					Description: "Whether or not this cache rule is enabled",
					Optional:    true,
					Default:     true,
				},
				"expire_policy": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "CACHE_CONTROL",
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(string)
						if !utilities.SliceContainsString(v, models.ValidExpirePolicies) {
							errs = append(errs, fmt.Errorf("%q must be one of (%v), got %s", key, models.ValidExpirePolicies, val))
						}
						return warns, errs
					},
				},
				"expire_seconds": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  31536000,
				},
				"force_bypass_cache": {
					Type:     schema.TypeBool,
					Default:  false,
					Optional: true,
				},
				"honor_must_revalidate": {
					Type:     schema.TypeBool,
					Default:  false,
					Optional: true,
				},
				"honor_no_cache": {
					Type:     schema.TypeBool,
					Default:  false,
					Optional: true,
				},
				"honor_no_store": {
					Type:     schema.TypeBool,
					Default:  false,
					Optional: true,
				},
				"honor_private": {
					Type:     schema.TypeBool,
					Default:  false,
					Optional: true,
				},
				"honor_smax_age": {
					Type:     schema.TypeBool,
					Default:  false,
					Optional: true,
				},
				"http_headers": {
					Type:     schema.TypeString,
					Default:  "*",
					Optional: true,
				},
				"must_revalidate_to_no_cache": {
					Type:     schema.TypeBool,
					Default:  false,
					Optional: true,
				},
				"no_cache_behavior": {
					Type:     schema.TypeString,
					Default:  "spec",
					Optional: true,
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(string)
						if !utilities.SliceContainsString(v, models.ValidNoCacheBehaviors) {
							errs = append(errs, fmt.Errorf("%q must be one of (%v), got %s", key, models.ValidNoCacheBehaviors, val))
						}
						return warns, errs
					},
				},
				"update_http_headers_on_304_response": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"default_cache_behavior": {
					Type:     schema.TypeString,
					Default:  "ttl",
					Optional: true,
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(string)
						if !utilities.SliceContainsString(v, models.ValidCacheBehaviors) {
							errs = append(errs, fmt.Errorf("%q must be one of (%v), got %s", key, models.ValidCacheBehaviors, val))
						}
						return warns, errs
					},
				},
				"max_age_zero_to_no_cache": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"bypass_cache_identifier": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"content_type_filter": {
					Type:     schema.TypeString,
					Default:  "*",
					Optional: true,
				},
				"header_filter": {
					Type:     schema.TypeString,
					Default:  "*",
					Optional: true,
				},
				"method_filter": {
					Type:     schema.TypeString,
					Default:  "*",
					Optional: true,
				},
				"path_filter": {
					Type:     schema.TypeString,
					Default:  "*",
					Optional: true,
				},
				"status_code_match": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}

	// Request Modifications - weighted
	requestModificationsSchema := &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Edge rules targeting the origin response/request",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"weight": {
					Type:        schema.TypeInt,
					Description: "Indicates the position in the ordered stack of this CDN cache rule",
					Required:    true,
				},
				"enabled": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
				},
				"add_headers": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"flow_control": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "next",
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(string)
						valid := []string{"next", "break"}
						if !utilities.SliceContainsString(strings.ToLower(v), valid) {
							errs = append(errs, fmt.Errorf("%q must be one of (next, break), got %s", key, val))
						}
						return warns, errs
					},
					StateFunc: func(val interface{}) string {
						return strings.ToLower(val.(string))
					},
				},
				"url_pattern": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"url_rewrite": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"header_pattern": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"header_rewrite": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}

	compressionSchema := &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Fields for compression on delivery",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type: schema.TypeBool,
					//Default:  true,
					Optional: true,
				},
				"gzip": {
					Type:     schema.TypeString,
					Optional: true,
					//Default:     "txt,js,htm,html,css",
					Description: "File suffixes to compress if requested",
				},
				"level": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "The level of compression used for gzip",
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(int)
						if v < 0 || v > 6 {
							errs = append(errs, fmt.Errorf("%q must be greater than or equal to 0 or no greater than 6, got %s", key, val))
						}
						return warns, errs
					},
				},
				"mime": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Mime types to be used with gzip compression",
				},
			},
		},
	}

	staticHeaderSchema := &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Fields for static header injection",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"weight": {
					Type:        schema.TypeInt,
					Description: "Indicates the position in the ordered stack of this static header rule",
					Required:    true,
				},
				"enabled": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
				},
				"http": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Headers to add to the response to the client",
				},
				"origin_pull": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Headers to add to origin request",
				},
				"client_request": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Headers to add to the request from the client to the CDN",
				},
				"method_filter": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "HTTP Method filter for this rule",
				},
				"path_filter": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "URL Path filter for this rule",
				},
				"header_filter": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Header filter for this rule",
				},
				"client_response_code_filter": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Client response Status code this rule applies to, can use globs like 2* for all 200s or 3* for all 300s etc",
				},
			},
		},
	}

	httpMethodSchema := &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Selectively enable additional HTTP methods you'd like the CDN to process",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
				},
				"passthru": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "*",
				},
			},
		},
	}

	responseHeaderSchema := &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Enable and bypass certain origin headers that affect delivery",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"http": {
					Type:        schema.TypeString,
					Description: "Force download by user agent",
					Deprecated:  "This key is used for legacy sites, you should use client response modifications instead",
					Optional:    true,
				},
				"enable_etag": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enables the etag header on responses from the CDN",
				},
			},
		},
	}

	// Delivery
	deliverySchema := &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Delivery related features, set of sets",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"compression":     compressionSchema,
				"static_header":   staticHeaderSchema,
				"http_methods":    httpMethodSchema,
				"response_header": responseHeaderSchema,
				// delivery_behaviors => http_methods, custom_http_response_headers
				// rate_limiting => bandwidth_rate_limiting, pattern_based_rate_limiting
				// force_downloads => disposition_by_http_header, disposition_by_url
				// custom_mime_types
				// edge_responses => cache_rule
				// media_delivery => pseudo_flash_streaming, pseudo_streaming, reserved_query_string_params
				// error_redirects => exceptions_to_redirect, response_code_redirection
				// dynamic_files => robots_configuration
				//
			},
		},
	}

	return &schema.Resource{
		Create: resourceConfigurationCreate,
		Read:   resourceConfigurationRead,
		Update: resourceConfigurationUpdate,
		Delete: resourceConfigurationDelete,
		Exists: resourceConfigurationExists,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// TODO: Importer - Ensure this actually works and expand as needed
				accountHash, hostHash, scopeID, err := ResourceConfigurationParseHashID(d.Id())
				if err != nil {
					return nil, err
				}
				d.Set("account_hash", accountHash)
				d.Set("host_hash", hostHash)
				d.SetId(scopeID)

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"account_hash": &schema.Schema{
				Description: "The destination account hash where the origin will be created",
				Type:        schema.TypeString,
				Required:    true,
			},
			"host_hash": &schema.Schema{
				Description: "The hash code of the parent host this scope is being attached to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"hostnames": &schema.Schema{
				Description: "Hostnames to be associated with this scope",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// TODO: The default DNS associated by highwinds can skew the diff
					// that determines if this field is sane, this is an attempt
					// to suppress that and may need reviewed
					// As of right now this means the default CDN dns gets removed
					// unless added to the tf schema manually
					debug.Log("hostnames", "OLD: %v | NEW: %v | STATE: %s", old, new, d.Get("hostnames"))
					if strings.Contains(old, "hwcdn.net") && new == "" {
						return true
					}
					if old == new {
						return true
					}

					return false
				},
				Optional: true,
			},
			"scope":                     scopeSchema,
			"origin_pull_host":          originHostSchema,
			"stale_cache_extension":     originPullCacheExtensionSchema,
			"cache_policy":              originPullPolicySchema,
			"origin_request_edge_rule":  requestModificationsSchema,
			"origin_response_edge_rule": requestModificationsSchema,
			"client_request_edge_rule":  requestModificationsSchema,
			"client_response_edge_rule": requestModificationsSchema,
			"delivery":                  deliverySchema,
		},
	}
}
