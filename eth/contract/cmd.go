package contract

import (
	"github.com/dtome123/go-bcwe3/eth/provider"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
)

type Cmd struct {
	address  string
	caller   *bind.BoundContract
	provider provider.Provider
}
