package main

import (
	"encoding/json"
	"encoding/xml"

	"github.com/dgrijalva/jwt-go"
)

type WebappConfig struct {
	FileEditPermissions bool `json:"file_edit_permissions"`
}

func (c *WebappConfig) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	b, _ := json.Marshal(c)
	_ = json.Unmarshal(b, &out)
	return out
}

// WopiToken is the token used for WOPI authentication.
// When a user wants to open a file with Collabora Online this token is passed to Collabora Online
// Collabora Online will use this token when it loads/saves a file
type WopiToken struct {
	UserID string `json:"userId"`
	FileID string `json:"fileId"`
	jwt.StandardClaims
}

// WopiDiscovery represents the XML from <WOPI>/hosting/discovery
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

// WopiCheckFileInfo is the required response from http:// wopi.readthedocs.io/projects/wopirest/en/latest/files/CheckFileInfo.html#checkfileinfo
type WopiCheckFileInfo struct {
	// The string name of the file, including extension, without a path. Used for display in user interface (UI), and determining the extension of the file.
	BaseFileName string `json:"BaseFileName"`

	// The size of the file in bytes, expressed as a long, a 64-bit signed integer.
	Size int64 `json:"Size"`

	// A string that uniquely identifies the owner of the file.
	OwnerID string `json:"OwnerId"`

	// A string value uniquely identifying the user currently accessing the file.
	UserID string `json:"UserId"`

	// The name visible to other users while editing collaboratively.
	UserFriendlyName string `json:"UserFriendlyName"`

	// User permissions
	UserCanWrite bool `json:"UserCanWrite"`

	// Enables/disables the "Save As" acton in the File menu
	UserCanNotWriteRelative bool `json:"UserCanNotWriteRelative"`
}

// WopiFile is used top map file extension with the action & url
type WopiFile struct {
	URL    string // WOPI url to view/edit the file
	Action string // edit or view
}

// ClientFileInfo contains file information sent to the client
type ClientFileInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Action    string `json:"action"` // view or edit
}
