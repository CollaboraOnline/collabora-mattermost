module github.com/CollaboraOnline/collabora-mattermost

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/mattermost/mattermost-plugin-api v0.1.2-0.20221110071900-f8b73bc6795e
	// mmgoget: github.com/mattermost/mattermost-server/v6@v7.7.0 is replaced by -> github.com/mattermost/mattermost-server/v6@ea08d47f60
	github.com/mattermost/mattermost-server/v6 v6.0.0-20230113170349-ea08d47f6051
	github.com/pkg/errors v0.9.1
)
