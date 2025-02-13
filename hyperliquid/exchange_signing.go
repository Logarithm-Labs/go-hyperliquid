package hyperliquid

import (
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func (api *ExchangeAPI) Sign(request *SignRequest) (byte, [32]byte, [32]byte, error) {
	signer := NewSigner(api.keyManager)
	v, r, s, err := signer.Sign(request)
	if err != nil {
		api.debug("Error SignInner: %s", err)
		return 0, [32]byte{}, [32]byte{}, err
	}
	return v, r, s, nil
}

func (api *ExchangeAPI) SignUserSignableAction(action any, payloadTypes []apitypes.Type, primaryType string) (byte, [32]byte, [32]byte, error) {
	message, err := StructToMap(action)
	if err != nil {
		return 0, [32]byte{}, [32]byte{}, err
	}
	// Remove unnecessary fields for signing
	delete(message, "type")
	delete(message, "signatureChainId")

	signRequest := &SignRequest{
		DomainName:  "HyperliquidSignTransaction",
		PrimaryType: primaryType,
		DType:       payloadTypes,
		DTypeMsg:    message,
		IsMainNet:   api.IsMainnet(),
	}
	return api.Sign(signRequest)
}

func (api *ExchangeAPI) SignL1Action(action any, timestamp uint64) (byte, [32]byte, [32]byte, error) {
	srequest, err := api.BuildEIP712Message(action, timestamp, api.VaultAddress())
	if err != nil {
		api.debug("Error building EIP712 message: %s", err)
		return 0, [32]byte{}, [32]byte{}, err
	}
	return api.Sign(srequest)
}

func (api *ExchangeAPI) BuildEIP712Message(action any, timestamp uint64, vaultAddress string) (*SignRequest, error) {
	hash, err := buildActionHash(action, vaultAddress, timestamp)
	if err != nil {
		return nil, err
	}
	message := buildMessage(hash.Bytes(), api.IsMainnet())
	srequest := &SignRequest{
		DomainName:  "Exchange",
		PrimaryType: "Agent",
		DType: []apitypes.Type{
			{
				Name: "source",
				Type: "string",
			},
			{
				Name: "connectionId",
				Type: "bytes32",
			},
		},
		DTypeMsg:  message,
		IsMainNet: api.IsMainnet(),
	}
	return srequest, nil
}

func (api *ExchangeAPI) SignWithdrawAction(action WithdrawAction) (byte, [32]byte, [32]byte, error) {
	types := []apitypes.Type{
		{
			Name: "hyperliquidChain",
			Type: "string",
		},
		{
			Name: "destination",
			Type: "string",
		},
		{
			Name: "amount",
			Type: "string",
		},
		{
			Name: "time",
			Type: "uint64",
		},
	}
	return api.SignUserSignableAction(action, types, "HyperliquidTransaction:Withdraw")
}
