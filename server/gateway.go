package server

import (
	"asura-gateway/log"
	"asura-gateway/request"
	"asura-gateway/response"
	"asura-gateway/util"
	"context"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func GateWayHandler(c *gin.Context) {

	req := &request.GwRequest{}
	err := req.ParseRequest(c.Request)
	if err != nil{
		c.JSON(http.StatusBadRequest,err.Error())
		return
	}

	result, err := GateWayProcess(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func GateWayProcess(req *request.GwRequest)(*response.GwResponse,error){
	resp,err := ParallelRequest(req)
	return resp,err
}


func ParallelRequest(allReq *request.GwRequest) (*response.GwResponse,error) {
	var reqCnt = len(allReq.RequestArray)
	RespChannel := make(chan *response.SingleResponse,5)

	resp := &response.GwResponse{}

	for _, r := range allReq.RequestArray {
		go func(req request.SingleReq,respChannel chan *response.SingleResponse) {
			var result = make(chan *response.SingleResponse)
			ctx,cancel := context.WithTimeout(context.Background(),1*time.Second)
			defer cancel()

			var s response.SingleResponse
			if req.Post != "" {
				go DoQuery(ctx,"POST", req, allReq.Header,result)
			} else if req.Post == "" {
				go DoQuery(ctx,"GET", req, allReq.Header,result)
			}

			select {
			case <- ctx.Done():
				s.Code = http.StatusInternalServerError
				s.Msg = "request time over limit"
				s.Url = req.Url
				log.Info("request time over limit,method:%v,request:%v",req.Url,req.Header)
				respChannel<-&s
				return
				case data :=<-result:
					respChannel<-data
					return
			}
		}(r,RespChannel)
	}


	for i:=0;i<reqCnt;i++{
		select {
		case data := <-RespChannel:
			resp.Data = append(resp.Data,data)
		}
	}
	return resp,nil
}

func DoQuery(ctx context.Context,method string, r request.SingleReq, p http.Header,respChannel chan  *response.SingleResponse) {
	currResp := &response.SingleResponse{Code: 500, Url: r.Url}
	if err := util.CheckProtocol(r.Url); err != nil {
		currResp.Msg = err.Error()
		respChannel<-currResp
		return
	}
	request, err := http.NewRequest(method, r.Url, nil)
	if err != nil {
		currResp.Msg = err.Error()
		respChannel<-currResp
		return
	}


	FillHeader(r, p, request)

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		currResp.Msg = err.Error()
		respChannel<-currResp
		return
		}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		currResp.Code = 500
		currResp.Msg = err.Error()
		respChannel<-currResp
		return
		}
	currResp.Code = resp.StatusCode
	currResp.Msg = "ok"
	currResp.Data = string(body)
	respChannel<- currResp
}

func FillHeader(r request.SingleReq, p http.Header, req *http.Request) {
	for key,value := range r.Header {
		req.Header[key] = value
	}

	for key,value := range p {
		if _,ok := req.Header[key];!ok {
			valStr := ""
			for _,v := range value{
				valStr = valStr + v + ";"
			}
			req.Header.Set(key,valStr)
		}
	}
}
