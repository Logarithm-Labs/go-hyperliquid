package hyperliquid

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

// IInfoAPI is an interface for the /info service.
type IInfoAPI interface {
	IClient // Base client interface

	// INFO API ENDPOINTS
	GetAllMids() (*map[string]string, error)
	GetOpenOrders(address string) (*[]Order, error)
	GetAccountOpenOrders() (*[]Order, error)
	GetUserFills(address string) (*[]OrderFill, error)
	GetAccountFills() (*[]OrderFill, error)
	GetUserRateLimits(address string) (*float64, error)
	GetL2BookSnapshot(coin string) (*L2BookSnapshot, error)
	GetCandleSnapshot(coin string, interval string, startTime int64, endTime int64) (*CandleSnapshot, error)

	// PERPETUALS INFO API ENDPOINTS
	GetMeta() (*Meta, error)
	GetUserState(address string) (*UserState, error)
	GetAccountState() (*UserState, error)
	GetFundingUpdates(address string, startTime int64, endTime int64) (*[]FundingUpdate, error)
	GetAccountFundingUpdates(startTime int64, endTime int64) (*[]FundingUpdate, error)
	GetNonFundingUpdates(address string, startTime int64, endTime int64) (*[]NonFundingUpdate, error)
	GetAccountNonFundingUpdates(startTime int64, endTime int64) (*[]NonFundingUpdate, error)
	GetHistoricalFundingRates() (*[]HistoricalFundingRate, error)

	// Additional helper functions
	GetMartketPx(coin string) (float64, error)
	BuildMetaMap() (map[string]AssetInfo, error)
	GetWithdrawals(address string) (*[]Withdrawal, error)
	GetAccountWithdrawals() (*[]Withdrawal, error)
	GetUserRole() (*UserRole, error)
}

type InfoAPI struct {
	Client
	baseEndpoint string
	spotMeta     map[string]AssetInfo
}

// NewInfoAPI returns a new instance of the InfoAPI struct.
// It sets the base endpoint to "/info" and the client to the NewClient function.
// The isMainnet parameter is used to set the network type.
func NewInfoAPI(isMainnet bool) *InfoAPI {
	api := InfoAPI{
		baseEndpoint: "/info",
		Client:       *NewClient(isMainnet),
	}
	spotMeta, err := api.BuildSpotMetaMap()
	if err != nil {
		api.SetDebugActive()
		api.debug("Error building meta map: %s", err)
	}
	api.spotMeta = spotMeta
	return &api
}

// Endpoint returns the base endpoint for the InfoAPI.
func (api *InfoAPI) Endpoint() string {
	return api.baseEndpoint
}

// Retrieve mids for all actively traded coins
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint#retrieve-mids-for-all-actively-traded-coins
func (api *InfoAPI) GetAllMids() (*map[string]string, error) {
	request := InfoRequest{
		Type: "allMids",
	}
	return MakeUniversalRequest[map[string]string](api, request)
}

// Retrieve spot meta and asset contexts
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint/spot#retrieve-spot-asset-contexts
func (api *InfoAPI) GetAllSpotPrices() (*map[string]string, error) {
	request := InfoRequest{
		Type: "spotMetaAndAssetCtxs",
	}
	response, err := MakeUniversalRequest[SpotMetaAndAssetCtxsResponse](api, request)
	if err != nil {
		return nil, err
	}

	marketsData, ok := response[1].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid markets data format")
	}

	result := make(map[string]string)

	marketBytes, err := json.Marshal(marketsData)
	if err != nil {
		return nil, err
	}

	var markets []Market
	if err := json.Unmarshal(marketBytes, &markets); err != nil {
		return nil, err
	}

	for _, market := range markets {
		result[market.Coin] = market.MidPx
	}

	return &result, nil
}

