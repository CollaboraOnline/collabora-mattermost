package main

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func GenerateEncryptionPassword() string {
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
	return b.String()
}

func (p *Plugin) getHTTPClient() *http.Client {
	config := p.getConfiguration()
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	if config.DisableCertificateVerification {
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{Transport: customTransport}
	return client
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
