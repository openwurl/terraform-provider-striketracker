package highwinds

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/pkg/debug"
	"github.com/openwurl/wurlwind/striketracker/models"
)

// buildNewConfigurationFromState builds a configuration model from terraform state
func buildNewConfigurationFromState(d *schema.ResourceData) (*models.NewHostConfiguration, error) {
	// Pull state and process
	scopeMapRaw := d.Get("scope").(map[string]interface{})
	newConfigModel, err := models.NewHostConfigurationFromState(scopeMapRaw)
	if err != nil {
		return nil, err
	}

	return newConfigModel, nil
}

// buildConfigurationFromState builds a configuration model from terraform state
func buildConfigurationFromState(d *schema.ResourceData) (*models.Configuration, error) {
	//	var err error

	config := &models.Configuration{}

	if v, ok := d.GetOk("scope"); ok {
		config.Scope = models.StructFromMap(&models.Scope{}, v.(map[string]interface{})).(*models.Scope)
		//return nil, fmt.Errorf("Erroring on purpose: %s", "scope")
	}

	if v, ok := d.GetOk("hostnames"); ok {
		config.Hostname = models.ScopeHostnameFromInterfaceSlice(v.([]interface{}))
	}

	if v, ok := d.GetOk("delivery"); ok {
		config.IngestSchema(expandSetOfMaps(v))
	}

	if v, ok := d.GetOk("origin"); ok {
		config.IngestSchema(expandSetOfMaps(v))
	}

	debug.Log("STATE", "%v", spew.Sprintf("%v", config.Scope))

	return config, config.Validate()
}

// ingestState updates terraform state from a configuration model
func ingestState(d *schema.ResourceData, config *models.Configuration) []error {
	var errs []error

	// Set scope details
	err := d.Set("scope", models.MapFromStruct(config.Scope))
	if err != nil {
		errs = append(errs, err)
	}

	// Set hostnames
	err = d.Set("hostnames", config.HostnamesFromModel())
	if err != nil {
		errs = append(errs, err)
	}

	// Should only be used for schema.Set's of schema.Sets
	content := config.ExtractSchema() //[0].(map[string]interface{})

	debug.Log("DELIVERY", "%v", spew.Sprintf("%v", content["delivery"]))

	// TODO: returned as map not interface, take a look if this is intended
	err = d.Set("delivery", []interface{}{content["delivery"]})
	if err != nil {
		errs = append(errs, err)
	}

	err = d.Set("origin", []interface{}{content["origin"]})
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}