// Retrieve a user's open orders
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint#retrieve-a-users-open-orders
func (api *InfoAPI) GetOpenOrders(address string) (*[]Order, error) {
	request := InfoRequest{
		User: address,
		Type: "openOrders",
	}
	return MakeUniversalRequest[[]Order](api, request)
}

// Retrieve a account's order history
// The same as GetOpenOrders but user is set to the account address
// Check AccountAddress() or SetAccountAddress() if there is a need to set the account address
func (api *InfoAPI) GetAccountOpenOrders() (*[]Order, error) {
	return api.GetOpenOrders(api.AccountAddress())
}

// Retrieve a user's fills
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint#retrieve-a-users-fills
func (api *InfoAPI) GetUserFills(address string) (*[]OrderFill, error) {
	request := InfoRequest{
		User: address,
		Type: "userFills",
	}
	return MakeUniversalRequest[[]OrderFill](api, request)
}

// Retrieve a account's fill history
// The same as GetUserFills but user is set to the account address
// Check AccountAddress() or SetAccountAddress() if there is a need to set the account address
func (api *InfoAPI) GetAccountFills() (*[]OrderFill, error) {
	return api.GetUserFills(api.AccountAddress())
}

// Query user rate limits
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint#query-user-rate-limits
func (api *InfoAPI) GetUserRateLimits(address string) (*RatesLimits, error) {
	request := InfoRequest{
		User: address,
		Type: "userRateLimit",
	}
	return MakeUniversalRequest[RatesLimits](api, request)
}

// Query account rate limits
// The same as GetUserRateLimits but user is set to the account address
// Check AccountAddress() or SetAccountAddress() if there is a need to set the account address
func (api *InfoAPI) GetAccountRateLimits() (*RatesLimits, error) {
	return api.GetUserRateLimits(api.AccountAddress())
}

// L2 Book snapshot
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint#l2-book-snapshot
func (api *InfoAPI) GetL2BookSnapshot(coin string) (*L2BookSnapshot, error) {
	request := InfoRequest{
		Type: "l2Book",
		Coin: coin,
	}
	return MakeUniversalRequest[L2BookSnapshot](api, request)
}

// Candle snapshot (Only the most recent 5000 candles are available)
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint#candle-snapshot
func (api *InfoAPI) GetCandleSnapshot(coin string, interval string, startTime int64, endTime int64) (*[]CandleSnapshot, error) {
	request := CandleSnapshotRequest{
		Type: "candleSnapshot",
		Req: CandleSnapshotSubRequest{
			Coin:      coin,
			Interval:  interval,
			StartTime: startTime,
			EndTime:   endTime,
		},
	}
	return MakeUniversalRequest[[]CandleSnapshot](api, request)
}

// Retrieve perpetuals metadata
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint/perpetuals#retrieve-perpetuals-metadata
func (api *InfoAPI) GetMeta() (*Meta, error) {
	request := InfoRequest{
		Type: "meta",
	}
	return MakeUniversalRequest[Meta](api, request)
}

// Retrieve spot metadata
func (api *InfoAPI) GetSpotMeta() (*SpotMeta, error) {
	request := InfoRequest{
		Type: "spotMeta",
	}
	return MakeUniversalRequest[SpotMeta](api, request)
}

// Retrieve user's perpetuals account summary
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint/perpetuals#retrieve-users-perpetuals-account-summary
func (api *InfoAPI) GetUserState(address string) (*UserState, error) {
	request := UserStateRequest{
		User: address,
		Type: "clearinghouseState",
	}
	return MakeUniversalRequest[UserState](api, request)
}

// GetAccountState retrieve account's perpetuals account summary
// The same as GetUserState but user is set to the account address
// Check AccountAddress() or SetAccountAddress() if there is a need to set the account address
func (api *InfoAPI) GetAccountState() (*UserState, error) {
	return api.GetUserState(api.AccountAddress())
}

