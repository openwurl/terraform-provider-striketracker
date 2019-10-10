package highwinds

import (
	"encoding/json"
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
			"origin":             originSchema,
			"origin_pull_policy": originPullPolicySchema,
			"scope":              scopeSchema,
		},
	}
}

/*
	Create
*/
func resourceConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	c := m.(*striketracker.Client)
	conf := configuration.New(c)
	accountHash := d.Get("account_hash").(string)
	hostHash := d.Get("host_hash").(string)

	ctx, cancel := getContext()
	defer cancel()

	// Build our model to send
	newConfigurationScope := buildCreateScopeConfiguration(d)

	devLog("Scope Config: %v", newConfigurationScope)
	devLog("Scope Config Pullhost: %v", newConfigurationScope.OriginPullHost)

	test, _ := json.Marshal(newConfigurationScope)
	devLog("BEFORE: %s", string(test))

	// Send model
	returnedModel, err := conf.Create(ctx, accountHash, hostHash, newConfigurationScope)
	if returnedModel != nil {
		if returnedModel.ID != 0 {
			d.SetId(fmt.Sprintf("%d", returnedModel.ID))
		}
	}
	if err != nil {
		return err
	}

	test2, _ := json.Marshal(returnedModel)
	devLog("AFTER: %s", string(test2))

	d.Partial(false)

	return resourceConfigurationUpdate(d, m)
}

/*
	Update
*/
func resourceConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
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

	devLog("Preparing to update configuration %s/%s/%s", accountHash, hostHash, scopeID)

	// Build our model to send
	newConfigurationScope := buildScopeConfiguration(d)

	devLog("Updating configuration %s/%s/%d", accountHash, hostHash, scopeID)
	// Ship object
	returnedModel, err := conf.Update(ctx, accountHash, hostHash, scopeID, newConfigurationScope)
	if err != nil {
		return err
	}
	if returnedModel == nil {
		return fmt.Errorf("Something went wrong updating the scope %s, returned model is nil", d.Id())
	}

	return resourceConfigurationRead(d, m)
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

	devLog("Reading configuration %s/%s/%d", accountHash, hostHash, scopeID)

	// Fetch resource
	configModel, err := conf.Get(ctx, accountHash, hostHash, scopeID)
	if err != nil {
		return err
	}
	if configModel == nil {
		return fmt.Errorf("Resource %s does not exist", d.Id())
	}

	devLog("Setting configuration state %s/%s/%d", accountHash, hostHash, scopeID)

	if configModel.Platform == "" || configModel.ID == 0 || configModel.Path == "" {
		return fmt.Errorf("Scope is nil on %s/%s/%d", accountHash, hostHash, scopeID)
	}

	// Set state
	d.Set("scope", buildHostScopeInterface(configModel))
	//if err := d.Set("scope", configModel.Scope); err != nil {
	//	return fmt.Errorf("error setting scope for resource %s: %s", d.Id(), err)
	//}

	devLog("Done setting configuration state %s/%s/%d", accountHash, hostHash, scopeID)

	return nil
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
	Exists
*/
func resourceConfigurationExists(d *schema.ResourceData, m interface{}) (bool, error) {
	return false, nil
}

/*
	Helpers
*/

