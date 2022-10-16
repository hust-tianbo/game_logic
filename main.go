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
	mux.HandleFunc("/get_box_info", func(w http.ResponseWriter, r *http.Request) {
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
	mux.HandleFunc("/acquire_box", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req logic.AcquireBoxReq
		json.Unmarshal(body, &req)
		var rsp logic.AcquireBoxRsp
		defer func() {
			log.Debugf("[GetHttpServerMux]deal log:%+v,%+v", req, rsp)
		}()

		rsp = logic.AcquireBox(&req)
		resBytes, _ := json.Marshal(rsp)
		w.Write([]byte(resBytes))
	})
	mux.HandleFunc("/acquire_box_check", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req logic.AcquireBoxCheckReq
		json.Unmarshal(body, &req)
		var rsp logic.AcquireBoxCheckRsp
		defer func() {
			log.Debugf("[GetHttpServerMux]deal log:%+v,%+v", req, rsp)
		}()

		rsp = logic.AcquireBoxCheck(&req)
		resBytes, _ := json.Marshal(rsp)
		w.Write([]byte(resBytes))
	})
	mux.HandleFunc("/acquire_prize", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req logic.AcquirePrizeReq
		json.Unmarshal(body, &req)
		var rsp logic.AcquirePrizeRsp
		defer func() {
			log.Debugf("[GetHttpServerMux]deal log:%+v,%+v", req, rsp)
		}()

		rsp = logic.AcquirePrize(&req)
		resBytes, _ := json.Marshal(rsp)
		w.Write([]byte(resBytes))
	})
	mux.HandleFunc("/get_user_box_list", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req logic.GetUserBoxListReq
		json.Unmarshal(body, &req)
		var rsp logic.GetUserBoxListRsp
		defer func() {
			log.Debugf("[GetHttpServerMux]deal log:%+v,%+v", req, rsp)
		}()

		rsp = logic.GetUserBoxList(&req)
		resBytes, _ := json.Marshal(rsp)
		w.Write([]byte(resBytes))
	})
	return mux
}
