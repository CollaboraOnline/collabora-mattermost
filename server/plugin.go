package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/plugin"
)

//FileInfo contains file informaton sent to the client
type FileInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Action    string `json:"action"` //view or edit
}

//Plugin required by plugin
type Plugin struct {
	plugin.MattermostPlugin
	configurationLock sync.RWMutex
	configuration     *configuration
}

//OnActivate is called when the plugin is activated
func (p *Plugin) OnActivate() error {
	GenerateEncryptionPassword(p)
	return nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	//send file info to client (name, extension and id) for each file
	//body contains an array with file ids in JSON format
	if r.URL.Path == "/fileInfo" {
		p.parseFileIds(w, r)
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

func (p *Plugin) returnFileInfoForWOPI(w http.ResponseWriter, r *http.Request) {

	splittedURL := strings.Split(r.URL.Path, "/")
	fileID := splittedURL[len(splittedURL)-1]

	wopiToken, isValid := DecodeToken(getAccessTokenFromURI(r.RequestURI), p)
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
		UserFriendlyName:        user.Nickname,
		UserCanWrite:            true,
		UserCanNotWriteRelative: true,
	}

	responseJSON, _ := json.Marshal(wopiFileInfo)

	w.Header().Set("Content-Type", "application/json")
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

	//MATTERMOST_SERVER_URL will be replaced on the client side by Javascript
	//could not find how to get it here :)
	url := WOPIFiles[strings.ToLower(file.Extension)].Url + "WOPISrc=MATTERMOST_SERVER_URL/plugins/" + manifest.ID + "/wopi/files/" + fileID
	token := EncodeToken(r.Header.Get("Mattermost-User-Id"), fileID, p)

	response := struct {
		URL         string `json:"url"`
		AccessToken string `json:"access_token"` //client will pass this token as a POST parameter to Collabora Online when loading the iframe
	}{url, token}

	responseJSON, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}

func (p *Plugin) parseWopiRequests(w http.ResponseWriter, r *http.Request) {

	splittedURL := strings.Split(r.URL.Path, "/")
	fileID := splittedURL[len(splittedURL)-2] //the segment before last segment is the file url

	wopiToken, isValid := DecodeToken(getAccessTokenFromURI(r.RequestURI), p)
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
		p.API.LogError("Error occured when retrieving file info: " + fileInfoError.Error())
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
		if err != nil {
			p.API.LogError("Error occured when creating new file: ", fileCreateError.Error())
			return
		}

		body, bodyReadError := ioutil.ReadAll(r.Body)
		if bodyReadError != nil {
			p.API.LogError("Error occured when reading body:", bodyReadError.Error())
			return
		}

		_, fileSaveError := f.Write(body)
		if err != nil {
			p.API.LogError("Error occured when writing contents to file: " + fileSaveError.Error())
			f.Close()
			return
		}

		fileCloseError := f.Close()
		if err != nil {
			p.API.LogError("Error occured when closing the file: " + fileCloseError.Error())
			return
		}
	}
}

func (p *Plugin) parseFileIds(w http.ResponseWriter, r *http.Request) {

	//extract fileIds array from body
	body, bodyReadError := ioutil.ReadAll(r.Body)
	if bodyReadError != nil {
		p.API.LogError("Error when reading body: ", bodyReadError.Error())
		return
	}
	var fileIds []string
	_ = json.Unmarshal(body, &fileIds)

	//create an array with more detailed file info for each file
	files := make([]FileInfo, len(fileIds))
	for i := 0; i < len(fileIds); i++ {
		fileInfo, fileInfoError := p.API.GetFileInfo(fileIds[i])
		if fileInfoError != nil {
			p.API.LogError("Error when retrieving file info: ", fileInfoError.Error())
		}
		if value, ok := WOPIFiles[strings.ToLower(fileInfo.Extension)]; ok {
			files[i] = FileInfo{fileInfo.Id, fileInfo.Name, fileInfo.Extension, value.Action}
		}
	}

	responseJSON, _ := json.Marshal(files)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}

func (p *Plugin) serveAsset(w http.ResponseWriter, r *http.Request) {

	splittedURL := strings.Split(r.URL.Path, "/")
	fileName := splittedURL[len(splittedURL)-1] //last segment is the file name

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

//Because the access_token get's removed from Query parameters by Mattermost before
//it reaches the plugin HTTP request parser, it should be manually extracted from the URI
func getAccessTokenFromURI(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	m, _ := url.ParseQuery(u.RawQuery)
	return m["access_token"][0]
}
