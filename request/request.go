package request

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type GwRequest struct {
	startTime    int64
	interval     int64
	RequestArray []SingleReq `json:"request_array"`
	Header		http.Header
}

type SingleReq struct {
	startTime int64
	interval  int64
	ip        string
	body      string

	Url   string `json:"url"`
	Get   string `json:"get"`
	Post  string `json:"post"`
	Body  string `json:"body"`
	Header http.Header `json:"header"`
}

func (myreq *GwRequest) StartTime() {
	myreq.startTime = time.Now().UnixNano()
}

func (myreq *GwRequest) EndTime() {
	myreq.interval = time.Now().UnixNano() - myreq.startTime
}

func (myreq *GwRequest) ParseRequest(r *http.Request) error {
	myreq.StartTime()

	myreq.Header = r.Header

	body := []byte{}
	body,err := ioutil.ReadAll(r.Body)
	if err != nil{
		return errors.New("read request body error " + err.Error())
	}

	if len(body) <= 0 {
		return errors.New("read 0 byte from body")
	}

	if err := json.Unmarshal(body, &myreq.RequestArray); err != nil {
		return errors.New("request param gw is not json data," + err.Error())
	}
	return nil
}
