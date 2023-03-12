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
		p.client.Log.Error("Failed to encode WOPI token.", "Error", err.Error())
		return ""
	}
	return signedString
}

// DecodeToken decodes a token string and returns WopiToken and isValid
func (p *Plugin) DecodeToken(tokenString string) (WopiToken, bool) {
	config := p.getConfiguration()
	wopiToken := WopiToken{}
	_, err := jwt.ParseWithClaims(tokenString, &wopiToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.EncryptionKey), nil
	})

	if err != nil {
		p.client.Log.Error("Failed to decode WOPI token.", "Error", err.Error())
		return WopiToken{}, false
	}

	return wopiToken, true
}

// GetWopiTokenFromURI decodes a token string from the URI
// returns WopiToken and error
func (p *Plugin) GetWopiTokenFromURI(uri string) (WopiToken, error) {
	parsedURL, parseURLErr := url.Parse(uri)
	if parseURLErr != nil {
		return WopiToken{}, errors.Wrap(parseURLErr, "failed to parse URI: "+uri)
	}

	// this token is in query parameters as Mattermost otherwise removes the access_token before it reaches the plugin HTTP request parser.
	queryValues := parsedURL.Query()
	token := queryValues.Get("access_token")
	if token == "" {
		return WopiToken{}, errors.New("failed to retrieve token from URI: " + uri)
	}

	wopiToken, isValid := p.DecodeToken(token)
	if !isValid {
		return WopiToken{}, errors.New("collaboraOnline called the plugin with an invalid token")
	}

	return wopiToken, nil
}
