package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/striketracker"
	"github.com/openwurl/wurlwind/striketracker/models"
	"github.com/openwurl/wurlwind/striketracker/services/certificates"
)

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Update: resourceCertificateUpdate,
		Delete: resourceCertificateDelete,
		Exists: resourceCertificateExists,
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
			"ca_bundle": &schema.Schema{
				Description: "The text of the certificate's CA bundle",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"certificate": &schema.Schema{
				Description: "The text of the x.509 certificate",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"certificate_information": &schema.Schema{
				Description: "The text of the certificate's CA bundle",
				Type:        schema.TypeMap,
				Computed:    true,
			},
			"ciphers": &schema.Schema{
				Description: "The ciper list which should be used during the SSL handshake",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"common_name": &schema.Schema{
				Description: "The primary hostname for which this certificate can serve traffic",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_date": &schema.Schema{
				Description: "The date at which this certificate was uploaded",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"expiration_date": &schema.Schema{
				Description: "The time at which this certificate is no longer valid",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"fingerprint": &schema.Schema{
				Description: "The cryptographic hash of the certificate used for uniqueness checking",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"account_hash": &schema.Schema{
				Description: "The destination account hash where the origin will be created",
				Type:        schema.TypeString,
				Required:    true,
			},
			"issuer": &schema.Schema{
				Description: "The organization which issued the certificate",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"key": &schema.Schema{
				Description: "The text of the x.509 private key",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"requester": &schema.Schema{
				Description: "The user which uploaded the certificate",
				Type:        schema.TypeMap,
				Computed:    true,
			},
			"trusted": &schema.Schema{
				Description: "Whether or not this certificate passes CA validation",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"updated_date": &schema.Schema{
				Description: "The date this certificate was last updated",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

/*
	Read
*/
func resourceCertificateRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*striketracker.Client)
	accountHash := d.Get("account_hash").(string)

	cs := certificates.New(c)

	ctx, cancel := getContext()
	defer cancel()

	certificateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	certResource, err := cs.Get(ctx, accountHash, certificateID)
	if err != nil {
		return err
	}

	d.Set("ca_bundle", certResource.CABundle)
	d.Set("certificate", certResource.Certificate)
	d.Set("certificate_information", certResource.CertificateInformation)
	d.Set("ciphers", certResource.Ciphers)
	d.Set("common_name", certResource.CommonName)
	d.Set("created_date", certResource.CreatedDate)
	d.Set("expiration_date", certResource.ExpirationDate)
	d.Set("fingerprint", certResource.Fingerprint)
	d.Set("issuer", certResource.Issuer)
	d.Set("key", certResource.Key)
	d.Set("requester", certResource.Requester)
	d.Set("trusted", certResource.Trusted)
	d.Set("updated_date", certResource.UpdatedDate)

	return nil
}

/*
	Create
*/
func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	c := m.(*striketracker.Client)
	accountHash := d.Get("account_hash").(string)

	cs := certificates.New(c)

	ctx, cancel := getContext()
	defer cancel()

	certificate := &models.Certificate{
		CABundle:    d.Get("ca_bundle").(string),
		Certificate: d.Get("certificate").(string),
		Ciphers:     d.Get("ciphers").(string),
		CommonName:  d.Get("common_name").(string),
		Trusted:     d.Get("trusted").(bool),
		Key:         d.Get("key").(string),
	}

	returnedCertificate, err := cs.Upload(ctx, accountHash, certificate)
	if returnedCertificate != nil {
		if returnedCertificate.ID != 0 {
			d.SetId(fmt.Sprintf("%d", returnedCertificate.ID))
		}
	}
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", returnedCertificate.ID))
	d.Partial(false)

	return resourceCertificateRead(d, m)
}

/*
	Update
*/
func resourceCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	c := m.(*striketracker.Client)
	accountHash := d.Get("account_hash").(string)

	cs := certificates.New(c)

	ctx, cancel := getContext()
	defer cancel()

	certificateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	certificate := &models.Certificate{
		ID:          certificateID,
		CABundle:    d.Get("ca_bundle").(string),
		Certificate: d.Get("certificate").(string),
		Ciphers:     d.Get("ciphers").(string),
		CommonName:  d.Get("common_name").(string),
		Trusted:     d.Get("trusted").(bool),
		Key:         d.Get("key").(string),
	}

	returnedCertificate, err := cs.Update(ctx, accountHash, certificate)
	if returnedCertificate != nil {
		if returnedCertificate.ID != 0 {
			d.SetId(fmt.Sprintf("%d", returnedCertificate.ID))
		}
	}
	if err != nil {
		return err
	}

	d.Partial(false)

	return resourceCertificateRead(d, m)
}

/*
	Delete
*/
func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	c := m.(*striketracker.Client)
	accountHash := d.Get("account_hash").(string)

	cs := certificates.New(c)

	ctx, cancel := getContext()
	defer cancel()

	certificateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	err = cs.Delete(ctx, accountHash, certificateID)
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
func resourceCertificateExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*striketracker.Client)
	accountHash := d.Get("account_hash").(string)

	cs := certificates.New(c)

	ctx, cancel := getContext()
	defer cancel()

	certificateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	certResource, err := cs.Get(ctx, accountHash, certificateID)
	if err != nil {
		return false, err
	}

	if certResource != nil {
		if certResource.Code == ErrCodeNotFound && strings.Contains(certResource.Error, ErrNotFound) {
			err = fmt.Errorf("Resource does not exist")
			return false, certResource.Err(err)
		}
	}

	return true, nil
}
