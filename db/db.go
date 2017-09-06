package db

import "github.com/xy02/alipay/pb"
import "gopkg.in/mgo.v2/bson"

//QRTradeDoc 交易文档
type QRTradeDoc struct {
	ObjectID    bson.ObjectId  `bson:"_id,omitempty"` //数据库内部ID号，包含创建时间
	IDType      pb.IDType      `bson:"id_type"`       //ID类型
	ID          []byte         `bson:"id"`            //交易号
	Subject     string         `bson:"subject"`       //主题
	AmountInFen int64          `bson:"amount_in_fen"` //交易总价 单位：分
	QRCode      string         `bson:"qr_code"`       //预付款二维码地址
	Status      pb.TradeStatus `bson:"status"`        //交易状态
	Detail      pb.TradeDetail `bson:"detail"`        //交易详情
}
