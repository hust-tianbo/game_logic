package wxpay

import (
	"context"

	"github.com/hust-tianbo/game_logic/config"
	"github.com/hust-tianbo/go_lib/log"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
)

func PreOrder(pay_id string, amount int, openid string) (*jsapi.PrepayWithRequestPaymentResponse, error) {
	client := getWechatPaySvr()

	cfg := config.GetConfig()
	svc := jsapi.JsapiApiService{Client: client}
	ctx := context.Background()
	req := jsapi.PrepayRequest{
		Appid:       core.String(cfg.AppID),
		Mchid:       core.String(cfg.MchID),
		Description: core.String("兑换盒子"),
		OutTradeNo:  core.String(pay_id),
		Attach:      core.String("自定义数据说明"),
		NotifyUrl:   core.String(cfg.NotifyUrl),
		Amount:      &jsapi.Amount{Total: core.Int64(int64(amount))},
		Payer:       &jsapi.Payer{Openid: core.String(openid)},
	}
	resp, _, err := svc.PrepayWithRequestPayment(ctx, req)

	if err == nil {
		log.Debugf("[PreOrder]result is:%+v,%+v,%+v,%+v,%+v", openid, pay_id, amount, req, resp)
		return resp, nil
	} else {
		log.Errorf("[PreOrder]failed:%+v,%+v,%+v,%+v,%+v", openid, pay_id, amount, req, err)
		return nil, err
	}
}
