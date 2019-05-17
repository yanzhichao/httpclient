package main

import (
     "httpclient/httpclient"
)

func main() {
     conf := httpclient.APIConf{
          Endpoint: "https://192.168.245.129:8443",
          Cacert: "D:/cet36/好雨/ca/ca.pem ",
          Cert: "D:/cet36/好雨/ca/client.pem",
          CertKey: "D:/cet36/好雨/ca/client.key.pem",
     }

     client, _ := httpclient.NewHttpClient(conf)

     var body string
     client.DoRequest("/v2/cluster", "GET",nil,&body)

     println(body)

}