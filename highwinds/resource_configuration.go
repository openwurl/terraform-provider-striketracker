package highwinds

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/pkg/utilities"
	"github.com/openwurl/wurlwind/striketracker"
	"github.com/openwurl/wurlwind/striketracker/models"
	"github.com/openwurl/wurlwind/striketracker/services/configuration"
)

func resourceConfiguration() *schema.Resource {
	originPullHost := &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "Define the primary and secondary origin pull host",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"primary": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "The primary origin for this configuration",
				},
				"secondary": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "The secondary/failover origin for this configuration",
				},
				"path": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The path to pull from at the specified origin",
					// TODO: Validate is a path
				},
			},
		},
	}

	originPullPolicy := &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Policy rules for managing origin pull cache behavior",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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
			},
		},
	}

	cacheControl := &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Cache control rules",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"max_age": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  31536000,
				},
				"synchronize_max_age": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
				},
			},
		},
	}

	compression := &schema.Schema{
		Type:        schema.TypeMap,
		Optional:    true,
		Description: "GZIP and mime rules",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"gzip": {
					Type:     schema.TypeString,
					Default:  "txt,js,htm,html,css",
					Optional: true,
				},
				"level": {
					Type:     schema.TypeInt,
					Default:  1,
					Optional: true,
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(int)
						if v <= 0 {
							errs = append(errs, fmt.Errorf("Minimum compression level if enabled is 1, got %d", v))
						}
						if v > 6 {
							errs = append(errs, fmt.Errorf("Highest compression level is 6, got: %d", v))
						}
						return warns, errs
					},
				},
				"mime": {
					Type:     schema.TypeString,
					Default:  "text/*",
					Optional: true,
				},
			},
		},
	}

	staticHeader := &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Static headers to cache from origin pull",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"http": {
					Type:        schema.TypeString,
					Description: "Header to cache from origin pull, ex. Access-Control-Allow-Origin",
					Optional:    true,
				},
				"origin_pull": {
					Type:        schema.TypeString,
					Description: "Header to pass, ex. Host: %client.request.host%",
					Optional:    true,
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
				accountHash, resourceID, err := ResourceImportParseHashID(d.Id())
				if err != nil {
					return nil, err
				}
				d.Set("account_hash", accountHash)
				d.SetId(resourceID)

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
			"root_id": &schema.Schema{
				Description: "The ID of the configuration if it is the root scope",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"name": &schema.Schema{
				Description: "The name of this configuration scope",
				Type:        schema.TypeString,
				Required:    true,
			},
			"path": &schema.Schema{
				Description: "The path of this configuration scope",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/",
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "CDS",
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
			"enable_origin_pull_logs": &schema.Schema{
				Description: "Whether or not to enable logging for origin pulls",
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
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
			"origin_pull_policy": originPullPolicy, // is []
			"enable_file_segmentation": &schema.Schema{
				Description: "Whether or not to pull origin files in segments",
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
			},
			"gzip_origin_pull": &schema.Schema{
				Description: "Whether or not to request gzip'd content on origin pull",
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
			},
			"origin_persistent_connections": &schema.Schema{
				Description: "Whether or not to maintain persistent connections to the origin",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"origin_pull": &schema.Schema{
				Description: "The redirect action for delivering origin requests to the requester",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "proxy",
				// TODO: Validation of all possible types (proxy, or follow)
			},
			"cache_control": cacheControl, // is []
			"cache_key_modification": &schema.Schema{
				Description: "Whether or not to normalize key path to lower case",
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
			},
			"compression":   compression,
			"static_header": staticHeader, // is []
			"http_methods": &schema.Schema{
				Description: "What to pass through in the methods",
				Type:        schema.TypeString,
				Default:     "*",
				Optional:    true,
			},
			"access_logs": &schema.Schema{
				Description: "Whether or not to enable access logs",
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
			},
			"origin_pull_host": originPullHost,
		},
	}
}

/*
	Create
*/
func resourceConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	if id := d.Get("root_id").(int); id != 0 {
		d.SetId(fmt.Sprintf("%d", id))
		d.Partial(false)
		return resourceConfigurationRead(d, m)
	}

	c := m.(*striketracker.Client)
	conf := configuration.New(c)
	accountHash := d.Get("account_hash").(string)
	hostHash := d.Get("host_hash").(string)

	ctx, cancel := getContext()
	defer cancel()

	newConfiguration := &models.Configuration{
		AccessLogs: &models.AccessLogs{
			Enabled: d.Get("access_logs").(bool),
		},
		Scope: models.Scope{
			Name:     d.Get("name").(string),
			Platform: d.Get("platform").(string),
			Path:     d.Get("path").(string),
		},
	}

	hostnamesList := d.Get("hostnames").([]interface{})
	hostnameList := *buildHostnameList(&hostnamesList)
	if len(hostnameList) > 0 {
		for _, hostname := range hostnameList {
			newConfiguration.Hostname = append(newConfiguration.Hostname, &models.ConfigurationHostname{
				Domain: hostname,
			})
		}
	}

	returnedModel, err := conf.Create(ctx, accountHash, hostHash, newConfiguration)
	if returnedModel != nil {
		if returnedModel.ID != 0 {
			d.SetId(fmt.Sprintf("%d", returnedModel.ID))
		}
	}
	if err != nil {
		return err
	}

	d.Partial(false)

	return resourceConfigurationRead(d, m)
}

/*
	Update
*/
func resourceConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	c := m.(*striketracker.Client)
	conf := configuration.New(c)
	accountHash := d.Get("account_hash").(string)
	hostHash := d.Get("host_hash").(string)

	ctx, cancel := getContext()
	defer cancel()

	newConfiguration := &models.Configuration{
		AccessLogs: &models.AccessLogs{
			Enabled: d.Get("access_logs").(bool),
		},
		Scope: models.Scope{
			Name:     d.Get("name").(string),
			Platform: d.Get("platform").(string),
			Path:     d.Get("path").(string),
		},
	}

	hostnamesList := d.Get("hostnames").([]interface{})
	hostnameList := *buildHostnameList(&hostnamesList)
	if len(hostnameList) > 0 {
		for _, hostname := range hostnameList {
			newConfiguration.Hostname = append(newConfiguration.Hostname, &models.ConfigurationHostname{
				Domain: hostname,
			})
		}
	}

	returnedModel, err := conf.Create(ctx, accountHash, hostHash, newConfiguration)
	if returnedModel != nil {
		if returnedModel.Scope.ID != 0 {
			d.SetId(fmt.Sprintf("%d", returnedModel.Scope.ID))
		}
	}
	if err != nil {
		return err
	}

	d.Partial(false)
	return resourceConfigurationRead(d, m)
}

