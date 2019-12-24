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

	// TODO: Need to make this less of a spiderweb somehow
	if v, ok := d.GetOk("delivery"); ok {
		delivery := expandDeliverySet(v)
		if compression, ok := delivery["compression"]; ok {
			config.Compression = expandDeliveryCompression(compression)
		}
		if staticHeader, ok := delivery["static_header"]; ok {
			// TODO: Must implement a weighting like OriginPullPolicy, order matters
			config.StaticHeader = expandDeliveryStaticHeaders(staticHeader)
		}

	}
	// delivery is a set []interface{}
	// filled with map[string]interface{}
	// filled with sets []interface

	// cast delivery to []interface{}
	/*
				// get compression slice from the map
				if compression, ok := deliverySlice["compression"].(*schema.Set); ok {

					// get the first item in the set list (max items 1)
					compressionSet := compression.List()[0]

					// cast compression to a map
					if compressionSlice, ok := compressionSet.(map[string]interface{}); ok {

						// handle compression
						debug.Log("COMPRESSION", fmt.Sprintf("%v", compressionSlice))
						config.Compression = expandDeliveryCompression(compressionSlice)
						// map[enabled:true gzip:m3u8,ts level:2 mime:text/*,application/x-mpegUR,vnd.apple.mpegURL,video/MP2T]
					}

				}
		if test, ok := v.(*schema.Set); ok {
			thing := test.List()
			if _, ok2 := thing[0].(map[string]interface{}); ok2 {

			} else {
				return nil, fmt.Errorf("failed to cast set to slice interface")
			}
		} else {
			return nil, fmt.Errorf("failed to extract Set list: %v", v)
		}
	*/
	/*
		if deliverySlice, dok := v.(*schema.Set).List()[0].([]interface{}); dok {

			// iterate all items in delivery []interface{}
			for _, thisDeliveryItem := range deliverySlice {

				// cast item to map[string]interface{}
				if thisField, fok := thisDeliveryItem.(map[string]interface{}); fok {

					// extract single field (compression), cast its slice,
					// pick first element (max items: 1), and cast that to map
					if compression, cok := thisField["compression"].([]interface{})[0].(map[string]interface{}); cok {
						// deal with compression
						debug.Log("TESTING COMPRESSION", "%v", fmt.Sprintf("%v", compression))

					} else {
						return nil, fmt.Errorf("failed to find/process compression field")
					}

				} else {
					return nil, fmt.Errorf("failed to parse map item in delivery slice")
				}

			}

		} else {
			return nil, fmt.Errorf("failed to parse delivery slice parent")
		}
	*/
	/*
		for _, field := range deliverySlice {
			thisField, fok := field.(map[string]interface{})
			if !fok {
				return nil, fmt.Errorf("Failed to parse fields in delivery")
			}

			if compression, cok := thisField["compression"]; cok {

			}

		}
	*/
	// extract
	// compression
	// static header
	// http methods
	// gzip origin pull
	//}

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

	/*
		Delivery
	*/
	err = d.Set("delivery", compressDeliverySet(config))
	// needs to be fully prepacked
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}
