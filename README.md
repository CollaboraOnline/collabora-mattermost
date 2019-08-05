# Collabora Online

This plugin integrates [Collabora Online](https://www.collaboraoffice.com/collabora-online/) with Mattermost so users can view or edit files directly in Mattermost.

## Installation

You can get the latest version on the release page.
Upload & install it via [System Console](https://about.mattermost.com/default-plugin-uploads)


## Configuring

After installing the plugin you should go to plugin's settings in System Console and set the Collabora Online WOPI address.

## Notes

This plugin will set the post type "custom_post_with_file" for each post that contains attached files.

# Development

To make the plugin, you need to install a development version of Mattermost [https://developers.mattermost.com/contribute/server/developer-setup/](https://developers.mattermost.com/contribute/server/developer-setup/)

## Make & manual installation

Build the plugin:
```
make
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/com.collaboraonline.mattermost-x.y.z.tar.gz
```

## Make & auto installation

There is a build target to automate deploying and enabling the plugin to your server, but it requires configuration and [http](https://httpie.org/) to be installed:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=username
export MM_ADMIN_PASSWORD=password
make deploy
```


