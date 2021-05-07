package main

import (
	"sync"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

//Plugin required by plugin
type Plugin struct {
	plugin.MattermostPlugin
	configurationLock sync.RWMutex
	configuration     *configuration
}

//OnActivate is called when the plugin is activated
func (p *Plugin) OnActivate() error {
	p.GenerateEncryptionPassword()
	return nil
}
