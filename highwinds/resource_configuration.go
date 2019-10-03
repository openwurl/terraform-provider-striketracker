package highwinds

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/pkg/utilities"
	"github.com/openwurl/wurlwind/striketracker/models"
)

func resourceConfiguration() *schema.Resource {
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
			"parent_host": &schema.Schema{
				Description: "The hash code of the parent host this scope is being attached to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": &schema.Schema{
				Description: "The name of this configuration scope",
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
			// "origin_pull_policy": NEEDS TO BE RESOURCE BUT WITHOUT REMOTE RESOURCE?
			"enable_file_segmentation": &schema.Schema{
				Description: "Whether or not to pull origin files in segments",
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
			},
		},
	}
}

/*
	Create
*/
func resourceConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	return resourceConfigurationRead(d, m)
}

/*
	Update
*/
func resourceConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceConfigurationRead(d, m)
}

/*
	Delete
*/
func resourceConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
	Read
*/
func resourceConfigurationRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
	Exists
*/
func resourceConfigurationExists(d *schema.ResourceData, m interface{}) (bool, error) {
	return false, nil
}
