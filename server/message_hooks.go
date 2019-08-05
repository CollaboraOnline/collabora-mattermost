package main

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// MessageWillBePosted is invoked when a message is posted by a user before it is committed to the
// database. If you also want to act on edited posts, see MessageWillBeUpdated. Return values
// should be the modified post or nil if rejected and an explanation for the user.
//
// If you don't need to modify or reject posts, use MessageHasBeenPosted instead.
//
// Note that this method will be called for posts created by plugins, including the plugin that created the post.
//
// This demo implementation rejects posts in the demo channel, as well as posts that @-mention
// the demo plugin user.
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {

	if len(post.FileIds) > 0 {
		post.Type = "custom_post_with_file"
	}

	return post, ""
}
