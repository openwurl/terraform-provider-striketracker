package main

import (
	"fmt"
	"strconv"

	"github.com/openwurl/wurlwind/striketracker"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/striketracker/models"
	"github.com/openwurl/wurlwind/striketracker/services/origin"
)

// TODO
/*
	Fix partial states
	return reads on resource

*/

func resourceOrigin() *schema.Resource {
	return &schema.Resource{
		Create: resourceOriginCreate,
		Read:   resourceOriginRead,
		Update: resourceOriginUpdate,
		Delete: resourceOriginDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"origin_id": &schema.Schema{
				Description: "The computed ID of the origin",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"name": &schema.Schema{
				Description: "The name of the origin",
				Type:        schema.TypeString,
				Required:    true,
			},
			"hostname": &schema.Schema{
				Description: "The hostname or IP of the origin",
				Type:        schema.TypeString,
				Required:    true,
			},
			"port": &schema.Schema{
				Description: "The port of the origin (80, 443, 8080, 1936)",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"account_hash": &schema.Schema{
				Description: "The destination account hash where the origin will be created",
				Type:        schema.TypeString,
				Required:    true,
			},
			"authentication_type": &schema.Schema{
				Description: "Authentication type, NONE or BASIC",
				Type:        schema.TypeString,
				Optional:    true,
				//Default:     "NONE",
			},
			"certificate_cn": &schema.Schema{
				Description: "Common name to validate if VerifyCertificate",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"error_cache_ttl_seconds": &schema.Schema{
				Description: "DNS Timeout for origin request",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"max_connections_per_edge": &schema.Schema{
				Description: "Maximum concurrent connections any single edge will make",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"max_connections_per_edge_enabled": &schema.Schema{
				Description: "Limit maximum connections from edges",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"maximum_origin_pull_seconds": &schema.Schema{
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"max_retry_count": &schema.Schema{
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"origin_cache_headers": &schema.Schema{
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"origin_default_keep_alive": &schema.Schema{
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"origin_pull_headers": &schema.Schema{
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"origin_pull_neg_linger": &schema.Schema{
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"path": &schema.Schema{
				Description: "The path at the origin to request",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"request_timeout_seconds": &schema.Schema{
				Description: "Timeout for request to origin",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"secure_port": &schema.Schema{
				Description: "Port for use with TLS connections to origin",
				Type:        schema.TypeInt,
				Optional:    true,
				//Default:     443,
			},
			"verify_certificate": &schema.Schema{
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
	}
}

func resourceOriginCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*striketracker.Client)
	s := origin.New(c)
	accountHash := d.Get("account_hash").(string)
	origin := &models.Origin{
		Name:                         d.Get("name").(string),
		Hostname:                     d.Get("hostname").(string),
		Port:                         d.Get("port").(int),
		Path:                         d.Get("path").(string),
		AuthenticationType:           d.Get("authentication_type").(string),
		CertificateCN:                d.Get("certificate_cn").(string),
		ErrorCacheTTLSeconds:         d.Get("error_cache_ttl_seconds").(int),
		MaxConnectionsPerEdge:        d.Get("max_connections_per_edge").(int),
		MaxConnectionsPerEdgeEnabled: d.Get("max_connections_per_edge_enabled").(bool),
		MaximumOriginPullSeconds:     d.Get("maximum_origin_pull_seconds").(int),
		MaxRetryCount:                d.Get("max_retry_count").(int),
		OriginCacheHeaders:           d.Get("origin_cache_headers").(string),
		OriginDefaultKeepAlive:       d.Get("origin_default_keep_alive").(int),
		OriginPullHeaders:            d.Get("origin_pull_headers").(string),
		OriginPullNegLinger:          d.Get("origin_pull_neg_linger").(string),
		RequestTimeoutSeconds:        d.Get("request_timeout_seconds").(int),
		SecurePort:                   d.Get("secure_port").(int),
		VerifyCertificate:            d.Get("verify_certificate").(bool),
	}

	returnedModel, err := s.Create(accountHash, origin)
	if err != nil {
		d.Partial(true)
		return err
	}

	d.SetId(fmt.Sprintf("%d", returnedModel.ID))
	//d.Set("origin_id", returnedModel.ID)
	return resourceOriginRead(d, m)
}

func resourceOriginRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceOriginUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*striketracker.Client)
	originID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Origin ID %s is an invalid origin ID: %v", d.Id(), err)
	}

	s := origin.New(c)
	accountHash := d.Get("account_hash").(string)
	origin := &models.Origin{
		ID:                           originID,
		Name:                         d.Get("name").(string),
		Hostname:                     d.Get("hostname").(string),
		Port:                         d.Get("port").(int),
		Path:                         d.Get("path").(string),
		AuthenticationType:           d.Get("authentication_type").(string),
		CertificateCN:                d.Get("certificate_cn").(string),
		ErrorCacheTTLSeconds:         d.Get("error_cache_ttl_seconds").(int),
		MaxConnectionsPerEdge:        d.Get("max_connections_per_edge").(int),
		MaxConnectionsPerEdgeEnabled: d.Get("max_connections_per_edge_enabled").(bool),
		MaximumOriginPullSeconds:     d.Get("maximum_origin_pull_seconds").(int),
		MaxRetryCount:                d.Get("max_retry_count").(int),
		OriginCacheHeaders:           d.Get("origin_cache_headers").(string),
		OriginDefaultKeepAlive:       d.Get("origin_default_keep_alive").(int),
		OriginPullHeaders:            d.Get("origin_pull_headers").(string),
		OriginPullNegLinger:          d.Get("origin_pull_neg_linger").(string),
		RequestTimeoutSeconds:        d.Get("request_timeout_seconds").(int),
		SecurePort:                   d.Get("secure_port").(int),
		VerifyCertificate:            d.Get("verify_certificate").(bool),
	}

	returnedModel, err := s.Update(accountHash, origin)
	if returnedModel.ID != 0 {
		d.SetId(fmt.Sprintf("%d", returnedModel.ID))
	}
	if err != nil {
		d.Partial(true)
		return err
	}

	return resourceOriginRead(d, m)
}

func resourceOriginDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*striketracker.Client)

	s := origin.New(c)
	accountHash := d.Get("account_hash").(string)
	//originID := d.Get("origin_id").(int)
	originID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Origin ID %s is an invalid origin ID: %v", d.Id(), err)
	}
	err = s.Delete(accountHash, originID)
	if err != nil {
		d.Partial(true)
		return err
	}
	d.SetId("")
	return nil
}
