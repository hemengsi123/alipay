package server

import (
	"log"
	"testing"

	"github.com/oklog/ulid"
	"github.com/xy02/alipay/pb"
	"github.com/xy02/alipay/trade"
	"github.com/xy02/utils"
	mgo "gopkg.in/mgo.v2"
)

var server pb.AlipayServer

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
		NotifyURL:    "http://140.206.154.90:2222/test",
	})
	if err != nil {
		panic(err)
	}
	server = &Server{
		tradeCollection: s.DB("testAlipay").C("trade"),
		alipayClient:    client,
	}
}

func TestServer_PrecreateTrade(t *testing.T) {
	// id := []byte("ABC")
	id := utils.NewULID()
	trade, err := server.PrecreateTrade(nil, &pb.PrecreateParam{
		TradeId:     id[:],
		IdType:      pb.IDType_ULID,
		Subject:     "xx_product",
		AmountInFen: 2,
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
	id := ulid.MustParse("01BSTK5GGQE3S14ZA5SB91DSQY")
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
	id := ulid.MustParse("01BSTK5GGQE3S14ZA5SB91DSQY")
	trade, err := server.RefreshQR(nil, &pb.RefreshQRParam{
		TradeId: id[:],
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(*trade)
}
