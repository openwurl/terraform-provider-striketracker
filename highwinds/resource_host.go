package highwinds

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/striketracker"
	"github.com/openwurl/wurlwind/striketracker/models"
	"github.com/openwurl/wurlwind/striketracker/services/hosts"
)

func resourceHost() *schema.Resource {
	scopeList := &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    false,
		Description: "The scopes that have been attached to this service",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"platform": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The platform this scope operates on",
				},
				"path": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The path this scope routes",
				},
				"id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The ID of the scope that is attached",
				},
			},
		},
	}

	servicesList := &schema.Schema{
		Description: "The enabled services for the host",
		Type:        schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
		Required: true,
	}

	return &schema.Resource{
		Create: resourceHostCreate,
		Read:   resourceHostRead,
		Update: resourceHostUpdate,
		Delete: resourceHostDelete,
		Exists: resourceHostExists,
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
			"name": &schema.Schema{
				Description: "The name of the host",
				Type:        schema.TypeString,
				Required:    true,
			},
			"hash_code": &schema.Schema{
				Description: "The hash code pointer to the host",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"services": servicesList,
			"scopes":   scopeList,
			"type": &schema.Schema{
				Description: "The type of host",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"root_scope_id": &schema.Schema{
				Description: "The ID of the root CDS scope",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    false,
			},
		},
	}
}

/*
	Create
*/
func resourceHostCreate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	c := m.(*striketracker.Client)
	h := hosts.New(c)
	accountHash := d.Get("account_hash").(string)

	host := &models.Host{
		Name: d.Get("name").(string),
	}
	servicesList := d.Get("services").([]interface{})
	serviceList := *buildServiceList(&servicesList)

	if len(serviceList) > 0 {
		for _, service := range servicesList {
			host.Services = append(host.Services, &models.DeliveryService{
				ID: service.(int),
			})
		}
	}

	ctx, cancel := getContext()
	defer cancel()

	returnedModel, err := h.Create(ctx, accountHash, host)
	if returnedModel != nil {
		if returnedModel.HashCode != "" {
			d.SetId(returnedModel.HashCode)
			d.Set("root_scope_id", fmt.Sprintf("%d", returnedModel.GetCDSScope().ID))
		}
	}
	if err != nil {
		return err
	}

	d.Partial(false)

	return resourceHostRead(d, m)
}

/*
	Update
*/
func resourceHostUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	c := m.(*striketracker.Client)
	h := hosts.New(c)
	accountHash := d.Get("account_hash").(string)

	host := &models.Host{
		Name: d.Get("name").(string),
	}
	servicesList := d.Get("services").([]interface{})
	serviceList := *buildServiceList(&servicesList)

	if len(serviceList) > 0 {
		for _, service := range servicesList {
			host.Services = append(host.Services, &models.DeliveryService{
				ID: service.(int),
			})
		}
	}

	ctx, cancel := getContext()
	defer cancel()

	returnedModel, err := h.Update(ctx, accountHash, d.Id(), host)
	if returnedModel != nil {
		if returnedModel.HashCode != "" {
			d.SetId(returnedModel.HashCode)
			d.Set("root_scope_id", fmt.Sprintf("%d", returnedModel.GetCDSScope().ID))
		}
	}
	if err != nil {
		return err
	}

	d.Partial(false)
	return resourceHostRead(d, m)
}

/*
	Delete
*/
func resourceHostDelete(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	c := m.(*striketracker.Client)
	h := hosts.New(c)
	accountHash := d.Get("account_hash").(string)

	ctx, cancel := getContext()
	defer cancel()

	err := h.Delete(ctx, accountHash, d.Id())
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
func resourceHostRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*striketracker.Client)
	h := hosts.New(c)
	accountHash := d.Get("account_hash").(string)

	ctx, cancel := getContext()
	defer cancel()

	devLog("Fetching host %s", d.Id())

	hostResource, err := h.Get(ctx, accountHash, d.Id())
	if err != nil {
		return err
	}

	if hostResource == nil {
		return fmt.Errorf("Resource %s does not exist", d.Id())
	}

	d.Set("root_scope_id", fmt.Sprintf("%d", hostResource.GetCDSScope().ID))
	d.Set("name", hostResource.Name)
	d.Set("hash_code", hostResource.HashCode)
	d.Set("services", hostResource.Services)

	scopesList := buildScopesList(hostResource.Scopes)
	if err := d.Set("scopes", scopesList); err != nil {
		return fmt.Errorf("error setting scopes on %s: %v", hostResource.Name, err)
	}
	d.Set("type", hostResource.Type)

	return nil
}

/*
	Exists
*/
func resourceHostExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*striketracker.Client)
	h := hosts.New(c)
	accountHash := d.Get("account_hash").(string)

	ctx, cancel := getContext()
	defer cancel()

	hostResource, err := h.Get(ctx, accountHash, d.Id())
	if err != nil {
		return false, nil
	}

	if hostResource == nil {
		return false, nil
	}

	return true, nil
}

func buildServiceList(terraformServiceList *[]interface{}) *[]int {
	hostScopeList := make([]int, len(*terraformServiceList))
	for i, serviceID := range *terraformServiceList {
		hostScopeList[i] = serviceID.(int)
	}
	return &hostScopeList
}

func buildScopesList(scopes []*models.Scope) []map[string]string {
	scopesList := []map[string]string{}
	for _, scope := range scopes {
		sc := map[string]string{}
		sc["id"] = fmt.Sprintf("%d", scope.ID)
		sc["platform"] = scope.Platform
		sc["path"] = scope.Path
		scopesList = append(scopesList, sc)
	}
	return scopesList
}
