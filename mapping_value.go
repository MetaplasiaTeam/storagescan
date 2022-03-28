package storagescan

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

type MappingValueI interface {
	Key(k string) interface{}
}

type MappingValue struct {
	baseSlotIndex common.Hash

	keyTyp SolidityTyp

	valueTyp Variable

	f GetValueStorageAtFunc
}

// slotIndex = abi.encode(key,slot)
func (m MappingValue) Key(k string) interface{} {
	var keyByte []byte
	switch m.keyTyp {
	case UintTy:
		keyByte = encodeUintString(k)
	case IntTy:
		keyByte = encodeIntString(k)
	case BytesTy:
		keyByte = encodeByteString(k)
	case StringTy:
		keyByte = []byte(k)
	case AddressTy:
		keyByte = encodeHexString(k)
	default:
		panic("invalid Key-type!")

	}

	slotIndex := crypto.Keccak256Hash(keyByte, m.baseSlotIndex.Bytes())
	rv := reflect.ValueOf(m.valueTyp).Elem()
	rs := rv.FieldByName("SlotIndex")
	rs.Set(reflect.ValueOf(slotIndex))

	return rv.Interface().(Variable).Value(m.f)

}

func encodeHexString(v string) []byte {
	return common.HexToHash(v).Bytes()
}

func encodeByteString(v string) []byte {
	if strings.Contains(v, "0x") {
		return common.RightPadBytes(common.FromHex(v), 32)
	} else {
		return common.RightPadBytes([]byte(v), 32)
	}

}

func encodeUintString(v string) []byte {
	if strings.Contains(v, "0x") {
		return encodeHexString(v)
	} else {
		bn := new(big.Int)
		bn.SetString(v, 10)
		return common.BigToHash(bn).Bytes()
	}

}

func encodeIntString(c string) []byte {
	intVar, err := strconv.ParseInt(c, 0, 64)
	if err != nil {
		panic(err)
	}
	if intVar < 0 {
		// 取反加1
		bs := common.BigToHash(big.NewInt(intVar)).Bytes()
		ub := make([]byte, 0)
		for _, tb := range bs {
			ub = append(ub, ^tb)
		}
		rb := new(big.Int).SetBytes(ub)
		return rb.Add(rb, big.NewInt(1)).Bytes()
	} else {
		return common.BigToHash(big.NewInt(intVar)).Bytes()
	}

}
