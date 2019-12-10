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
	var err error

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

	// Append stale cache extension (origin pull cache extension)
	if d.HasChange("stale_cache_extension") {
		sce := d.Get("stale_cache_extension").(*schema.Set).List()
		config.OriginPullCacheExtensionFromState(sce[0].(map[string]interface{}))
		d.SetPartial("stale_cache_extension")
	}

	// Append Origin Pull Policy (cache_policy)
	if d.HasChange("cache_policy") {
		cp := d.Get("cache_policy").(*schema.Set).List()
		err = config.OriginPullPolicyFromState(cp)
		if err != nil {
			return nil, err
		}
		d.SetPartial("cache_policy")
	}

	// Append Origin Request Edge Rule
	if d.HasChange("origin_request_edge_rule") {
		orer := d.Get("origin_request_edge_rule").(*schema.Set).List()
		err = config.OriginRequestModificationFromState(orer)
		if err != nil {
			return nil, err
		}
		d.SetPartial("origin_request_edge_rule")
	}

	// Append Origin Response Edge Rule
	if d.HasChange("origin_response_edge_rule") {
		orer := d.Get("origin_response_edge_rule").(*schema.Set).List()
		err = config.OriginResponseModificationFromState(orer)
		if err != nil {
			return nil, err
		}
		d.SetPartial("origin_response_edge_rule")
	}

	// Append Client Request Edge Rule
	if d.HasChange("client_request_edge_rule") {
		orer := d.Get("client_request_edge_rule").(*schema.Set).List()
		err = config.ClientRequestModificationFromState(orer)
		if err != nil {
			return nil, err
		}
		d.SetPartial("client_request_edge_rule")
	}

	// Append Client Response Edge Rule
	if d.HasChange("client_response_edge_rule") {
		orer := d.Get("client_response_edge_rule").(*schema.Set).List()
		err = config.ClientResponseModificationFromState(orer)
		if err != nil {
			return nil, err
		}
		d.SetPartial("client_response_edge_rule")
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
