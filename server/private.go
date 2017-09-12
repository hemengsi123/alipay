package server

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/oklog/ulid"

	"github.com/xy02/alipay/db"
	"github.com/xy02/alipay/pb"
)

func stringifyID(id []byte, idType pb.IDType) string {
	switch idType {
	case pb.IDType_ULID:
		tmp := &ulid.ULID{}
		if err := tmp.UnmarshalBinary(id); err == nil {
			return tmp.String()
		}
		fallthrough
	case pb.IDType_UTF8:
		return string(id)
	default:
		//hex
		return hex.EncodeToString(id)
	}
}

func parseDoc2Trade(doc *db.TradeDoc) *pb.Trade {
	return &pb.Trade{
		Id:            doc.ID,
		IdType:        doc.IDType,
		Subject:       doc.Subject,
		AmountInFen:   doc.AmountInFen,
		QrCode:        doc.QRCode,
		Status:        doc.Status,
		StatusChanges: getStatusChanges(doc.StatusChanges),
		Detail:        getDetail(doc.Detail),
		CreatedAt:     getCreatedAt(doc.ObjectID),
	}
}

func getStatusChanges(changes []db.StatusChange) []*pb.StatusChange {
	if changes == nil {
		return nil
	}
	result := []*pb.StatusChange{}
	for _, change := range changes {
		data := &pb.StatusChange{
			Status: change.Status,
			SyncedAt: &timestamp.Timestamp{
				Seconds: change.SyncAt.Unix(),
				Nanos:   int32(change.SyncAt.Nanosecond()),
			},
		}
		result = append(result, data)
	}
	return result
}

func getDetail(detailDoc db.Detail) *pb.TradeDetail {
	buf, err := json.Marshal(detailDoc)
	if err != nil {
		panic(err)
	}
	detail := &pb.TradeDetail{}
	if err := json.Unmarshal(buf, detail); err != nil {
		panic(err)
	}
	return detail
}

func getCreatedAt(objID bson.ObjectId) *timestamp.Timestamp {
	var createdAt time.Time
	if objID == "" {
		createdAt = time.Now()
	} else {
		createdAt = objID.Time()
	}
	return &timestamp.Timestamp{
		Seconds: createdAt.Unix(),
		Nanos:   int32(createdAt.Nanosecond()),
	}
}
