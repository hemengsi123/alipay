package main

import (
	"net/http"

	"github.com/xy02/alipay/pb"
	ser "github.com/xy02/alipay/server"
	"github.com/xy02/alipay/trade"

	mgo "gopkg.in/mgo.v2"
)

var server pb.AlipayServer
var server2 *ser.Server

func init() {
	s, _ := mgo.Dial("localhost")
	client, err := trade.NewAlipayClient(trade.AlipayConfig{
		GatewayURL: trade.GatewayDevURL,
		PriKey:     "./testKeys/xy_pri_key.pem",
		PriKeyPwd:  "",
		// AlipayPubKey: "../testKeys/alipay_dev_pub_key.pem",
		AlipayPubKey: "./testKeys/alipay_dev_rsa2_pub.pem",
		AlipayAppID:  "2016072800108822",
		SignType:     trade.RSA2,
		NotifyURL:    "http://140.206.154.90:2345/test",
	})
	if err != nil {
		panic(err)
	}
	server2 = &ser.Server{}
	server2.SetAlipayClient(client)
	server2.SetTradeCollection(s.DB("testAlipay").C("trade"))
	server = server2
}
func main() {
	http.HandleFunc("/test", server2.GetNotificationHandler())
	http.ListenAndServe(":2345", nil)
}
