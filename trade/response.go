package trade

type (
	//AlipayTradePrecreateResponse 阿里创建预付的结果中的应答
	AlipayTradePrecreateResponse struct {
		Code       string `json:"code" bson:"code"`
		Msg        string `json:"msg" bson:"msg"`
		SubCode    string `json:"sub_code" bson:"sub_code"`
		SubMsg     string `json:"sub_msg" bson:"sub_msg"`
		OutTradeNo string `json:"out_trade_no" bson:"out_trade_no"`
		QrCode     string `json:"qr_code" bson:"qr_code"`
	}
	//AlipayTradeQueryResponse 查询结果中的应答
	//TradeStatus交易状态：WAIT_BUYER_PAY（交易创建，等待买家付款）、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）、TRADE_SUCCESS（交易支付成功）、TRADE_FINISHED（交易结束，不可退款）
	AlipayTradeQueryResponse struct {
		Code           string     `bson:"code" json:"code"`
		Msg            string     `bson:"msg" json:"msg"`
		SubCode        string     `bson:"sub_code" json:"sub_code"`
		SubMsg         string     `bson:"sub_msg" json:"sub_msg"`
		TradeNo        string     `bson:"trade_no" json:"trade_no"`
		OutTradeNo     string     `bson:"out_trade_no" json:"out_trade_no"`
		OpenID         string     `bson:"open_id" json:"open_id"`
		BuyerLogonID   string     `bson:"buyer_logon_id" json:"buyer_logon_id"`
		TradeStatus    string     `bson:"trade_status" json:"trade_status"`
		TotalAmount    string     `bson:"total_amount" json:"total_amount"`
		ReceiptAmount  string     `bson:"receipt_amount" json:"receipt_amount"`
		BuyerPayAmount string     `bson:"buyer_pay_amount" json:"buyer_pay_amount"`
		PointAmount    string     `bson:"point_amount" json:"point_amount"`
		InvoiceAmount  string     `bson:"invoice_amount" json:"invoice_amount"`
		SendPayDate    string     `bson:"send_pay_date" json:"send_pay_date"`
		AlipayStoreID  string     `bson:"alipay_store_id" json:"alipay_store_id"`
		StoreID        string     `bson:"store_id" json:"store_id"`
		TerminalID     string     `bson:"terminal_id" json:"terminal_id"`
		StoreName      string     `bson:"store_name" json:"store_name"`
		BuyerUserID    string     `bson:"buyer_user_id" json:"buyer_user_id"`
		FundBillList   []FundBill `bson:"fund_bill_list" json:"fund_bill_list"`
		//		Discount_goods_detail           string `json:"discount_goods_detail"`
		//		Industry_sepc_detail         string `json:"industry_sepc_detail"`
	}
	//FundBill ...
	FundBill struct {
		Amount      string `json:"amount" bson:"amount"`
		FundChannel string `json:"fund_channel" bson:"fund_channel"`
		RealAmount  string `json:"real_amount" bson:"real_amount"`
	}
)
