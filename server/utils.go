package main

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func (p *Plugin) getHTTPClient() *http.Client {
	config := p.getConfiguration()
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	if config.DisableCertificateVerification {
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{Transport: customTransport}
	return client
}

//Because the access_token get's removed from Query parameters by Mattermost before
//it reaches the plugin HTTP request parser, it should be manually extracted from the URI
func getAccessTokenFromURI(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse uri")
	}
	m, parseErr := url.ParseQuery(u.RawQuery)
	if parseErr != nil {
		return "", errors.Wrap(parseErr, "failed to parse raw query")
	}
	return m["access_token"][0], nil
}
