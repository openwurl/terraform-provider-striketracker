package highwinds

import (
	"strconv"

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
func buildNewConfigurationFromState(d *schema.ResourceData) (*models.Configuration, error) {

	return nil, nil
}

// buildConfigurationFromState builds a configuration for updating a scope
func buildConfigurationFromState(d *schema.ResourceData) (*models.Configuration, error) {

	return nil, nil
}

// REPLACE ME
func buildCreateScopeConfiguration(d *schema.ResourceData) *models.ConfigurationCreate {
	// Pull scope resource from HCL and process the interface
	scopeMapRaw := d.Get("scope").(map[string]interface{})

	scopeMap := buildHostScopeList(scopeMapRaw)

	// Weird bugfix because default isn't appearing in state sometimes
	if scopeMap["platform"] == "" {
		scopeMap["platform"] = "CDS"
	}

	// Create base model
	newConfigurationScope := &models.ConfigurationCreate{
		Name:     scopeMap["name"],
		Platform: scopeMap["platform"],
		Path:     scopeMap["path"],
		//OriginPullHost: originHost,
	}

	// Append hostnames to model
	hostnamesList := d.Get("hostnames").([]interface{})
	hostnameList := *buildHostnameList(&hostnamesList)
	if len(hostnameList) > 0 {
		for _, hostname := range hostnameList {
			newConfigurationScope.Hostname = append(newConfigurationScope.Hostname, &models.ConfigurationHostname{
				Domain: hostname,
			})
		}
	}

	return newConfigurationScope
}

// REPLACE ME
func buildScopeConfiguration(d *schema.ResourceData) *models.Configuration {
	// Pull scope resource from HCL and process the interface
	scopeMapRaw := d.Get("scope").(map[string]interface{})

	scopeMap := buildHostScopeList(scopeMapRaw)

	// Weird bugfix because default isn't appearing in state sometimes
	if scopeMap["platform"] == "" {
		scopeMap["platform"] = "CDS"
	}

	// Create base model
	newConfigurationScope := &models.Configuration{
		Scope: &models.Scope{
			Name:     scopeMap["name"],
			Platform: scopeMap["platform"],
			Path:     scopeMap["path"],
		},
	}

	/*
		Attach origin pull protocol
	*/
	originRaw := d.Get("origin").(map[string]interface{})
	originHost := buildOriginMap(originRaw)
	if ra := originHost["redirect_action"]; ra != "" {
		newConfigurationScope.OriginPull = &models.OriginPull{
			RedirectAction: ra,
		}
	} else {
		devLog("Origin pull redirect action not defined")
	}

	// Attach origin redirect action
	if opp := originHost["origin_pull_protocol"]; opp != "" {
		newConfigurationScope.OriginPullProtocol = &models.OriginPullProtocol{
			Protocol: opp,
		}
	} else {
		devLog("Origin pull protocol not defined")
	}

	/*
		Attach primary/secondary origin host
	*/
	originModel := &models.OriginPullHost{}
	if ohPrimary, err := strconv.Atoi(originHost["primary"]); err == nil {
		originModel.Primary = ohPrimary
	} else {
		devLog("Failed to parse primary: %v", err)
	}
	if ohSecondary, err := strconv.Atoi(originHost["secondary"]); err == nil {
		originModel.Secondary = ohSecondary
	} else {
		devLog("Failed to parse secondary: %v", err)
	}
	if ohPath := originHost["path"]; ohPath != "" {
		originModel.Path = ohPath
	}
	newConfigurationScope.OriginPullHost = originModel

	/*
		Append hostnames to model
	*/
	hostnamesList := d.Get("hostnames").([]interface{})
	hostnameList := *buildHostnameList(&hostnamesList)
	if len(hostnameList) > 0 {
		for _, hostname := range hostnameList {
			newConfigurationScope.Hostname = append(newConfigurationScope.Hostname, &models.ConfigurationHostname{
				Domain: hostname,
			})
		}
	}

	/*
		Append origin pull policies to the model
	*/
	originPullPolicyListRaw := d.Get("origin_pull_policy").([]interface{})
	originPullPolicyList := models.BuildOriginPullPoliciesList(&originPullPolicyListRaw)
	if len(originPullPolicyList) > 0 {
		for _, policy := range originPullPolicyList {
			newConfigurationScope.OriginPullPolicy = append(newConfigurationScope.OriginPullPolicy, policy)
		}
	}

	/*
		Add Edge rules
	*/
	cReqMod := d.Get("client_request_edge_rule").([]interface{})
	if len(cReqMod) > 0 {
		cReqs := models.BuildClientRequestModification(cReqMod)
		newConfigurationScope.ClientRequestModification = cReqs
	}

	cRespMod := d.Get("client_response_edge_rule").([]interface{})
	if len(cRespMod) > 0 {
		cReqs := models.BuildClientResponseModification(cRespMod)
		newConfigurationScope.ClientResponseModification = cReqs
	}
	oReqMod := d.Get("origin_request_edge_rule").([]interface{})
	if len(oReqMod) > 0 {
		cReqs := models.BuildOriginRequestModification(oReqMod)
		newConfigurationScope.OriginRequestModification = cReqs
	}
	oRespMod := d.Get("origin_response_edge_rule").([]interface{})
	if len(oRespMod) > 0 {
		cReqs := models.BuildOriginResponseModification(oRespMod)
		newConfigurationScope.OriginResponseModification = cReqs
	}

	return newConfigurationScope
}