func buildOriginPullPoliciesList(terraformPullPolicyList *[]interface{}) []*models.OriginPullPolicy {
	policylist := []*models.OriginPullPolicy{}
	// extract policies
	for _, policy := range *terraformPullPolicyList {
		devLog("POLICY: %v", policy)
		thisMap := policy.(map[string]interface{})

		newPolicy := &models.OriginPullPolicy{
			Enabled:                        thisMap["enabled"].(bool),
			ExpirePolicy:                   thisMap["expire_policy"].(string),
			ExpireSeconds:                  thisMap["expire_seconds"].(int),
			ForceBypassCache:               thisMap["force_bypass_cache"].(bool),
			HonorMustRevalidate:            thisMap["honor_must_revalidate"].(bool),
			HonorNoCache:                   thisMap["honor_no_cache"].(bool),
			HonorNoStore:                   thisMap["honor_no_store"].(bool),
			HonorPrivate:                   thisMap["honor_private"].(bool),
			HonorSMaxAge:                   thisMap["honor_smax_age"].(bool),
			HTTPHeaders:                    thisMap["http_headers"].(string),
			MustRevalidateToNoCache:        thisMap["must_revalidate_to_no_cache"].(bool),
			NoCacheBehavior:                thisMap["no_cache_behavior"].(string),
			UpdateHTTPHeadersOn304Response: thisMap["update_http_headers_on_304_response"].(bool),
			DefaultCacheBehavior:           thisMap["default_cache_behavior"].(string),
			MaxAgeZeroToNoCache:            thisMap["max_age_zero_to_no_cache"].(bool),
			ContentTypeFilter:              thisMap["content_type_filter"].(string),
			HeaderFilter:                   thisMap["header_filter"].(string),
			MethodFilter:                   thisMap["method_filter"].(string),
			PathFilter:                     thisMap["path_filter"].(string),
		}
		policylist = append(policylist, newPolicy)
	}

	return policylist
}

func buildHostnameList(terraformHostnameList *[]interface{}) *[]string {
	hostnames := make([]string, len(*terraformHostnameList))
	for i, hostname := range *terraformHostnameList {
		hostnames[i] = hostname.(string)
	}
	return &hostnames
}

func buildHostScopeList(scopes map[string]interface{}) map[string]string {
	sc := make(map[string]string)
	for scopeKey, scopeVal := range scopes {
		sc[scopeKey] = scopeVal.(string)
	}
	return sc
}

func buildHostScopeInterface(scope *models.Configuration) map[string]interface{} {
	scopeList := make(map[string]interface{})

	scopeList["id"] = string(scope.ID)
	scopeList["platform"] = scope.Platform
	scopeList["path"] = scope.Path

	devLog("Done building scope list: %v", scopeList)

	return scopeList
}

func buildOriginMap(originMap map[string]interface{}) map[string]string {
	raw := make(map[string]string)

	for mapKey, mapVal := range originMap {
		raw[mapKey] = mapVal.(string)
	}

	return raw
}

func buildCreateScopeConfiguration(d *schema.ResourceData) *models.ConfigurationCreate {
	// Pull scope resource from HCL and process the interface
	scopeMapRaw := d.Get("scope").(map[string]interface{})

	scopeMap := buildHostScopeList(scopeMapRaw)

	// Weird bugfix because default isn't appearing in state sometimes
	if scopeMap["platform"] == "" {
		scopeMap["platform"] = "CDS"
	}

	// Create base model
	newConfigurationScope := &models.ConfigurationCreate{
		Name:     scopeMap["name"],
		Platform: scopeMap["platform"],
		Path:     scopeMap["path"],
		//OriginPullHost: originHost,
	}

	originRaw := d.Get("origin").(map[string]interface{})
	originHost := buildOriginMap(originRaw)

	// Attach origin pull protocol
	if ra := originHost["redirect_action"]; ra != "" {
		newConfigurationScope.OriginPull = &models.OriginPull{
			RedirectAction: ra,
		}
	} else {
		devLog("Origin pull redirect action not defined")
	}

	// Attach origin redirect action
	if opp := originHost["origin_pull_protocol"]; opp != "" {
		newConfigurationScope.OriginPullProtocol = &models.OriginPullProtocol{
			Protocol: opp,
		}
	} else {
		devLog("Origin pull protocol not defined")
	}

	// Attach primary/secondary origin host
	originModel := &models.OriginPullHost{}
	if ohPrimary, err := strconv.Atoi(originHost["primary"]); err == nil {
		originModel.Primary = ohPrimary
	} else {
		devLog("Failed to parse primary: %v", err)
	}
	if ohSecondary, err := strconv.Atoi(originHost["secondary"]); err == nil {
		originModel.Secondary = ohSecondary
	} else {
		devLog("Failed to parse secondary: %v", err)
	}
	if ohPath := originHost["path"]; ohPath != "" {
		originModel.Path = ohPath
	}
	newConfigurationScope.OriginPullHost = originModel

	// Append hostnames to model
	hostnamesList := d.Get("hostnames").([]interface{})
	hostnameList := *buildHostnameList(&hostnamesList)
	if len(hostnameList) > 0 {
		for _, hostname := range hostnameList {
			newConfigurationScope.Hostname = append(newConfigurationScope.Hostname, &models.ConfigurationHostname{
				Domain: hostname,
			})
		}
	}

	// Append origin pull policies to the model
	originPullPolicyListRaw := d.Get("origin_pull_policy").([]interface{})
	originPullPolicyList := buildOriginPullPoliciesList(&originPullPolicyListRaw)
	if len(originPullPolicyList) > 0 {
		for _, policy := range originPullPolicyList {
			newConfigurationScope.OriginPullPolicy = append(newConfigurationScope.OriginPullPolicy, policy)
		}
	}

	return newConfigurationScope
}

