package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"

	root "github.com/CollaboraOnline/collabora-mattermost"
)

// Plugin struct is required by the plugin.
type Plugin struct {
	plugin.MattermostPlugin

	client *pluginapi.Client
	router *mux.Router

	manifest model.Manifest

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *Configuration

	siteURL string
}

// OnActivate is called when the plugin is activated.
func (p *Plugin) OnActivate() error {
	p.client = pluginapi.NewClient(p.API, p.Driver)
	p.router = p.InitAPI()
	p.manifest = root.Manifest

	if err := p.registerSiteURL(); err != nil {
		return errors.Wrap(err, "could not register site URL")
	}

	return nil
}

func (p *Plugin) OnDeactivate() error {
	return nil
}

// registerSiteURL fetches the site URL and sets it in the plugin object.
func (p *Plugin) registerSiteURL() error {
	siteURL := p.client.Configuration.GetConfig().ServiceSettings.SiteURL
	if siteURL == nil || *siteURL == "" {
		return errors.New("could not fetch siteURL")
	}
	p.siteURL = *siteURL
	return nil
}

// OnConfigurationChange is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	// if running the plugin for the first time, then run onActivate.
	if p.client == nil {
		if err := p.OnActivate(); err != nil {
			return err
		}
	}

	var configuration = new(Configuration)

	// Load the public configuration fields from the Mattermost server configuration.
	if loadConfigErr := p.client.Configuration.LoadPluginConfiguration(configuration); loadConfigErr != nil {
		return errors.Wrap(loadConfigErr, "failed to load plugin configuration")
	}

	configuration.ProcessConfiguration()

	if err := configuration.IsValid(); err != nil {
		return errors.Wrap(err, "failed to validate configuration")
	}

	if err := p.LoadWopiFileInfo(configuration.WOPIAddress); err != nil {
		return errors.Wrap(err, "could not load wopi file info")
	}

	if err := p.TestFileStoreConnection(); err != nil {
		return errors.Wrap(err, "could not connect to the filestore backend")
	}

	p.setConfiguration(configuration)

	// Send the updated config to the webapp.
	p.client.Frontend.PublishWebSocketEvent(WebsocketEventConfigUpdated, configuration.ToWebappConfig().ToMap(), &model.WebsocketBroadcast{})

	return nil
}

// UserHasLoggedIn is invoked after a user has logged in.
func (p *Plugin) UserHasLoggedIn(_ *plugin.Context, user *model.User) {
	// send the config to the webapp
	config := p.getConfiguration().ToWebappConfig()
	p.client.Frontend.PublishWebSocketEvent(WebsocketEventConfigUpdated, config.ToMap(), &model.WebsocketBroadcast{
		UserId: user.Id,
	})
}

// ServeHTTP handles HTTP requests for the plugin.
func (p *Plugin) ServeHTTP(_ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.client.Log.Debug("New plugin request:", "Host", r.Host, "RequestURI", r.RequestURI, "Method", r.Method)
	p.router.ServeHTTP(w, r)
}
