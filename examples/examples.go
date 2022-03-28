package examples

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"storagescan"
)


func RopstenTestNet() storagescan.GetValueStorageAtFunc {
	return storagescan.GenGetStorageValueFunc(context.Background(), "https://ropsten.infura.io/v3", common.HexToAddress("0x24302f327764f94c15d930f5ac70d362b4a156f9"))
}

// Int Type
// value is int64
func GetIntValue() {
	int1 := storagescan.SolidityInt{
		SlotIndex: common.HexToHash("0x0"),
		Length:    8,
		Offset:    0,
	}
	int1Value := int1.Value(RopstenTestNet())
	log.Printf("int1value:%v\n", int1Value)

	// int1 and int2 belong to the same slot, offset is 8,because int1 takes up 8 bit of space.
	int2 := storagescan.SolidityInt{
		SlotIndex: common.HexToHash("0x0"),
		Length:    128,
		Offset:    8,
	}

	int2Value := int2.Value(RopstenTestNet())
	log.Printf("int2value:%v\n", int2Value)

	int3 := storagescan.SolidityInt{
		SlotIndex: common.HexToHash("0x1"),
		Length:    256,
		Offset:    0,
	}

	int3Value := int3.Value(RopstenTestNet())
	log.Printf("int3value:%v\n", int3Value)

}

// value is uint64,beyond uint64-max is represented as a string
func GetUintValue() {
	uint1 := storagescan.SolidityUint{
		SlotIndex: common.HexToHash("0x2"),
		Length:    8,
		Offset:    0,
	}
	uint1Value := uint1.Value(RopstenTestNet())
	log.Printf("uint1value:%v\n", uint1Value)

	uint2 := storagescan.SolidityUint{
		SlotIndex: common.HexToHash("0x2"),
		Length:    128,
		Offset:    8,
	}
	uint2Value := uint2.Value(RopstenTestNet())
	log.Printf("uint2Value:%v\n", uint2Value)

	uint3 := storagescan.SolidityUint{
		SlotIndex: common.HexToHash("0x3"),
		Length:    256,
		Offset:    0,
	}
	uint3Value := uint3.Value(RopstenTestNet())
	log.Printf("uint3Value:%v\n", uint3Value)

}

// value is bool type
func GetBoolValue() {
	bool1 := storagescan.SolidityBool{
		SlotIndex: common.HexToHash("0x4"),
		Offset:    0,
	}
	bool1Value := bool1.Value(RopstenTestNet())
	log.Printf("bool1Value:%v\n", bool1Value)

	// bool1 takes up 8 bit of space
	bool2 := storagescan.SolidityBool{
		SlotIndex: common.HexToHash("0x4"),
		Offset:    8,
	}
	bool2Value := bool2.Value(RopstenTestNet())
	log.Printf("bool2Value:%v\n", bool2Value)

}

func GetStringValue() {
	string1 := storagescan.SolidityString{
		SlotIndex: common.HexToHash("0x5"),
	}
	string1Value := string1.Value(RopstenTestNet())
	log.Printf("string1Value:%v\n", string1Value)

	string2 := storagescan.SolidityString{
		SlotIndex: common.HexToHash("0x6"),
	}
	string2Value := string2.Value(RopstenTestNet())
	log.Printf("string2Value:%v\n", string2Value)

}


// value is string
func GetFixedBytesValue() {
	// b1 takes up 1*8 bit of space
	b1 := storagescan.SolidityBytes{
		SlotIndex: common.HexToHash("0x7"),
		Length:    8,
		Offset:    0,
	}
	b1Value := b1.Value(RopstenTestNet())
	log.Printf("b1Value:%v\n", b1Value)

	// b2 takes up 8*8=64 bit of space
	b2 := storagescan.SolidityBytes{
		SlotIndex: common.HexToHash("0x7"),
		Length:    64,
		Offset:    8,
	}
	b2Value := b2.Value(RopstenTestNet())
	log.Printf("b2Value:%v\n", b2Value)

	b3 := storagescan.SolidityBytes{
		SlotIndex: common.HexToHash("0x8"),
		Length:    256,
		Offset:    0,
	}
	b3Value := b3.Value(RopstenTestNet())
	log.Printf("b3Value:%v\n", b3Value)

}

// value is common.Address
func GetAddressValue() {
	addr1 := storagescan.SolidityAddress{
		SlotIndex: common.HexToHash("0x9"),
		Offset:    0,
	}
	addr1Value := addr1.Value(RopstenTestNet())
	log.Printf("addr1Value:%v\n", addr1Value)

}

