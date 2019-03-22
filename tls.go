package pay

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"sync"
)

var tlsConfig *tls.Config
var tlsSyncOnce sync.Once

func getTLSConfig() (*tls.Config, error) {
	tlsSyncOnce.Do(func() {
		cert, err := tls.LoadX509KeyPair(GetConf().CertPath, GetConf().KeyPath)
		if err != nil {
			panic("cert load fail")
		}
		caData, err := ioutil.ReadFile(GetConf().CaPath)
		if err != nil {
			panic("ca load fail")
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
		}
	})
	return tlsConfig, nil
}

//ClientDo 请求封装
func ClientDo(r *http.Request) (*http.Response, error) {
	config, err := getTLSConfig()
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		TLSClientConfig: config,
	}
	client := http.Client{
		Transport: transport,
	}
	return client.Do(r)
}
