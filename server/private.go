package server

import (
	"errors"

	"github.com/xy02/alipay/pb"
	"github.com/xy02/utils"
)

var (
	errIDType = errors.New("invalid id type")
)

func makeID(idType pb.IDType) ([]byte, string, error) {
	switch idType {
	case pb.IDType_ULID:
		id := utils.NewULID()
		return id[:], id.String(), nil
	}
	return nil, "", errIDType
}