// GetUserStateSpot retrieve's a user's spot account summary
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint/spot#retrieve-a-users-token-balances
func (api *InfoAPI) GetUserStateSpot(address string) (*UserStateSpot, error) {
	request := UserStateRequest{
		User: address,
		Type: "spotClearinghouseState",
	}
	return MakeUniversalRequest[UserStateSpot](api, request)
}

// Retrieve account's spot account summary
// The same as GetUserStateSpot but user is set to the account address
// Check AccountAddress() or SetAccountAddress() if there is a need to set the account address
func (api *InfoAPI) GetAccountStateSpot() (*UserStateSpot, error) {
	return api.GetUserStateSpot(api.AccountAddress())
}

// GetFundingUpdates retrieves user's funding history between startTime and endTime (Unix timestamps).
// Returns chronological funding payments for perpetual positions. See API docs for details.
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint/perpetuals#retrieve-funding-updates
func (api *InfoAPI) GetFundingUpdates(address string, startTime int64, endTime int64) (*[]FundingUpdate, error) {
	request := InfoRequest{
		User:      address,
		Type:      "userFunding",
		StartTime: startTime,
		EndTime:   endTime,
	}
	return MakeUniversalRequest[[]FundingUpdate](api, request)
}

// GetAccountFundingUpdates retrieve's an account's funding history
// The same as GetFundingUpdates but user is set to the account address
// Check AccountAddress() or SetAccountAddress() if there is a need to set the account address
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint/perpetuals#retrieve-a-users-funding-history-or-non-funding-ledger-updates
func (api *InfoAPI) GetAccountFundingUpdates(startTime int64, endTime int64) (*[]FundingUpdate, error) {
	return api.GetFundingUpdates(api.AccountAddress(), startTime, endTime)
}

// Retrieve a user's funding history or non-funding ledger updates
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint/perpetuals#retrieve-a-users-funding-history-or-non-funding-ledger-updates
func (api *InfoAPI) GetNonFundingUpdates(address string, startTime int64, endTime int64) (*[]NonFundingUpdate, error) {
	request := InfoRequest{
		User:      address,
		Type:      "userNonFundingLedgerUpdates",
		StartTime: startTime,
		EndTime:   endTime,
	}
	return MakeUniversalRequest[[]NonFundingUpdate](api, request)
}

// Retrieve account's funding history or non-funding ledger updates
// The same as GetNonFundingUpdates but user is set to the account address
// Check AccountAddress() or SetAccountAddress() if there is a need to set the account address
func (api *InfoAPI) GetAccountNonFundingUpdates(startTime int64, endTime int64) (*[]NonFundingUpdate, error) {
	return api.GetNonFundingUpdates(api.AccountAddress(), startTime, endTime)
}

// Retrieve a user's role ("missing", "user", "agent", "vault", or "subAccount")
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint#query-a-users-role
func (api *InfoAPI) GetUserRole(address string) (*UserRoleResponse, error) {
	request := InfoRequest{
		User: address,
		Type: "userRole",
	}
	return MakeUniversalRequest[UserRoleResponse](api, request)
}

// Retrieve historical funding rates
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint/perpetuals#retrieve-historical-funding-rates
func (api *InfoAPI) GetHistoricalFundingRates(coin string, startTime int64, endTime int64) (*[]HistoricalFundingRate, error) {
	request := InfoRequest{
		Type:      "fundingHistory",
		Coin:      coin,
		StartTime: startTime,
		EndTime:   endTime,
	}
	return MakeUniversalRequest[[]HistoricalFundingRate](api, request)
}

