package main

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// MessageWillBePosted will set a post type for each post that contains at least one file
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {

	if len(post.FileIds) > 0 {
		post.Type = "custom_post_with_file"
	}

	return post, ""
}
