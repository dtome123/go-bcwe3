package erc1155

import "errors"

var (
	ErrParseABI               = errors.New("failed to parse ABI")
	ErrPackBalanceOf          = errors.New("failed to pack balanceOf")
	ErrCallBalanceOf          = errors.New("failed to call balanceOf")
	ErrUnpackBalanceOf        = errors.New("failed to unpack balanceOf result")
	ErrEmptyBalanceOfResponse = errors.New("empty response from balanceOf")

	ErrPackTokenOfOwnerByIndex   = errors.New("failed to pack tokenOfOwnerByIndex")
	ErrCallTokenOfOwnerByIndex   = errors.New("failed to call tokenOfOwnerByIndex")
	ErrUnpackTokenOfOwnerByIndex = errors.New("failed to unpack tokenOfOwnerByIndex result")
	ErrEmptyTokenOfOwnerResponse = errors.New("empty response from tokenOfOwnerByIndex")
)
