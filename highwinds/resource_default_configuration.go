package highwinds

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func defaultResourceConfiguration() *schema.Resource {
	drc := resourceConfiguration()
	drc.Create = resourceDefaultConfigurationCreate
	drc.Delete = resourceDefaultConfigurationDelete

	return drc
}

func resourceDefaultConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	// validate it exists
	// set ID
	// run an update
	return resourceConfigurationUpdate(d, m)
}

func resourceDefaultConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[WARN] Cannot destroy Default Scope Configuration. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
