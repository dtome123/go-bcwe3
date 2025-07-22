package contract

import (
	"fmt"
	"math/big"
)

type ContractResult struct {
	Value interface{}
}

type ContractResults []ContractResult

func (r ContractResults) Index(index int) ContractResult {
	return r[index]
}

func (r ContractResults) Len() int {
	return len(r)
}

func (v ContractResult) String() string {
	return fmt.Sprintf("%v", v.Value)
}

func (v ContractResult) AsString() (string, error) {
	str, ok := v.Value.(string)
	if !ok {
		return "", fmt.Errorf("value is not a string")
	}
	return str, nil
}

func (v ContractResult) AsUint64() (uint64, error) {
	num, ok := v.Value.(uint64)
	if !ok {
		return 0, fmt.Errorf("value is not uint64")
	}
	return num, nil
}

func (v ContractResult) AsUnit8() (uint8, error) {
	num, ok := v.Value.(uint8)
	if !ok {
		return 0, fmt.Errorf("value is not uint8")
	}
	return num, nil
}

func (v ContractResult) AsUint16() (uint16, error) {
	num, ok := v.Value.(uint16)
	if !ok {
		return 0, fmt.Errorf("value is not uint16")
	}
	return num, nil
}

func (v ContractResult) AsUint32() (uint32, error) {
	num, ok := v.Value.(uint32)
	if !ok {
		return 0, fmt.Errorf("value is not uint32")
	}
	return num, nil
}

func (v ContractResult) AsInt64() (int64, error) {
	num, ok := v.Value.(int64)
	if !ok {
		return 0, fmt.Errorf("value is not int64")
	}
	return num, nil
}

func (v ContractResult) AsInt32() (int32, error) {
	num, ok := v.Value.(int32)
	if !ok {
		return 0, fmt.Errorf("value is not int32")
	}
	return num, nil
}

func (v ContractResult) AsInt16() (int16, error) {
	num, ok := v.Value.(int16)
	if !ok {
		return 0, fmt.Errorf("value is not int16")
	}
	return num, nil
}

func (v ContractResult) AsInt8() (int8, error) {
	num, ok := v.Value.(int8)
	if !ok {
		return 0, fmt.Errorf("value is not int8")
	}
	return num, nil
}

func (v ContractResult) AsFloat64() (float64, error) {
	num, ok := v.Value.(float64)
	if !ok {
		return 0, fmt.Errorf("value is not float64")
	}
	return num, nil
}

func (v ContractResult) AsBool() (bool, error) {
	num, ok := v.Value.(bool)
	if !ok {
		return false, fmt.Errorf("value is not bool")
	}
	return num, nil
}

func (v ContractResult) AsBytes() ([]byte, error) {
	num, ok := v.Value.([]byte)
	if !ok {
		return nil, fmt.Errorf("value is not []byte")
	}
	return num, nil
}

func (v ContractResult) AsBigInt() (*big.Int, error) {
	num, ok := v.Value.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("value is not *big.Int")
	}
	return num, nil
}

func (v ContractResult) AsBigFloat() (*big.Float, error) {
	num, ok := v.Value.(*big.Float)
	if !ok {
		return nil, fmt.Errorf("value is not *big.Float")
	}
	return num, nil
}
