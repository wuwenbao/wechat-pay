package pay

import "fmt"

type ReturnError struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息
}

type ResultError struct {
	ResultCode string `xml:"result_code"`  //业务结果
	ErrCode    string `xml:"err_code"`     //错误代码
	ErrCodeDes string `xml:"err_code_des"` //错误代码描述
}

type ResponseError struct {
	ReturnError
	ResultError
}

//Fail 通知失败
func Fail(msg error) string {
	f := `<xml><return_code><![CDATA[FAIL]]></return_code>return_msg><![CDATA[%s]]></return_msg></xml>`
	return fmt.Sprintf(f, msg)
}

//Success 通知成功
func Success() string {
	msg := `<xml><return_code><![CDATA[SUCCESS]]></return_code>return_msg><![CDATA[OK]]></return_msg></xml>`
	return msg
}