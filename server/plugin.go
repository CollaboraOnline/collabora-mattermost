package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

// Plugin required by plugin
type Plugin struct {
	plugin.MattermostPlugin
	router            *mux.Router
	configurationLock sync.RWMutex
	configuration     *configuration
}

// OnActivate is called when the plugin is activated
func (p *Plugin) OnActivate() error {
	p.router = p.InitAPI()
	return nil
}

// OnConfigurationChange is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	// Load the public configuration fields from the Mattermost server configuration.
	if loadConfigErr := p.API.LoadPluginConfiguration(configuration); loadConfigErr != nil {
		return errors.Wrap(loadConfigErr, "failed to load plugin configuration")
	}

	if err := configuration.ProcessConfiguration(); err != nil {
		p.API.LogError("Error in ProcessConfiguration.", "Error", err.Error())
		return errors.Wrap(err, "failed to process configuration")
	}

	if err := configuration.IsValid(); err != nil {
		return errors.Wrap(err, "failed to validate configuration")
	}

	if err := p.LoadWopiFileInfo(configuration.WOPIAddress); err != nil {
		return errors.Wrap(err, "could not load wopi file info")
	}

	p.setConfiguration(configuration)

	// send the updated config to the webapp
	p.API.PublishWebSocketEvent(WebsocketEventConfigUpdated, configuration.ToWebappConfig().ToMap(), &model.WebsocketBroadcast{})

	return nil
}

// UserHasLoggedIn is invoked after a user has logged in.
func (p *Plugin) UserHasLoggedIn(c *plugin.Context, user *model.User) {
	// send the config to the webapp
	config := p.getConfiguration().ToWebappConfig()
	p.API.PublishWebSocketEvent(WebsocketEventConfigUpdated, config.ToMap(), &model.WebsocketBroadcast{})
}

// ServeHTTP handles HTTP requests for the plugin.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.API.LogDebug("New plugin request:", "Host", r.Host, "RequestURI", r.RequestURI, "Method", r.Method)
	p.router.ServeHTTP(w, r)
}
