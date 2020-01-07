package highwinds

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/striketracker/models"
)

/*
	Helpers
*/

// expandSetOfMaps expands a tf set of maps into its first level map
func expandSetOfMaps(raw interface{}) map[string]interface{} {
	if deliverySet, ok := raw.(*schema.Set); ok {
		set := deliverySet.List()[0]
		if deliverySlice, ok := set.(map[string]interface{}); ok {
			return deliverySlice
		}
	}
	return nil
}

// getMapFromZeroedSet returns the (*schema.Set).List()[0] map if it is one
func getMapFromZeroedSet(set interface{}) map[string]interface{} {
	if v, ok := set.(*schema.Set); ok {
		if len(v.List()) < 1 {
			return nil
		}
		if sv, ok := v.List()[0].(map[string]interface{}); ok {
			return sv
		}
	}
	return nil
}

// getMapFromInterface
func getMapFromInterface(set interface{}) map[string]interface{} {
	if v, ok := set.(map[string]interface{}); ok {
		return v
	}
	return nil
}

// getSliceIfaceFromSet returns the (*schema.Set).List() if it is one
func getSliceIfaceFromSet(set interface{}) []interface{} {
	if v, ok := set.(*schema.Set); ok {
		return v.List()
	}
	return nil
}

// expandWeightedList returns an ordered list by it's weight field
func expandWeightedList(set []interface{}) ([]interface{}, error) {
	orderedList := make([]interface{}, len(set))

	for _, obj := range set {
		objCast := obj.(map[string]interface{})
		i := objCast["weight"].(int)

		if orderedList[i] != nil {
			return nil, fmt.Errorf("Weight %d used multiple times, did you set a weight?", i)
		}

		orderedList[i] = objCast
	}

	return orderedList, nil
}

/*
	Origin Pull Policy requires expansion/compression due to being actual
	list sets with weighted application
*/

func expandOriginPullPolicies(set []interface{}) ([]*models.OriginPullPolicy, error) {
	orderedList := make([]interface{}, len(set))
	ret := make([]*models.OriginPullPolicy, 0)
	orderedList, err := expandWeightedList(set)
	if err != nil {
		return nil, err
	}

	for _, policy := range orderedList {
		ret = append(ret, models.StructFromMap(&models.OriginPullPolicy{}, policy.(map[string]interface{})).(*models.OriginPullPolicy))
	}

	return ret, nil

}

func compressOriginPullPolicies(model []*models.OriginPullPolicy) []interface{} {
	originPullPolicyIface := make([]interface{}, 0)
	for index, policy := range model {
		thisPolicy := models.MapFromStruct(policy)
		thisPolicy["weight"] = index
		originPullPolicyIface = append(originPullPolicyIface, thisPolicy)
	}
	return originPullPolicyIface
}

/*
	Request & Response Modifications requires expansion/compression due to being actual
	list sets with weighted application
*/

// OriginRequestModification

func expandOriginRequestModification(set []interface{}) ([]*models.OriginRequestModification, error) {
	ret := make([]*models.OriginRequestModification, 0)

	orderedList, err := expandWeightedList(set)
	if err != nil {
		return nil, err
	}

	for _, modification := range orderedList {
		ret = append(ret, models.StructFromMap(&models.OriginRequestModification{}, modification.(map[string]interface{})).(*models.OriginRequestModification))
	}
	return ret, nil
}

func compressOriginRequestModification(model []*models.OriginRequestModification) []interface{} {
	originRequestModificationIface := make([]interface{}, 0)
	for index, policy := range model {
		thisPolicy := models.MapFromStruct(policy)
		thisPolicy["weight"] = index
		originRequestModificationIface = append(originRequestModificationIface, thisPolicy)
	}
	return originRequestModificationIface
}

// OriginResponseModification

func expandOriginResponseModification(set []interface{}) ([]*models.OriginResponseModification, error) {
	ret := make([]*models.OriginResponseModification, 0)

	orderedList, err := expandWeightedList(set)
	if err != nil {
		return nil, err
	}

	for _, modification := range orderedList {
		ret = append(ret, models.StructFromMap(&models.OriginResponseModification{}, modification.(map[string]interface{})).(*models.OriginResponseModification))
	}
	return ret, nil
}

func compressOriginResponseModification(model []*models.OriginResponseModification) []interface{} {
	originResponseModificationIface := make([]interface{}, 0)
	for index, policy := range model {
		thisPolicy := models.MapFromStruct(policy)
		thisPolicy["weight"] = index
		originResponseModificationIface = append(originResponseModificationIface, thisPolicy)
	}
	return originResponseModificationIface
}

