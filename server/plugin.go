package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	//send file info to client (name, extension and id) for each file
	//body contains an array with file ids in JSON format
	if r.URL.Path == "/fileInfo" {
		p.parseFileIds(w, r)
	}

	//send url and token with witch client will load Collabora in iframe
	if r.URL.Path == "/collaboraURL" {
		p.collaboraURL(w, r)
	}

	//https://<WOPI host URL>/<...>/wopi/files/<id>/contents gets/saves a file, used by Collabora Online
	if strings.Contains(r.URL.Path, "/wopi/files/") && strings.Contains(r.URL.Path, "contents") {
		p.parseWopiRequests(w, r)
	}

	//https://<WOPI host URL>/<...>/wopi/files/<id> returns file info, used by Collabora Online
	if strings.Contains(r.URL.Path, "/wopi/files/") && !strings.Contains(r.URL.Path, "contents") {
		p.returnFileInfoForWOPI(w, r)
	}
}

func (p *Plugin) returnFileInfoForWOPI(w http.ResponseWriter, r *http.Request) {

	splittedURL := strings.Split(r.URL.Path, "/")
	fileID := splittedURL[len(splittedURL)-1]

	wopiToken, isValid := DecodeToken(getAccessToken(r.RequestURI))
	if !isValid || wopiToken.FileID != fileID {
		fmt.Println("Invalid token")
		return
	}

	user, userErr := p.API.GetUser(wopiToken.UserID)
	if userErr != nil {
		fmt.Println("Error retrieving user")
		return
	}

	fileInfo, err := p.API.GetFileInfo(fileID)
	if err != nil {
		fmt.Println("Error retrieving file info " + fileID)
		return
	}

	post, postErr := p.API.GetPost(fileInfo.PostId)
	if postErr != nil {
		fmt.Println("Could not retrive file's post")
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
		//PostMessageOrigin       string `json:"PostMessageOrigin"`
	}{
		BaseFileName:            fileInfo.Name,
		Size:                    fileInfo.Size,
		OwnerID:                 post.UserId,
		UserID:                  user.Id,
		UserFriendlyName:        user.Nickname,
		UserCanWrite:            true,
		UserCanNotWriteRelative: true,
		//PostMessageOrigin:       "http://localhost:8065",
	}

	responseJSON, _ := json.Marshal(wopiFileInfo)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}

func (p *Plugin) collaboraURL(w http.ResponseWriter, r *http.Request) {
	fileIDs, ok := r.URL.Query()["file_id"]
	fileID := fileIDs[0]

	if !ok {
		p.API.LogError("file_id query parameter missing!")
		return
	}

	url := "http://localhost:9980/loleaflet/1e4154c/loleaflet.html?WOPISrc=http://localhost:8065/plugins/" + manifest.ID + "/wopi/files/" + fileID
	token := EncodeToken(r.Header.Get("Mattermost-User-Id"), fileID)

	response := struct {
		URL         string `json:"url"`
		AccessToken string `json:"access_token"` //will be used when submitting the form on client to load Collabora in iFrame
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

	wopiToken, isValid := DecodeToken(getAccessToken(r.RequestURI))
	if !isValid || wopiToken.FileID != fileID {
		fmt.Println("Invalid token")
		return
	}

	file, err := p.API.GetFile(fileID)
	if err != nil {
		fmt.Println("Error retrieving file info " + fileID)
	}

	//collabora online asks for the file
	if r.Method == http.MethodGet {
		if _, err := w.Write(file); err != nil {
			p.API.LogError("failed to write status", "err", err.Error())
		}
	}

	//collabora online wants to save the file
	if r.Method == http.MethodPost {
		fmt.Println(fileID)
		fileInfo, err1 := p.API.GetFileInfo(fileID)

		if err1 != nil {
			fmt.Println("Error retrieving file info on save")
			return
		}

		fmt.Println(fileInfo.Path)
		f, err := os.Create("./data/" + fileInfo.Path)
		if err != nil {
			fmt.Println(err)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error reading body: %v", err)
			return
		}

		l, err := f.Write(body)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		fmt.Println(l, " bytes written successfully!")
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (p *Plugin) parseFileIds(w http.ResponseWriter, r *http.Request) {

	//extract fileIds array from body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	var fileIds []string
	_ = json.Unmarshal(body, &fileIds)

	//create an array with more detailed file info
	files := make([]FileInfo, len(fileIds))
	for i := 0; i < len(fileIds); i++ {
		fileInfo, err := p.API.GetFileInfo(fileIds[i])
		if err != nil {
			p.API.LogError("failed to write status", "err", err.Error())
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

//Because the access_token get's removed from Query parameters by mattermost before reaching plugin server,
//it should be manually extracted from the URI
func getAccessToken(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	m, _ := url.ParseQuery(u.RawQuery)
	return m["access_token"][0]
}
