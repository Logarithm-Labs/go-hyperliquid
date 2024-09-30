package hyperliquid

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/vmihailenco/msgpack/v5"
)

// SignRequest is the implementation of EIP-712 typed data
type SignRequest struct {
	PrimaryType string
	DType       []apitypes.Type
	DTypeMsg    map[string]interface{}
	IsMainNet   bool
	DomainName  string
}

func (request *SignRequest) getChainId() *math.HexOrDecimal256 {
	if request.DomainName == "HyperliquidSignTransaction" {
		if request.IsMainNet {
			return math.NewHexOrDecimal256(int64(ARBITRUM_CHAIN_ID))
		}
		return math.NewHexOrDecimal256(int64(ARBITRUM_TESTNET_CHAIN_ID))
	}
	return math.NewHexOrDecimal256(int64(HYPERLIQUID_CHAIN_ID))
}

func (request *SignRequest) GetTypes() apitypes.Types {
	types := apitypes.Types{
		request.PrimaryType: request.DType,
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
	}
	return types
}

func (request *SignRequest) GetDomain() apitypes.TypedDataDomain {
	return apitypes.TypedDataDomain{
		Name:              request.DomainName,
		Version:           "1",
		ChainId:           request.getChainId(),
		VerifyingContract: VERIFYING_CONTRACT,
	}
}

type Signer struct {
	manager *PKeyManager
}

func NewSigner(manager *PKeyManager) Signer {
	return Signer{
		manager: manager,
	}
}

func (signer *Signer) Sign(request *SignRequest) (byte, [32]byte, [32]byte, error) {
	return signer.signInternal(SignRequestToEIP712TypedData(request))
}

// signInternal signs the typed data and returns the signature in VRS format
func (signer *Signer) signInternal(message apitypes.TypedData) (byte, [32]byte, [32]byte, error) {
	pkey := signer.manager.PrivateECDSA()
	bytes, _, err := apitypes.TypedDataAndHash(message)
	if err != nil {
		log.Printf("Error hashing typed data: %s", err)
		return 0, [32]byte{}, [32]byte{}, err
	}
	signature, err := crypto.Sign(bytes, pkey)
	if err != nil {
		log.Printf("Error signing typed data: %s", err)
		return 0, [32]byte{}, [32]byte{}, err
	}
	return SignatureToVRS(signature)
}

func SignRequestToEIP712TypedData(request *SignRequest) apitypes.TypedData {
	return apitypes.TypedData{
		Domain:      request.GetDomain(),
		Types:       request.GetTypes(),
		PrimaryType: request.PrimaryType,
		Message:     request.DTypeMsg,
	}
}

func SignatureToVRS(sig []byte) (byte, [32]byte, [32]byte, error) {
	var v byte
	var r [32]byte
	var s [32]byte
	v = sig[64] + 27
	copy(r[:], sig[:32])
	copy(s[:], sig[32:64])
	return v, r, s, nil
}

// Create a hash of an action (json object)
func buildActionHash(action any, vaultAd string, nonce uint64) (common.Hash, error) {
	data, err := msgpack.Marshal(action)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error while marshaling action: %s", err)
	}
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(nonce))
	data = ArrayAppend(data, nonceBytes)

	if vaultAd == "" {
		data = ArrayAppend(data, []byte("\x00"))
	} else {
		data = ArrayAppend(data, []byte("\x01"))
		data = ArrayAppend(data, HexToBytes(vaultAd))
	}
	result := crypto.Keccak256Hash(data)
	return result, nil
}

func getNetSource(isMainnet bool) string {
	if isMainnet {
		return "a"
	} else {
		return "b"
	}
}

// Build a message to sign
func buildMessage(hash []byte, isMainnet bool) apitypes.TypedDataMessage {
	source := getNetSource(isMainnet)
	return apitypes.TypedDataMessage{
		"source":       source,
		"connectionId": hash,
	}
}
