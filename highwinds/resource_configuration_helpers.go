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
	}

	if v, ok := d.GetOk("hostnames"); ok {
		config.Hostname = models.ScopeHostnameFromInterfaceSlice(v.([]interface{}))
	}

	if v, ok := d.GetOk("origin_pull_host"); ok {
		config.OriginPullHost = models.StructFromMap(&models.OriginPullHost{}, getMapFromZeroedSet(v)).(*models.OriginPullHost)
	}

	if v, ok := d.GetOk("stale_cache_extension"); ok {
		config.OriginPullCacheExtension = models.StructFromMap(&models.OriginPullCacheExtension{}, v.(*schema.Set).List()[0].(map[string]interface{})).(*models.OriginPullCacheExtension)
	}

	// weighted set
	if v, ok := d.GetOk("cache_policy"); ok {
		models, err := expandOriginPullPolicies(getSliceIfaceFromSet(v))
		if err != nil {
			return nil, err
		}
		config.OriginPullPolicy = models
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

	// Delivery is a complex set
	if v, ok := d.GetOk("delivery"); ok {
		delivery := expandDeliverySet(v)

		debug.Log("DEBUG CHECK", "SETTING MODEL")

		if compression, ok := delivery["compression"]; ok {
			config.Compression = models.StructFromMap(&models.Compression{}, getMapFromZeroedSet(compression)).(*models.Compression)
		}

		if staticHeader, ok := delivery["static_header"]; ok {
			sh, err := expandDeliveryStaticHeader(getSliceIfaceFromSet(staticHeader))
			if err != nil {
				return nil, err
			}
			config.StaticHeader = sh
		}

		if httpMethods, ok := delivery["http_methods"]; ok {
			config.HTTPMethods = models.StructFromMap(&models.HTTPMethods{}, getMapFromZeroedSet(httpMethods)).(*models.HTTPMethods)
		}

		if responseHeader, ok := delivery["response_header"]; ok {
			config.ResponseHeader = models.StructFromMap(&models.ResponseHeader{}, getMapFromZeroedSet(responseHeader)).(*models.ResponseHeader)
		}

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

	// Delivery is a complex set
	err = d.Set("delivery", compressDeliverySet(config))
	if err != nil {
		errs = append(errs, err)
	}

	// Set origin_pull_host - unweighted set
	err = d.Set("origin_pull_host", []interface{}{models.MapFromStruct(config.OriginPullHost)})
	if err != nil {
		errs = append(errs, err)
	}

	// Set stale_cache_extension - unweighted set
	err = d.Set("stale_cache_extension", []interface{}{models.MapFromStruct(config.OriginPullCacheExtension)})
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

	return errs
}
