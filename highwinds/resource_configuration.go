package highwinds

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/pkg/utilities"
	"github.com/openwurl/wurlwind/striketracker/models"
)

/*
	TODO:
		* The rest  of the configuration payload
		* Read should set all of the known state
		* Finish reverse engineering API
*/

func resourceConfiguration() *schema.Resource {

	scopeSchema := &schema.Schema{
		Type:        schema.TypeMap,
		Optional:    true,
		Description: "Fields concerning the identity of this scope configuration",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeInt,
					Description: "The ID of this scoped configuration",
					Computed:    true,
					Optional:    false,
				},
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

	originPullPolicySchema := &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Policy rules for managing origin pull cache behavior",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
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
					Default:  true,
					Optional: true,
				},
				"honor_no_cache": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
				},
				"honor_no_store": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
				},
				"honor_private": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
				},
				"honor_smax_age": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
				},
				"http_headers": {
					Type:     schema.TypeString,
					Default:  "*",
					Optional: true,
				},
				"must_revalidate_to_no_cache": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
				},
				"no_cache_behavior": {
					Type:     schema.TypeString,
					Default:  "spec",
					Optional: true,
					// TODO: Validation spec or legacy
				},
				"update_http_headers_on_304_response": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
				},
				"default_cache_behavior": {
					Type:     schema.TypeString,
					Default:  "ttl",
					Optional: true,
				},
				"max_age_zero_to_no_cache": {
					Type:     schema.TypeBool,
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
			},
		},
	}

	originSchema := &schema.Schema{
		Type:        schema.TypeMap,
		Optional:    true,
		Description: "Fields concerning the configuration of the origin",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"primary": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The primary origin for this configuration",
				},
				"secondary": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The secondary/failover origin for this configuration",
				},
				"path": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The path to pull from at the specified origin",
					// TODO: Validate is a path
				},
				"origin_pull_protocol": &schema.Schema{
					Description: "The protocol to use for pulling from this origin. (http, https, or match)",
					Type:        schema.TypeString,
					Required:    true,
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(string)
						if !utilities.SliceContainsString(strings.ToLower(v), models.ValidPullProtocols) {
							errs = append(errs, fmt.Errorf("%q must be one of (http, https, or match), got %s", key, val))
						}
						return warns, errs
					},
					StateFunc: func(val interface{}) string {
						return strings.ToLower(val.(string))
					},
				},
				"redirect_action": &schema.Schema{
					Description: "The protocol to use for pulling from this origin. (proxy, follow)",
					Type:        schema.TypeString,
					Required:    true,
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(string)
						valid := []string{"proxy", "follow"}
						if !utilities.SliceContainsString(strings.ToLower(v), valid) {
							errs = append(errs, fmt.Errorf("%q must be one of (proxy, follow), got %s", key, val))
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

	requestModifications := &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Edge rules targeting the origin request",
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
				// TODO: These needs expanded to account for host hash as well
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
				Description: "Hostnames to be associated with this configuration",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				// TODO: Validation
				Optional: true,
			},
			"dns_hostnames": &schema.Schema{
				Description: "Non default hostnames for use in DNS",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"origin":                    originSchema,
			"origin_pull_policy":        originPullPolicySchema,
			"scope":                     scopeSchema,
			"origin_request_edge_rule":  requestModifications,
			"origin_response_edge_rule": requestModifications,
			"client_request_edge_rule":  requestModifications,
			"client_response_edge_rule": requestModifications,
		},
	}
}
