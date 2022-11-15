package logic

import "github.com/hust-tianbo/go_lib/log"

type ResourceData struct {
	OriginalType   string `json:"original_type"`
	Algorithm      string `json:"algorithm"`
	Ciphertext     string `json:"ciphertext"`
	AssociatedData string `json:"associated_data"`
	Nonce          string `json:"nonce"`
}

type NotifyReq struct {
	ID           string       `json:"id"`
	CreateTime   string       `json:"create_time"`
	ResourceType string       `json:"resource_type"`
	EventType    string       `json:"event_type"`
	Summary      string       `json:"summary"`
	Resource     ResourceData `json:"reource"`
}

type NotifyRsp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Notify(req *NotifyReq) *NotifyRsp {
	rsp := &NotifyRsp{Code: "SUCCESS"}
	log.Debugf("[Notify]req:%+v", req)
	return rsp
}
