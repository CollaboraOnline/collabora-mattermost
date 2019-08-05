package main

import (
	"log"

	jwt "github.com/dgrijalva/jwt-go"
)

var tokenSignPassword = "aUniquePassword"

//WOPIToken is the token used for WOPI authnetication.
//When a user wants to open a file with Collabora Online this token is passed to Collabora Online
//Collabora Online will use this token when it loads/saves a file
type WOPIToken struct {
	UserID string `json:"userId"`
	FileID string `json:"fileId"`
	jwt.StandardClaims
}

//EncodeToken creates a token for WOPI
func EncodeToken(userID string, fileID string) string {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &WOPIToken{
		UserID: userID,
		FileID: fileID,
	})
	tokenstring, err := token.SignedString([]byte(tokenSignPassword))
	if err != nil {
		log.Fatalln(err)
	}
	return tokenstring
}

//DecodeToken decodes a token string an returns WOPIToken and isValid
func DecodeToken(tokenString string) (WOPIToken, bool) {

	wopiToken := WOPIToken{}
	_, err := jwt.ParseWithClaims(tokenString, &wopiToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSignPassword), nil
	})

	if err != nil {
		return WOPIToken{}, false
	}

	return wopiToken, true
}
