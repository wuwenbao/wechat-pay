package order

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/wuwenbao/wechat-go/pay"
	"io/ioutil"
	"net/http"
)

type unifyResponse struct {
	pay.ResponseError
	Appid      string `xml:"appid"`       //公众账号ID
	MchId      string `xml:"mch_id"`      //商户号
	DeviceInfo string `xml:"device_info"` //设备号
	NonceStr   string `xml:"nonce_str"`   //随机字符串
	Sign       string `xml:"sign"`        //签名`
	Openid     string `xml:"openid"`      //用户标识
	PrepayId   string `xml:"prepay_id"`   //预支付交易会话标识
	TradeType  string `xml:"trade_type"`  //交易类型
	CodeUrl    string `xml:"code_url"`    //二维码链接
}

type unifyParam struct {
	Appid          string    `xml:"appid"`
	MchId          string    `xml:"mch_id"`
	DeviceInfo     string    `xml:"device_info,omitempty"`
	NonceStr       string    `xml:"nonce_str"`
	Sign           string    `xml:"sign"`
	SignType       string    `xml:"sign_type,omitempty"`
	Body           string    `xml:"body"`
	Detail         pay.CDATA `xml:"detail,omitempty"`
	Attach         string    `xml:"attach,omitempty"`
	OutTradeNo     string    `xml:"out_trade_no"`
	FeeType        string    `xml:"fee_type,omitempty"`
	TotalFee       int       `xml:"total_fee"`
	SpbillCreateIp string    `xml:"spbill_create_ip"`
	TimeStart      string    `xml:"time_start,omitempty"`
	TimeExpire     string    `xml:"time_expire,omitempty"`
	GoodsTag       string    `xml:"goods_tag,omitempty"`
	NotifyUrl      string    `xml:"notify_url"`
	TradeType      string    `xml:"trade_type"`
	ProductId      string    `xml:"product_id,omitempty"`
	LimitPay       string    `xml:"limit_pay,omitempty"`
	Openid         string    `xml:"openid,omitempty"`
	SceneInfo      pay.CDATA `xml:"scene_info,omitempty"`
}

func (u *unifyParam) SetDetail(str string) *unifyParam {
	u.Detail = pay.CDATA{Text: str}
	return u
}

func (u *unifyParam) SetSceneInfo(str string) *unifyParam {
	u.SceneInfo = pay.CDATA{Text: str}
	return u
}

//Unify 统一下单
func (u *unifyParam) Unify() (*unifyResponse, error) {
	//数据签名
	u.NonceStr = pay.RandomStr(10)
	u.Sign = pay.SignCheck(u)

	bts, err := xml.Marshal(struct {
		XMLName xml.Name `xml:"xml"`
		*unifyParam
	}{unifyParam: u})
	if err != nil {
		return nil, err
	}
	fmt.Println(string(bts))

	resp, err := http.Post(
		`https://api.mch.weixin.qq.com/pay/unifiedorder`,
		"", bytes.NewBuffer(bts),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := new(unifyResponse)
	err = xml.Unmarshal(all, response)
	if err != nil {
		return nil, err
	}
	if response.ReturnCode == "FAIL" {
		return nil, errors.New(response.ReturnMsg)
	}
	if response.ResultCode == "FAIL" {
		return nil, errors.New(response.ErrCode)
	}
	return response, nil
}

//NewJSAPI 微信网页端
func NewJSAPI() (jsApi *unifyParam) {
	order := NewUnify()
	order.TradeType = "JSAPI"
	return order
}

func NewUnify() *unifyParam {
	order := &unifyParam{}
	order.Appid = pay.GetConf().Appid
	order.MchId = pay.GetConf().MchId
	order.NotifyUrl = pay.GetConf().NotifyUrl
	return order
}
