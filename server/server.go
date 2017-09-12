package server

import (
	"context"
	"encoding/json"

	"github.com/xy02/alipay/db"
	"github.com/xy02/alipay/pb"
	"github.com/xy02/alipay/trade"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Server 服务器
type Server struct {
	tradeCollection *mgo.Collection
	alipayClient    *trade.AlipayClient
}

//PrecreateTrade 预创建交易
func (s *Server) PrecreateTrade(ctx context.Context, param *pb.PrecreateParam) (*pb.Trade, error) {
	outTradeNo := stringifyID(param.TradeId, param.IdType)
	//创建交易
	jsonBytes, err := s.alipayClient.PrecreateTrade(trade.PrecreateParam{
		Subject:     param.Subject,
		TotalAmount: param.AmountInFen,
		OutTradeNo:  outTradeNo,
		NotifyURL:   param.NotifyUrl,
	})
	if err != nil {
		return nil, err
	}
	//获取trade detail
	detailDoc := db.Detail{}
	if err := json.Unmarshal(jsonBytes, &detailDoc); err != nil {
		return nil, err
	}
	//持久化
	tradeDoc := &db.TradeDoc{
		ID:          param.TradeId,
		IDType:      param.IdType,
		Subject:     param.Subject,
		AmountInFen: param.AmountInFen,
		QRCode:      detailDoc.QRCode,
		Status:      pb.TradeStatus_PRECREATE,
		Detail:      detailDoc,
	}
	if err := s.tradeCollection.Insert(tradeDoc); err != nil {
		return nil, err
	}
	return parseDoc2Trade(tradeDoc), nil
}

//QueryTrade 查询
func (s *Server) QueryTrade(ctx context.Context, param *pb.QueryParam) (*pb.Trade, error) {
	//查数据库
	tradeDoc := &db.TradeDoc{}
	if err := s.tradeCollection.Find(bson.M{
		db.ID: param.TradeId,
	}).One(tradeDoc); err != nil {
		return nil, err
	}
	if tradeDoc.Status == pb.TradeStatus_FINISHED || tradeDoc.Status == pb.TradeStatus_CLOSED {
		//交易已经结束
		return parseDoc2Trade(tradeDoc), nil
	}
	outTradeNo := stringifyID(param.TradeId, param.IdType)
	jsonBytes, err := s.alipayClient.QueryTrade(outTradeNo)
	if err != nil {
		return nil, err
	}
	//获取trade detail
	detailDoc := db.Detail{}
	if err := json.Unmarshal(jsonBytes, &detailDoc); err != nil {
		return nil, err
	}
	if detailDoc.TradeStatus != tradeDoc.Detail.TradeStatus {
		//有变化
		tradeDoc.Detail = detailDoc
		if err := s.tradeCollection.Update(bson.M{db.ID: param.TradeId}, tradeDoc); err != nil {
			return nil, err
		}
	}
	return parseDoc2Trade(tradeDoc), nil
}

//RefreshQR 刷新QR
func (s *Server) RefreshQR(ctx context.Context, param *pb.RefreshQRParam) (*pb.Trade, error) {

	return nil, nil
}
