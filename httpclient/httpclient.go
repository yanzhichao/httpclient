package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"net/http"
)
var region HttpAPI
//Region region api
type HttpAPI interface {
  DoRequest(path, method string, requestBody io.Reader,responseBody *string,) (int, error)
}
//APIConf region api config
type APIConf struct {
	Endpoint  string `yaml:"endpoints"`
	Token     string   `yaml:"token"`
	AuthType  string   `yaml:"auth_type"`
	Cacert    string   `yaml:"client-ca-file"`
	Cert      string   `yaml:"tls-cert-file"`
	CertKey   string   `yaml:"tls-private-key-file"`
}

type httpClient struct {
	APIConf
	Client *http.Client
}

//NewRegion NewRegion
func NewHttpClient(c APIConf) (HttpAPI, error) {
	if region == nil {
		re := &httpClient{
			APIConf: c,
		}
		if c.Cacert != "" && c.Cert != "" && c.CertKey != "" {
			pool := x509.NewCertPool()
			caCrt, err := ioutil.ReadFile(c.Cacert)
			if err != nil {
				println("read ca file err: %s", err)
				return nil, err
			}
			pool.AppendCertsFromPEM(caCrt)

			cliCrt, err := tls.LoadX509KeyPair(c.Cert, c.CertKey)
			if err != nil {
				println("Loadx509keypair err: %s", err)
				return nil, err
			}
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      pool,
					Certificates: []tls.Certificate{cliCrt},
				},
			}
			re.Client = &http.Client{
				Transport: tr,
			}
		} else {
			re.Client = http.DefaultClient
		}
		region = re
	}
	return region, nil
}


func (r *httpClient) GetEndpoint() string {
	return r.Endpoint
}

//DoRequest do request
func (r *httpClient) DoRequest(path, method string, requestBody io.Reader,responseBody *string,) (int, error) {

	request, err := http.NewRequest(method, r.GetEndpoint()+path, requestBody)
	if err != nil {
		println(err.Error())
		return 500, err

	}
	//request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Type", "application/json")
	if r.Token != "" {
		request.Header.Set("Authorization", "Token "+r.Token)
	}
	res, err := r.Client.Do(request)
	if err != nil {
		return 500, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if responseBody != nil {
		re, err := ioutil.ReadAll(res.Body)
		if err != nil {
			println(err.Error())
			return 500, err
		}
		*responseBody = string(re)
	}
	return res.StatusCode, err
}


