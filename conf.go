package pay

type Conf struct {
	Appid           string
	MchId           string
	MchKey          string
	NotifyUrl       string
	RefundNotifyUrl string
	CertPath        string
	KeyPath         string
	CaPath          string
}

var config *Conf

func SetConf(conf *Conf) {
	config = conf
}

func GetConf() *Conf {
	return config
}
