package highwinds

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/openwurl/wurlwind/pkg/debug"
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
	newConfigurationScope, err := buildNewConfigurationFromState(d)
	if err != nil {
		return err
	}

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
	Create
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

	debug.Log("Update", "Preparing to update configuration %s/%s/%d", accountHash, hostHash, scopeID)

	// Build our model to send
	newConfigurationScope, err := buildConfigurationFromState(d)
	if err != nil {
		return fmt.Errorf("Error building config from state: %v", err.Error())
	}

	debug.Log("Update", "Updating configuration %s/%s/%d", accountHash, hostHash, scopeID)
	// Ship object
	returnedModel, err := conf.Update(ctx, accountHash, hostHash, scopeID, newConfigurationScope)
	if err != nil {
		return err
	}
	if returnedModel == nil {
		return fmt.Errorf("Something went wrong updating the scope %s, returned model is nil", d.Id())
	}

	d.Partial(false)

	return resourceConfigurationRead(d, m)
}

/*
	Create
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

	debug.Log("Read", "Reading configuration %s/%s/%d", accountHash, hostHash, scopeID)

	// Fetch resource
	configModel, err := conf.Get(ctx, accountHash, hostHash, scopeID)
	if err != nil {
		return err
	}
	if configModel == nil {
		return fmt.Errorf("Resource %s does not exist", d.Id())
	}

	debug.Log("Read", "Setting configuration state %s/%s/%d", accountHash, hostHash, scopeID)

	if configModel.Platform == "" || configModel.ID == 0 || configModel.Path == "" {
		return ErrScopeIsNil(accountHash, hostHash, scopeID)
	}

	// Ingest the remote state and update our local state
	errs := ErrSetState(ingestState(d, configModel))
	if errs != nil {
		return errs
	}

	debug.Log("Read", "Done setting configuration state %s/%s/%d", accountHash, hostHash, scopeID)

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

// TODO: exists
/*
	Exists
*/
func resourceConfigurationExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*striketracker.Client)
	conf := configuration.New(c)
	accountHash := d.Get("account_hash").(string)
	hostHash := d.Get("host_hash").(string)
	scopeID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}
	ctx, cancel := getContext()
	defer cancel()

	debug.Log("Exists", "Reading configuration %s/%s/%d", accountHash, hostHash, scopeID)

	// Fetch resource
	configModel, err := conf.Get(ctx, accountHash, hostHash, scopeID)
	if err != nil {
		return false, err
	}
	if configModel == nil {
		return false, fmt.Errorf("Resource %s does not exist", d.Id())
	}
	return true, nil
}
