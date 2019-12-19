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
	var err error

	config := &models.Configuration{
		Scope:                    expandScopeModel(d.Get("scope").(map[string]interface{})),
		Hostname:                 models.ScopeHostnameFromInterfaceSlice(d.Get("hostnames").([]interface{})),
		OriginPullHost:           expandOriginPullHost(d.Get("origin_pull_host").(*schema.Set).List()[0].(map[string]interface{})),
		OriginPullCacheExtension: expandOriginPullCacheExtension(d.Get("stale_cache_extension").(*schema.Set).List()[0].(map[string]interface{})),
	}

	if v, ok := d.GetOk("cache_policy"); ok {
		err = config.OriginPullPolicyFromState(v.(*schema.Set).List())
		if err != nil {
			return nil, err
		}
	}

	if v, ok := d.GetOk("origin_request_edge_rule"); ok {
		err = config.OriginRequestModificationFromState(v.(*schema.Set).List())
		if err != nil {
			return nil, err
		}
	}
	if v, ok := d.GetOk("origin_response_edge_rule"); ok {
		err = config.OriginResponseModificationFromState(v.(*schema.Set).List())
		if err != nil {
			return nil, err
		}
	}
	if v, ok := d.GetOk("client_request_edge_rule"); ok {
		err = config.ClientRequestModificationFromState(v.(*schema.Set).List())
		if err != nil {
			return nil, err
		}
	}
	if v, ok := d.GetOk("client_response_edge_rule"); ok {
		err = config.ClientResponseModificationFromState(v.(*schema.Set).List())
		if err != nil {
			return nil, err
		}
	}

	/*
		Delivery
	*/

	/*
		// Append Compression
		if d.HasChange("delivery") {
			test1 := d.Get("delivery").(*schema.Set).List()
			dcer := d.Get("delivery.compression").(map[string]interface{})
			debug.Log("THING!!!!", "%v", spew.Sprintf("%v", test1))
			debug.Log("THING!!!!", "%v", spew.Sprintf("%v", dcer))

		}
	*/

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

	// Set stale_cache_extension
	err = d.Set("stale_cache_extension", config.OriginPullCacheExtensionFromModel())
	if err != nil {
		errs = append(errs, err)
	}

	// Set cache_policy (origin pull policy)
	err = d.Set("cache_policy", config.OriginPullPolicyFromModel())
	if err != nil {
		errs = append(errs, err)
	}

	// Set origin_request_edge_rule
	err = d.Set("origin_request_edge_rule", config.OriginRequestModificationFromModel())
	if err != nil {
		errs = append(errs, err)
	}

	// Set origin_response_edge_rule
	err = d.Set("origin_response_edge_rule", config.OriginResponseModificationFromModel())
	if err != nil {
		errs = append(errs, err)
	}

	// Set client_request_edge_rule
	err = d.Set("client_request_edge_rule", config.ClientRequestModificationFromModel())
	if err != nil {
		errs = append(errs, err)
	}

	// Set client_response_edge_rule
	err = d.Set("client_response_edge_rule", config.ClientResponseModificationFromModel())
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}
