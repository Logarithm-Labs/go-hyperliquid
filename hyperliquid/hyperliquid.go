package hyperliquid

type IHyperliquid interface {
	IExchangeAPI
	IInfoAPI
}

type Hyperliquid struct {
	ExchangeAPI
	InfoAPI
}

// HyperliquidClientConfig represents the configuration options for the Hyperliquid client.
// It allows configuring the network type, private key, and account address settings.
//
// The configuration options include:
//   - IsMainnet: Set to true for mainnet and false for testnet
//   - PrivateKey: Optional key for authenticated endpoints (can be empty for public endpoints)
//   - AccountAddress: Default account address used by the API (modifiable via SetAccountAddress)
type HyperliquidClientConfig struct {
	IsMainnet      bool
	PrivateKey     string
	AccountAddress string
}

// NewHyperliquid creates a new Hyperliquid API client.
func NewHyperliquid(config *HyperliquidClientConfig) *Hyperliquid {
	var defaultConfig *HyperliquidClientConfig
	if config == nil {
		defaultConfig = &HyperliquidClientConfig{
			IsMainnet:      true,
			PrivateKey:     "",
			AccountAddress: "",
		}
	} else {
		defaultConfig = config
	}
	exchangeAPI := NewExchangeAPI(defaultConfig.IsMainnet)
	exchangeAPI.SetPrivateKey(defaultConfig.PrivateKey)
	exchangeAPI.SetAccountAddress(defaultConfig.AccountAddress)
	infoAPI := NewInfoAPI(defaultConfig.IsMainnet)
	infoAPI.SetAccountAddress(defaultConfig.AccountAddress)
	hl := &Hyperliquid{
		ExchangeAPI: *exchangeAPI,
		InfoAPI:     *infoAPI,
	}
	hl.UpdateAccountVaultAddress()
	return hl
}

func (h *Hyperliquid) SetDebugActive() {
	h.ExchangeAPI.SetDebugActive()
	h.InfoAPI.SetDebugActive()
}

func (h *Hyperliquid) SetPrivateKey(privateKey string) error {
	err := h.ExchangeAPI.SetPrivateKey(privateKey)
	if err != nil {
		return err
	}
	return nil
}

func (h *Hyperliquid) SetAccountAddress(accountAddress string) {
	h.ExchangeAPI.SetAccountAddress(accountAddress)
	h.InfoAPI.SetAccountAddress(accountAddress)
}

func (h *Hyperliquid) UpdateAccountVaultAddress() {
	h.UpdateVaultAddress(h.AccountAddress())
}

func (h *Hyperliquid) UpdateVaultAddress(address string) error {
	role, err := h.InfoAPI.GetUserRole(address)
	if err != nil {
		return err
	}
	h.ExchangeAPI.SetUserRole(role.UserRole.Role)
	return nil
}
func (h *Hyperliquid) AccountAddress() string {
	return h.ExchangeAPI.AccountAddress()
}

func (h *Hyperliquid) IsMainnet() bool {
	return h.ExchangeAPI.IsMainnet()
}
