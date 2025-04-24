package main

import (
	"fmt"
	"log"

	"github.com/dtome123/go-bcwe3/eth/verifier"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func main() {

	// legacy
	message := "hello"
	// Địa chỉ người ký (để xác minh)
	address := "0x7556989c2A60E60F0c66A2b9D77079BC9F189037"

	// Chữ ký hex từ client (personal_sign)
	sigHex := "0xbd51be0700eb411813268a5dc6e893fc8aa2326e82e9c405b30cd9acc65881b31f9a5b2aed01f854bd8c12a8cbfcceecb1175aedb18b7540bfc5c1dc4eb21ac71b" // 65 bytes (130 hex chars)

	err := verifier.Verify(&verifier.VerifyRequest{
		SignatureType: verifier.SignaturePersonalSign,
		Payload:       message,
		Signature:     sigHex,
		ExpectedAddr:  address,
	})

	if err != nil {
		log.Println("Verify error:", err)
	} else {
		log.Println("✅ Signature is valid")
	}

	// Địa chỉ người ký (để xác minh)

	// EIP712
	sigHex = "0xd1059fd1ce46b05c3f0287ba60f6eb7e3e8d10d70fd89acf7e93e873de537b760efd7caf781803e947c07f90f543e200f0e8e044b7ec5f2e75a1f7457df27ae31b" // 65 bytes (130 hex chars)

	typedData := &apitypes.TypedData{
		Types: apitypes.Types{
			"Person": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "wallet", Type: "address"},
			},
		},
		PrimaryType: "Person",
		Domain: apitypes.TypedDataDomain{
			Name:              "MyDApp",
			Version:           "1",
			ChainId:           math.NewHexOrDecimal256(17000), // ví dụ chainId
			VerifyingContract: "0x0000000000000000000000000000000000000000",
		},
		Message: apitypes.TypedDataMessage{
			"name":   address,
			"wallet": address,
		},
	}

	err = verifier.Verify(&verifier.VerifyRequest{
		SignatureType: verifier.SignatureTypedData,
		Payload:       typedData,
		Signature:     sigHex,
		ExpectedAddr:  address,
	})
	if err != nil {
		fmt.Println("❌ Verify failed:", err)
	} else {
		log.Println("✅ Signature is valid")
	}

}
