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

type closeResponse struct {
	pay.ResponseError
	Appid      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	OutTradeNo string `xml:"out_trade_no"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
}

type closeParam struct {
	Appid      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	OutTradeNo string `xml:"out_trade_no"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	SignType   string `xml:"sign_type"`
}

func (c *closeParam) Close() (*closeResponse, error) {
	//数据签名
	c.NonceStr = pay.RandomStr(10)
	c.Sign = pay.SignCheck(c)

	bts, err := xml.Marshal(struct {
		XMLName xml.Name `xml:"xml"`
		*closeParam
	}{closeParam: c})
	if err != nil {
		return nil, err
	}
	fmt.Println(string(bts))

	resp, err := http.Post(
		`https://api.mch.weixin.qq.com/pay/closeorder`,
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
	response := new(closeResponse)
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

func NewClose(outTradeNo string) *closeParam {
	param := &closeParam{
		Appid:      pay.GetConf().Appid,
		MchId:      pay.GetConf().MchId,
		OutTradeNo: outTradeNo,
	}
	return param
}

func Close(outTradeNo string) (*closeResponse, error) {
	return NewClose(outTradeNo).Close()
}
