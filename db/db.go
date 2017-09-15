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
	Detail        Detail         `bson:"detail"`                   //交易详情
	AppID         string         `bson:"app_id"`                   //支付宝应用id，这个是重要的补充，所有订单都属于一个appID
	StatusChanges []StatusChange `bson:"status_changes,omitempty"` //交易的同步记录
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
	//以下是回掉通知才有的属性
	// app_id string //支付宝应用id
	// gmt_create string //创建交易的时间
	//gmt_payment
	// seller_id string //卖方id
	// seller_email string
	//buyer_id string //同buyer_user_id
	//notify_type //trade_status_sync
	//notify_time //发现同一个交易会重复通知多次
	//open_id
}

//GetPBStatus 获取状态
func (detail *Detail) GetPBStatus() pb.TradeStatus {
	switch detail.TradeStatus {
	case "":
		return pb.TradeStatus_PRECREATE
	case "WAIT_BUYER_PAY":
		return pb.TradeStatus_WAIT
	case "TRADE_CLOSED":
		return pb.TradeStatus_CLOSED
	case "TRADE_SUCCESS":
		return pb.TradeStatus_SUCCESS
	case "TRADE_FINISHED":
		return pb.TradeStatus_FINISHED
	default:
		return pb.TradeStatus_UNKNOWN
	}
}

//FundBill 交易支付使用的资金渠道
type FundBill struct {
	Amount      string `json:"amount" bson:"amount,omitempty"`
	FundChannel string `json:"fund_channel" bson:"fund_channel,omitempty"`
	RealAmount  string `json:"real_amount" bson:"real_amount,omitempty"`
	FundType    string `json:"fund_type" bson:"fund_type,omitempty"`
}

// //PrecreatedDoc 交易的预创建记录，n:1 TradeDoc
// type PrecreatedDoc struct {
// 	ObjectID bson.ObjectId `bson:"_id,omitempty"` //数据库内部ID号，包含创建时间
// 	TradeID  []byte        `bson:"trade_id"`      //交易号
// 	QrCode   string        `bson:"qr_code"`       //预付款二维码地址
// }
