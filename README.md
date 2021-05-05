# Collabora Online Mattermost plugin

This plugin enables Mattermost users to preview and collaboratively edit documents (simple text files, word, spreadsheet and presentation documents) via a [Collabora Online](https://www.collaboraoffice.com/collabora-online/) Server using the WOPI protocol.

![Demo](https://www.collaboraoffice.com/wp-content/uploads/2019/09/plugin_video.gif)

## How does it work?

The plugin will parse the messages that contain attachments and will display a list with the files that can be viewed or edited. The file is automatically saved when editing, so you don’t have to do it manually.

**Note**: This plugin will set the post type "custom_post_with_file" for each post that contains attached files.

Below you can see a list with the supported formats:

- OpenDocument format: .ods, .odt, .odp, .odg etc.
- Microsoft: .doc, .docx, .xls, .xlsx, .ppt, .pptx etc.
- Others: .txt, .csv, .gif, .jpg, .jpeg, .png, .svg, .pdf etc.

Collabora Online uses a WOPI-like protocol (client) to access the files on your Mattermost server (host). You can read more about it on https://wopi.readthedocs.io. Hence, you will also need a Collabora Online instance to use the plugin.
You can build your own, or conveniently use a version of our [CODE edition](https://www.collaboraoffice.com/code/).

## Installation

1. You can get the latest version on the [releases page](https://github.com/CollaboraOnline/collabora-mattermost/releases/latest).
1. Upload & install it via Mattermost [System Console](https://about.mattermost.com/default-plugin-uploads)
1. After installing the plugin you should go to the plugin's settings in System Console and set the Collabora Online address (more about this below).
   A page refresh may be required for the plugin’s settings to appear in the System Console.

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

2. Starting the CODE docker container for the first time will take a while!
   `docker run` returns quickly, but you will not be able to actually use it for the next 5 minutes or so.
   Look at `docker logs -f code` for details.

3. Replace the string `my-local-ipaddress` with your local IP address!
   Each `.` of the address must be accompanied by the double-backslash `\\`, thus avoiding misinterpretations.

## Building the plugin

- Make sure you have following components installed:
    - Go - v1.16 - [Getting Started](https://golang.org/doc/install)
      > **Note:** If you have installed Go to a custom location, make sure the `$GOROOT` variable is set properly. Refer [Installing to a custom location](https://golang.org/doc/install#install).
    - Make

- Build the plugin:
    ```
    make
    ```

- This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:
    ```
    dist/com.collaboraonline.mattermost-x.y.z.tar.gz
    ```
- This plugin contains both a server and web app portion.
  Read the Mattermost documentation about the [Developer Workflow](https://developers.mattermost.com/extend/plugins/developer-workflow/)
  and [Developer Setup](https://developers.mattermost.com/extend/plugins/developer-setup/) for more information about developing and extending plugins.

## Make & auto installation

There is a build target to automate deploying and enabling the plugin to your server, but it requires configuration and [http](https://httpie.org/) to be installed:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=username
export MM_ADMIN_PASSWORD=password
make deploy
```

## Troubleshooting

- Q. Failed to read document from storage. Please contact your storage server administrator.  
  A. Make sure you are running both Mattermost and Collabora Server with the same protocol (http/https).
