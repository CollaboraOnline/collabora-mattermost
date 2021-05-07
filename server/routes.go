package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	//send file info to client (name, extension and id) for each file
	//body contains an array with file ids in JSON format
	if r.URL.Path == "/fileInfo" {
		p.parseFileIDs(w, r)
	}

	//send list with file extensions and actions associated with these files
	if r.URL.Path == "/wopiFileList" {
		p.returnWopiFileList(w, r)
	}

	//send URL and token which the client will use to load Collabora Online in the iframe
	if r.URL.Path == "/collaboraURL" {
		p.returnCollaboraOnlineFileURL(w, r)
	}

	//https://<WOPI host URL>/<...>/wopi/files/<id>/contents gets/saves a file, used by Collabora Online
	if strings.Contains(r.URL.Path, "/wopi/files/") && strings.Contains(r.URL.Path, "contents") {
		p.parseWopiRequests(w, r)
	}

	//https://<WOPI host URL>/<...>/wopi/files/<id> returns file info, used by Collabora Online
	if strings.Contains(r.URL.Path, "/wopi/files/") && !strings.Contains(r.URL.Path, "contents") {
		p.returnFileInfoForWOPI(w, r)
	}

	//for serving assets from the assets/folder to the client side of the plugin
	if strings.Contains(r.URL.Path, "/assets/") {
		p.serveAsset(w, r)
	}
}

