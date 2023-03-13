package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
)

const (
	HeaderMattermostUserID = "Mattermost-User-Id"
)

// InitAPI initializes the REST API.
func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.Use(p.withRecovery)

	p.handleStaticFiles(r)
	s := r.PathPrefix("/api/v1").Subrouter()

	// Add the custom plugin routes here
	s.HandleFunc("/config", p.getWebappConfig).Methods(http.MethodGet)
	s.HandleFunc("/files/{fileID:[A-Za-z0-9_-]+}/access", handleAuthRequired(p.handleSaveFilePermissions)).Methods(http.MethodPost)
	s.HandleFunc("/channels/{channelID:[A-Za-z0-9_-]+}/files/new", handleAuthRequired(p.createFileFromTemplate)).Methods(http.MethodPost).Queries("name", "{name}", "ext", "{ext}")
	s.HandleFunc("/fileInfo", handleAuthRequired(p.getClientFileInfos)).Methods(http.MethodGet)
	s.HandleFunc("/wopiFileList", handleAuthRequired(p.returnWopiFileList)).Methods(http.MethodGet)
	s.HandleFunc("/collaboraURL", handleAuthRequired(p.returnCollaboraOnlineFileURL)).Methods(http.MethodGet)
	s.HandleFunc("/wopi/files/{fileID:[a-z0-9]+}", p.getWopiFileInfo).Methods(http.MethodGet)
	s.HandleFunc("/wopi/files/{fileID:[a-z0-9]+}/contents", p.getWopiFileContents).Methods(http.MethodGet)
	s.HandleFunc("/wopi/files/{fileID:[a-z0-9]+}/edit", p.getWopiFileInfoEditable).Methods(http.MethodGet)
	s.HandleFunc("/wopi/files/{fileID:[a-z0-9]+}/edit/contents", p.getWopiFileContents).Methods(http.MethodGet)
	s.HandleFunc("/wopi/files/{fileID:[a-z0-9]+}/edit/contents", p.saveWopiFileContents).Methods(http.MethodPost)

	// 404 handler
	r.Handle("{anything:.*}", http.NotFoundHandler())
	return r
}

func (p *Plugin) getAPIBaseURL() string {
	return p.siteURL + "/plugins/" + p.manifest.Id + "/api/v1"
}

func returnStatusOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	m := make(map[string]string)
	m[model.STATUS] = model.StatusOk
	_, _ = w.Write([]byte(model.MapToJSON(m)))
}

// withRecovery allows recovery from panics.
func (p *Plugin) withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				p.client.Log.Error("Recovered from a panic",
					"url", r.URL.String(),
					"error", x,
					"stack", string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// handleStaticFiles handles the static files under the assets directory.
func (p *Plugin) handleStaticFiles(r *mux.Router) {
	bundlePath, err := p.client.System.GetBundlePath()
	if err != nil {
		p.client.Log.Warn("Failed to get bundle path.", "Error", err.Error())
		return
	}

	// This will serve static files from the 'assets' directory under '/static/<filename>'
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(bundlePath, "assets")))))
}

// handleAuthRequired verifies if provided request is performed by a logged-in Mattermost user.
func handleAuthRequired(handleFunc func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(HeaderMattermostUserID)
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		handleFunc(w, r)
	}
}

func (p *Plugin) getWebappConfig(w http.ResponseWriter, r *http.Request) {
	var config = p.getConfiguration().ToWebappConfig()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(config); err != nil {
		p.client.Log.Warn("failed to write response", "error", err.Error())
	}
}

