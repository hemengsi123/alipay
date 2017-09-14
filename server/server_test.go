package server

import (
	"log"
	"net/http"
	"testing"

	"github.com/oklog/ulid"
	"github.com/xy02/alipay/pb"
	"github.com/xy02/alipay/trade"
	"github.com/xy02/utils"
	mgo "gopkg.in/mgo.v2"
)

var server pb.AlipayServer
var server2 *Server

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
		NotifyURL:    "http://140.206.154.90:2345/test",
	})
	if err != nil {
		panic(err)
	}
	server2 = &Server{
		tradeCollection: s.DB("testAlipay").C("trade"),
		alipayClient:    client,
	}
	server = server2
}

func TestHandler(t *testing.T) {
	http.HandleFunc("/test", server2.GetNotificationHandler())
	http.ListenAndServe(":2345", nil)
}

func TestServer_PrecreateTrade(t *testing.T) {
	// id := []byte("ABC")
	id := utils.NewULID()
	trade, err := server.PrecreateTrade(nil, &pb.PrecreateParam{
		TradeId:     id[:],
		IdType:      pb.IDType_ULID,
		Subject:     "xx_product",
		AmountInFen: 1,
	})
	if err != nil {
		t.Error(err)
	} else {
		log.Println(*trade)
		// log.Println((*trade).Detail)
	}
}

func TestServer_QueryTrade(t *testing.T) {
	// id := []byte("ABC")
	id := ulid.MustParse("01BSZREFBYK65K6K3YSGC0MNSH")
	trade, err := server.QueryTrade(nil, &pb.QueryParam{
		TradeId: id[:],
	})
	if err != nil {
		t.Error(err)
	} else {
		log.Println(*trade)
		// log.Println((*trade).Detail)
	}
}

func TestServer_RefreshQR(t *testing.T) {
	id := ulid.MustParse("01BSZREFBYK65K6K3YSGC0MNSH")
	trade, err := server.RefreshQR(nil, &pb.RefreshQRParam{
		TradeId: id[:],
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(*trade)
}
