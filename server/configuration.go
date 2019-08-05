package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
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
			Action struct {
				Text   string `xml:",chardata"`
				Ext    string `xml:"ext,attr"`
				Name   string `xml:"name,attr"`
				Urlsrc string `xml:"urlsrc,attr"`
			} `xml:"action"`
		} `xml:"app"`
	} `xml:"net-zone"`
}

//WOPIData contains the XML from <WOPI>/hosting/discovery
var WOPIData WopiDiscovery

//WOPIAddress stores current WOPI address. Used to check if we need to reload WOPI data OnConfigurationChange
var wopiAddress string

//WOPIFileInfo is used top map file extension with the action & url
type WOPIFileInfo struct {
	Url    string //WOPI url to view/edit the file
	Action string //edit or view
}

//WOPIFiles maps file extension with file action & url
var WOPIFiles map[string]WOPIFileInfo

// Clone deep copies the configuration
func (c *configuration) Clone() *configuration {
	return &configuration{WOPIAddress: c.WOPIAddress}
}

//OnConfigurationChange will request new WOPI data
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	// Load the public configuration fields from the Mattermost server configuration.
	if loadConfigErr := p.API.LoadPluginConfiguration(configuration); loadConfigErr != nil {
		return errors.Wrap(loadConfigErr, "failed to load plugin configuration")
	}

	//retrieve the new WOPI data from <WOPI>/hosting/discovery
	if wopiAddress != configuration.WOPIAddress {
		wopiAddress = configuration.WOPIAddress
		fmt.Println("WOPI address changed. Load new WOPI file info.")
		resp, err := http.Get(configuration.WOPIAddress + "hosting/discovery")
		if err != nil {
			fmt.Println("WOPI request error. Please check the WOPI address.", err.Error())
			return errors.New("wopo request error, please check WOPI address")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		xml.Unmarshal(body, &WOPIData)

		WOPIFiles = make(map[string]WOPIFileInfo)
		for i := 0; i < len(WOPIData.NetZone.App); i++ {
			WOPIFiles[strings.ToLower(WOPIData.NetZone.App[i].Action.Ext)] = WOPIFileInfo{WOPIData.NetZone.App[i].Action.Urlsrc, WOPIData.NetZone.App[i].Action.Name}
		}
		fmt.Println("WOPI file info loaded successfully!")
	}

	return nil
}
