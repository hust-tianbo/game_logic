package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hust-tianbo/go_lib/log"
)

type CheckAuthReq struct {
	PersonID      string `json:"personid"`
	Code          string `json:"code"`           // 平台返回的code码
	InternalToken string `json:"internal_token"` // 如果已经有内部票据，则携带
	ReturnWXToken bool   `json:"return_wxtoken"`
}

type CheckAuthRsp struct {
	Ret           int    `json:"ret"`            // 错误码
	Msg           string `json:"msg"`            // 错误信息
	InternalToken string `json:"internal_token"` // 内部票据
	PersonID      string `json:"personid"`       // 内部id
	WXToken       string `json:"wxtoken"`
}

const CheckAuthUrl string = "http://127.0.0.1:50052/check_auth"

func CheckToken(pid string, token string) bool {
	req := CheckAuthReq{
		PersonID:      pid,
		InternalToken: token,
		ReturnWXToken: true,
	}

	bytesData, _ := json.Marshal(req)

	resp, err := http.Post(CheckAuthUrl, "application/json;charset=utf-8", bytes.NewBuffer([]byte(bytesData)))
	if err != nil {
		log.Errorf("[CheckToken]req failed:%+|%+v", CheckAuthUrl, err)
		return false
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[CheckToken]read content failed:%+v|%+v", pid, err)
		return false
	}

	var res CheckAuthRsp
	json.Unmarshal(body, &res)

	return res.Ret == RetSuccess
}

func CheckTokenDirect(code string) (*CheckAuthRsp, bool) {
	req := CheckAuthReq{
		Code:          code,
		ReturnWXToken: true,
	}

	bytesData, _ := json.Marshal(req)

	resp, err := http.Post(CheckAuthUrl, "application/json;charset=utf-8", bytes.NewBuffer([]byte(bytesData)))
	if err != nil {
		log.Errorf("[CheckToken]req failed:%+|%+v", CheckAuthUrl, err)
		return nil, false
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[CheckToken]read content failed:%+v|%+v", code, err)
		return nil, false
	}

	var res CheckAuthRsp
	json.Unmarshal(body, &res)

	return &res, res.Ret == RetSuccess
}
