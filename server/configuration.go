package main

import (
	"encoding/xml"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var (
	//WOPIData contains the XML from <WOPI>/hosting/discovery
	WOPIData WopiDiscovery

	//WOPIFiles maps file extension with file action & url
	WOPIFiles map[string]WOPIFileInfo
)

// configuration captures the plugin's external configuration as exposed in the Mattermost server
// configuration, as well as values computed from the configuration. Any public fields will be
// deserialized from the Mattermost server configuration in OnConfigurationChange.
//
// As plugins are inherently concurrent (hooks being called asynchronously), and the plugin
// configuration can change at any time, access to the configuration must be synchronized. The
// strategy used in this plugin is to guard a pointer to the configuration, and clone the entire
// struct whenever it changes. You may replace this with whatever strategy you choose.
//
// If you add non-reference types to your configuration struct, be sure to rewrite Clone as a deep
// copy appropriate for your types.
type configuration struct {
	WOPIAddress                    string
	DisableCertificateVerification bool
}

// Clone deep copies the configuration
func (c *configuration) Clone() *configuration {
	return &configuration{WOPIAddress: c.WOPIAddress}
}

// ProcessConfiguration processes the config.
func (c *configuration) ProcessConfiguration() error {
	// trim trailing slash or spaces from the WOPI address, if needed
	c.WOPIAddress = strings.TrimSpace(c.WOPIAddress)
	c.WOPIAddress = strings.Trim(c.WOPIAddress, "/")

	return nil
}

// IsValid checks if all needed fields are set.
func (c *configuration) IsValid() error {
	if !strings.HasPrefix(c.WOPIAddress, "http") {
		return errors.New("please provide the WOPIAddress")
	}

	return nil
}

// getConfiguration retrieves the active configuration under lock, making it safe to use
// concurrently. The active configuration may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

// setConfiguration replaces the active configuration under lock.
//
// Do not call setConfiguration while holding the configurationLock, as sync.Mutex is not
// reentrant. In particular, avoid using the plugin API entirely, as this may in turn trigger a
// hook back into the plugin. If that hook attempts to acquire this lock, a deadlock may occur.
//
// This method panics if setConfiguration is called with the existing configuration. This almost
// certainly means that the configuration was modified without being cloned and may result in
// an unsafe access.
func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		// Ignore assignment if the configuration struct is empty. Go will optimize the
		// allocation for same to point at the same memory address, breaking the check
		// above.
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}

		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

// OnConfigurationChange is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	// Load the public configuration fields from the Mattermost server configuration.
	if loadConfigErr := p.API.LoadPluginConfiguration(configuration); loadConfigErr != nil {
		return errors.Wrap(loadConfigErr, "failed to load plugin configuration")
	}

	if err := configuration.ProcessConfiguration(); err != nil {
		p.API.LogError("Error in ProcessConfiguration.", "Error", err.Error())
		return err
	}

	if err := configuration.IsValid(); err != nil {
		p.API.LogError("Error in Validating Configuration.", "Error", err.Error())
		return err
	}

	if err := p.loadWopiFileInfo(configuration.WOPIAddress); err != nil {
		return errors.Wrap(err, "could not load wopi file info")
	}

	p.setConfiguration(configuration)

	return nil
}

// loadWopiFileInfo loads the WOPI file data
func (p *Plugin) loadWopiFileInfo(wopiAddress string) error {
	client := p.getHTTPClient()
	resp, err := client.Get(wopiAddress + "/hosting/discovery")
	if err != nil {
		p.API.LogError("WOPI request error. Please check the WOPI address.", err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.API.LogError("WOPI request error. Failed to read WOPI request body. Please check the WOPI address.", err.Error())
		return err
	}

	if err := xml.Unmarshal(body, &WOPIData); err != nil {
		p.API.LogError("WOPI request error. Failed to unmarshal WOPI XML. Please check the WOPI address.", err.Error())
		return err
	}

	WOPIFiles = make(map[string]WOPIFileInfo)
	for i := 0; i < len(WOPIData.NetZone.App); i++ {
		for j := 0; j < len(WOPIData.NetZone.App[i].Action); j++ {
			ext := strings.ToLower(WOPIData.NetZone.App[i].Action[j].Ext)
			if ext == "" || ext == "png" || ext == "jpg" || ext == "jpeg" || ext == "gif" {
				continue
			}
			WOPIFiles[strings.ToLower(ext)] = WOPIFileInfo{WOPIData.NetZone.App[i].Action[j].URLSrc, WOPIData.NetZone.App[i].Action[j].Name}
		}
	}

	p.API.LogInfo("WOPI file info loaded successfully!")
	return nil
}
