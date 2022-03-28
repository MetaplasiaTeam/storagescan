package storagescan

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"reflect"
	"strings"
)

type StructValueI interface {
	Field(f string) interface{}
	String() string
}

type StructValue struct {
	baseSlotIndex common.Hash

	filedValueMap map[string]struct {
		V Variable
		I uint64
	}

	f GetValueStorageAtFunc
}

func (s StructValue) Field(fd string) interface{} {
	filedValue, ok := s.filedValueMap[fd]
	if !ok {
		return nil
	}

	slotIndex := new(big.Int)
	slotIndex.Add(s.baseSlotIndex.Big(), big.NewInt(int64(filedValue.I)))
	rv := reflect.ValueOf(filedValue.V).Elem()
	rs := rv.FieldByName("SlotIndex")
	rs.Set(reflect.ValueOf(common.BigToHash(slotIndex)))

	return rv.Interface().(Variable).Value(s.f)

}

func (s StructValue) String() string {
	var fSting string
	for filedName := range s.filedValueMap {
		fSting += fmt.Sprintf("%v:%v ", filedName, s.Field(filedName))
	}
	return "struct{" + strings.TrimRight(fSting, " ") + "}"

}