// handleSaveFilePermissions allows setting file permissions when the edit permissions system console setting is enabled.
func (p *Plugin) handleSaveFilePermissions(w http.ResponseWriter, r *http.Request) {
	conf := p.getConfiguration()
	if !conf.FileEditPermissions {
		p.client.Log.Error("handleSaveFilePermissions: Edit permissions feature is disabled in system console.")
		http.Error(w, "Edit permissions feature is disabled in system console.", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	fileID := params["fileID"]
	userID := r.Header.Get(HeaderMattermostUserID)
	permissionQuery := r.URL.Query().Get("permission")
	if _, ok := AllowedFilePermissions[permissionQuery]; !ok {
		p.client.Log.Error("Invalid permission query param.", "permissionQuery", permissionQuery)
		http.Error(w, "Invalid permission query param.", http.StatusBadRequest)
		return
	}

	if err := p.setFilePermissions(fileID, userID, permissionQuery, false); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnStatusOK(w)
}

func (p *Plugin) setFilePermissions(fileID, userID, permission string, bypassFileOwnerCheck bool) error {
	fileInfo, fileInfoError := p.client.File.GetInfo(fileID)
	if fileInfoError != nil {
		p.client.Log.Error("Error when retrieving file info: ", fileInfoError.Error(), "fileID", fileID)
		return errors.Wrap(fileInfoError, "error when retrieving file info, invalid fileID")
	}

	post, postError := p.client.Post.GetPost(fileInfo.PostId)
	if postError != nil {
		p.client.Log.Error("Error occurred when retrieving post info for file: " + postError.Error())
		return errors.Wrap(postError, "error when retrieving post for file")
	}

	if !bypassFileOwnerCheck && post.UserId != userID {
		p.client.Log.Error("User does not have access to change file permissions.")
		return errors.New("only the file owner can change file permissions")
	}

	filePermissionsKey := p.GetFilePermissionsKey(fileID)
	post.AddProp(filePermissionsKey, permission)
	if postErr := p.client.Post.UpdatePost(post); postErr != nil {
		p.client.Log.Error("Failed to update post", "Error", postErr.Error())
		return errors.Wrap(postError, "error when saving post with updated permissions")
	}

	return nil
}

// createFileFromTemplate creates a new file from template in the given channel.
func (p *Plugin) createFileFromTemplate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	channelID := params["channelID"]
	_, channelErr := p.client.Channel.Get(channelID)
	if channelErr != nil {
		p.client.Log.Error("Invalid or missing channel ID: ", channelErr.Error(), "channelID", channelID)
		http.Error(w, channelErr.Error(), http.StatusBadRequest)
		return
	}

	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		http.Error(w, "missing filename", http.StatusBadRequest)
		return
	}

	fileExt := r.URL.Query().Get("ext")
	if fileExt == "" {
		http.Error(w, "missing file extension", http.StatusBadRequest)
		return
	}

	templateName, templateFound := TemplateFromExt[fileExt]
	if !templateFound {
		p.client.Log.Warn("no template found for file extension: " + fileExt)
		http.Error(w, "template not found for provided file extension", http.StatusBadRequest)
		return
	}

	bundlePath, err := p.client.System.GetBundlePath()
	if err != nil {
		p.client.Log.Warn("Failed to get bundle path.", "Error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateDir := filepath.Join(bundlePath, "assets", "templates")
	templatePath := path.Join(templateDir, templateName)

	templateFileReader, err := os.Open(templatePath)
	if err != nil {
		p.client.Log.Error("Failed to get the template content.", "Error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer templateFileReader.Close()

	fileInfo, appErr := p.client.File.Upload(templateFileReader, fileName+"."+fileExt, channelID)
	if appErr != nil {
		p.client.Log.Error("Failed to upload the template file.", "Error", appErr.Error())
		http.Error(w, appErr.Error(), http.StatusInternalServerError)
		return
	}

	post := &model.Post{
		ChannelId: channelID,
		UserId:    r.Header.Get(HeaderMattermostUserID),
		FileIds:   model.StringArray{fileInfo.Id},
	}

	if appErr := p.client.Post.CreatePost(post); appErr != nil {
		p.client.Log.Error("Failed to create post with the template file.", "Error", appErr.Error())
		http.Error(w, appErr.Error(), http.StatusInternalServerError)
		return
	}

	returnStatusOK(w)
}

// getClientFileInfos sends the ClientFileInfo (name, extension and id) for each file to the client.
// The response body contains an array with file ids in JSON format.
func (p *Plugin) getClientFileInfos(w http.ResponseWriter, r *http.Request) {
	// extract fileIDs array from request body
	var fileIDs []string
	err := json.NewDecoder(r.Body).Decode(&fileIDs)
	if err != nil {
		p.client.Log.Error("failed to decode fileIDs", "err", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create an array with more detailed file info for each file
	files := make([]ClientFileInfo, 0, len(fileIDs))
	for _, fileID := range fileIDs {
		fileInfo, fileInfoError := p.client.File.GetInfo(fileID)
		if fileInfoError != nil {
			p.client.Log.Error("Error when retrieving file info: ", fileInfoError.Error())
			continue
		}
		if value, ok := WopiFiles[strings.ToLower(fileInfo.Extension)]; ok {
			file := ClientFileInfo{
				fileInfo.Id,
				fileInfo.Name,
				fileInfo.Extension,
				value.Action,
			}
			files = append(files, file)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		p.client.Log.Warn("failed to write response", "error", err.Error())
	}
}

// returnWopiFileList returns the list with file extensions and actions associated with these files.
func (p *Plugin) returnWopiFileList(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(WopiFiles); err != nil {
		p.client.Log.Warn("failed to write response", "error", err.Error())
	}
}

// returnCollaboraOnlineFileURL returns the URL and token that the client will use to
// load Collabora Online in the iFrame.
func (p *Plugin) returnCollaboraOnlineFileURL(w http.ResponseWriter, r *http.Request) {
	// retrieve fileID and file info
	fileID := r.URL.Query().Get("file_id")
	if fileID == "" {
		p.client.Log.Error("Failed to retrieve file. `file_id` query parameter missing!")
		http.Error(w, "missing file_id parameter", http.StatusBadRequest)
		return
	}

	fileInfo, fileError := p.client.File.GetInfo(fileID)
	if fileError != nil {
		p.client.Log.Error("Failed to retrieve file. Error: ", fileError.Error())
		http.Error(w, "Invalid fileID. Error: "+fileError.Error(), http.StatusBadRequest)
		return
	}

	post, postError := p.client.Post.GetPost(fileInfo.PostId)
	if postError != nil {
		p.client.Log.Error("Error occurred when retrieving post info for file: " + postError.Error())
		http.Error(w, postError.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Header.Get(HeaderMattermostUserID)
	conf := p.getConfiguration()
	existingFilePermission := post.GetProp(p.GetFilePermissionsKey(fileID))

	// initialize file permission if not already exists
	if conf.FileEditPermissions && existingFilePermission == nil {
		// If the edit permissions feature is enabled,
		// set the default permission to allow only the owner to edit
		// skip the file owner check to handle the scenario
		// when some user A has uploaded a file(and has not opened it with collabora), which is opened for first time by user B

		if err := p.setFilePermissions(fileID, userID, PermissionOwner, true); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	wopiURL := WopiFiles[strings.ToLower(fileInfo.Extension)].URL + "WOPISrc=" + (p.getAPIBaseURL() + "/wopi/files/" + fileID)
	wopiToken := p.EncodeToken(userID, fileID)

	response := struct {
		URL         string `json:"url"`
		AccessToken string `json:"access_token"` // client will pass this token as a POST parameter to Collabora Online when loading the iFrame
	}{wopiURL, wopiToken}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		p.client.Log.Warn("failed to write response", "error", err.Error())
	}
}

// getWopiFileContents is used by Collabora Online server to get the contents of a file.
func (p *Plugin) getWopiFileContents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileID := params["fileID"]

	wopiToken, tokenErr := p.GetWopiTokenFromURI(r.RequestURI)
	if tokenErr != nil || wopiToken.FileID != fileID {
		p.client.Log.Error(fmt.Sprintf("Invalid token. Error: %v", tokenErr))
		http.Error(w, "Invalid token.", http.StatusBadRequest)
		return
	}

	fileInfo, fileInfoError := p.client.File.GetInfo(fileID)
	if fileInfoError != nil {
		p.client.Log.Error("Error occurred when retrieving file info: " + fileInfoError.Error())
		http.Error(w, fileInfoError.Error(), http.StatusInternalServerError)
		return
	}

	post, postError := p.client.Post.GetPost(fileInfo.PostId)
	if postError != nil {
		p.client.Log.Error("Error occurred when retrieving post info for file: " + postError.Error())
		http.Error(w, postError.Error(), http.StatusInternalServerError)
		return
	}

	// check if user has access to the channel where the file was sent
	if !p.client.User.HasPermissionToChannel(wopiToken.UserID, post.ChannelId, model.PermissionReadChannel) {
		p.client.Log.Error("User: " + wopiToken.UserID + " does not have the appropriate permissions: PermissionReadChannel. Channel: " + post.ChannelId)
		http.Error(w, "You do not have the appropriate permissions.", http.StatusForbidden)
		return
	}

	fileContent, getFileErr := p.client.File.Get(fileID)
	if getFileErr != nil {
		p.client.Log.Error("Error retrieving file info, fileID: " + fileID)
		http.Error(w, getFileErr.Error(), http.StatusInternalServerError)
		return
	}

	// send file to Collabora Online
	_, _ = io.Copy(w, fileContent)
}

// saveWopiFileContents is used by Collabora Online server to save the updated contents of a file.
func (p *Plugin) saveWopiFileContents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileID := params["fileID"]

	wopiToken, tokenErr := p.GetWopiTokenFromURI(r.RequestURI)
	if tokenErr != nil || wopiToken.FileID != fileID {
		p.client.Log.Error(fmt.Sprintf("Invalid token. Error: %v", tokenErr))
		http.Error(w, "Invalid token.", http.StatusBadRequest)
		return
	}

	fileInfo, fileInfoError := p.client.File.GetInfo(fileID)
	if fileInfoError != nil {
		p.client.Log.Error("Error occurred when retrieving file info: " + fileInfoError.Error())
		http.Error(w, fileInfoError.Error(), http.StatusInternalServerError)
		return
	}

	post, postError := p.client.Post.GetPost(fileInfo.PostId)
	if postError != nil {
		p.client.Log.Error("Error occurred when retrieving post info for file: " + postError.Error())
		http.Error(w, postError.Error(), http.StatusInternalServerError)
		return
	}

	conf := p.getConfiguration()
	filePermission := post.GetProp(p.GetFilePermissionsKey(fileID))
	canChannelEdit := !conf.FileEditPermissions || filePermission == PermissionChannel
	canOwnerOnlyEdit := conf.FileEditPermissions && filePermission == PermissionOwner
	canCurrentUserEdit := (canChannelEdit && p.client.User.HasPermissionToChannel(wopiToken.UserID, post.ChannelId, model.PermissionReadChannel)) || (canOwnerOnlyEdit && post.UserId == wopiToken.UserID)

	if !canCurrentUserEdit {
		p.client.Log.Error("User does not have the appropriate permissions to edit the file.", "channelID", post.ChannelId, "userID", wopiToken.UserID)
		http.Error(w, "You do not have the appropriate permissions.", http.StatusForbidden)
		return
	}

	// save the file received from Collabora Online using the appropriate file backend.
	if _, err := p.WriteFile(r.Body, fileInfo.Path); err != nil {
		p.client.Log.Error("Failed to save the updated file contents.", "Error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnStatusOK(w)
}

// generateWopiFileInfo generates the file information, used by Collabora Online
// see: http:// wopi.readthedocs.io/projects/wopirest/en/latest/files/CheckFileInfo.html#checkfileinfo
func (p *Plugin) generateWopiFileInfo(wopiToken WopiToken, userCanEdit bool) (*WopiCheckFileInfo, error) {
	user, userErr := p.client.User.Get(wopiToken.UserID)
	if userErr != nil {
		p.client.Log.Error("Error retrieving user. Token UserID is corrupted or the user doesn't exist.", "TokenUserID", wopiToken.UserID, "Error", userErr.Error())
		return nil, userErr
	}

	fileInfo, fileInfoErr := p.client.File.GetInfo(wopiToken.FileID)
	if fileInfoErr != nil {
		p.client.Log.Error("Error retrieving file info", "FileID", wopiToken.FileID, "Error", fileInfoErr.Error())
		return nil, fileInfoErr
	}

	post, postErr := p.client.Post.GetPost(fileInfo.PostId)
	if postErr != nil {
		p.client.Log.Error("Error retrieving file's post.", "PostID", fileInfo.PostId, "Error", postErr.Error())
		return nil, postErr
	}

	wopiFileInfo := &WopiCheckFileInfo{
		BaseFileName:            fileInfo.Name,
		Size:                    fileInfo.Size,
		OwnerID:                 post.UserId,
		UserID:                  user.Id,
		UserFriendlyName:        user.GetDisplayName(model.ShowFullName),
		UserCanWrite:            userCanEdit,
		UserCanNotWriteRelative: true,
	}

	return wopiFileInfo, nil
}

// getWopiFileInfo returns the file information, used by Collabora Online.
func (p *Plugin) getWopiFileInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileID := params["fileID"]

	wopiToken, tokenErr := p.GetWopiTokenFromURI(r.RequestURI)
	if tokenErr != nil || wopiToken.FileID != fileID {
		p.client.Log.Error(fmt.Sprintf("Invalid token. Error: %v", tokenErr))
		http.Error(w, "Invalid token.", http.StatusBadRequest)
		return
	}

	wopiFileInfo, wopiFileInfoErr := p.generateWopiFileInfo(wopiToken, false)
	if wopiFileInfoErr != nil {
		http.Error(w, wopiFileInfoErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wopiFileInfo); err != nil {
		p.client.Log.Warn("failed to write response", "error", err.Error())
	}
}

// getWopiFileInfoEditable returns the file information, used by Collabora Online
// with editable set to true.
// TODO: refactor this to same function as getWopiFileInfo.
func (p *Plugin) getWopiFileInfoEditable(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileID := params["fileID"]

	wopiToken, tokenErr := p.GetWopiTokenFromURI(r.RequestURI)
	if tokenErr != nil || wopiToken.FileID != fileID {
		p.client.Log.Error(fmt.Sprintf("Invalid token. Error: %v", tokenErr))
		http.Error(w, "Invalid token.", http.StatusBadRequest)
		return
	}

	wopiFileInfo, wopiFileInfoErr := p.generateWopiFileInfo(wopiToken, true)
	if wopiFileInfoErr != nil {
		http.Error(w, wopiFileInfoErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wopiFileInfo); err != nil {
		p.client.Log.Warn("failed to write response", "error", err.Error())
	}
}
