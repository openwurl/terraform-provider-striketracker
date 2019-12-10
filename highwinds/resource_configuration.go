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

	// Cache Policy
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

	// Request Modifications
	// TODO: This will need expanded, there are more fields
	requestModifications := &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Edge rules targeting the origin response/request",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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
			"origin_request_edge_rule":  requestModifications,
			"origin_response_edge_rule": requestModifications,
			"client_request_edge_rule":  requestModifications,
			"client_response_edge_rule": requestModifications,
		},
	}
}
