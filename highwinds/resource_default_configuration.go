package highwinds

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/pkg/debug"
	"github.com/openwurl/wurlwind/striketracker"
	"github.com/openwurl/wurlwind/striketracker/services/hosts"
)

func defaultResourceConfiguration() *schema.Resource {
	drc := resourceConfiguration()
	drc.Create = resourceDefaultConfigurationCreate
	drc.Delete = resourceDefaultConfigurationDelete

	return drc
}

func resourceDefaultConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	// Fetch defined host
	c := m.(*striketracker.Client)
	h := hosts.New(c)
	accountHash := d.Get("account_hash").(string)
	hostHash := d.Get("host_hash").(string)

	debug.Log("Fetch", "Fetching %s/%s", accountHash, hostHash)

	ctx, cancel := getContext()
	defer cancel()

	hostResource, err := h.Get(ctx, accountHash, hostHash)
	if err != nil {
		return err
	}

	if hostResource == nil {
		return fmt.Errorf("Resource %s does not exist", d.Id())
	}

	// Extract root scope

	rootScope := hostResource.GetCDSScope()
	if rootScope == nil {
		return fmt.Errorf("Could not fetch Root Scope on parent host: %v", rootScope.Name)
	}

	// Set ID from root scope
	d.SetId(rootScope.GetIDString())

	// Return an update on the resource
	debug.Log("Update", "Updating instead of creating %s/%s", accountHash, hostHash)

	return resourceConfigurationUpdate(d, m)
}

func resourceDefaultConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[WARN] Cannot destroy Default Scope Configuration. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
