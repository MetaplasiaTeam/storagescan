package storagescan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type StorageLayout struct {
	Storage []Storage              `json:"storage"`
	Types   map[string]StorageType `json:"types"`
}

type Storage struct {
	AstId    uint   `json:"astId"`
	Contract string `json:"contract"`
	Label    string `json:"label"`
	Offset   uint64 `json:"offset"`
	Slot     string `json:"slot"`
	Type     string `json:"type"`
}

type StorageType struct {
	Base          string    `json:"base"`
	Encoding      string    `json:"encoding"`
	Label         string    `json:"label"`
	Members       []Storage `json:"members"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	NumberOfBytes string    `json:"numberOfBytes"`
}

type Contract struct {
	Address common.Address `json:"address"`

	RPCNode string `json:"rpc_node"`
	// key is variable name, value is variable type
	Variables map[string]Variable `json:"variables"`

	StorageLayout StorageLayout `json:"storage_layout"`
}

type VariableDesc struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func NewContract(address common.Address, rpcNode string) *Contract {
	return &Contract{
		Address:   address,
		RPCNode:   rpcNode,
		Variables: map[string]Variable{},
	}
}

func (c Contract) ParseByStorageLayout(layOutJson string) (err error) {
	err = json.Unmarshal([]byte(layOutJson), &c.StorageLayout)
	if err != nil {
		err = fmt.Errorf("parse storage layout error: %v", err)
		return
	}
	for _, s := range c.StorageLayout.Storage {
		variableName := s.Label
		offset := s.Offset * 8
		sb := new(big.Int)
		sb.SetString(s.Slot, 10)
		slotIndex := common.BigToHash(sb)

		v := c.getVariableByVariableType(s.Type)
		reflect.ValueOf(v).Elem().FieldByName("SlotIndex").Set(reflect.ValueOf(slotIndex))
		if v.Len() < 256 && offset != 0 {
			reflect.ValueOf(v).Elem().FieldByName("Offset").Set(reflect.ValueOf(uint(offset)))
		}
		c.Variables[variableName] = v

	}
	return
}

func (c Contract) GetVariableValue(name string) interface{} {
	return c.Variables[name].Value(GenGetStorageValueFunc(context.Background(), c.RPCNode, c.Address))
}

func (c Contract) GetAllVariables() []VariableDesc {
	var variables []VariableDesc
	for k, v := range c.Variables {
		variables = append(variables, VariableDesc{
			Name: k,
			Type: v.Typ().String(),
		})
	}
	// sort by name
	sort.Slice(variables, func(i, j int) bool {
		return variables[i].Name < variables[j].Name
	})
	return variables
}

func (c Contract) getVariableByVariableType(vt string) Variable {
	if vtForm, ok := c.StorageLayout.Types[vt]; ok {
		switch vtForm.Encoding {
		case "bytes":
			// string
			return &SolidityString{}
		case "inplace":
			if vtForm.Base != "" {
				// array
				arrayRegexp := regexp.MustCompile(`(.*)\[(.*)\]`)
				arrayMatch := arrayRegexp.FindStringSubmatch(vtForm.Label)
				arraySize, _ := strconv.ParseUint(arrayMatch[2], 10, 64)

				return &SolidityArray{
					UnitLength: arraySize,
					UnitTyp:    c.getVariableByVariableType(vtForm.Base),
				}
			}
			// bytes1,uint256,int1
			regExp := regexp.MustCompile(`(bytes|uint|int)(\d+)`)
			isMatch := regExp.MatchString(vtForm.Label)
			if isMatch {
				// bytes1,uint256,int1
				subMatch := regExp.FindStringSubmatch(vtForm.Label)
				length, _ := strconv.ParseUint(subMatch[2], 10, 64)
				switch subMatch[1] {
				case "bytes":
					return &SolidityBytes{
						Length: uint(length * 8),
					}
				case "uint":
					return &SolidityUint{
						Length: uint(length),
					}
				case "int":
					return &SolidityInt{
						Length: uint(length),
					}
				}
			} else {
				// bool,address,struct
				if vtForm.Label == "address" {
					return &SolidityAddress{}
				}

				if vtForm.Label == "bool" {
					return &SolidityBool{}
				}
				// enum
				if strings.HasPrefix(vtForm.Label, "enum") {
					bytesLen, _ := strconv.ParseUint(vtForm.NumberOfBytes, 10, 64)
					return &SolidityUint{
						Length: uint(bytesLen) * 8,
					}
				}

				if strings.HasPrefix(vtForm.Label, "struct") {
					filedValueMap := make(map[string]Variable)
					for _, m := range vtForm.Members {
						offset := m.Offset * 8
						sb := new(big.Int)
						sb.SetString(m.Slot, 10)
						slotIndex := common.BigToHash(sb)

						mv := c.getVariableByVariableType(m.Type)

						reflect.ValueOf(mv).Elem().FieldByName("SlotIndex").Set(reflect.ValueOf(slotIndex))
						if mv.Len() < 256 && offset != 0 {
							reflect.ValueOf(mv).Elem().FieldByName("Offset").Set(reflect.ValueOf(uint(offset)))
						}
						filedValueMap[m.Label] = mv
					}

					return &SolidityStruct{
						FiledValueMap: filedValueMap,
					}

				}

			}

		case "mapping":
			return &SolidityMapping{
				KeyTyp:   c.getVariableByVariableType(vtForm.Key).Typ(),
				ValueTyp: c.getVariableByVariableType(vtForm.Value),
			}

		case "dynamic_array":
			return &SoliditySlice{
				UnitTyp: c.getVariableByVariableType(vtForm.Base),
			}

		}

	}
	return nil
}
