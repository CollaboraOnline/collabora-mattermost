# Collabora Online Mattermost plugin

This plugin enables Mattermost users to preview and collaboratively edit documents (simple text files, word, spreadsheet and presentation documents) via a [Collabora Online](https://www.collaboraoffice.com/collabora-online/) Server using the WOPI protocol.

![Demo](https://www.collaboraoffice.com/wp-content/uploads/2019/09/plugin_video.gif)

## How does it work?

The plugin will parse the messages that contain attachments and will display a list with the files that can be viewed or edited. The file is automatically saved when editing, so you don’t have to do it manually.

Below you can see a list with the supported formats:

- OpenDocument format: .ods, .odt, .odp, .odg etc.
- Microsoft: .doc, .docx, .xls, .xlsx, .ppt, .pptx etc.
- Others: .txt, .csv, .gif, .jpg, .jpeg, .png, .svg, .pdf etc.

Please note that files like .pdf, .jpg, .svg, and others can only be viewed and not edited.
  
Collabora Online uses a WOPI-like protocol (client) to access the files on your Mattermost server (host). You can read more about it on https://wopi.readthedocs.io. Hence, you will also need a Collabora Online instance to use the plugin.
You can build your own, or conveniently use a version of our [CODE edition](https://www.collaboraoffice.com/code/).

## Installation

1. You can get the latest version on the [releases page](https://github.com/CollaboraOnline/collabora-mattermost/releases/latest).
1. Upload this file in the Mattermost **System Console > Plugins > Management** page to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#custom-plugins).
1. After installing the plugin, you should go to the plugin's settings in System Console and set the Collabora Online address (more about this below).
   A page refresh may be required for the plugin’s settings to appear in the System Console.


### System Console Settings

- **Collabora Online URL**:
  The URL (and port) of the Collabora Online server that provides the editing functionality as a WOPI client. Collabora Online should use the same protocol (http:// or https://) as the server installation. Naturally, https:// is recommended.

- **Disable certificate verification**:
  You must enable this setting and accept the local ssl certificate in your browser to be able to preview and edit files when using a self-signed certificate for CollaboraOnline server.

- **Token Encryption Key**:
  The plugin internally generates and passes an access token to Collabora Online that is used later by it to do various operations.
  This setting is the key used to encrypt/decrypt such tokens and must be generated once before starting the plugin for the first time.

## Development

You can use the self-hosted Collabora Online Server i.e. the [CODE](https://www.collaboraoffice.com/code/) docker image.

```sh
# pull collabora image
docker pull collabora/code

# Run with http url: http://localhost:9980/
docker run -t -d -p 127.0.0.1:9980:9980 -p [::1]:9980:9980 -e 'domain=my\\.-local\\.-ip\\.address' -e "username=admin" -e "password=secret" --restart always --cap-add MKNOD -e "extra_params=--o:ssl.enable=false" --name=code collabora/code

# Run with https url: https://localhost:9980/
docker run -t -d -p 127.0.0.1:9980:9980 -p [::1]:9980:9980 -e 'domain=my\\.-local\\.-ip\\.address' -e "username=admin" -e "password=secret" --restart always --cap-add MKNOD --name=code collabora/code
```

**Additional Notes**: 

1. You need to make sure SSL is enabled (with local ssl certificate accepted in your browser) or disabled for both Mattermost and CODE.
   If one setup has https enabled and the other not, it will not work.

1. Starting the CODE docker container for the first time will take a while!
   `docker run` returns quickly, but you will not be able to actually use it for the next 5 minutes or so.
   Look at `docker logs -f code` for details.

1. Replace the string `my-local-ipaddress` with your local IP address!
   Each `.` of the address must be accompanied by the double-backslash `\\`, thus avoiding misinterpretations.

1. The `domain` environment variable should point to the Mattermost server's IP address and not the Collabora Server.

1. If you are using a self-signed certificate with mattermost running over `https`, you must enable the `Disable certificate verification` system console setting 
   and accept the local ssl certificate in your browser to be able to preview files.

## Building the plugin

- Make sure you have following components installed:
    - Go - v1.16 - [Getting Started](https://golang.org/doc/install)
      > **Note:** If you have installed Go to a custom location, make sure the `$GOROOT` variable is set properly. Refer [Installing to a custom location](https://golang.org/doc/install#install).
    - NodeJS - v14.17 and NPM - [Downloading and installing Node.js and npm](https://docs.npmjs.com/getting-started/installing-node).
    - Make

- Note that this project uses [Go modules](https://github.com/golang/go/wiki/Modules). Be sure to locate the project outside of `$GOPATH`.
To learn more about plugins, see [plugin documentation](https://developers.mattermost.com/extend/plugins/).

- Build your plugin:
    ```
    make dist
    ```

- This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:
    ```
    dist/com.collaboraonline.mattermost-x.y.z.tar.gz
    ```

- This plugin contains both a server and web app portion.
  Read the Mattermost documentation about the [Developer Workflow](https://developers.mattermost.com/extend/plugins/developer-workflow/)
  and [Developer Setup](https://developers.mattermost.com/extend/plugins/developer-setup/) for more information about developing and extending plugins.

## Make & auto installation

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options.

### Deploying with Local Mode

If your Mattermost server is running locally, you can enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode) to streamline deploying your plugin. Edit your server configuration as follows:

```json
{
    "ServiceSettings": {
        ...
        "EnableLocalMode": true,
        "LocalModeSocketLocation": "/var/tmp/mattermost_local.socket"
    }
}
```

and then deploy your plugin:
```
make deploy
```

You may also customize the Unix socket path:
```
export MM_LOCALSOCKETPATH=/var/tmp/alternate_local.socket
make deploy
```

If developing a plugin with a webapp, watch for changes and deploy those automatically:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=<mm-admin-auth-token>
make watch
```

### Deploying with credentials

Alternatively, you can authenticate with the server's API with credentials:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```

or with a [personal access token](https://docs.mattermost.com/developer/personal-access-tokens.html):
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=<mm-admin-auth-token>
make deploy
```

## Troubleshooting

- Q. Failed to read document from storage. Please contact your storage server administrator.  
  A. Make sure you are running both Mattermost and Collabora Server with the same protocol (http/https).
     Check your Mattermost logs for more information.

- Q. CollaboraOnline Server URL in the system console does not get updated.
  A. You may need to disable and re-enable the plugin for the server URL (or other system console settings) changes to take effect.
