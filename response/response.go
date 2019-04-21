package response

import (
	"encoding/json"
	"strconv"
)

type GwResponse struct {
	Code   int              `json:"code"`
	Sysmsg string           `json:"sysmsg"`
	Data   []*SingleResponse `json:"data"`
}

type SingleResponse struct {
	Code int
	Msg  string
	Url  string
	Data string
}

func (gw *GwResponse) FillResponse() []byte {
	var reStr = ""
	gw.Code = 200
	gw.Sysmsg = "ok"
	r, err := json.Marshal(gw.Data)
	reStr = "{\"code\":" + strconv.Itoa(gw.Code) + ",\"sysmsg\":\"" + gw.Sysmsg + "\",\"data\":" + string(r) + "}"
	if err != nil {
		return []byte(err.Error())
	}
	return []byte(reStr)
}
