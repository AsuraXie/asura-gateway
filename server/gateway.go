package server

import (
	"asura-gateway/config"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func StartGateWay() {
	http.HandleFunc("/metrics", MetricsHandler)
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/gw", gwHandler)
	if err := http.ListenAndServe(getListenAddr(), nil); err != nil {
		fmt.Println("start gateway error ", err.Error())
	}
}

func getListenAddr() string {
	port := strconv.Itoa(config.ServerConfig.ListenPort)

	var listenAddr string
	if strings.Contains(config.ServerConfig.ListenAddress, "*") {
		listenAddr = ":" + port
	} else {
		listenAddr = config.ServerConfig.ListenAddress + ":" + port
	}

	fmt.Println("gateway is start ,listen on " + listenAddr)
	return listenAddr
}

func gwHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write([]byte("Method Error"))
		return
	}

	req := GwRequest{}
	if err := req.ParseRequest(*r); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	resp := ParallelRequest(req)
	w.Write(resp.FillResponse())
	req.EndTime()
	MetricRequestPoint(req)
}

func ParallelRequest(allReq GwRequest) GwResponse {
	var wg sync.WaitGroup
	resp := GwResponse{}
	writeLock := sync.Mutex{}
	for _, r := range allReq.RequestArray {
		wg.Add(1)
		go func(r SingleReq) {
			var s = SingleResponse{}
			if r.Post == "" {
				s = DoGet(r)
			} else if r.Post != "" {
				s = DoPost(r)
			}
			writeLock.Lock()
			resp.Data = append(resp.Data, s)
			writeLock.Unlock()
			wg.Done()
		}(r)
	}
	wg.Wait()
	return resp
}

func DoGet(r SingleReq) SingleResponse {
	currResp := SingleResponse{Code: 500, Url: r.Url}
	if err := CheckParam(&r); err != nil {
		currResp.Msg = err.Error()
		return currResp
	}
	request, err := http.NewRequest("GET", r.Url, nil)
	if err != nil {
		currResp.Msg = err.Error()
		return currResp
	}

	ctx, cancel := context.WithCancel(context.Background())
	request.WithContext(ctx)
	timer := time.AfterFunc(2*time.Second, func() {
		fmt.Println("time limit execed,1 second")
		cancel()
	})

	FillHeader(r, request)
	client := http.Client{}
	resp, err := client.Do(request)
	timer.Stop()
	if err != nil {
		currResp.Msg = err.Error()
		return currResp
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(resp)
	fmt.Print("Get" + string(body))
	if err != nil {
		currResp.Code = 500
		currResp.Msg = err.Error()
		return currResp
	}
	currResp.Code = resp.StatusCode
	currResp.Msg = "ok"
	currResp.Data = string(body)
	return currResp
}

func DoPost(r SingleReq) SingleResponse {
	currResp := SingleResponse{Code: 500, Url: r.Url}
	if err := CheckParam(&r); err != nil {
		currResp.Msg = err.Error()
		return currResp
	}

	request, err := http.NewRequest("POST", r.Url, strings.NewReader(r.Post))
	if err != nil {
		currResp.Msg = err.Error()
		return currResp
	}

	ctx, cancel := context.WithCancel(context.Background())
	request.WithContext(ctx)
	timer := time.AfterFunc(2*time.Second, func() {
		fmt.Println("time limit execed,1 second")
		cancel()
	})

	FillHeader(r, request)
	client := http.Client{}
	resp, err := client.Do(request)
	timer.Stop()

	if err != nil {
		currResp.Msg = err.Error()
		return currResp
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		currResp.Msg = err.Error()
		return currResp
	}
	fmt.Println(string(body))
	currResp.Code = resp.StatusCode
	currResp.Msg = "ok"
	currResp.Data = string(body)
	return currResp
}

func FillHeader(r SingleReq, req *http.Request) {
	if r.ContentType != "" {
		req.Header.Set("Content-Type", r.ContentType)
	}

	if r.Host != "" {
		req.Header.Set("Host", r.Host)
	} else {
		req.Header.Set("Host", GetHostFromUrl(r.Url))
	}

	if r.Cookie != "" {
		req.Header.Set("Cookie", r.Cookie)
	}
}

func GetHostFromUrl(url string) string {
	var urlArray = strings.Split(url, "?")
	return urlArray[0]
}

func CheckParam(r *SingleReq) error {
	if !strings.Contains(r.Url, "http") {
		return errors.New("miss protocl http")
	}
	return nil
}
