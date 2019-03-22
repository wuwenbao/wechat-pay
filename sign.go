package pay

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"time"
)

type CDATA struct {
	Text string `xml:",cdata"`
}

//SignCheck 数据有效检查
func SignCheck(signType interface{}) string {
	valueOf := reflect.ValueOf(signType)
	typeOf := reflect.TypeOf(signType)

	switch typeOf.Kind() {
	case reflect.Ptr:
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return ""
	}

	storeKv := make(map[string]string)
	var sortK []string
	for i := 0; i < typeOf.NumField(); i++ {
		tag, ok := typeOf.Field(i).Tag.Lookup("xml")
		if !ok || tag == "" {
			continue
		}
		tags := strings.Split(tag, ",")
		if len(tags) < 0 {
			continue
		}
		s := valueOf.Field(i).String()
		if s == "" || s == "0" {
			continue
		}
		key := strings.TrimSpace(tags[0])

		if key == "sign" {
			continue
		}

		switch valueOf.Field(i).Interface().(type) {
		case CDATA:
			storeKv[key] = valueOf.Field(i).Interface().(CDATA).Text
		default:
			storeKv[key] = fmt.Sprintf("%v", valueOf.Field(i).Interface())
		}
		sortK = append(sortK, key)
	}

	sort.Strings(sortK)

	var buf bytes.Buffer
	for index, val := range sortK {
		if index != 0 {
			buf.WriteString("&")
		}
		buf.WriteString(val + "=" + storeKv[val])
	}
	buf.WriteString("&key=" + GetConf().MchKey)
	hs := Md5(buf.Bytes())
	return strings.ToUpper(hs)
}

//RandomStr 随机生成字符串
func RandomStr(length int) string {
	str := `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	bts := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bts[r.Intn(len(bts))])
	}
	return string(result)
}
