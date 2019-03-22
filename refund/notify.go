package refund

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"github.com/wuwenbao/wechat-go/pay"
	"io"
	"io/ioutil"
	"net/http"
)

type notifyResponse struct {
	pay.ReturnError
	Appid    string `xml:"appid"`
	MchId    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	ReqInfo  string `xml:"req_info"`
}

type NotifyReqInfo struct {
	TransactionId       string `xml:"transaction_id"`
	OutTradeNo          string `xml:"out_trade_no"`
	RefundId            string `xml:"refund_id"`
	OutRefundNo         string `xml:"out_refund_no"`
	TotalFee            int    `xml:"total_fee"`
	SettlementTotalFee  int    `xml:"settlement_total_fee"`
	RefundFee           int    `xml:"refund_fee"`
	SettlementRefundFee int    `xml:"settlement_refund_fee"`
	RefundStatus        string `xml:"refund_status"`
	SuccessTime         string `xml:"success_time"`
	RefundRecvAccout    string `xml:"refund_recv_accout"`
	RefundAccount       string `xml:"refund_account"`
	RefundRequestSource string `xml:"refund_request_source"`
}

func notifySign(r io.Reader) (*NotifyReqInfo, error) {
	bts, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	notify := new(notifyResponse)
	err = xml.Unmarshal(bts, notify)
	if err != nil {
		return nil, err
	}
	if notify.ReturnCode == "FAIL" {
		return nil, errors.New(notify.ReturnMsg)
	}
	src, err := base64.StdEncoding.DecodeString(notify.ReqInfo)
	if err != nil {
		return nil, err
	}
	key := pay.Md5([]byte(pay.GetConf().MchKey))
	dst, err := pay.AesEBCDecrypt(src, []byte(key))
	if err != nil {
		return nil, err
	}
	reqInfo := new(NotifyReqInfo)
	err = xml.Unmarshal(dst, reqInfo)
	if err != nil {
		return nil, err
	}
	return reqInfo, nil
}

//Notify
func Notify(w http.ResponseWriter, r *http.Request, notifyHandle func(notify *NotifyReqInfo) error) {
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
