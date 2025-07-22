package erc165

import (
	"context"
)

type ERC165 interface {
	SupportInterface(ctx context.Context, contractAddr string, interfaceIdBytes [4]byte) (bool, error)
}
