package trade

import (
	"encoding/json"
	"log"
	"testing"
)

var client *AlipayClient

func init() {
	c, err := NewAlipayClient(AlipayConfig{
		GatewayURL: GatewayDevURL,
		PriKey:     "../testKeys/xy_pri_key.pem",
		PriKeyPwd:  "",
		// AlipayPubKey: "../testKeys/alipay_dev_pub_key.pem",
		AlipayPubKey: "../testKeys/alipay_dev_rsa2_pub.pem",
		AlipayAppID:  "2016072800108822",
		SignType:     RSA2,
	})
	if err != nil {
		panic(err)
	}
	client = c
}

func TestQuery(t *testing.T) {
	log.Println(1 == 1, 2)
	id := "01BSTK5GGQE3S14ZA5SB91DSQY"
	buf, err := client.QueryTrade(id)
	if err != nil {
		t.Fatal(err)
	}
	res := &AlipayTradeQueryResponse{}
	if err := json.Unmarshal(buf, res); err != nil {
		panic(err)
	}
	log.Println(res)
}
