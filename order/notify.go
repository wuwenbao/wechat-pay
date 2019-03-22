package order

import (
	"encoding/xml"
	"errors"
	"github.com/wuwenbao/wechat-go/pay"
	"io"
	"io/ioutil"
	"net/http"
)

//NotifyResponse
type NotifyResponse struct {
	pay.ResponseError
	Appid              string `xml:"appid"`
	MchId              string `xml:"mch_id"`
	DeviceInfo         string `xml:"device_info"`
	NonceStr           string `xml:"nonce_str"`
	Sign               string `xml:"sign"`
	SignType           string `xml:"sign_type"`
	Openid             string `xml:"openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	BankType           string `xml:"bank_type"`
	TotalFee           string `xml:"total_fee"`
	SettlementTotalFee string `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            string `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	CouponFee          string `xml:"coupon_fee"`
	CouponCount        string `xml:"coupon_count"`
	TransactionId      string `xml:"transaction_id"`
	OutTradeNo         string `xml:"out_trade_no"`
	Attach             string `xml:"attach"`
	TimeEnd            string `xml:"time_end"`
}

//notifySign
func notifySign(r io.Reader) (*NotifyResponse, error) {
	bts, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	notify := new(NotifyResponse)
	err = xml.Unmarshal(bts, notify)
	if err != nil {
		return nil, err
	}
	if notify.ReturnCode == "FAIL" {
		return nil, errors.New(notify.ReturnMsg)
	}
	if notify.ResultCode == "FAIL" {
		return nil, errors.New(notify.ErrCode)
	}
	if notify.Sign != pay.SignCheck(notify) {
		return nil, errors.New("sign error")
	}
	return notify, nil
}

//Notify
func Notify(w http.ResponseWriter, r *http.Request, notifyHandle func(notify *NotifyResponse) error) {
	notify, err := notifySign(r.Body)
	if err != nil {
		w.Write([]byte(pay.Fail(err)))
		return
	}
	err = notifyHandle(notify)
	if err != nil {
		w.Write([]byte(pay.Fail(err)))
		return
	}
	w.Write([]byte(pay.Success()))
}
