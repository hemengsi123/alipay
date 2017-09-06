package server

import (
	"log"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/xy02/alipay/pb"
	"github.com/xy02/alipay/trade"
)

var server *Server

func init() {
	s, _ := mgo.Dial("localhost")
	client, err := trade.NewAlipayClient(trade.AlipayConfig{
		GatewayURL: trade.GatewayDevURL,
		PriKey:     "../testKeys/xy_pri_key.pem",
		PriKeyPwd:  "",
		// AlipayPubKey: "../testKeys/alipay_dev_pub_key.pem",
		AlipayPubKey: "../testKeys/alipay_dev_rsa2_pub.pem",
		AlipayAppID:  "2016072800108822",
		SignType:     trade.RSA2,
	})
	if err != nil {
		panic(err)
	}
	server = &Server{
		tradeCollection: s.DB("testAlipay").C("trade"),
		alipayClient:    client,
	}
}

func TestServer_CreateQRTrade(t *testing.T) {
	trade, err := server.CreateQRTrade(nil, &pb.CreateQRParam{
		IdType:      pb.IDType_ULID,
		Subject:     "xx product",
		AmountInFen: 1,
		NotifyUrl:   "http://140.206.154.90:2222/test",
	})
	if err != nil {
		t.Error(err)
	} else {
		log.Println(*trade)
		// log.Println((*trade).Detail)
	}

}
