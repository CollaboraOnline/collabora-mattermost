module github.com/CollaboraOnline/collabora-mattermost

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/gorilla/mux v1.8.0
	github.com/mattermost/mattermost-server/v5 v5.39.0
	github.com/pkg/errors v0.9.1
)

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.1+incompatible
