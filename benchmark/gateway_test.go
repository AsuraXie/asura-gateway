package benchmark

import (
	"asura-gateway/request"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func Test_gateway(t *testing.T){
	req := []request.SingleReq{}
	req = append(req,BuildPostParam())
	req = append(req,BuildGetParam())
	DoPost(req)
}

func BuildPostParam() request.SingleReq{
	postReq := request.SingleReq{}
	postReq.Url = "http://localhost:8089/postmethod"
	postReq.Post = "wokao"
	postReq.Header = map[string][]string {
		"User-Agent":[]string{"hello world"},
		"Content-Type":[]string{"application/json","utf-8"},
		"Accept":[]string{"*/*"},
		"PostUserId":[]string{"123"},
		"Host":[]string{"localhostpost"},
	}
	return postReq
}

func BuildGetParam() request.SingleReq {
	getReq := request.SingleReq{}
	getReq.Url = "http://localhost:8089/getmethod"
	getReq.Get = "id=1"
	getReq.Header = map[string][]string {
		"Content-Type":[]string{"text/html"},
		"Accept":[]string{"*/*"},
		"GetUserId":[]string{"321"},
		"Host":[]string{"localhosthostget"},
	}
	return getReq
}

func DoPost(reqParam interface{}){
	client := &http.Client{}
	reqParamBytes,err := json.Marshal(reqParam)
	if err != nil{
		fmt.Println("error",err)
		return
	}
	req,err := http.NewRequest("POST","http://localhost:12345/gw",bytes.NewReader(reqParamBytes))

	if err != nil{
		fmt.Println("err")
		os.Exit(-1)
	}
	PostHeader(req)
	resp,err := client.Do(req)
	if err !=nil{
		fmt.Println(err)
		os.Exit(-1)
	}

	defer resp.Body.Close()

	var body []byte
	body,err = ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println("post ",resp.StatusCode,string(body))

	if resp.StatusCode != http.StatusOK {
		fmt.Println("post ",resp.StatusCode,string(body))
	}
}

func PostHeader(req *http.Request) {
	req.Header.Set("Accept","*/*")
	req.Header.Set("User-Agent","Putong/3.4.2.1 Android/26 Xiaomi/MIX+2")
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Host","localhostglobal")
}