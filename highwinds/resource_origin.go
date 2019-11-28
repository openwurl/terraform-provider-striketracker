package highwinds

import (
	"fmt"
	"strconv"

	"github.com/openwurl/wurlwind/striketracker"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/striketracker/models"
	"github.com/openwurl/wurlwind/striketracker/services/origin"
)

func resourceOrigin() *schema.Resource {
	return &schema.Resource{
		Create: resourceOriginCreate,
		Read:   resourceOriginRead,
		Update: resourceOriginUpdate,
		Delete: resourceOriginDelete,
		Exists: resourceOriginExists,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				accountHash, resourceID, err := ResourceImportParseHashID(d.Id())
				if err != nil {
					return nil, err
				}
				d.Set("account_hash", accountHash)
				d.SetId(resourceID)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
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
				Default:     "NONE",
			},
			"certificate_cn": &schema.Schema{
				Description: "Common name to validate if VerifyCertificate",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"error_cache_ttl_seconds": &schema.Schema{
				Description: "DNS Timeout for origin request",
				Type:        schema.TypeInt,
				Default:     100,
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
				Default:     1,
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
				Default:     "/",
			},
			"request_timeout_seconds": &schema.Schema{
				Description: "Timeout for request to origin",
				Type:        schema.TypeInt,
				Default:     15,
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

/*
	Create
*/
func resourceOriginCreate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

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

	ctx, cancel := getContext()
	defer cancel()

	returnedModel, err := s.Create(ctx, accountHash, origin)
	if returnedModel != nil {
		if returnedModel.ID != 0 {
			d.SetId(fmt.Sprintf("%d", returnedModel.ID))
		}
	}
	if err != nil {
		return err
	}

	d.Partial(false)

	return resourceOriginRead(d, m)
}

/*
	Read
*/
func resourceOriginRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*striketracker.Client)
	s := origin.New(c)
	accountHash := d.Get("account_hash").(string)

	ctx, cancel := getContext()
	defer cancel()

	originID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	originResource, err := s.Get(ctx, accountHash, originID)
	if err != nil {
		return err
	}

	d.Set("name", originResource.Name)
	d.Set("hostname", originResource.Hostname)
	d.Set("port", originResource.Port)
	d.Set("path", originResource.Path)
	d.Set("authentication_type", originResource.AuthenticationType)
	d.Set("certificate_cn", originResource.CertificateCN)
	d.Set("error_cache_ttl_seconds", originResource.ErrorCacheTTLSeconds)
	d.Set("max_connections_per_edge", originResource.MaxConnectionsPerEdge)
	d.Set("max_connections_per_edge_enabled", originResource.MaxConnectionsPerEdgeEnabled)
	d.Set("maximum_origin_pull_seconds", originResource.MaximumOriginPullSeconds)
	d.Set("max_retry_count", originResource.MaxRetryCount)
	d.Set("origin_cache_headers", originResource.OriginCacheHeaders)
	d.Set("origin_default_keep_alive", originResource.OriginDefaultKeepAlive)
	d.Set("origin_pull_headers", originResource.OriginPullHeaders)
	d.Set("origin_pull_neg_linger", originResource.OriginPullNegLinger)
	d.Set("request_timeout_seconds", originResource.RequestTimeoutSeconds)
	d.Set("secure_port", originResource.SecurePort)
	d.Set("verify_certificate", originResource.VerifyCertificate)

	return nil
}

/*
	Update
*/
func resourceOriginUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
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

	ctx, cancel := getContext()
	defer cancel()

	returnedModel, err := s.Update(ctx, accountHash, origin)
	if returnedModel != nil {
		if returnedModel.ID != 0 {
			d.SetId(fmt.Sprintf("%d", returnedModel.ID))
		}
	}
	if err != nil {
		return err
	}

	d.Partial(false)
	return resourceOriginRead(d, m)
}

/*
	Delete
*/
func resourceOriginDelete(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	c := m.(*striketracker.Client)

	s := origin.New(c)
	accountHash := d.Get("account_hash").(string)
	originID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Origin ID %s is an invalid origin ID: %v", d.Id(), err)
	}

	ctx, cancel := getContext()
	defer cancel()

	err = s.Delete(ctx, accountHash, originID)
	if err != nil {
		return err
	}
	d.Partial(false)
	d.SetId("")
	return nil
}

/*
	Exists
*/
func resourceOriginExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*striketracker.Client)
	accountHash := d.Get("account_hash").(string)

	cs := origin.New(c)

	ctx, cancel := getContext()
	defer cancel()

	originID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	originResource, err := cs.Get(ctx, accountHash, originID)
	if err != nil {
		return false, err
	}

	if originResource == nil {
		return false, fmt.Errorf(striketracker.ErrNotFound)
	}

	return true, nil
}
