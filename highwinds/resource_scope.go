package highwinds

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceScope() *schema.Resource {
	return &schema.Resource{
		Create: resourceScopeCreate,
		Read:   resourceScopeRead,
		Update: resourceScopeUpdate,
		Delete: resourceScopeDelete,
		Exists: resourceScopeExists,
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
			"parent_host": &schema.Schema{
				Description: "The hash code of the parent host this scope is being attached to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"path": &schema.Schema{
				Description: "The name of this configuration scope",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

/*
	Create
*/
func resourceScopeCreate(d *schema.ResourceData, m interface{}) error {
	return resourceScopeRead(d, m)
}

/*
	Update
*/
func resourceScopeUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceScopeRead(d, m)
}

/*
	Delete
*/
func resourceScopeDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
	Read
*/
func resourceScopeRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
	Exists
*/
func resourceScopeExists(d *schema.ResourceData, m interface{}) (bool, error) {
	return false, nil
}
