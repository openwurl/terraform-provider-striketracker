package highwinds

import "github.com/openwurl/wurlwind/striketracker/models"

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

	return o

}
