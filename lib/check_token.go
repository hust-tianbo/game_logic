package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hust-tianbo/go_lib/log"
)

type CheckAuthRsp struct {
	Ret           int    `json:"ret"`            // 错误码
	Msg           string `json:"msg"`            // 错误信息
	InternalToken string `json:"internal_token"` // 内部票据
	PersonID      string `json:"personid"`       // 内部id
}

func CheckToken(pid string, token string) bool {
	url := fmt.Sprintf("https://127.0.0.1:50052/check_auth?personid=%+v&internal_token=%+v", pid, token)

	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("[CheckToken]req failed:%+|%+v", url, err)
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
