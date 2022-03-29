package storagescan

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type SolidityTyp uint8

type GetValueStorageAtFunc func(s common.Hash) []byte

// GenGetStorageValueFunc this is a wrapper for the storage at function
func GenGetStorageValueFunc(ctx context.Context, rpcNode string, contractAddr common.Address) GetValueStorageAtFunc {
	return func(s common.Hash) []byte {
		cli, err := ethclient.DialContext(ctx, rpcNode)
		if err != nil {
			return nil
		}
		var value []byte
		value, err = cli.StorageAt(ctx, contractAddr, s, nil)
		if err != nil {
			return nil
		}
		return value
	}
}

type Variable interface {
	Typ() SolidityTyp

	Value(f GetValueStorageAtFunc) interface{}
}

// Type enumerator
const (
	IntTy SolidityTyp = iota
	UintTy
	BoolTy
	StringTy
	SliceTy
	ArrayTy
	MappingTy
	AddressTy
	BytesTy
	StructTy
)

type SolidityInt struct {
	SlotIndex common.Hash

	Length uint

	Offset uint
}

func (s SolidityInt) Typ() SolidityTyp {
	return IntTy
}

func (s SolidityInt) Value(f GetValueStorageAtFunc) interface{} {
	v := f(s.SlotIndex)

	vb := common.BytesToHash(v).Big()
	vb.Rsh(vb, s.Offset)

	// get mask for length
	mask := new(big.Int)
	mask.SetBit(mask, int(s.Length), 1).Sub(mask, big.NewInt(1))

	// get value by mask
	vb.And(vb, mask)

	// signBit is 0 if the value is positive and 1 if it is negative
	signBit := new(big.Int)
	signBit.Rsh(vb, s.Length-1)
	if signBit.Uint64() == 0 {
		return vb.Uint64()

	} else {
		// flip the bits
		vb.Sub(vb, big.NewInt(1))
		r := make([]byte, 0)
		for _, b := range vb.Bytes() {
			r = append(r, ^b)
		}
		// convert back to big int
		return -new(big.Int).SetBytes(r).Int64()
	}

}

type SolidityUint struct {
	SlotIndex common.Hash

	Length uint

	Offset uint
}

func (s SolidityUint) Typ() SolidityTyp {
	return UintTy
}

func (s SolidityUint) Value(f GetValueStorageAtFunc) interface{} {
	v := f(s.SlotIndex)
	vb := common.BytesToHash(v).Big()
	vb.Rsh(vb, s.Offset)

	mask := new(big.Int)
	mask.SetBit(mask, int(s.Length), 1).Sub(mask, big.NewInt(1))

	vb.And(vb, mask)

	// if vb > uint64 max, return hex string, else return uint64
	if vb.Cmp(big.NewInt(0).SetUint64(1<<64-1)) > 0 {
		return common.BigToHash(vb).Hex()
	} else {
		return vb.Uint64()
	}

}

type SolidityAddress struct {
	SlotIndex common.Hash

	Offset uint
}

func (s SolidityAddress) Typ() SolidityTyp {
	return AddressTy
}

func (s SolidityAddress) Value(f GetValueStorageAtFunc) interface{} {
	v := f(s.SlotIndex)
	vb := common.BytesToHash(v).Big()
	vb.Rsh(vb, s.Offset)

	lengthOffset := new(big.Int)
	lengthOffset.SetBit(lengthOffset, 160, 1).Sub(lengthOffset, big.NewInt(1))

	vb.And(vb, lengthOffset)

	return common.BytesToAddress(vb.Bytes())
}

type SolidityBool struct {
	SlotIndex common.Hash

	Offset uint
}

func (s SolidityBool) Typ() SolidityTyp {
	return BoolTy

}

func (s SolidityBool) Value(f GetValueStorageAtFunc) interface{} {
	v := f(s.SlotIndex)
	vb := common.BytesToHash(v).Big()
	vb.Rsh(vb, s.Offset)

	lengthOffset := new(big.Int)
	lengthOffset.SetBit(lengthOffset, 8, 1).Sub(lengthOffset, big.NewInt(1))

	vb.And(vb, lengthOffset)
	return vb.Uint64() == 1

}

