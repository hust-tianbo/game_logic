package wxpay

import (
	"context"

	"github.com/hust-tianbo/go_lib/log"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
)

func PreOrder(pay_id string, amount int, openid string) (*jsapi.PrepayWithRequestPaymentResponse, error) {
	client := getWechatPaySvr()

	svc := jsapi.JsapiApiService{Client: client}
	ctx := context.Background()
	resp, _, err := svc.PrepayWithRequestPayment(ctx, jsapi.PrepayRequest{
		Appid:       core.String(appID),
		Mchid:       core.String(mchID),
		Description: core.String("兑换盒子"),
		OutTradeNo:  core.String(pay_id),
		Attach:      core.String("自定义数据说明"),
		NotifyUrl:   core.String(""),
		Amount:      &jsapi.Amount{Total: core.Int64(int64(amount))},
		Payer:       &jsapi.Payer{Openid: core.String(openid)},
	})

	if err == nil {
		log.Debugf("[PreOrder]result is:%+v,%+v,%+v,%+v", openid, pay_id, amount, resp)
		return resp, nil
	} else {
		log.Errorf("[PreOrder]failed:%+v,%+v,%+v,%+v", openid, pay_id, amount, err)
		return nil, err
	}
}
