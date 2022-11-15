package wxpay

import (
	"context"

	"github.com/hust-tianbo/game_logic/config"
	"github.com/hust-tianbo/go_lib/log"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

func getWechatPaySvr() *core.Client {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	cfg := config.GetConfig()
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(cfg.PrivatePath)
	if err != nil {
		log.Fatal("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(cfg.MchID, cfg.SerialNumber, mchPrivateKey, cfg.ApiV3),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err:%s", err)
	}

	return client
}

func CheckOrder(transaction_id string) (*payments.Transaction, error) {
	client := getWechatPaySvr()

	cfg := config.GetConfig()
	svc := jsapi.JsapiApiService{Client: client}
	ctx := context.Background()
	resp, _, err := svc.QueryOrderById(ctx, jsapi.QueryOrderByIdRequest{
		TransactionId: core.String(transaction_id),
		Mchid:         core.String(cfg.MchID),
	})

	if err == nil {
		log.Debugf("[CheckOrder]transaction_id is:%+v,%+v", transaction_id, resp)
		return resp, nil
	} else {
		log.Errorf("[CheckOrder]failed :%+v,%+v", transaction_id, err)
		return nil, err
	}
}
