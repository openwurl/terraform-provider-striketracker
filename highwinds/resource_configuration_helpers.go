package highwinds

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/striketracker/models"
)

/*
	Helpers
*/

// buildHostnameList converts the hostname state object to a list of strings
func buildHostnameList(terraformHostnameList *[]interface{}) *[]string {
	hostnames := make([]string, len(*terraformHostnameList))
	for i, hostname := range *terraformHostnameList {
		hostnames[i] = hostname.(string)
	}
	return &hostnames
}

func buildHostScopeList(scopes map[string]interface{}) map[string]string {
	sc := make(map[string]string)
	for scopeKey, scopeVal := range scopes {
		sc[scopeKey] = scopeVal.(string)
	}
	return sc
}

func buildOriginMap(originMap map[string]interface{}) map[string]string {
	raw := make(map[string]string)

	for mapKey, mapVal := range originMap {
		raw[mapKey] = mapVal.(string)
	}

	return raw
}

// ingestRemoteState picks apart the model from the remote and sets terraform state accordingly
func ingestRemoteState(d *schema.ResourceData, config *models.Configuration) []error {
	var errs []error

	// Set scope details
	err := d.Set("scope", config.BuildHostScopeInterface())
	if err != nil {
		errs = append(errs, err)
	}

	// Set hostnames both complete and pared down
	err = d.Set("hostnames", config.HostnamesAsStringSlice())
	if err != nil {
		errs = append(errs, err)
	}
	err = d.Set("dns_hostnames", config.ActionableHostnamesAsStringSlice())
	if err != nil {
		errs = append(errs, err)
	}

	// If OriginPull is nil we still need an empty string set
	if config.OriginPull == nil {
		config.OriginPull = &models.OriginPull{
			RedirectAction: "",
		}
	}

	// Set Origin configuration
	err = d.Set("origin", config.BuildOriginInterface())
	if err != nil {
		errs = append(errs, err)
	}

	// Set origin pull policies
	err = d.Set("origin_pull_policy", config.BuildOriginPullPoliciesInterface())
	if err != nil {
		errs = append(errs, err)
	}

	// Set log bools
	err = d.Set("logs", config.GetLogState())
	if err != nil {
		errs = append(errs, err)
	}

	// Set browser cache control
	err = d.Set("cache_control", config.BuildCacheControlInterface())
	if err != nil {
		errs = append(errs, err)
	}

	// TODO: Redo delivery
	// Set delivery
	//err = d.Set("delivery", config.BuildDeliveryMap())
	//if err != nil {
	//	errs = append(errs, err)
	//}

	// Set cachekeys
	err = d.Set("cache_keys", config.BuildCacheKeyMap())
	if err != nil {
		errs = append(errs, err)
	}

	// Set Access an Origin Pull logs
	err = d.Set("logs", config.BuildLogMap())
	if err != nil {
		errs = append(errs, err)
	}

	// Set EdgeRules
	err = d.Set("client_request_edge_rule", config.ClientRequestMap())
	if err != nil {
		errs = append(errs, err)
	}
	err = d.Set("client_response_edge_rule", config.ClientResponseMap())
	if err != nil {
		errs = append(errs, err)
	}
	err = d.Set("origin_request_edge_rule", config.OriginRequestMap())
	if err != nil {
		errs = append(errs, err)
	}
	err = d.Set("origin_response_edge_rule", config.OriginResponseMap())
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}

// buildNewConfigurationFromState builds a configuration only used to create new scopes with
func buildNewConfigurationFromState(d *schema.ResourceData) (*models.ConfigurationCreate, error) {
	// Pull scope resource from HCL and process the interface into base config
	scopeMapRaw := d.Get("scope").(map[string]interface{})
	newConfigScope, err := models.NewCreateConfigurationFromScope(scopeMapRaw)
	if err != nil {
		return nil, err
	}

	//  Append hostnames to model
	hostnamesList := d.Get("hostnames").([]interface{})
	newConfigScope.AppendHostnames(hostnamesList)

	return newConfigScope, nil
}

// buildConfigurationFromState builds a configuration for updating a scope
func buildConfigurationFromState(d *schema.ResourceData) (*models.Configuration, error) {
	// Pull scope resource from HCL and process the interface into base config
	scopeMapRaw := d.Get("scope").(map[string]interface{})
	newConfigScope, err := models.ConfigurationFromScope(scopeMapRaw)
	if err != nil {
		return nil, err
	}

	// Attach Compression, HTTPMethods, StaticHeader from delivery map
	deliveryMap := d.Get("delivery").([]interface{})
	newConfigScope.IngestDeliveryMap(deliveryMap)
	return nil, fmt.Errorf("failing on purpose")

	// Attach CacheKeyModification from cache_keys map
	cacheKeysMap := d.Get("cache_keys").(map[string]interface{})
	newConfigScope.IngestCacheKeyMap(cacheKeysMap)

	// Attach AccessLogs and OriginPullLogs from logs map
	logsMap := d.Get("logs").(map[string]interface{})
	newConfigScope.IngestLogMap(logsMap)

	// Attach hostnames
	hostnamesList := d.Get("hostnames").([]interface{})
	newConfigScope.IngestHostnames(hostnamesList)

	// Attach origin details map
	originMap := d.Get("origin").(map[string]interface{})
	newConfigScope.IngestOriginMap(originMap)

	// Attach origin pull policy
	originPullPolicyListRaw := d.Get("origin_pull_policy").([]interface{})
	newConfigScope.IngestOriginPullPolicies(originPullPolicyListRaw)

	// attach cache control
	cacheControlListRaw := d.Get("cache_control").([]interface{})
	newConfigScope.IngestCacheControl(cacheControlListRaw)

	// Attach edge rule policies
	cReqMod := d.Get("client_request_edge_rule").([]interface{})
	if len(cReqMod) > 0 {
		newConfigScope.IngestClientRequestModification(cReqMod)
	}
	cRespMod := d.Get("client_response_edge_rule").([]interface{})
	if len(cRespMod) > 0 {
		newConfigScope.IngestClientResponseModification(cRespMod)
	}
	oReqMod := d.Get("origin_request_edge_rule").([]interface{})
	if len(oReqMod) > 0 {
		newConfigScope.IngestOriginRequestModification(oReqMod)
	}
	oRespMod := d.Get("origin_response_edge_rule").([]interface{})
	if len(oRespMod) > 0 {
		newConfigScope.IngestOriginResponseModification(oRespMod)
	}

	return newConfigScope, nil
}
