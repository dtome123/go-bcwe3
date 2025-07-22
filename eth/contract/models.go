package contract

import (
	"fmt"
	"math/big"
)

type ViewSingleResult struct {
	Value interface{}
}

type ViewResults []ViewSingleResult

func (r ViewResults) Index(index int) ViewSingleResult {
	return r[index]
}

func (r ViewResults) Len() int {
	return len(r)
}

func (v ViewSingleResult) String() string {
	return fmt.Sprintf("%v", v.Value)
}

func (v ViewSingleResult) AsString() (string, error) {
	str, ok := v.Value.(string)
	if !ok {
		return "", fmt.Errorf("value is not a string")
	}
	return str, nil
}

func (v ViewSingleResult) AsUint64() (uint64, error) {
	num, ok := v.Value.(uint64)
	if !ok {
		return 0, fmt.Errorf("value is not uint64")
	}
	return num, nil
}

func (v ViewSingleResult) AsUnit8() (uint8, error) {
	num, ok := v.Value.(uint8)
	if !ok {
		return 0, fmt.Errorf("value is not uint8")
	}
	return num, nil
}

func (v ViewSingleResult) AsUint16() (uint16, error) {
	num, ok := v.Value.(uint16)
	if !ok {
		return 0, fmt.Errorf("value is not uint16")
	}
	return num, nil
}

func (v ViewSingleResult) AsUint32() (uint32, error) {
	num, ok := v.Value.(uint32)
	if !ok {
		return 0, fmt.Errorf("value is not uint32")
	}
	return num, nil
}

func (v ViewSingleResult) AsInt64() (int64, error) {
	num, ok := v.Value.(int64)
	if !ok {
		return 0, fmt.Errorf("value is not int64")
	}
	return num, nil
}

func (v ViewSingleResult) AsInt32() (int32, error) {
	num, ok := v.Value.(int32)
	if !ok {
		return 0, fmt.Errorf("value is not int32")
	}
	return num, nil
}

func (v ViewSingleResult) AsInt16() (int16, error) {
	num, ok := v.Value.(int16)
	if !ok {
		return 0, fmt.Errorf("value is not int16")
	}
	return num, nil
}

func (v ViewSingleResult) AsInt8() (int8, error) {
	num, ok := v.Value.(int8)
	if !ok {
		return 0, fmt.Errorf("value is not int8")
	}
	return num, nil
}

func (v ViewSingleResult) AsFloat64() (float64, error) {
	num, ok := v.Value.(float64)
	if !ok {
		return 0, fmt.Errorf("value is not float64")
	}
	return num, nil
}

func (v ViewSingleResult) AsBool() (bool, error) {
	num, ok := v.Value.(bool)
	if !ok {
		return false, fmt.Errorf("value is not bool")
	}
	return num, nil
}

func (v ViewSingleResult) AsBytes() ([]byte, error) {
	num, ok := v.Value.([]byte)
	if !ok {
		return nil, fmt.Errorf("value is not []byte")
	}
	return num, nil
}

func (v ViewSingleResult) AsBigInt() (*big.Int, error) {
	num, ok := v.Value.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("value is not *big.Int")
	}
	return num, nil
}

func (v ViewSingleResult) AsBigFloat() (*big.Float, error) {
	num, ok := v.Value.(*big.Float)
	if !ok {
		return nil, fmt.Errorf("value is not *big.Float")
	}
	return num, nil
}