func (p *Plugin) parseFileIDs(w http.ResponseWriter, r *http.Request) {
	//extract fileIDs array from body
	body, bodyReadError := ioutil.ReadAll(r.Body)
	if bodyReadError != nil {
		p.API.LogError("Error when reading body: ", bodyReadError.Error())
		return
	}
	var fileIDs []string
	_ = json.Unmarshal(body, &fileIDs)

	//create an array with more detailed file info for each file
	files := make([]CollaboraFileInfo, 0, len(fileIDs))
	for _, fileID := range fileIDs {
		fileInfo, fileInfoError := p.API.GetFileInfo(fileID)
		if fileInfoError != nil {
			p.API.LogError("Error when retrieving file info: ", fileInfoError.Error())
			continue
		}
		if value, ok := WOPIFiles[strings.ToLower(fileInfo.Extension)]; ok {
			file := CollaboraFileInfo{
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
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}

func (p *Plugin) returnWopiFileList(w http.ResponseWriter, r *http.Request) {
	responseJSON, _ := json.Marshal(WOPIFiles)
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}

func (p *Plugin) returnCollaboraOnlineFileURL(w http.ResponseWriter, r *http.Request) {
	//retrieve fileID and file info
	queryFileID, ok := r.URL.Query()["file_id"]
	if !ok {
		p.API.LogError("file_id query parameter missing!")
		return
	}
	fileID := queryFileID[0]
	file, fileError := p.API.GetFileInfo(fileID)
	if fileError != nil {
		p.API.LogError("Failed to retrieve file. Error: ", fileError.Error())
		return
	}

	wopiURL := WOPIFiles[strings.ToLower(file.Extension)].URL + "WOPISrc=" + *p.API.GetConfig().ServiceSettings.SiteURL + "/plugins/" + manifest.Id + "/wopi/files/" + fileID
	wopiToken := EncodeToken(r.Header.Get("Mattermost-User-Id"), fileID, p)

	response := struct {
		URL         string `json:"url"`
		AccessToken string `json:"access_token"` //client will pass this token as a POST parameter to Collabora Online when loading the iframe
	}{wopiURL, wopiToken}

	responseJSON, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}

func (p *Plugin) parseWopiRequests(w http.ResponseWriter, r *http.Request) {
	splitURL := strings.Split(r.URL.Path, "/")
	fileID := splitURL[len(splitURL)-2] //the segment before last segment is the file url

	token, tokenErr := getAccessTokenFromURI(r.RequestURI)
	if tokenErr != nil {
		p.API.LogError("Error retrieving token from URI: "+r.RequestURI, "Error", tokenErr.Error())
		return
	}

	wopiToken, isValid := DecodeToken(token, p)
	if !isValid || wopiToken.FileID != fileID {
		p.API.LogError("Invalid token.")
		return
	}

	fileContent, err := p.API.GetFile(fileID)
	if err != nil {
		p.API.LogError("Error retrieving file info, fileID: " + fileID)
		return
	}

	fileInfo, fileInfoError := p.API.GetFileInfo(fileID)
	if fileInfoError != nil {
		p.API.LogError("Error occurred when retrieving file info: " + fileInfoError.Error())
		return
	}

	postInfo, postInfoError := p.API.GetPost(fileInfo.PostId)
	if postInfoError != nil {
		p.API.LogError("Error occurred when retrieving post info for file: " + postInfoError.Error())
		return
	}

	//check if user has access to the channel where the file was sent
	//p.API.HasPermissionToChannel(userID,channelID) was returning false for some reason...
	members, channelMembersError := p.API.GetChannelMembersByIds(postInfo.ChannelId, []string{wopiToken.UserID})
	if channelMembersError != nil {
		p.API.LogError("Error occurred when retrieving channel members: " + channelMembersError.Error())
	}
	if members == nil {
		p.API.LogError("User doesn't have access to the channel where the file was sent")
		return
	}

	//send file to Collabora Online
	if r.Method == http.MethodGet {
		if _, err := w.Write(fileContent); err != nil {
			p.API.LogError("failed to write status", "err", err.Error())
		}
	}

	//save file received from Collabora Online
	if r.Method == http.MethodPost {
		f, fileCreateError := os.Create("./data/" + fileInfo.Path)
		if fileCreateError != nil {
			p.API.LogError("Error occurred when creating new file: ", fileCreateError.Error())
			return
		}

		body, bodyReadError := ioutil.ReadAll(r.Body)
		if bodyReadError != nil {
			p.API.LogError("Error occurred when reading body:", bodyReadError.Error())
			return
		}

		_, fileSaveError := f.Write(body)
		if fileSaveError != nil {
			p.API.LogError("Error occurred when writing contents to file: " + fileSaveError.Error())
			f.Close()
			return
		}

		fileCloseError := f.Close()
		if fileCloseError != nil {
			p.API.LogError("Error occurred when closing the file: " + fileCloseError.Error())
			return
		}
	}
}

func (p *Plugin) returnFileInfoForWOPI(w http.ResponseWriter, r *http.Request) {
	splitURL := strings.Split(r.URL.Path, "/")
	fileID := splitURL[len(splitURL)-1]

	token, tokenErr := getAccessTokenFromURI(r.RequestURI)
	if tokenErr != nil {
		p.API.LogError("Error retrieving token from URI:"+r.RequestURI, "Error", tokenErr.Error())
		return
	}

	wopiToken, isValid := DecodeToken(token, p)
	if !isValid || wopiToken.FileID != fileID {
		p.API.LogError("Collabora Online called the plugin with an invalid token.")
		return
	}

	user, userErr := p.API.GetUser(wopiToken.UserID)
	if userErr != nil {
		p.API.LogError("Error retrieving user. Token UserID is corrupted or the user doesn't exist.")
		return
	}

	fileInfo, err := p.API.GetFileInfo(fileID)
	if err != nil {
		p.API.LogError("Error retrieving file info, fileId: " + fileID)
		return
	}

	post, postErr := p.API.GetPost(fileInfo.PostId)
	if postErr != nil {
		p.API.LogError("Error retrieving file's post, postId: " + fileInfo.PostId)
		return
	}

	wopiFileInfo := struct {
		BaseFileName            string `json:"BaseFileName"`
		Size                    int64  `json:"Size"`
		OwnerID                 string `json:"OwnerId"`
		UserID                  string `json:"UserId"`
		UserFriendlyName        string `json:"UserFriendlyName"`
		UserCanWrite            bool   `json:"UserCanWrite"`
		UserCanNotWriteRelative bool   `json:"UserCanNotWriteRelative"`
	}{
		BaseFileName:            fileInfo.Name,
		Size:                    fileInfo.Size,
		OwnerID:                 post.UserId,
		UserID:                  user.Id,
		UserFriendlyName:        user.GetFullName(),
		UserCanWrite:            true,
		UserCanNotWriteRelative: true,
	}

	responseJSON, _ := json.Marshal(wopiFileInfo)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}

func (p *Plugin) serveAsset(w http.ResponseWriter, r *http.Request) {
	splitURL := strings.Split(r.URL.Path, "/")
	fileName := splitURL[len(splitURL)-1] //last segment is the file name

	bundlePath, bundlePathError := p.API.GetBundlePath()
	if bundlePathError != nil {
		p.API.LogError("Error when getting bundle path: " + bundlePathError.Error())
		return
	}

	asset, assetError := ioutil.ReadFile(filepath.Join(bundlePath, "assets", fileName))
	if assetError != nil {
		p.API.LogError("Error when loading assets: " + assetError.Error())
		return
	}

	if _, err := w.Write(asset); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}
