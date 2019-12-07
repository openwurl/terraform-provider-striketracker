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
	config := models.NewConfiguration()

	// Pull state and process
	if d.HasChange("scope") {
		scopeMapRaw := d.Get("scope").(map[string]interface{})
		config.ScopeFromState(scopeMapRaw)
		d.SetPartial("scope")
	}

	//  Append hostnames to model
	if d.HasChange("hostnames") {
		hostnamesList := d.Get("hostnames").([]interface{})
		config.HostnamesFromState(hostnamesList)
		d.SetPartial("hostnames")
	}

	// Append OriginHost to model
	if d.HasChange("origin_pull_host") {
		originHost := d.Get("origin_pull_host").(*schema.Set).List()
		config.OriginHostFromState(originHost[0].(map[string]interface{}))
		d.SetPartial("origin_pull_host")
	}

	// Append stale cache extension
	if d.HasChange("stale_cache_extension") {
		sce := d.Get("stale_cache_extension").(*schema.Set).List()
		config.OriginPullCacheExtensionFromState(sce[0].(map[string]interface{}))
		d.SetPartial("stale_cache_extension")
	}

	debug.Log("STATE", "%v", spew.Sprintf("%v", config.Scope))

	return config, config.Validate()
}

// ingestState updates terraform state from a configuration model
func ingestState(d *schema.ResourceData, config *models.Configuration) []error {
	var errs []error

	// Set scope details
	err := d.Set("scope", config.ScopeFromModel())
	if err != nil {
		errs = append(errs, err)
	}

	// Set hostnames
	err = d.Set("hostnames", config.HostnamesFromModel())
	if err != nil {
		errs = append(errs, err)
	}

	// Set origin_pull_host
	err = d.Set("origin_pull_host", config.OriginHostFromModel())
	if err != nil {
		errs = append(errs, err)
	}

	// set stale_cache_extension
	err = d.Set("stale_cache_extension", config.OriginPullCacheExtensionFromModel())
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}
