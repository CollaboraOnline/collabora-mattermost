package main

import (
	"crypto/tls"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/shared/filestore"
)

func (p *Plugin) getFileBackend() (filestore.FileBackend, error) {
	license := p.API.GetLicense()
	serverConfig := p.API.GetUnsanitizedConfig()
	backend, err := filestore.NewFileBackend(serverConfig.FileSettings.ToFileBackendSettings(license != nil && *license.Features.Compliance))
	if err != nil {
		return nil, err
	}
	return backend, nil
}

func (p *Plugin) WriteFile(fr io.Reader, path string) (int64, error) {
	backend, err := p.getFileBackend()
	if err != nil {
		return 0, err
	}

	result, nErr := backend.WriteFile(fr, path)
	if nErr != nil {
		return result, nErr
	}
	return result, nil
}

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

func (p *Plugin) GetHTTPClient() *http.Client {
	config := p.getConfiguration()
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	if config.DisableCertificateVerification {
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{Transport: customTransport}
	return client
}