func buildScopeConfiguration(d *schema.ResourceData) *models.Configuration {
	// Pull scope resource from HCL and process the interface
	scopeMapRaw := d.Get("scope").(map[string]interface{})

	scopeMap := buildHostScopeList(scopeMapRaw)

	// Weird bugfix because default isn't appearing in state sometimes
	if scopeMap["platform"] == "" {
		scopeMap["platform"] = "CDS"
	}

	// Create base model
	newConfigurationScope := &models.Configuration{
		Scope: &models.Scope{
			Name:     scopeMap["name"],
			Platform: scopeMap["platform"],
			Path:     scopeMap["path"],
		},
		//OriginPullHost: originHost,
	}

	originRaw := d.Get("origin").(map[string]interface{})
	originHost := buildOriginMap(originRaw)

	// Attach origin pull protocol
	if ra := originHost["redirect_action"]; ra != "" {
		newConfigurationScope.OriginPull = &models.OriginPull{
			RedirectAction: ra,
		}
	} else {
		devLog("Origin pull redirect action not defined")
	}

	// Attach origin redirect action
	if opp := originHost["origin_pull_protocol"]; opp != "" {
		newConfigurationScope.OriginPullProtocol = &models.OriginPullProtocol{
			Protocol: opp,
		}
	} else {
		devLog("Origin pull protocol not defined")
	}

	// Attach primary/secondary origin host
	originModel := &models.OriginPullHost{}
	if ohPrimary, err := strconv.Atoi(originHost["primary"]); err == nil {
		originModel.Primary = ohPrimary
	} else {
		devLog("Failed to parse primary: %v", err)
	}
	if ohSecondary, err := strconv.Atoi(originHost["secondary"]); err == nil {
		originModel.Secondary = ohSecondary
	} else {
		devLog("Failed to parse secondary: %v", err)
	}
	if ohPath := originHost["path"]; ohPath != "" {
		originModel.Path = ohPath
	}
	newConfigurationScope.OriginPullHost = originModel

	// Append hostnames to model
	hostnamesList := d.Get("hostnames").([]interface{})
	hostnameList := *buildHostnameList(&hostnamesList)
	if len(hostnameList) > 0 {
		for _, hostname := range hostnameList {
			newConfigurationScope.Hostname = append(newConfigurationScope.Hostname, &models.ConfigurationHostname{
				Domain: hostname,
			})
		}
	}

	// Append origin pull policies to the model
	originPullPolicyListRaw := d.Get("origin_pull_policy").([]interface{})
	originPullPolicyList := buildOriginPullPoliciesList(&originPullPolicyListRaw)
	if len(originPullPolicyList) > 0 {
		for _, policy := range originPullPolicyList {
			newConfigurationScope.OriginPullPolicy = append(newConfigurationScope.OriginPullPolicy, policy)
		}
	}

	return newConfigurationScope
}
