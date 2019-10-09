package highwinds

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
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
			"scope": scopeSchema,
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

	// Pull scope resource from HCL and process the interface
	scopeMapRaw := d.Get("scope").(map[string]interface{})

	scopeMap := buildHostScopeList(scopeMapRaw)

	// Weird bugfix because default isn't appearing in state sometimes
	if scopeMap["platform"] == "" {
		scopeMap["platform"] = "CDS"
	}

	newConfigurationScope := &models.ConfigurationCreate{
		Name:     scopeMap["name"],
		Platform: scopeMap["platform"],
		Path:     scopeMap["path"],
	}

	devLog("Scope Config: %v", newConfigurationScope)

	returnedModel, err := conf.Create(ctx, accountHash, hostHash, newConfigurationScope)
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

	// Pull scope resource from HCL and process the interface
	scopeMapRaw := d.Get("scope").(map[string]interface{})

	scopeMap := buildHostScopeList(scopeMapRaw)

	// Weird bugfix because default isn't appearing in state
	if scopeMap["platform"] == "" {
		scopeMap["platform"] = "CDS"
	}

	// Build object
	newConfigurationScope := &models.Configuration{
		Scope: &models.Scope{
			Name:     scopeMap["name"],
			Platform: scopeMap["platform"],
			Path:     scopeMap["path"],
		},
	}

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
