package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hust-tianbo/game_logic/config"

	"github.com/hust-tianbo/game_logic/internal/logic"

	"github.com/hust-tianbo/go_lib/log"
)

const (
	port = ":50053"
)

func main() {
	//modbus.get()
	log.Debugf("begin logic server")
	config.InitConfig()

	// 注册http接口
	mux := GetHttpServerMux()
	http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}),
	)

	log.Debugf("end logic server")
}

func GetHttpServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req logic.GetBoxInfoReq
		json.Unmarshal(body, &req)
		var rsp logic.GetBoxInfoRsp
		defer func() {
			log.Debugf("[GetHttpServerMux]deal log:%+v,%+v", req, rsp)
		}()

		rsp = logic.GetBoxInfo(req)
		resBytes, _ := json.Marshal(rsp)
		w.Write([]byte(resBytes))
	})
	return mux
}
