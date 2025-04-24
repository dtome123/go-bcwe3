package verifier

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var (
	ErrInvalidSignature       = errors.New("invalid signature")
	ErrUnsupportedType        = errors.New("unsupported signature type")
	ErrInvalidSignatureLength = errors.New("invalid signature length")
	ErrRecoverPublicKey       = errors.New("failed to recover public key")
	ErrInvalidPayloadType     = errors.New("invalid payload type")
	ErrPayloadIsEmpty         = errors.New("payload is empty")
)

type SignatureType int

const (
	SignaturePersonalSign SignatureType = iota
	SignatureTypedData
)

type VerifyRequest struct {
	SignatureType SignatureType
	Payload       any
	Signature     string
	ExpectedAddr  string
}

func Verify(req *VerifyRequest) error {
	switch req.SignatureType {
	case SignaturePersonalSign:
		if msg, ok := req.Payload.(string); ok {
			return verifyLegacySignature(msg, req.Signature, req.ExpectedAddr)
		}

		return ErrInvalidPayloadType
	case SignatureTypedData:

		switch v := any(req.Payload).(type) {
		case apitypes.TypedData:
			return verifyEIP712Signature(v, req.Signature, req.ExpectedAddr)
		case *apitypes.TypedData:
			if v == nil {
				return ErrPayloadIsEmpty
			}
			return verifyEIP712Signature(*v, req.Signature, req.ExpectedAddr)
		default:
			return ErrInvalidPayloadType
		}

	default:
		return ErrUnsupportedType
	}
}

func verifyLegacySignature(message string, sigHex string, address string) error {
	// Decode hex signature
	sigHex = strings.TrimPrefix(sigHex, "0x")
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		return ErrInvalidSignature
	}

	if len(sig) != 65 {
		return ErrInvalidSignatureLength
	}

	// Adjust recovery ID if necessary
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	// Hash message with Ethereum prefix
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))
	msgHash := crypto.Keccak256([]byte(prefix + message))

	// Recover public key
	pubKey, err := crypto.SigToPub(msgHash, sig)
	if err != nil {
		return ErrRecoverPublicKey
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	normalizedExpected := strings.ToLower(common.HexToAddress(address).Hex())
	normalizedActual := strings.ToLower(recoveredAddr.Hex())

	if normalizedExpected != normalizedActual {
		return ErrInvalidSignature
	}

	return nil
}

func verifyEIP712Signature(typedData apitypes.TypedData, sigHex string, expectedAddrHex string) error {

	typedData.Types["EIP712Domain"] = []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return err
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return err
	}

	// EIP-712 hash
	data := crypto.Keccak256(
		[]byte("\x19\x01"),
		domainSeparator,
		typedDataHash,
	)

	sig, err := hex.DecodeString(strings.TrimPrefix(sigHex, "0x"))
	if err != nil {
		return err
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	pubKey, err := crypto.SigToPub(data, sig)
	if err != nil {
		return err
	}

	if crypto.PubkeyToAddress(*pubKey).Hex() != expectedAddrHex {
		return ErrInvalidSignature
	}

	return nil
}
