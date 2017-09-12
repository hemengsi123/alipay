package server

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

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
	detail := pb.TradeDetail{}
	if err := json.Unmarshal(jsonBytes, &detail); err != nil {
		return nil, err
	}
	detailDoc := db.Detail{}
	if err := json.Unmarshal(jsonBytes, &detailDoc); err != nil {
		return nil, err
	}
	//持久化
	doc := &db.TradeDoc{
		ID:          param.TradeId,
		IDType:      param.IdType,
		Subject:     param.Subject,
		AmountInFen: param.AmountInFen,
		QRCode:      detailDoc.QRCode,
		Status:      pb.TradeStatus_PRECREATE,
		Detail:      detailDoc,
	}
	if err := s.tradeCollection.Insert(doc); err != nil {
		return nil, err
	}
	//return
	now := time.Now()
	return &pb.Trade{
		Id:          param.TradeId,
		IdType:      param.IdType,
		Subject:     param.Subject,
		AmountInFen: param.AmountInFen,
		QrCode:      detailDoc.QRCode,
		Detail:      &detail,
		Status:      pb.TradeStatus_PRECREATE,
		CreatedAt: &timestamp.Timestamp{
			Seconds: now.Unix(),
		},
	}, nil
}

//QueryTrade 查询
func (s *Server) QueryTrade(ctx context.Context, param *pb.QueryParam) (*pb.Trade, error) {
	return nil, nil
}

//RefreshQR 刷新QR
func (s *Server) RefreshQR(ctx context.Context, param *pb.RefreshQRParam) (*pb.Trade, error) {
	return nil, nil
}
