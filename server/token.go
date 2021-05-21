package main

import (
	"net/url"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	kvEncryptionPasswordKey = "encryptionPassword"
)

var (
	//if the plugin fails to save password to KV, this fallback password will be used
	fallbackPassword = ""
)

//EnsureEncryptionPassword generates a password for encrypting the tokens, if it does not exist
//This method is called from plugin.go, and will generate a password only the first time when the plugin is loaded
func (p *Plugin) EnsureEncryptionPassword() {
	password := GenerateEncryptionPassword()
	ok, err := p.KVEnsure(kvEncryptionPasswordKey, []byte(password))
	if err != nil {
		p.API.LogError("Failed to set an encryption password for the plugin, fallback password will be used.", "Error", err.Error())
		fallbackPassword = password
		return
	}
	if !ok {
		p.API.LogWarn("Skipped write since already set by another plugin instance.")
	}
}

func (p *Plugin) getEncryptionPassword() []byte {
	//if the fallbackPassword is set this means the plugin cannot read from KV pair
	if fallbackPassword != "" {
		return []byte(fallbackPassword)
	}

	tokenSignPasswordByte, _ := p.API.KVGet(kvEncryptionPasswordKey)
	return tokenSignPasswordByte
}

//EncodeToken creates a token for WOPI
func (p *Plugin) EncodeToken(userID string, fileID string) string {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &WopiToken{
		UserID: userID,
		FileID: fileID,
	})
	signedString, err := token.SignedString(p.getEncryptionPassword())
	if err != nil {
		p.API.LogError("Failed to encode WOPI token", "Error", err.Error())
		return ""
	}
	return signedString
}

//DecodeToken decodes a token string an returns WopiToken and isValid
func (p *Plugin) DecodeToken(tokenString string) (WopiToken, bool) {
	wopiToken := WopiToken{}
	_, err := jwt.ParseWithClaims(tokenString, &wopiToken, func(token *jwt.Token) (interface{}, error) {
		return p.getEncryptionPassword(), nil
	})

	if err != nil {
		p.API.LogError("Failed to decode WOPI token", "Error", err.Error())
		return WopiToken{}, false
	}

	return wopiToken, true
}

// getAccessTokenFromURI extracts the access_token from the URI
// We need to do this manually as Mattermost removes the access_token before it reaches the plugin HTTP request parser
func getAccessTokenFromURI(uri string) (string, error) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse uri")
	}
	urlValues, parseErr := url.ParseQuery(parsedURL.RawQuery)
	if parseErr != nil {
		return "", errors.Wrap(parseErr, "failed to parse raw query")
	}
	return urlValues.Get("access_token"), nil
}

// GetWopiTokenFromURI decodes a token string from the URI
// returns WopiToken and error
func (p *Plugin) GetWopiTokenFromURI(uri string) (WopiToken, error) {
	token, tokenErr := getAccessTokenFromURI(uri)
	if tokenErr != nil {
		return WopiToken{}, errors.Wrap(tokenErr, "failed to retrieve token from URI: "+uri)
	}

	wopiToken, isValid := p.DecodeToken(token)
	if !isValid {
		return WopiToken{}, errors.New("collaboraOnline called the plugin with an invalid token")
	}

	return wopiToken, nil
}
