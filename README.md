# wechat-pay

### 使用

```go

package main

import (
	"encoding/xml"
	"log"
	"fmt"
	"github.com/wuwenbao/wechat-go/pay/order"
	"net/http"
	"github.com/wuwenbao/wechat-go/pay/refund"
	"github.com/wuwenbao/wechat-go/pay"
)

func main() {
	conf := &pay.Conf{
		Appid:           "",
		MchId:           "",
		MchKey:          "",
		NotifyUrl:       "",
		RefundNotifyUrl: "",
		CertPath:        "",
		KeyPath:         "",
		CaPath:          "",
	}
	pay.SetConf(conf)

	jsApi := order.NewJSAPI()
	jsApi.Body = "body"
	jsApi.TotalFee = 1000000
	jsApi.Openid = "openid"
	jsApi.SpbillCreateIp = "0.0.0.0"
	jsApi.OutTradeNo = "out_trade_no"
	jsApi.SetDetail("商品详情")
	jsApi.SetSceneInfo("场景信息")
	prepay, err := jsApi.Unify()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prepay)


    //PayNotifyHandle 支付通知逻辑
    func PayNotifyHandle(w http.ResponseWriter, r *http.Request) {
        order.Notify(w, r, func(notify *order.NotifyResponse) error {
            fmt.Println(notify)
            return nil
        })
    }

    //RefundNotifyHandle 退款通知逻辑
    func RefundNotifyHandle(w http.ResponseWriter, r *http.Request) {
        refund.Notify(w, r, func(notify *refund.NotifyReqInfo) error {
            fmt.Println(notify)
            return nil
        })
    }
}

```
