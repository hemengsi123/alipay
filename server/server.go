package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/net/context"

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

func (s *Server) SetTradeCollection(c *mgo.Collection) {
	s.tradeCollection = c
}
func (s *Server) SetAlipayClient(c *trade.AlipayClient) {
	s.alipayClient = c
}

//GetNotificationHandler 返回通知回掉处理器
func (s *Server) GetNotificationHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// body, err := ioutil.ReadAll(r.Body)
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}
		// log.Println(string(body))
		//验签
		if err := s.alipayClient.VerifyNotification(r.Form); err != nil {
			// log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}
		// log.Println(r.Form)
		fundBillList := []db.FundBill{}
		strFundList := strings.Replace(r.Form.Get("fund_bill_list"), "fundChannel", "fund_channel", 1)
		if err := json.Unmarshal([]byte(strFundList), &fundBillList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}
		detailDoc := db.Detail{
			TradeNo:        r.Form.Get("trade_no"),
			OutTradeNo:     r.Form.Get("out_trade_no"),
			BuyerLogonID:   r.Form.Get("buyer_logon_id"),
			TradeStatus:    r.Form.Get("trade_status"),
			TotalAmount:    r.Form.Get("total_amount"),
			ReceiptAmount:  r.Form.Get("receipt_amount"),
			BuyerPayAmount: r.Form.Get("buyer_pay_amount"),
			PointAmount:    r.Form.Get("point_amount"),
			InvoiceAmount:  r.Form.Get("invoice_amount"),
			SendPayDate:    r.Form.Get("gmt_payment"),
			StoreID:        r.Form.Get("store_id"),
			TerminalID:     r.Form.Get("terminal_id"),
			StoreName:      r.Form.Get("store_name"),
			BuyerUserID:    r.Form.Get("buyer_id"),
			FundBillList:   fundBillList,
		}
		//查数据库
		tradeDoc := &db.TradeDoc{}
		if err := s.tradeCollection.Find(bson.M{
			db.OutTradeNo: detailDoc.OutTradeNo,
		}).One(tradeDoc); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}
		status := detailDoc.GetPBStatus()
		if tradeDoc.Status != status {
			//有变化
			tradeDoc.Status = status
			tradeDoc.Detail = detailDoc
			if err := s.tradeCollection.Update(bson.M{db.OutTradeNo: detailDoc.OutTradeNo}, tradeDoc); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				panic(err)
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

//PrecreateTrade 预创建交易
func (s *Server) PrecreateTrade(ctx context.Context, param *pb.PrecreateParam) (*pb.Trade, error) {
	outTradeNo := stringifyID(param.TradeId, param.IdType)
	//创建交易
	jsonBytes, err := s.alipayClient.PrecreateTrade(trade.PrecreateParam{
		Subject:     param.Subject,
		TotalAmount: param.AmountInFen,
		OutTradeNo:  outTradeNo,
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
		AppID:       s.alipayClient.GetAppID(),
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
	outTradeNo := stringifyID(param.TradeId, tradeDoc.IDType)
	jsonBytes, err := s.alipayClient.QueryTrade(outTradeNo)
	if err != nil {
		return nil, err
	}
	//获取trade detail
	detailDoc := db.Detail{}
	if err := json.Unmarshal(jsonBytes, &detailDoc); err != nil {
		return nil, err
	}
	status := detailDoc.GetPBStatus()
	if tradeDoc.Status != status {
		//有变化
		tradeDoc.Status = status
		tradeDoc.Detail = detailDoc
		if err := s.tradeCollection.Update(bson.M{db.ID: param.TradeId}, tradeDoc); err != nil {
			return nil, err
		}
	}
	return parseDoc2Trade(tradeDoc), nil
}

//RefreshQR 刷新QR
func (s *Server) RefreshQR(ctx context.Context, param *pb.RefreshQRParam) (*pb.Trade, error) {
	//查数据库
	tradeDoc := &db.TradeDoc{}
	if err := s.tradeCollection.Find(bson.M{
		db.ID: param.TradeId,
	}).One(tradeDoc); err != nil {
		return nil, err
	}
	if tradeDoc.Status != pb.TradeStatus_PRECREATE && tradeDoc.Status != pb.TradeStatus_WAIT {
		//交易尚未支付
		return parseDoc2Trade(tradeDoc), nil
	}
	outTradeNo := stringifyID(param.TradeId, tradeDoc.IDType)
	//创建预交易
	jsonBytes, err := s.alipayClient.PrecreateTrade(trade.PrecreateParam{
		Subject:     tradeDoc.Subject,
		TotalAmount: tradeDoc.AmountInFen,
		OutTradeNo:  outTradeNo,
	})
	if err != nil {
		return nil, err
	}
	//更新 QR
	detailDoc := db.Detail{}
	if err := json.Unmarshal(jsonBytes, &detailDoc); err != nil {
		return nil, err
	}
	if detailDoc.QRCode != "" {
		tradeDoc.Status = detailDoc.GetPBStatus()
		tradeDoc.QRCode = detailDoc.QRCode
		tradeDoc.Detail = detailDoc
		if err := s.tradeCollection.Update(bson.M{db.ID: param.TradeId}, tradeDoc); err != nil {
			return nil, err
		}
	}
	return parseDoc2Trade(tradeDoc), nil
}

//WatchTrade 监控交易
func (s *Server) WatchTrade(param *pb.WatchParam, stream pb.Alipay_WatchTradeServer) error {
	return nil
}
