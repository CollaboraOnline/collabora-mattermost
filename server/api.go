package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"

	root "github.com/CollaboraOnline/collabora-mattermost"
)

const (
	HeaderMattermostUserID = "Mattermost-User-Id"
)

// InitAPI initializes the REST API
func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.Use(p.withRecovery)

	p.handleStaticFiles(r)
	s := r.PathPrefix("/api/v1").Subrouter()

	// Add the custom plugin routes here
	s.HandleFunc("/fileInfo", handleAuthRequired(p.parseFileIDs)).Methods(http.MethodGet)
	s.HandleFunc("/wopiFileList", handleAuthRequired(p.returnWopiFileList)).Methods(http.MethodGet)
	s.HandleFunc("/collaboraURL", handleAuthRequired(p.returnCollaboraOnlineFileURL)).Methods(http.MethodGet)
	s.HandleFunc("/wopi/files/{fileID:[a-z0-9]+}", p.returnWopiFileInfo).Methods(http.MethodGet)
	s.HandleFunc("/wopi/files/{fileID:[a-z0-9]+}/contents", p.getWopiFileContents).Methods(http.MethodGet)
	s.HandleFunc("/wopi/files/{fileID:[a-z0-9]+}/contents", p.saveWopiFileContents).Methods(http.MethodPost)

	// 404 handler
	r.Handle("{anything:.*}", http.NotFoundHandler())
	return r
}

func (p *Plugin) getBaseAPIURL() string {
	return *p.API.GetConfig().ServiceSettings.SiteURL + "/plugins/" + root.Manifest.Id + "/api/v1"
}

func returnStatusOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	m := make(map[string]string)
	m[model.STATUS] = model.STATUS_OK
	_, _ = w.Write([]byte(model.MapToJson(m)))
}

