package server

import (
	"context"
	"encoding/json"

	"github.com/xy02/alipay/db"
	"github.com/xy02/alipay/pb"
	"github.com/xy02/alipay/trade"

	mgo "gopkg.in/mgo.v2"
)

//Server 服务器
type Server struct {
	tradeCollection *mgo.Collection
	alipayClient    *trade.AlipayClient
}

//CreateQRTrade 创建交易
func (s *Server) CreateQRTrade(ctx context.Context, data *pb.CreateQRParam) (*pb.QRTrade, error) {
	id, idStr, err := makeID(data.IdType)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := s.alipayClient.PrecreateTrade(trade.PrecreateParam{
		Subject:     data.Subject,
		TotalAmount: data.AmountInFen,
		OutTradeNo:  idStr,
		NotifyURL:   data.NotifyUrl,
	})
	if err != nil {
		return nil, err
	}
	//获取trade
	detail := pb.TradeDetail{}
	if err := json.Unmarshal(jsonBytes, &detail); err != nil {
		return nil, err
	}
	//持久化
	doc := &db.QRTradeDoc{
		IDType:      data.IdType,
		ID:          id,
		Subject:     data.Subject,
		AmountInFen: data.AmountInFen,
		QRCode:      detail.QrCode,
		Status:      pb.TradeStatus_PRECREATE,
		Detail:      detail,
	}
	if err := s.tradeCollection.Insert(doc); err != nil {
		return nil, err
	}
	return &pb.QRTrade{
		IdType:      data.IdType,
		Id:          idStr,
		Subject:     data.Subject,
		AmountInFen: data.AmountInFen,
		QrCode:      detail.QrCode,
		Detail:      &detail,
		Status:      pb.TradeStatus_PRECREATE,
	}, nil
}

//QueryQRTrade 查询
func (s *Server) QueryQRTrade(ctx context.Context, data *pb.QueryQRParam) (*pb.QRTrade, error) {
	return nil, nil
}

//RefreshQR 刷新QR
func (s *Server) RefreshQR(ctx context.Context, data *pb.RefreshQRParam) (*pb.QRTrade, error) {
	return nil, nil
}

//QueryQRTrades 按交易号查询交易记录变化
func (s *Server) QueryQRTrades(data *pb.CreateQRParam, stream pb.Alipay_QueryQRTradesServer) (*pb.QRTrade, error) {
	return nil, nil
}