type SolidityString struct {
	SlotIndex common.Hash
}

func (s SolidityString) Typ() SolidityTyp {
	return StringTy
}

// Value calculate the string length of the current slot record
// the length of the string exceeds 31 bytes (0x1f), and the entire slot stores the length of the string*2+1
// the length of the string does not exceed 31 bytes, the rightmost bit of the entire slot stores the character length*2, and the leftmost stores the string content
// if the last digit is odd then it is a long string, otherwise it is a short  string
func (s SolidityString) Value(f GetValueStorageAtFunc) interface{} {
	data := f(s.SlotIndex)
	v := common.BytesToHash(data).Big()

	// get the last digit of v
	lastDigit := v.Bit(0)

	//  equal to 1 means it is a long string
	if lastDigit == 1 {
		// get the current string length bit
		length := new(big.Int)
		length.Sub(v, big.NewInt(1)).Div(length, big.NewInt(2)).Mul(length, big.NewInt(8))

		remainB := new(big.Int)
		remainB.Mod(length, big.NewInt(256))

		slotNum := new(big.Int)
		if remainB.Uint64() == 0 {
			slotNum.Div(length, big.NewInt(256))
		} else {
			slotNum.Div(length, big.NewInt(256)).Add(slotNum, big.NewInt(1))
		}

		firstSlotIndex := crypto.Keccak256Hash(s.SlotIndex.Bytes())

		value := f(firstSlotIndex)

		for i := int64(1); i < slotNum.Int64()-1; i++ {
			nextSlot := new(big.Int)
			nextSlot.Add(firstSlotIndex.Big(), big.NewInt(i))
			nextValue := f(common.BigToHash(nextSlot))
			value = append(value, nextValue...)
		}

		lastSlotIndex := new(big.Int)
		lastSlotIndex.Add(firstSlotIndex.Big(), big.NewInt(slotNum.Int64()-1))

		lastSlotValue := f(common.BigToHash(lastSlotIndex))

		if remainB.Uint64() == 0 {
			value = append(value, lastSlotValue...)
		} else {
			// move right to get the final value
			lastValueBig := common.BytesToHash(lastSlotValue).Big()
			lastValueBig.Rsh(lastValueBig, 256-uint(remainB.Uint64()))
			value = append(value, lastValueBig.Bytes()...)
		}

		return string(value)
	} else {

		length := new(big.Int)
		length.And(v, big.NewInt(0xff))
		length.Div(length, big.NewInt(2)).Mul(length, big.NewInt(8))

		v.Rsh(v, 256-uint(length.Uint64()))

		return string(v.Bytes())
	}

}

type SolidityBytes struct {
	SlotIndex common.Hash

	Length uint

	Offset uint
}

func (s SolidityBytes) Typ() SolidityTyp {
	return BytesTy
}

func (s SolidityBytes) Value(f GetValueStorageAtFunc) interface{} {
	v := f(s.SlotIndex)
	vb := common.BytesToHash(v).Big()
	vb.Rsh(vb, s.Offset)

	lengthOffset := new(big.Int)
	lengthOffset.SetBit(lengthOffset, int(s.Length), 1).Sub(lengthOffset, big.NewInt(1))

	vb.And(vb, lengthOffset)

	return string(common.TrimRightZeroes(vb.Bytes()))

}

// bytes = byte[] = uint8[]

type SoliditySlice struct {
	SlotIndex common.Hash

	UnitTyp Variable `json:"unit_typ"`
}

func (s SoliditySlice) Typ() SolidityTyp {
	return SliceTy
}

