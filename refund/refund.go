package refund

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/wuwenbao/wechat-go/pay"
	"io/ioutil"
	"net/http"
)

type refundResponse struct {
	pay.ResponseError
}

type refundParam struct {
	Appid         string `xml:"appid"`
	MchId         string `xml:"mch_id"`
	NonceStr      string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	SignType      string `xml:"sign_type,omitempty"`
	TransactionId string `xml:"transaction_id,omitempty"`
	OutTradeNo    string `xml:"out_trade_no,omitempty"`
	OutRefundNo   string `xml:"out_refund_no"`
	TotalFee      int    `xml:"total_fee"`
	RefundFee     int    `xml:"refund_fee"`
	RefundFeeType string `xml:"refund_fee_type,omitempty"`
	RefundDesc    string `xml:"refund_desc,omitempty"`
	RefundAccount string `xml:"refund_account,omitempty"`
	NotifyUrl     string `xml:"notify_url,omitempty"`
}

//Refund 退款操作
func (r *refundParam) Refund() (*refundResponse, error) {
	//数据签名
	r.NonceStr = pay.RandomStr(10)
	r.Sign = pay.SignCheck(r)
	bts, err := xml.Marshal(struct {
		XMLName xml.Name `xml:"xml"`
		*refundParam
	}{refundParam: r})
	if err != nil {
		return nil, err
	}
	fmt.Println(string(bts))
	request, err := http.NewRequest(http.MethodPost, `https://api.mch.weixin.qq.com/secapi/pay/refund`, bytes.NewBuffer(bts))
	if err != nil {
		return nil, err
	}
	resp, err := pay.ClientDo(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := new(refundResponse)
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

//NewRefund 实例化退款
func NewRefund() *refundParam {
	param := &refundParam{
		Appid:     pay.GetConf().Appid,
		MchId:     pay.GetConf().MchId,
		NotifyUrl: pay.GetConf().RefundNotifyUrl,
	}
	return param
}