func GetStructValue() {
	i := storagescan.SolidityStruct{
		SlotIndex: common.HexToHash("0xa"),
		FiledValueMap: map[string]struct {
			V storagescan.Variable
			I uint64
		}{
			"id": {
				V: &storagescan.SolidityUint{
					Length: 256,
				},
				I: 0,
			},
			"value": {
				V: &storagescan.SolidityString{},
				I: 1,
			},
		},
	}
	structIValue := i.Value(RopstenTestNet())
	log.Printf("i struct value:%v\n", structIValue)

	idFieldValue := structIValue.(storagescan.StructValueI).Field("id")
	log.Printf("id field value:%v\n", idFieldValue)

	valueFieldValue := structIValue.(storagescan.StructValueI).Field("value")
	log.Printf("'value' field value:%v\n", valueFieldValue)

}

func GetSliceValue() {
	slice1 := storagescan.SoliditySlice{
		SlotIndex: common.HexToHash("0xc"),
		UnitTyp: &storagescan.SolidityUint{
			Length: 8,
		},
	}
	slice1Value := slice1.Value(RopstenTestNet())
	log.Printf("slice1Value:%v\n", slice1Value)

	index1ValueOfSlice1 := slice1Value.(storagescan.SliceArrayValueI).Index(0)
	log.Printf("slice1 index 0 value:%v\n", index1ValueOfSlice1)

	slice5 := storagescan.SoliditySlice{
		SlotIndex: common.HexToHash("0x10"),
		UnitTyp: &storagescan.SolidityStruct{
			FiledValueMap: map[string]struct {
				V storagescan.Variable
				I uint64
			}{
				"id": {V: &storagescan.SolidityUint{
					Length: 256,
				}, I: 0},
				"value": {
					V: &storagescan.SolidityString{},
					I: 1,
				},

			},
		},
	}
	slice5Value := slice5.Value(RopstenTestNet())
	log.Printf("slice5Value:%v\n", slice5Value)


	index1ValueOfSlice5 := slice5Value.(storagescan.SliceArrayValueI).Index(1)

	log.Printf("slice5 index 1 value->%v\n", index1ValueOfSlice5)

}


// array usage is consistent with slice
func GetArrayValue() {
	array1 := storagescan.SolidityArray{
		SlotIndex: common.HexToHash("0x11"),
		UnitTyp:   &storagescan.SolidityUint{
			Length:    8,
		},
		UnitLength: 5,
	}
	array1Value := array1.Value(RopstenTestNet())
	log.Printf("array1Value:%v\n",array1Value)

	array2 := storagescan.SolidityArray{
		SlotIndex:  common.HexToHash("0x12"),
		UnitLength: 5,
		UnitTyp:    &storagescan.SolidityUint{
			Length:    256,
		},
	}
	array2Value := array2.Value(RopstenTestNet())
	log.Printf("array2Value:%v\n",array2Value)


	array5 := storagescan.SolidityArray{
		SlotIndex: common.HexToHash("0x1a"),
		UnitTyp: &storagescan.SolidityStruct{
			FiledValueMap: map[string]struct {
				V storagescan.Variable
				I uint64
			}{
				"id": {V: &storagescan.SolidityUint{
					Length: 256,
				}, I: 0},
				"value": {
					V: &storagescan.SolidityString{},
					I: 1,
				},

			},
		},
		UnitLength: 2,
	}
	array5Value := array5.Value(RopstenTestNet())
	log.Printf("array5Value:%v\n", array5Value)


	index1ValueOfArray5 := array5Value.(storagescan.SliceArrayValueI).Index(1)

	log.Printf("array5 index 1 value->%v\n", index1ValueOfArray5)

}


func GetMappingValue() {
	mapping1 := storagescan.SolidityMapping{
		SlotIndex: common.HexToHash("0x1e"),
		KeyTyp:    storagescan.UintTy,
		ValueTyp:  &storagescan.SolidityString{},
	}

	mapping1KeyValue := mapping1.Value(RopstenTestNet()).(storagescan.MappingValueI).Key("1")
	log.Printf("mapping1KeyValue:%v\n", mapping1KeyValue)


	mapping6 := storagescan.SolidityMapping{
		SlotIndex: common.HexToHash("0x23"),
		KeyTyp:    storagescan.UintTy,
		ValueTyp:  &storagescan.SolidityStruct{
			FiledValueMap: map[string]struct {
				V storagescan.Variable
				I uint64
			}{
				"id": {V: &storagescan.SolidityUint{
					Length: 256,
				}, I: 0},
				"value": {
					V: &storagescan.SolidityString{},
					I: 1,
				},

			},
		},
	}
	mapping6KeyValue := mapping6.Value(RopstenTestNet()).(storagescan.MappingValueI).Key("123")
	log.Printf("mapping6KeyValue:%v\n", mapping6KeyValue)




}


