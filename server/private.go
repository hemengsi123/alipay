package server

import (
	"encoding/hex"

	"github.com/oklog/ulid"

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