/*
	Delete
*/
func resourceConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	c := m.(*striketracker.Client)
	conf := configuration.New(c)
	accountHash := d.Get("account_hash").(string)
	hostHash := d.Get("host_hash").(string)
	scopeID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	ctx, cancel := getContext()
	defer cancel()

	err = conf.Delete(ctx, accountHash, hostHash, scopeID, false)
	if err != nil {
		return err
	}

	d.Partial(false)
	d.SetId("")
	return nil
}

/*
	Read
*/
func resourceConfigurationRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*striketracker.Client)
	conf := configuration.New(c)
	accountHash := d.Get("account_hash").(string)
	hostHash := d.Get("host_hash").(string)
	scopeID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	ctx, cancel := getContext()
	defer cancel()

	confResource, err := conf.Get(ctx, accountHash, hostHash, scopeID)
	if err != nil {
		return err
	}

	if confResource == nil {
		return fmt.Errorf("Resource %s does not exist", d.Id())
	}

	if &confResource.Scope == nil {
		return fmt.Errorf("Resource %s has problem with Scope", d.Id())
	}

	if confResource.AccessLogs == nil {
		return fmt.Errorf("Resource %s has problem with access_logs", d.Id())
	}

	devLog("%v", confResource.Name)

	d.Set("name", confResource.Name)
	d.Set("platform", confResource.Platform)
	d.Set("path", confResource.Path)
	d.Set("access_logs", confResource.AccessLogs.Enabled)
	//d.Set("path", confResource.Path)

	return nil
}

/*
	Exists
*/
func resourceConfigurationExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*striketracker.Client)
	conf := configuration.New(c)
	accountHash := d.Get("account_hash").(string)
	hostHash := d.Get("host_hash").(string)
	scopeID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}
	ctx, cancel := getContext()
	defer cancel()

	confResource, err := conf.Get(ctx, accountHash, hostHash, scopeID)
	if err != nil {
		return false, nil
	}

	if confResource == nil {
		return false, nil
	}

	return true, nil
}

/*
	Helpers
*/

func buildHostnameList(terraformHostnameList *[]interface{}) *[]string {
	hostnames := make([]string, len(*terraformHostnameList))
	for i, hostname := range *terraformHostnameList {
		hostnames[i] = hostname.(string)
	}
	return &hostnames
}

func hostnameListToObject(hostnameList *[]string) []*models.ConfigurationHostname {
	hostnames := make([]*models.ConfigurationHostname, len(*hostnameList))
	for _, hostname := range *hostnameList {
		hostnames = append(hostnames, &models.ConfigurationHostname{
			Domain: hostname,
		})
	}
	return hostnames
}
