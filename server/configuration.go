package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type configuration struct {
	WOPIAddress string
}

//WopiDiscovery represents the XML from <WOPI>/hosting/discovery
type WopiDiscovery struct {
	XMLName xml.Name `xml:"wopi-discovery"`
	Text    string   `xml:",chardata"`
	NetZone struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
		App  []struct {
			Text   string `xml:",chardata"`
			Name   string `xml:"name,attr"`
			Action []struct {
				Text   string `xml:",chardata"`
				Ext    string `xml:"ext,attr"`
				Name   string `xml:"name,attr"`
				URLSrc string `xml:"urlsrc,attr"`
			} `xml:"action"`
		} `xml:"app"`
	} `xml:"net-zone"`
}

//WOPIData contains the XML from <WOPI>/hosting/discovery
var WOPIData WopiDiscovery

//WOPIFileInfo is used top map file extension with the action & url
type WOPIFileInfo struct {
	URL    string //WOPI url to view/edit the file
	Action string //edit or view
}

//WOPIFiles maps file extension with file action & url
var WOPIFiles map[string]WOPIFileInfo

// Clone deep copies the configuration
func (c *configuration) Clone() *configuration {
	return &configuration{WOPIAddress: c.WOPIAddress}
}

// OnConfigurationChange is called when plugin's configuration changes
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	// Load the public configuration fields from the Mattermost server configuration.
	if loadConfigErr := p.API.LoadPluginConfiguration(configuration); loadConfigErr != nil {
		return errors.Wrap(loadConfigErr, "failed to load plugin configuration")
	}

	p.setConfiguration(configuration)

	return nil
}

// set the new configuration and load WOPI file data
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

	wopiAddress := configuration.WOPIAddress

	//append trailing slash to the WOPI address if needed
	if wopiAddress[len(wopiAddress)-1:] != "/" {
		wopiAddress = fmt.Sprintf("%s%s", wopiAddress, "/")
	}

	resp, err := http.Get(wopiAddress + "hosting/discovery")
	if err != nil {
		p.API.LogError("WOPI request error. Please check the WOPI address.", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.API.LogError("WOPI request error. Failed to read WOPI request body. Please check the WOPI address.", err.Error())
		return
	}

	if err := xml.Unmarshal(body, &WOPIData); err != nil {
		p.API.LogError("WOPI request error. Failed to unmarshal WOPI XML. Please check the WOPI address.", err.Error())
		return
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
}
