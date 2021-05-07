package main

import (
	"math/rand"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var kvEncryptionPassword = "encryptionPassword"
var fallbackPassword = "" //if the plugin fails to save password to KV, this fallback password will be used

//WOPIToken is the token used for WOPI authentication.
//When a user wants to open a file with Collabora Online this token is passed to Collabora Online
//Collabora Online will use this token when it loads/saves a file
type WOPIToken struct {
	UserID string `json:"userId"`
	FileID string `json:"fileId"`
	jwt.StandardClaims
}

//GenerateEncryptionPassword generates a password for encrypting the tokens
//This method is called from main, and will generate a password only the first time when the plugin is loaded
func (p *Plugin) GenerateEncryptionPassword() {
	currentPassword, readPasswordError := p.API.KVGet(kvEncryptionPassword)
	if readPasswordError != nil {
		p.API.LogError("Cannot retrieve encryption password")
	}
	if len(currentPassword) == 0 {
		rand.Seed(time.Now().UnixNano())
		chars := []rune(
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
				"abcdefghijklmnopqrstuvwxyz" +
				"0123456789",
		)
		length := 20
		var b strings.Builder
		for i := 0; i < length; i++ {
			b.WriteRune(chars[rand.Intn(len(chars))])
		}
		password := b.String()
		saved, writePasswordError := p.API.KVCompareAndSet(kvEncryptionPassword, nil, []byte(password))
		if writePasswordError != nil {
			p.API.LogError("Cannot set an encryption password for the plugin, fallback password will be used")
			fallbackPassword = password
		}
		if !saved {
			p.API.LogWarn("Skipped write since already set by another plugin instance")
		}
	}
}

func getEncryptionPassword(p *Plugin) []byte {
	//if the fallbackPassword is set this means the plugin cannot read from KV pair
	if fallbackPassword != "" {
		return []byte(fallbackPassword)
	}

	tokenSignPasswordByte, _ := p.API.KVGet(kvEncryptionPassword)
	return tokenSignPasswordByte
}

//EncodeToken creates a token for WOPI
func EncodeToken(userID string, fileID string, p *Plugin) string {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &WOPIToken{
		UserID: userID,
		FileID: fileID,
	})
	signedString, err := token.SignedString(getEncryptionPassword(p))
	if err != nil {
		p.API.LogError("Failed to encode WOPI token", "Error", err.Error())
		return ""
	}
	return signedString
}

//DecodeToken decodes a token string an returns WOPIToken and isValid
func DecodeToken(tokenString string, p *Plugin) (WOPIToken, bool) {
	wopiToken := WOPIToken{}
	_, err := jwt.ParseWithClaims(tokenString, &wopiToken, func(token *jwt.Token) (interface{}, error) {
		return getEncryptionPassword(p), nil
	})

	if err != nil {
		p.API.LogError("Failed to decode WOPI token", "Error", err.Error())
		return WOPIToken{}, false
	}

	return wopiToken, true
}
