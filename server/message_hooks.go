package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// MessageWillBePosted will set a post type for each post that contains at least one file
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	//change the post type only if it contains any files that can be viewed/edited with Collabora Online
	changePostType := false
	for _, fileID := range post.FileIds {
		fileInfo, fileInfoError := p.API.GetFileInfo(fileID)
		if fileInfoError != nil {
			p.API.LogError("Could not retrieve file info on message post")
			continue
		}
		if _, ok := WopiFiles[strings.ToLower(fileInfo.Extension)]; ok {
			changePostType = true
		}
	}

	if changePostType {
		post.Type = "custom_post_with_file"
	}

	return post, ""
}
