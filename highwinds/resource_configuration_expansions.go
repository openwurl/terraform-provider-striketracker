package highwinds

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/pkg/debug"
	"github.com/openwurl/wurlwind/striketracker/models"
)

func expandScopeModel(m map[string]interface{}) *models.Scope {
	s := &models.Scope{}

	if v, ok := m["platform"]; ok {
		s.Platform = v.(string)
	}

	if v, ok := m["path"]; ok {
		s.Path = v.(string)
	}

	if v, ok := m["name"]; ok {
		s.Name = v.(string)
	}

	return s
}

func expandOriginPullHost(m map[string]interface{}) *models.OriginPullHost {
	o := &models.OriginPullHost{}

	if v, ok := m["primary"]; ok {
		o.Primary = v.(int)
	}

	if v, ok := m["secondary"]; ok {
		o.Secondary = v.(int)
	}

	if v, ok := m["path"]; ok {
		o.Path = v.(string)
	}

	return o
}

func expandOriginPullCacheExtension(m map[string]interface{}) *models.OriginPullCacheExtension {
	o := &models.OriginPullCacheExtension{}

	if v, ok := m["enabled"]; ok {
		o.Enabled = v.(bool)
	}

	if v, ok := m["expired_cache_extension"]; ok {
		v := v.(int)
		o.ExpiredCacheExtension = &v
	}

	if v, ok := m["origin_unreachable_cache_extension"]; ok {
		v := v.(int)
		o.OriginUnreachableCacheExtension = &v
	}

	debug.Log("STALE CACHE EXTENSION", fmt.Sprintf("%v", o))

	return o

}

/*
	Delivery
*/

func expandDeliverySet(raw interface{}) map[string]interface{} {
	if deliverySet, ok := raw.(*schema.Set); ok {
		set := deliverySet.List()[0]
		if deliverySlice, ok := set.(map[string]interface{}); ok {
			return deliverySlice
		}
	}
	return nil
}

//delivery = []interface{}[0] = map[string]interface{} ["static_header"] = []interface{}[x] = map[string]interface{}

func compressDeliverySet(c *models.Configuration) []interface{} {
	delivery := make([]interface{}, 0)

	deliverySet := make(map[string]interface{})

	if c.Compression != nil {
		deliverySet["compression"] = compressDeliveryCompression(c.Compression)
	}
	if c.StaticHeader != nil {
		deliverySet["static_header"] = compressDeliveryStaticHeaders(c.StaticHeader)
	}

	if len(deliverySet) > 0 {
		delivery = append(delivery, deliverySet)
		return delivery
	}

	return nil
}

/*
	Delivery Compression
*/

func expandDeliveryCompression(raw interface{}) *models.Compression {
	if compression, ok := raw.(*schema.Set); ok {
		compressionSet := compression.List()[0]
		if m, ok := compressionSet.(map[string]interface{}); ok {
			c := &models.Compression{}

			/*
				if v, ok := m["enabled"]; ok {
					c.Enabled = v.(bool)
				}

				if v, ok := m["gzip"]; ok {
					c.GZIP = v.(string)
				}

				if v, ok := m["level"]; ok {
					c.Level = v.(int)
				}

				if v, ok := m["mime"]; ok {
					c.Mime = v.(string)
				}
			*/
			c = models.StructFromMap(c, m).(*models.Compression)

			return c
		}
	}
	return nil
}

func compressDeliveryCompression(c *models.Compression) []interface{} {
	if c != nil {
		compression := make([]interface{}, 0)
		compression = append(compression, models.MapFromStruct(c))
		return compression
	}
	// we don't want to make empty things where there is nothing
	return nil
}

/*
	Delivery Static Headers
*/

func expandDeliveryStaticHeaders(raw interface{}) []*models.StaticHeader {
	// TODO: Must implement a weighting like OriginPullPolicy, order matters
	/* TODO FIELDS
	weight
	*/
	m := make([]*models.StaticHeader, 0)

	if sh, ok := raw.(*schema.Set); ok {
		shList := sh.List()
		for _, shSet := range shList {
			if thisStaticHeader, ok := shSet.(map[string]interface{}); ok {
				m = append(m, expandDeliveryStaticHeader(thisStaticHeader))
			}

		}
		if len(m) > 0 {
			return m
		}
	}

	return nil
}

func expandDeliveryStaticHeader(m map[string]interface{}) *models.StaticHeader {
	sh := &models.StaticHeader{}
	/*
		if v, ok := m["enabled"]; ok {
			sh.Enabled = v.(bool)
		}
		if v, ok := m["origin_pull"]; ok {
			sh.OriginPull = v.(string)
		}
		if v, ok := m["client_request"]; ok {
			sh.ClientRequest = v.(string)
		}
		if v, ok := m["http"]; ok {
			sh.HTTP = v.(string)
		}
	*/
	sh = models.StructFromMap(sh, m).(*models.StaticHeader)

	return sh
}

func compressDeliveryStaticHeaders(sh []*models.StaticHeader) []interface{} {
	// TODO: Must implement a weighting like OriginPullPolicy, order matters
	/* TODO FIELDS
	weight
	*/
	if len(sh) > 0 {
		staticHeader := make([]interface{}, 0)
		for _, header := range sh {
			if header != nil {
				//staticHeader = append(staticHeader, header.Map())
				staticHeader = append(staticHeader, models.MapFromStruct(header))
			}
		}
		return staticHeader

	}

	return nil
}
