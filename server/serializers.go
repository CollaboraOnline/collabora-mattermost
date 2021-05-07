package main

import "encoding/xml"

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

//WOPIFileInfo is used top map file extension with the action & url
type WOPIFileInfo struct {
	URL    string //WOPI url to view/edit the file
	Action string //edit or view
}

//CollaboraFileInfo contains file information sent to the client
type CollaboraFileInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Action    string `json:"action"` //view or edit
}
