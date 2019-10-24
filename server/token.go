package main

import (
	"log"
	"math/rand"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var kvEncryptionPassword = "encryptionPassword"
var fallbackPassword = "p345vnqcwqvc54v4cax" //if somehow the plugin can't retrieve the password from KV pair it will use this password

//WOPIToken is the token used for WOPI authnetication.
//When a user wants to open a file with Collabora Online this token is passed to Collabora Online
//Collabora Online will use this token when it loads/saves a file
type WOPIToken struct {
	UserID string `json:"userId"`
	FileID string `json:"fileId"`
	jwt.StandardClaims
}

//GenerateEncryptionPassword generates a password for encrypting the tokens
//This method is called from main, and will generate a password only the first time when the plugin is loaded
func GenerateEncryptionPassword(p *Plugin) {
	currentPassword, readPasswordError := p.API.KVGet(kvEncryptionPassword)
	if readPasswordError != nil {
		p.API.LogError("Cannot retrieve encryption password")
	}
	if len(currentPassword) == 0 {
		rand.Seed(time.Now().UnixNano())
		chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
			"abcdefghijklmnopqrstuvwxyz" +
			"0123456789")
		length := 20
		var b strings.Builder
		for i := 0; i < length; i++ {
			b.WriteRune(chars[rand.Intn(len(chars))])
		}
		password := b.String()
		saved, writePasswordError := p.API.KVCompareAndSet(kvEncryptionPassword, nil, []byte(password))
		if writePasswordError != nil {
			p.API.LogError("Cannot set an encryption password for the plugin")
		}
		if !saved {
			p.API.LogWarn("Skipped write since already set by another plugin instance")
		}
	}
}

//EncodeToken creates a token for WOPI
func EncodeToken(userID string, fileID string, p *Plugin) string {
	tokenSignPasswordByte, getPasswordError := p.API.KVGet(kvEncryptionPassword)
	if getPasswordError != nil {
		p.API.LogError("Cannot retrieve token signing password, eror: ", getPasswordError)
		tokenSignPasswordByte = []byte(fallbackPassword)
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &WOPIToken{
		UserID: userID,
		FileID: fileID,
	})
	tokenstring, err := token.SignedString(tokenSignPasswordByte)
	if err != nil {
		log.Fatalln(err)
	}
	return tokenstring
}

//DecodeToken decodes a token string an returns WOPIToken and isValid
func DecodeToken(tokenString string, p *Plugin) (WOPIToken, bool) {
	tokenSignPasswordByte, getPasswordError := p.API.KVGet(kvEncryptionPassword)
	if getPasswordError != nil {
		p.API.LogError("Cannot retrieve token signing password, eror: ", getPasswordError)
		tokenSignPasswordByte = []byte(fallbackPassword)
	}

	wopiToken := WOPIToken{}
	_, err := jwt.ParseWithClaims(tokenString, &wopiToken, func(token *jwt.Token) (interface{}, error) {
		return tokenSignPasswordByte, nil
	})

	if err != nil {
		return WOPIToken{}, false
	}

	return wopiToken, true
}
