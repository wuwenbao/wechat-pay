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

type queryResponse struct {
	pay.ResponseError
	Appid              string `xml:"appid"`
	MchId              string `xml:"mch_id"`
	NonceStr           string `xml:"nonce_str"`
	Sign               string `xml:"sign"`
	DeviceInfo         string `xml:"device_info"`
	Openid             string `xml:"openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	TradeState         string `xml:"trade_state"`
	BankType           string `xml:"bank_type"`
	TotalFee           int    `xml:"total_fee"`
	SettlementTotalFee int    `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            int    `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	CouponFee          int    `xml:"coupon_fee"`
	CouponCount        int    `xml:"coupon_count"`
	TransactionId      string `xml:"transaction_id"`
	OutTradeNo         string `xml:"out_trade_no"`
	Attach             string `xml:"attach"`
	TimeEnd            string `xml:"time_end"`
	TradeStateDesc     string `xml:"trade_state_desc"`
}

type queryParam struct {
	Appid         string `xml:"appid"`
	MchId         string `xml:"mch_id"`
	TransactionId string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	NonceStr      string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	SignType      string `xml:"sign_type"`
}

func (q *queryParam) Query() (*queryResponse, error) {
	//数据签名
	q.NonceStr = pay.RandomStr(10)
	q.Sign = pay.SignCheck(q)

	bts, err := xml.Marshal(struct {
		XMLName xml.Name `xml:"xml"`
		*queryParam
	}{queryParam: q})
	if err != nil {
		return nil, err
	}
	fmt.Println(string(bts))

	resp, err := http.Post(
		`https://api.mch.weixin.qq.com/pay/orderquery`,
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
	response := new(queryResponse)
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

func NewQuery(outTradeNo, transactionId string) *queryParam {
	param := &queryParam{
		Appid:         pay.GetConf().Appid,
		MchId:         pay.GetConf().MchId,
		OutTradeNo:    outTradeNo,
		TransactionId: transactionId,
	}
	return param
}

//QueryByOutTradeNo 商户订单号
func QueryByOutTradeNo(outTradeNo string) (*queryResponse, error) {
	return NewQuery(outTradeNo, "").Query()
}

//QueryByTransactionId 微信订单号
func QueryByTransactionId(transactionId string) (*queryResponse, error) {
	return NewQuery("", transactionId).Query()
}