// ClientResponseModification

func expandClientResponseModification(set []interface{}) ([]*models.ClientResponseModification, error) {
	ret := make([]*models.ClientResponseModification, 0)

	orderedList, err := expandWeightedList(set)
	if err != nil {
		return nil, err
	}

	for _, modification := range orderedList {
		ret = append(ret, models.StructFromMap(&models.ClientResponseModification{}, modification.(map[string]interface{})).(*models.ClientResponseModification))
	}
	return ret, nil
}

func compressClientResponseModification(model []*models.ClientResponseModification) []interface{} {
	clientResponseModificationIface := make([]interface{}, 0)
	for index, policy := range model {
		thisPolicy := models.MapFromStruct(policy)
		thisPolicy["weight"] = index
		clientResponseModificationIface = append(clientResponseModificationIface, thisPolicy)
	}
	return clientResponseModificationIface
}

// ClientRequestModification

func expandClientRequestModification(set []interface{}) ([]*models.ClientRequestModification, error) {
	ret := make([]*models.ClientRequestModification, 0)

	orderedList, err := expandWeightedList(set)
	if err != nil {
		return nil, err
	}

	for _, modification := range orderedList {
		ret = append(ret, models.StructFromMap(&models.ClientRequestModification{}, modification.(map[string]interface{})).(*models.ClientRequestModification))
	}
	return ret, nil
}

func compressClientRequestModification(model []*models.ClientRequestModification) []interface{} {
	clientRequestModificationIface := make([]interface{}, 0)
	for index, policy := range model {
		thisPolicy := models.MapFromStruct(policy)
		thisPolicy["weight"] = index
		clientRequestModificationIface = append(clientRequestModificationIface, thisPolicy)
	}
	return clientRequestModificationIface
}

/*
	Delivery
	This is a huge set of sets comprised of many submodels
*/

// compressDeliverySet packs several models into an []interface{} to be injected in a tf set
func compressDeliverySet(c *models.Configuration) []interface{} {
	// compression - can't avoid spelling out the fields here for now
	// no way to store object refs in json tags
	deliveryFirstLevel := make(map[string]interface{})
	deliveryFirstLevel["compression"] = []interface{}{models.MapFromStruct(c.Compression)}
	deliveryFirstLevel["static_header"] = compressDeliveryStaticHeader(c.StaticHeader)
	deliveryFirstLevel["http_methods"] = []interface{}{models.MapFromStruct(c.HTTPMethods)}
	deliveryFirstLevel["response_header"] = []interface{}{models.MapFromStruct(c.ResponseHeader)}
	deliveryFirstLevel["bandwidth_rate_limiting"] = []interface{}{models.MapFromStruct(c.BandwidthRateLimit)}
	deliveryFirstLevel["pattern_based_rate_limiting"] = []interface{}{models.MapFromStruct(c.BandwidthLimit)}

	if len(deliveryFirstLevel) > 0 {
		deliveryIface := make([]interface{}, 0)
		deliveryIface = append(deliveryIface, deliveryFirstLevel)
		return deliveryIface
	}

	return nil
}

// compressOriginSet packs several models into an []interface{} to be injected in a tf set
func compressOriginSet(c *models.Configuration) []interface{} {
	originFirstLevel := make(map[string]interface{})
	originFirstLevel["origin_pull_host"] = []interface{}{models.MapFromStruct(c.OriginPullHost)}

	if len(originFirstLevel) > 0 {
		originIface := make([]interface{}, 0)
		originIface = append(originIface, originFirstLevel)
		return originIface
	}
	return nil
}

// StaticHeader - weighted set

// expandDeliveryStaticHeader expands the weighted static header set
func expandDeliveryStaticHeader(set []interface{}) ([]*models.StaticHeader, error) {
	ret := make([]*models.StaticHeader, 0)

	orderedList, err := expandWeightedList(set)
	if err != nil {
		return nil, err
	}

	for _, modification := range orderedList {
		ret = append(ret, models.StructFromMap(&models.StaticHeader{}, modification.(map[string]interface{})).(*models.StaticHeader))
	}
	return ret, nil
}

func compressDeliveryStaticHeader(model []*models.StaticHeader) []interface{} {
	staticHeaderIface := make([]interface{}, 0)
	for index, policy := range model {
		thisPolicy := models.MapFromStruct(policy)
		thisPolicy["weight"] = index
		staticHeaderIface = append(staticHeaderIface, thisPolicy)
	}
	return staticHeaderIface
}
