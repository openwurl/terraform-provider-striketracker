package highwinds

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/striketracker"
	"github.com/openwurl/wurlwind/striketracker/services/configuration"
)

/*
	Create
*/
func resourceConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	c := m.(*striketracker.Client)
	conf := configuration.New(c)
	accountHash := d.Get("account_hash").(string)
	hostHash := d.Get("host_hash").(string)

	ctx, cancel := getContext()
	defer cancel()

	// Build our model to send
	newConfigurationScope := buildCreateScopeConfiguration(d)

	// Send model
	returnedModel, err := conf.Create(ctx, accountHash, hostHash, newConfigurationScope)
	if returnedModel != nil {
		if returnedModel.ID != 0 {
			d.SetId(fmt.Sprintf("%d", returnedModel.ID))
		}
	}
	if err != nil {
		return err
	}

	d.Partial(false)

	return resourceConfigurationUpdate(d, m)
}

/*
	Update
*/
func resourceConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*striketracker.Client)
	conf := configuration.New(c)
	accountHash := d.Get("account_hash").(string)
	hostHash := d.Get("host_hash").(string)
	scopeID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	ctx, cancel := getContext()
	defer cancel()

	devLog("Preparing to update configuration %s/%s/%s", accountHash, hostHash, scopeID)

	// Build our model to send
	newConfigurationScope := buildScopeConfiguration(d)

	devLog("Updating configuration %s/%s/%d", accountHash, hostHash, scopeID)
	// Ship object
	returnedModel, err := conf.Update(ctx, accountHash, hostHash, scopeID, newConfigurationScope)
	if err != nil {
		return err
	}
	if returnedModel == nil {
		return fmt.Errorf("Something went wrong updating the scope %s, returned model is nil", d.Id())
	}

	return resourceConfigurationRead(d, m)
}

/*
	Read
*/
func resourceConfigurationRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*striketracker.Client)
	conf := configuration.New(c)
	accountHash := d.Get("account_hash").(string)
	hostHash := d.Get("host_hash").(string)
	scopeID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	ctx, cancel := getContext()
	defer cancel()

	devLog("Reading configuration %s/%s/%d", accountHash, hostHash, scopeID)

	// Fetch resource
	configModel, err := conf.Get(ctx, accountHash, hostHash, scopeID)
	if err != nil {
		return err
	}
	if configModel == nil {
		return fmt.Errorf("Resource %s does not exist", d.Id())
	}

	devLog("Setting configuration state %s/%s/%d", accountHash, hostHash, scopeID)

	if configModel.Platform == "" || configModel.ID == 0 || configModel.Path == "" {
		return ErrScopeIsNil(accountHash, hostHash, scopeID)
	}

	errs := ErrSetState(ingestRemoteState(d, configModel))
	if errs != nil {
		return errs
	}

	devLog("Done setting configuration state %s/%s/%d", accountHash, hostHash, scopeID)

	return nil
}

/*
	Delete
*/
func resourceConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	c := m.(*striketracker.Client)
	conf := configuration.New(c)
	accountHash := d.Get("account_hash").(string)
	hostHash := d.Get("host_hash").(string)
	scopeID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	ctx, cancel := getContext()
	defer cancel()

	err = conf.Delete(ctx, accountHash, hostHash, scopeID, false)
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
func resourceConfigurationExists(d *schema.ResourceData, m interface{}) (bool, error) {
	return false, nil
}
