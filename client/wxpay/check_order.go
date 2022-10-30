package wxpay

import (
	"context"

	"github.com/hust-tianbo/go_lib/log"
	"github.com/wechatpay-apiv3/wecahtpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	mchID                      string = "190000****"                               // 商户号
	mchCertificateSerialNumber string = "3775B6A45ACD588826D15E583A95F5DD********" // 商户证书序列号
	mchAPIv3Key                string = "2ab9****************************"         // 商户APIv3密钥
	appID                      string = "wx*****"                                  // 小程序id
	notifyUrl                  string = ""                                         // 通知的链接
)

func getWechatPaySvr() *core.Client {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("/path/to/merchant/apiclient_key.pem")
	if err != nil {
		log.Fatal("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err:%s", err)
	}

	return client
}

func CheckOrder(transaction_id string) (*payments.Transaction, error) {
	client := getWechatPaySvr()

	svc := jsapi.JsapiApiService{Client: client}
	ctx := context.Background()
	resp, _, err := svc.QueryOrderById(ctx, jsapi.QueryOrderByIdRequest{
		TransactionId: core.String(transaction_id),
		Mchid:         core.String(mchID),
	})

	if err == nil {
		log.Debugf("[CheckOrder]transaction_id is:%+v,%+v", transaction_id, resp)
		return resp, nil
	} else {
		log.Errorf("[CheckOrder]failed :%+v,%+v", transaction_id, err)
		return nil, err
	}
}