func (s SoliditySlice) Value(f GetValueStorageAtFunc) interface{} {
	length := common.BytesToHash(f(s.SlotIndex)).Big().Uint64()
	valueSlotIndex := crypto.Keccak256Hash(s.SlotIndex.Bytes())

	switch s.UnitTyp.Typ() {
	case IntTy:
		si := s.UnitTyp.(*SolidityInt)
		return IntSliceValue{
			slotIndex:     valueSlotIndex,
			length:        length,
			uintBitLength: si.Length,
			f:             f,
		}
	case UintTy:
		su := s.UnitTyp.(*SolidityUint)
		return UintSliceValue{
			slotIndex:     valueSlotIndex,
			length:        length,
			uintBitLength: su.Length,
			f:             f,
		}
	case BytesTy:
		sb := s.UnitTyp.(*SolidityBytes)
		return BytesSliceValue{
			slotIndex:     valueSlotIndex,
			length:        length,
			uintBitLength: sb.Length,
			f:             f,
		}
	case StructTy:
		ss := s.UnitTyp.(*SolidityStruct)
		return StructSliceValue{
			slotIndex:     valueSlotIndex,
			length:        length,
			filedValueMap: ss.FiledValueMap,
			f:             f,
		}

	case BoolTy:
		return BoolSliceValue{
			slotIndex: valueSlotIndex,
			length:    length,
			f:         f,
		}
	case StringTy:
		return StringSliceValue{
			slotIndex: valueSlotIndex,
			length:    length,
			f:         f,
		}
	case AddressTy:
		return AddressSliceValue{
			slotIndex: valueSlotIndex,
			length:    length,
			f:         f,
		}

	}
	return nil

}

type SolidityArray struct {
	SlotIndex common.Hash

	UnitLength uint64 `json:"unit_length"`

	UnitTyp Variable `json:"unit_typ"`
}

func (s SolidityArray) Typ() SolidityTyp {
	return ArrayTy
}

func (s SolidityArray) Value(f GetValueStorageAtFunc) interface{} {
	switch s.UnitTyp.Typ() {
	case IntTy:
		si := s.UnitTyp.(*SolidityInt)
		return IntSliceValue{
			slotIndex:     s.SlotIndex,
			length:        s.UnitLength,
			uintBitLength: si.Length,
			f:             f,
		}
	case UintTy:
		su := s.UnitTyp.(*SolidityUint)
		return UintSliceValue{
			slotIndex:     s.SlotIndex,
			length:        s.UnitLength,
			uintBitLength: su.Length,
			f:             f,
		}
	case BytesTy:
		sb := s.UnitTyp.(*SolidityBytes)
		return BytesSliceValue{
			slotIndex:     s.SlotIndex,
			length:        s.UnitLength,
			uintBitLength: sb.Length,
			f:             f,
		}
	case StructTy:
		ss := s.UnitTyp.(*SolidityStruct)
		return StructSliceValue{
			slotIndex:     s.SlotIndex,
			length:        s.UnitLength,
			filedValueMap: ss.FiledValueMap,
			f:             f,
		}

	case BoolTy:
		return BoolSliceValue{
			length:    s.UnitLength,
			slotIndex: s.SlotIndex,
			f:         f,
		}
	case StringTy:
		return StringSliceValue{
			length:    s.UnitLength,
			slotIndex: s.SlotIndex,
			f:         f,
		}
	case AddressTy:
		return AddressSliceValue{
			length:    s.UnitLength,
			slotIndex: s.SlotIndex,
			f:         f,
		}

	}

	return nil

}

type SolidityStruct struct {
	SlotIndex common.Hash
	// field name and value mapping
	FiledValueMap map[string]struct {
		// field value type
		V Variable
		// the struct position slot index
		I uint64
	}
}

func (s SolidityStruct) Typ() SolidityTyp {
	return StructTy
}

func (s SolidityStruct) Value(f GetValueStorageAtFunc) interface{} {
	return StructValue{
		baseSlotIndex: s.SlotIndex,
		filedValueMap: s.FiledValueMap,
		f:             f,
	}

}

type SolidityMapping struct {
	SlotIndex common.Hash

	KeyTyp SolidityTyp

	ValueTyp Variable `json:"value_typ"`
}

func (s SolidityMapping) Typ() SolidityTyp {
	return MappingTy
}

func (s SolidityMapping) Value(f GetValueStorageAtFunc) interface{} {
	m := MappingValue{
		baseSlotIndex: s.SlotIndex,
		keyTyp:        s.KeyTyp,
		valueTyp:      s.ValueTyp,
		f:             f,
	}
	return m

}
