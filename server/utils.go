package main

import (
	"crypto/tls"
	"io"
	"net/http"

	"github.com/mattermost/mattermost-server/v6/shared/filestore"
)

func (p *Plugin) getFileBackend() (filestore.FileBackend, error) {
	license := p.client.System.GetLicense()
	insecure := p.client.Configuration.GetConfig().ServiceSettings.EnableInsecureOutgoingConnections
	serverConfig := p.client.Configuration.GetUnsanitizedConfig()
	backend, err := filestore.NewFileBackend(serverConfig.FileSettings.ToFileBackendSettings(license != nil && *license.Features.Compliance, insecure != nil && *insecure))
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

func (p *Plugin) GetHTTPClient() *http.Client {
	config := p.getConfiguration()
	customTransport := http.DefaultTransport.(*http.Transport).Clone()

	if config.SkipSSLVerify {
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{Transport: customTransport}
	return client
}

func (p *Plugin) GetFilePermissionsKey(fileID string) string {
	return p.manifest.Id + "_file_permissions_" + fileID
}
