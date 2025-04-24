package utils

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func GetFromAddressTx(tx *types.Transaction) string {
	var signer types.Signer
	chainID := tx.ChainId()

	switch tx.Type() {
	case types.LegacyTxType:
		if chainID == nil || chainID.Sign() == 0 {
			signer = types.HomesteadSigner{}
		} else {
			signer = types.NewEIP155Signer(chainID)
		}
	case types.AccessListTxType, types.DynamicFeeTxType:
		signer = types.NewLondonSigner(chainID)
	default:
		signer = types.LatestSignerForChainID(chainID)
	}

	fromAddress, err := types.Sender(signer, tx)
	if err != nil {
		return ""
	}

	return fromAddress.Hex()
}

func GetToAddressTx(tx *types.Transaction) string {

	if tx.To() == nil {
		return ""
	}

	return tx.To().Hex()
}
