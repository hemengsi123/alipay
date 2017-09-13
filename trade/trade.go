package trade

import (
	"encoding/json"
	"fmt"

	"github.com/xy02/utils"
)

type (
	//SignType 签名类型
	SignType string
	//Gateway 网关地址
	Gateway string

	//AlipayClient 阿里支付客户端
	AlipayClient struct {
		gatewayURL Gateway  //网关地址
		appID      string   //应用id号
		signType   SignType //RSA,RSA2
		signer     signer   //签名方案
		notifyURL  string   //回掉地址
	}
	//AlipayConfig 支付宝配置
	AlipayConfig struct {
		GatewayURL   Gateway //网关地址
		PriKey       string
		PriKeyPwd    string
		AlipayPubKey string
		AlipayAppID  string
		SignType     SignType
		NotifyURL    string
	}
	//PrecreateParam 预创建数据
	PrecreateParam struct {
		Subject     string //交易主题
		OutTradeNo  string //交易号
		TotalAmount int64  //总价
		// NotifyURL   string //回调地址
	}
)

//对外常量
const (
	GatewayDevURL Gateway  = "https://openapi.alipaydev.com/gateway.do"
	GatewayURL    Gateway  = "https://openapi.alipay.com/gateway.do"
	RSA           SignType = "RSA"
	RSA2          SignType = "RSA2"
)

//NewAlipayClient 创建支付宝客户端
func NewAlipayClient(config AlipayConfig) (*AlipayClient, error) {
	pk, err := utils.DecryptPrivateKey(config.PriKey, []byte(config.PriKeyPwd))
	if err != nil {
		return nil, err
	}
	pubKey, err := utils.RetrievePublicKey(config.AlipayPubKey)
	if err != nil {
		return nil, err
	}
	if config.GatewayURL == "" {
		config.GatewayURL = GatewayURL
	}
	if config.SignType == "" {
		config.SignType = RSA
	}
	var signer signer
	switch config.SignType {
	case RSA2:
		signer = rsa2Signer{
			key:    pk,
			pubKey: pubKey,
		}
	default:
		signer = rsaSigner{
			key:    pk,
			pubKey: pubKey,
		}
	}
	return &AlipayClient{
		gatewayURL: config.GatewayURL,
		appID:      config.AlipayAppID,
		signType:   config.SignType,
		signer:     signer,
		notifyURL:  config.NotifyURL,
	}, nil
}

//GetNotifyURL 返回通知地址
func (client *AlipayClient) GetNotifyURL() string {
	return client.notifyURL
}

//PrecreateTrade 预创建交易
func (client *AlipayClient) PrecreateTrade(param PrecreateParam) (json.RawMessage, error) {
	yuan := param.TotalAmount / 100
	fen := param.TotalAmount % 100
	totalAmount := fmt.Sprintf("%v.%v", yuan, fen)
	if fen < 10 {
		totalAmount = fmt.Sprintf("%v.0%v", yuan, fen)
	}
	bizContent, err := json.Marshal(preCreateContent{
		OutTradeNo:  param.OutTradeNo,
		TotalAmount: totalAmount,
		Subject:     param.Subject,
	})
	if err != nil {
		return nil, err
	}
	req := &preCreateRequest{
		// AppID:      client.appID,
		Method: alCreateMethod,
		// Charset:    alCharset,
		// SignType:   alSignType,
		// Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		BizContent: string(bizContent),
		NotifyURL:  client.notifyURL,
	}
	result := &preCreateResult{}
	if err := client.request(req, result); err != nil {
		return nil, err
	}
	//check sign
	if err := client.signer.verify(result.AlipayTradePrecreateResponse, result.Sign); err != nil {
		return nil, err
	}
	return result.AlipayTradePrecreateResponse, nil
}

//QueryTrade 查询交易,返回交易状态
func (client *AlipayClient) QueryTrade(outTradeNo string) (json.RawMessage, error) {
	bizContent, err := json.Marshal(queryContent{
		OutTradeNo: outTradeNo,
	})
	if err != nil {
		return nil, err
	}
	req := &queryRequest{
		// AppID:      client.appID,
		Method: alQueryMethod,
		// Charset:    alCharset,
		// SignType:   alSignType,
		// Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		BizContent: string(bizContent),
	}
	result := &queryResult{}
	if err := client.request(req, result); err != nil {
		return nil, err
	}
	//check sign
	if err := client.signer.verify(result.AlipayTradeQueryResponse, result.Sign); err != nil {
		return nil, err
	}
	return result.AlipayTradeQueryResponse, nil
}