// Helper function to get the market price of a given coin
// The coin parameter is the name of the coin
//
// Example:
//
//	api.GetMartketPx("BTC")
func (api *InfoAPI) GetMartketPx(coin string) (float64, error) {
	allMids, err := api.GetAllMids()
	if err != nil {
		return 0, err
	}
	parsed, err := strconv.ParseFloat((*allMids)[coin], 32)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

// Helper function to get the withdrawals of a given address
// By default returns last 90 days
func (api *InfoAPI) GetWithdrawals(address string) (*[]Withdrawal, error) {
	startTime, endTime := GetDefaultTimeRange()
	updates, err := api.GetNonFundingUpdates(address, startTime, endTime)
	if err != nil {
		return nil, err
	}
	var withdrawals []Withdrawal
	for _, update := range *updates {
		if update.Delta.Type == "withdraw" {
			withrawal := Withdrawal{
				Time:   update.Time,
				Hash:   update.Hash,
				Amount: update.Delta.Usdc,
				Fee:    update.Delta.Fee,
				Nonce:  update.Delta.Nonce,
			}
			withdrawals = append(withdrawals, withrawal)
		}
	}
	return &withdrawals, nil
}

// Helper function to get the withdrawals of the account address
// The same as GetWithdrawals but user is set to the account address
// Check AccountAddress() or SetAccountAddress() if there is a need to set the account address
func (api *InfoAPI) GetAccountWithdrawals() (*[]Withdrawal, error) {
	return api.GetWithdrawals(api.AccountAddress())
}

// Helper function to get the deposits of the given address
// By default returns last 90 days
func (api *InfoAPI) GetDeposits(address string) (*[]Deposit, error) {
	startTime, endTime := GetDefaultTimeRange()
	updates, err := api.GetNonFundingUpdates(address, startTime, endTime)
	if err != nil {
		return nil, err
	}
	var deposits []Deposit
	for _, update := range *updates {
		if update.Delta.Type == "deposit" {
			deposit := Deposit{
				Hash:   update.Hash,
				Amount: update.Delta.Usdc,
				Time:   update.Time,
			}
			deposits = append(deposits, deposit)
		}
	}
	return &deposits, nil
}

// Helper function to get the deposits of the account address
// The same as GetDeposits but user is set to the account address
// Check AccountAddress() or SetAccountAddress() if there is a need to set the account address
func (api *InfoAPI) GetAccountDeposits() (*[]Deposit, error) {
	return api.GetDeposits(api.AccountAddress())
}

// Helper function to build a map of asset names to asset info
// It is used to get the assetId for a given asset name
func (api *InfoAPI) BuildMetaMap() (map[string]AssetInfo, error) {
	metaMap := make(map[string]AssetInfo)
	result, err := api.GetMeta()
	if err != nil {
		log.Fatalf("Failed to get meta: %v", err)
	}
	for index, asset := range result.Universe {
		if asset.Name == "BTC" {
			metaMap["BTC"] = AssetInfo{
				SzDecimals: asset.SzDecimals,
				AssetID:    index,
			}
		}
		metaMap[asset.Name] = AssetInfo{
			SzDecimals: asset.SzDecimals,
			AssetID:    index,
		}
	}
	return metaMap, nil
}

// Helper function to build a map of asset names to asset info
// It is used to get the assetId for a given asset name
func (api *InfoAPI) BuildSpotMetaMap() (map[string]AssetInfo, error) {
	spotMeta, err := api.GetSpotMeta()
	if err != nil {
		return nil, err
	}

	tokenMap := make(map[int]struct {
		name        string
		szDecimals  int
		weiDecimals int
	}, len(spotMeta.Tokens))

	for _, token := range spotMeta.Tokens {
		tokenMap[token.Index] = struct {
			name        string
			szDecimals  int
			weiDecimals int
		}{token.Name, token.SzDecimals, token.WeiDecimals}
	}

	metaMap := make(map[string]AssetInfo)
	for _, universe := range spotMeta.Universe {
		for _, tokenId := range universe.Tokens {
			if tokenId == 0 {
				continue
			}
			if token, exists := tokenMap[tokenId]; exists {
				metaMap[token.name] = AssetInfo{
					SzDecimals:  token.szDecimals,
					WeiDecimals: token.weiDecimals,
					AssetID:     universe.Index,
					SpotName:    universe.Name,
				}
			}
		}
	}
	return metaMap, nil
}
