package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type GwRequest struct {
	startTime    int64
	interval     int64
	RequestArray []SingleReq `json:"request_array"`
}

type SingleReq struct {
	startTime int64
	interval  int64
	ip        string
	body      string

	Header      string `json:"header"`
	Url         string `json:"url"`
	Host        string `json:"host"`
	Get         string `json:"get"`
	Post        string `json:"post"`
	ContentType string `json:"content_type"`
	FileName    string `json:"file_name"`
	Cookie      string `json:"cookie"`
}

func (myreq *GwRequest) StartTime() {
	myreq.startTime = time.Now().UnixNano()
}

func (myreq *GwRequest) EndTime() {
	myreq.interval = time.Now().UnixNano() - myreq.startTime
}

func (myreq *GwRequest) ParseRequest(r http.Request) error {
	myreq.StartTime()
	postData := r.PostFormValue("gw")
	if err := json.Unmarshal([]byte(postData), &myreq.RequestArray); err != nil {
		return errors.New("request param gw is not json data," + err.Error())
	}
	return nil
}
