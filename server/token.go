package main

import (
	"net/url"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// EncodeToken creates a token for WOPI
func (p *Plugin) EncodeToken(userID string, fileID string) string {
	config := p.getConfiguration()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &WopiToken{
		UserID: userID,
		FileID: fileID,
	})
	signedString, err := token.SignedString([]byte(config.EncryptionKey))
	if err != nil {
		p.API.LogError("Failed to encode WOPI token.", "Error", err.Error())
		return ""
	}
	return signedString
}

// DecodeToken decodes a token string an returns WopiToken and isValid
func (p *Plugin) DecodeToken(tokenString string) (WopiToken, bool) {
	config := p.getConfiguration()
	wopiToken := WopiToken{}
	_, err := jwt.ParseWithClaims(tokenString, &wopiToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.EncryptionKey), nil
	})

	if err != nil {
		p.API.LogError("Failed to decode WOPI token.", "Error", err.Error())
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
