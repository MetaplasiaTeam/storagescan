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

	filedValueMap map[string]Variable

	f GetValueStorageAtFunc
}

func (s StructValue) Field(fd string) interface{} {
	filedValue, ok := s.filedValueMap[fd]
	if !ok {
		return nil
	}

	oldSlot := filedValue.Slot()

	slotIndex := new(big.Int)
	slotIndex.Add(s.baseSlotIndex.Big(), filedValue.Slot().Big())

	// convert the slotIndex to common.Hash and assign it to the SlotIndex field of filed Value.V, using reflection
	reflect.ValueOf(filedValue).Elem().FieldByName("SlotIndex").Set(reflect.ValueOf(common.BigToHash(slotIndex)))
	value := filedValue.Value(s.f)
	reflect.ValueOf(filedValue).Elem().FieldByName("SlotIndex").Set(reflect.ValueOf(oldSlot))
	return value

}

func (s StructValue) String() string {
	var fSting string
	for filedName := range s.filedValueMap {
		fSting += fmt.Sprintf("%v:%v ", filedName, s.Field(filedName))
	}
	return "struct{" + strings.TrimRight(fSting, " ") + "}"
}