// withRecovery allows recovery from panics
func (p *Plugin) withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				p.API.LogError("Recovered from a panic",
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
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		p.API.LogWarn("Failed to get bundle path.", "Error", err.Error())
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

// parseFileIDs sends the file info to the client (name, extension and id) for each file
// body contains an array with file ids in JSON format
func (p *Plugin) parseFileIDs(w http.ResponseWriter, r *http.Request) {
	//extract fileIDs array from body
	body, bodyReadError := ioutil.ReadAll(r.Body)
	if bodyReadError != nil {
		p.API.LogError("Error when reading body: ", bodyReadError.Error())
		http.Error(w, bodyReadError.Error(), http.StatusBadRequest)
		return
	}

	var fileIDs []string
	if err := json.Unmarshal(body, &fileIDs); err != nil {
		p.API.LogError("Failed to unmarshal request body: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//create an array with more detailed file info for each file
	files := make([]ClientFileInfo, 0, len(fileIDs))
	for _, fileID := range fileIDs {
		fileInfo, fileInfoError := p.API.GetFileInfo(fileID)
		if fileInfoError != nil {
			p.API.LogError("Error when retrieving file info: ", fileInfoError.Error())
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

	responseJSON, _ := json.Marshal(files)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseJSON)
}

// returnWopiFileList returns the list with file extensions and actions associated with these files
func (p *Plugin) returnWopiFileList(w http.ResponseWriter, _ *http.Request) {
	responseJSON, _ := json.Marshal(WopiFiles)
	_, _ = w.Write(responseJSON)
}

// returnCollaboraOnlineFileURL returns the URL and token that the client will use to
// load Collabora Online in the iframe
func (p *Plugin) returnCollaboraOnlineFileURL(w http.ResponseWriter, r *http.Request) {
	//retrieve fileID and file info
	fileID := r.URL.Query().Get("file_id")
	if fileID == "" {
		p.API.LogError("file_id query parameter missing!")
		http.Error(w, "missing file_id parameter", http.StatusBadRequest)
		return
	}

	file, fileError := p.API.GetFileInfo(fileID)
	if fileError != nil {
		p.API.LogError("Failed to retrieve file. Error: ", fileError.Error())
		http.Error(w, "Invalid fileID. Error: "+fileError.Error(), http.StatusBadRequest)
		return
	}

	wopiURL := WopiFiles[strings.ToLower(file.Extension)].URL + "WOPISrc=" + p.getBaseAPIURL() + "/wopi/files/" + fileID
	wopiToken := p.EncodeToken(r.Header.Get(HeaderMattermostUserID), fileID)

	response := struct {
		URL         string `json:"url"`
		AccessToken string `json:"access_token"` //client will pass this token as a POST parameter to Collabora Online when loading the iframe
	}{wopiURL, wopiToken}

	responseJSON, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseJSON)
}

// getWopiFileContents is used by Collabora Online server to get the contents of a file
func (p *Plugin) getWopiFileContents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileID := params["fileID"]

	wopiToken, tokenErr := p.GetWopiTokenFromURI(r.RequestURI)
	if tokenErr != nil || wopiToken.FileID != fileID {
		p.API.LogError(fmt.Sprintf("Invalid token. Error: %v", tokenErr))
		http.Error(w, "Invalid token.", http.StatusBadRequest)
		return
	}

	fileInfo, fileInfoError := p.API.GetFileInfo(fileID)
	if fileInfoError != nil {
		p.API.LogError("Error occurred when retrieving file info: " + fileInfoError.Error())
		http.Error(w, fileInfoError.Error(), http.StatusInternalServerError)
		return
	}

	post, postError := p.API.GetPost(fileInfo.PostId)
	if postError != nil {
		p.API.LogError("Error occurred when retrieving post info for file: " + postError.Error())
		http.Error(w, postError.Error(), http.StatusInternalServerError)
		return
	}

	//check if user has access to the channel where the file was sent
	if !p.API.HasPermissionToChannel(wopiToken.UserID, post.ChannelId, model.PERMISSION_READ_CHANNEL) {
		p.API.LogError("User: " + wopiToken.UserID + " does not have the appropriate permissions: PERMISSION_READ_CHANNEL. Channel: " + post.ChannelId)
		http.Error(w, "You do not have the appropriate permissions.", http.StatusForbidden)
		return
	}

	fileContent, getFileErr := p.API.GetFile(fileID)
	if getFileErr != nil {
		p.API.LogError("Error retrieving file info, fileID: " + fileID)
		http.Error(w, getFileErr.Error(), http.StatusInternalServerError)
		return
	}

	//send file to Collabora Online
	_, _ = w.Write(fileContent)
}

// saveWopiFileContents is used by Collabora Online server to save the updated contents of a file
func (p *Plugin) saveWopiFileContents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileID := params["fileID"]

	wopiToken, tokenErr := p.GetWopiTokenFromURI(r.RequestURI)
	if tokenErr != nil || wopiToken.FileID != fileID {
		p.API.LogError(fmt.Sprintf("Invalid token. Error: %v", tokenErr))
		http.Error(w, "Invalid token.", http.StatusBadRequest)
		return
	}

	fileInfo, fileInfoError := p.API.GetFileInfo(fileID)
	if fileInfoError != nil {
		p.API.LogError("Error occurred when retrieving file info: " + fileInfoError.Error())
		http.Error(w, fileInfoError.Error(), http.StatusInternalServerError)
		return
	}

	post, postError := p.API.GetPost(fileInfo.PostId)
	if postError != nil {
		p.API.LogError("Error occurred when retrieving post info for file: " + postError.Error())
		http.Error(w, postError.Error(), http.StatusInternalServerError)
		return
	}

	//check if user has access to the channel where the file was sent
	if !p.API.HasPermissionToChannel(wopiToken.UserID, post.ChannelId, model.PERMISSION_READ_CHANNEL) {
		p.API.LogError("User: " + wopiToken.UserID + " does not have the appropriate permissions: PERMISSION_READ_CHANNEL. Channel: " + post.ChannelId)
		http.Error(w, "You do not have the appropriate permissions.", http.StatusForbidden)
		return
	}

	//save file received from Collabora Online
	if _, err := p.WriteFile(r.Body, fileInfo.Path); err != nil {
		p.API.LogError("Failed to save the updated file contents.", "Error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnStatusOK(w)
}

// returnWopiFileInfo returns the file information, used by Collabora Online
// see: http://wopi.readthedocs.io/projects/wopirest/en/latest/files/CheckFileInfo.html#checkfileinfo
func (p *Plugin) returnWopiFileInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileID := params["fileID"]

	wopiToken, tokenErr := p.GetWopiTokenFromURI(r.RequestURI)
	if tokenErr != nil || wopiToken.FileID != fileID {
		p.API.LogError(fmt.Sprintf("Invalid token. Error: %v", tokenErr))
		http.Error(w, "Invalid token.", http.StatusBadRequest)
		return
	}

	user, userErr := p.API.GetUser(wopiToken.UserID)
	if userErr != nil {
		p.API.LogError("Error retrieving user. Token UserID is corrupted or the user doesn't exist.", "Error", userErr.Error())
		http.Error(w, userErr.Error(), http.StatusInternalServerError)
		return
	}

	fileInfo, fileInfoErr := p.API.GetFileInfo(fileID)
	if fileInfoErr != nil {
		p.API.LogError("Error retrieving file info, fileID: "+fileID, "Error", fileInfoErr.Error())
		http.Error(w, fileInfoErr.Error(), http.StatusInternalServerError)
		return
	}

	post, postErr := p.API.GetPost(fileInfo.PostId)
	if postErr != nil {
		p.API.LogError("Error retrieving file's post, postId: "+fileInfo.PostId, "Error", postErr.Error())
		http.Error(w, postErr.Error(), http.StatusInternalServerError)
		return
	}

	wopiFileInfo := WopiCheckFileInfo{
		BaseFileName:            fileInfo.Name,
		Size:                    fileInfo.Size,
		OwnerID:                 post.UserId,
		UserID:                  user.Id,
		UserFriendlyName:        user.GetDisplayName(model.SHOW_FULLNAME),
		UserCanWrite:            true,
		UserCanNotWriteRelative: true,
	}

	responseJSON, _ := json.Marshal(wopiFileInfo)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseJSON)
}
