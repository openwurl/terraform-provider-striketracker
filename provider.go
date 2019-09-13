package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/striketracker"
)

// Provider provides the striketracker functionality
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"authorization_header_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTHORIZATIONHEADERKEY", nil),
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPLICATIONID", "wurlwind-terraform"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"striketracker_origin":        resourceOrigin(),
			"striketracker_certificate":   resourceCertificate(),
			"striketracker_host":          resourceHost(),
			"striketracker_configuration": resourceConfiguration(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	c, err := striketracker.NewClientWithOptions(
		striketracker.WithAuthorizationHeaderToken(d.Get("authorization_header_key").(string)),
		striketracker.WithApplicationID(d.Get("application_id").(string)),
	)
	return c, err
}
