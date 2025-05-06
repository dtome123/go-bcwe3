package contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
)

func CallViewFunction[T any](contract *bind.BoundContract, method string, params ...interface{}) (T, error) {

	var zero T
	var result []interface{}

	if method == "" {
		return zero, fmt.Errorf("method name cannot be empty")
	}

	if contract == nil {
		return zero, fmt.Errorf("contract is nil")
	}

	err := contract.Call(nil, &result, method, params...)
	if err != nil {
		return zero, err
	}
	if len(result) == 0 {
		return zero, fmt.Errorf("no result returned for method %s", method)
	}

	val, ok := result[0].(T)
	if !ok {
		return zero, fmt.Errorf("cannot convert result to type %T", zero)
	}
	return val, nil
}
