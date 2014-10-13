package snapshot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SnapshotRequest struct {
	uri          string
	requestPath  string
	method       string
	pathSettings map[string]string
}

type listSnapshotsJSON struct {
	Snapshots []struct {
		Snapshot          string   `json:"snapshot"`
		Indices           []string `json:"indices"`
		State             []string `json:"state"`
		StartTime         string   `json:"start_time"`
		StartTimeInMillis int      `json:"start_time_in_millis"`
		EndTime           string   `json:"end_time"`
		EndTimeInMillis   int      `json:"end_time_in_millis"`
		DurationInMillis  int      `json:"DurationInMillis"`
		Failurs           []string `json:"failures"`
		Shards            struct {
			Total      int `json:"total"`
			Failed     int `json:"failed"`
			Successful int `json:"successful"`
		} `json: "shards"`
	} `json:"snapshots"`
}

var CreateSnapshotRequest SnapshotRequest = SnapshotRequest{
	"localhost:9200",
	"_snapshot/{{repo_name}}/{{snapshot_name}}",
	"PUT",
	map[string]string{},
}

var ListSnapshotsRequest SnapshotRequest = SnapshotRequest{
	"localhost:9200",
	"_snapshot/{{repo_name}}/_all",
	"GET",
	map[string]string{},
}

func (r *SnapshotRequest) setPath() {
	path := r.requestPath
	for name, value := range r.pathSettings {
		nameMark := fmt.Sprintf("{{%s}}", name)
		path = strings.Replace(path, nameMark, value, 1)
	}
	r.requestPath = path
}

func (r *SnapshotRequest) perform() (*http.Response, error) {
	r.setPath()
	client := &http.Client{}
	requestURL := fmt.Sprintf("%s/%s", r.uri, r.requestPath)
	req, err := http.NewRequest(r.method, requestURL, nil)
	if err != nil {
		return nil, err
	}
	response, connectionErr := client.Do(req)
	if connectionErr != nil {
		return nil, connectionErr
	}
	return response, nil
}

func createSnapshot(url, repoName, snapName string) {
	request := CreateSnapshotRequest
	request.uri = url
	request.pathSettings["repo_name"] = repoName
	request.pathSettings["snapshot_name"] = snapName
	request.perform()
}

func listSnapshots(url, repoName string) listSnapshotsJSON {
	request := ListSnapshotsRequest
	request.uri = url
	request.pathSettings["repo_name"] = repoName
	response, _ := request.perform()
	var js listSnapshotsJSON
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &js)
	return js
}
