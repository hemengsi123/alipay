package db

import (
	"time"

	"github.com/xy02/alipay/pb"
	"gopkg.in/mgo.v2/bson"
)

//TradeDoc 交易文档
type TradeDoc struct {
	ObjectID      bson.ObjectId  `bson:"_id,omitempty"`            //数据库内部ID号，包含创建时间
	ID            []byte         `bson:"id"`                       //交易号
	IDType        pb.IDType      `bson:"id_type"`                  //交易号类型
	Subject       string         `bson:"subject"`                  //主题
	AmountInFen   int64          `bson:"amount_in_fen,minsize"`    //交易总价 单位：分
	QRCode        string         `bson:"qr_code,omitempty"`        //预付款二维码地址，记录最后一个precreate的QR
	Status        pb.TradeStatus `bson:"status"`                   //最新交易状态
	StatusChanges []StatusChange `bson:"status_changes,omitempty"` //交易的同步记录
	Detail        Detail         `bson:"detail,omitempty"`         //交易详情
}

//StatusChange 交易的状态记录，n:1 TradeDoc
type StatusChange struct {
	SyncAt time.Time      `bson:"sync_at"` //同步时间
	Status pb.TradeStatus `bson:"status"`  //交易状态
}

//Detail 交易详情
type Detail struct {
	QRCode         string     `bson:"qr_code,omitempty" json:"qr_code"`
	TradeNo        string     `bson:"trade_no,omitempty" json:"trade_no"`
	OutTradeNo     string     `bson:"out_trade_no,omitempty" json:"out_trade_no"`
	BuyerLogonID   string     `bson:"buyer_logon_id,omitempty" json:"buyer_logon_id"`
	TradeStatus    string     `bson:"trade_status,omitempty" json:"trade_status"`
	TotalAmount    string     `bson:"total_amount,omitempty" json:"total_amount"`
	ReceiptAmount  string     `bson:"receipt_amount,omitempty" json:"receipt_amount"`
	BuyerPayAmount string     `bson:"buyer_pay_amount,omitempty" json:"buyer_pay_amount"`
	PointAmount    string     `bson:"point_amount,omitempty" json:"point_amount"`
	InvoiceAmount  string     `bson:"invoice_amount,omitempty" json:"invoice_amount"`
	SendPayDate    string     `bson:"send_pay_date,omitempty" json:"send_pay_date"`
	StoreID        string     `bson:"store_id,omitempty" json:"store_id"`
	TerminalID     string     `bson:"terminal_id,omitempty" json:"terminal_id"`
	StoreName      string     `bson:"store_name,omitempty" json:"store_name"`
	BuyerUserID    string     `bson:"buyer_user_id,omitempty" json:"buyer_user_id"`
	FundBillList   []FundBill `bson:"fund_bill_list,omitempty" json:"fund_bill_list"`
}

//FundBill 交易支付使用的资金渠道
type FundBill struct {
	Amount      string `json:"amount" bson:"amount"`
	FundChannel string `json:"fund_channel" bson:"fund_channel"`
	RealAmount  string `json:"real_amount" bson:"real_amount"`
	FundType    string `json:"fund_type" bson:"fund_type"`
}

// //PrecreatedDoc 交易的预创建记录，n:1 TradeDoc
// type PrecreatedDoc struct {
// 	ObjectID bson.ObjectId `bson:"_id,omitempty"` //数据库内部ID号，包含创建时间
// 	TradeID  []byte        `bson:"trade_id"`      //交易号
// 	QrCode   string        `bson:"qr_code"`       //预付款二维码地址
// }
