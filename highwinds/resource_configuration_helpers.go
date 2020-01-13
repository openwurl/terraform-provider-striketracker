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

	// weighted set
	if v, ok := d.GetOk("cache_policy"); ok {
		models, err := expandOriginPullPolicies(getSliceIfaceFromSet(v))
		if err != nil {
			return nil, err
		}
		config.OriginPullPolicy = models
	}

	if v, ok := d.GetOk("stale_cache_extension"); ok {
		//v.(*schema.Set).List()[0].(map[string]interface{})
		config.OriginPullCacheExtension = models.StructFromMap(&models.OriginPullCacheExtension{}, expandSetOfMaps(v)).(*models.OriginPullCacheExtension)
	}

	// weighted set
	if v, ok := d.GetOk("origin_request_edge_rule"); ok {
		models, err := expandOriginRequestModification(getSliceIfaceFromSet(v))
		if err != nil {
			return nil, err
		}
		config.OriginRequestModification = models
	}

	// weighted set
	if v, ok := d.GetOk("origin_response_edge_rule"); ok {
		models, err := expandOriginResponseModification(getSliceIfaceFromSet(v))
		if err != nil {
			return nil, err
		}
		config.OriginResponseModification = models
	}

	// weighted set
	if v, ok := d.GetOk("client_request_edge_rule"); ok {
		models, err := expandClientRequestModification(getSliceIfaceFromSet(v))
		if err != nil {
			return nil, err
		}
		config.ClientRequestModification = models
	}

	// weighted set
	if v, ok := d.GetOk("client_response_edge_rule"); ok {
		models, err := expandClientResponseModification(getSliceIfaceFromSet(v))
		if err != nil {
			return nil, err
		}
		config.ClientResponseModification = models
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

	err = d.Set("delivery", []interface{}{content["delivery"]})
	if err != nil {
		errs = append(errs, err)
	}

	err = d.Set("origin", []interface{}{content["origin"]})
	if err != nil {
		errs = append(errs, err)
	}

	// Set cache_policy (OriginPullPolicy) - weighted
	err = d.Set("cache_policy", compressOriginPullPolicies(config.OriginPullPolicy))
	if err != nil {
		errs = append(errs, err)
	}

	// Set Edge Rules
	err = d.Set("origin_request_edge_rule", compressOriginRequestModification(config.OriginRequestModification))
	if err != nil {
		errs = append(errs, err)
	}

	err = d.Set("origin_response_edge_rule", compressOriginResponseModification(config.OriginResponseModification))
	if err != nil {
		errs = append(errs, err)
	}

	err = d.Set("client_request_edge_rule", compressClientRequestModification(config.ClientRequestModification))
	if err != nil {
		errs = append(errs, err)
	}

	err = d.Set("client_response_edge_rule", compressClientResponseModification(config.ClientResponseModification))
	if err != nil {
		errs = append(errs, err)
	}

	// Set stale_cache_extension - unweighted set
	err = d.Set("stale_cache_extension", []interface{}{models.MapFromStruct(config.OriginPullCacheExtension)})
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}
