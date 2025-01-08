package hyperliquid

import (
	"os"
	"testing"
)

func GetInfoAPI() *InfoAPI {
	api := NewInfoAPI(true)
	if GLOBAL_DEBUG {
		api.SetDebugActive()
	}
	// It should be active account to pass all tests
	// like GetAccountFills, GetAccountWithdrawals, etc.
	TEST_ADDRESS := os.Getenv("TEST_ADDRESS")
	if TEST_ADDRESS == "" {
		panic("Set TEST_ADDRESS in .env file")
	}
	api.SetAccountAddress(TEST_ADDRESS)
	return api
}

func TestInfoAPI_AccountAddress(t *testing.T) {
	api := GetInfoAPI()
	address := api.AccountAddress()
	targetAddress := os.Getenv("TEST_ADDRESS")
	if targetAddress == "" {
		t.Errorf("Set TEST_ADDRESS in .env file")
	}
	if address != targetAddress {
		t.Errorf("AccountAddress() = %v, want %v", address, targetAddress)
	}
}

func TestInfoAPI_Endpoint(t *testing.T) {
	api := GetInfoAPI()
	res := api.Endpoint()
	if res != "/info" {
		t.Errorf("Endpoint() = %v, want %v", res, "/info")
	}
}

func TestInfoAPI_GetAllMids(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetAllMids()
	if err != nil {
		t.Errorf("GetAllMids() error = %v", err)
	}

	// Check BTC and ETH are in the map
	if _, ok := (*res)["BTC"]; !ok {
		t.Errorf("GetAllMids() doesnt return %v, want %v", res, "BTC")
	}
	if _, ok := (*res)["ETH"]; !ok {
		t.Errorf("GetAllMids() doesnt return %v, want %v", res, "ETH")
	}
	t.Logf("GetAllMids() = %v", res)
}

func TestInfoAPI_GetAccountFills(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetAccountFills()
	if err != nil {
		t.Errorf("GetAccountFills() error = %v", err)
	}
	if len(*res) == 0 {
		t.Errorf("GetAccountFills() len = %v, want > %v", res, 0)
	}
	res0 := (*res)[0]
	t.Logf("res0 = %+v", res0)
	if res0.Px == 0 {
		t.Errorf("res0.Px = %v, want > %v", res0.Px, 0)
	}
	if res0.Sz == 0 {
		t.Errorf("res0.Sz = %v, want > %v", res0.Sz, 0)
	}
	if res0.Fee == 0 {
		t.Errorf("res0.Fee = %v, want > %v", res0.Fee, 0)
	}
	t.Logf("GetAccountFills() = %v", res)
}

func TestInfoAPI_GetAccountRateLimits(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetAccountRateLimits()
	if err != nil {
		t.Errorf("GetAccountRateLimits() error = %v", err)
	}
	if res.CumVlm == 0 {
		t.Errorf("GetAccountRateLimits() len = %v, want > %v", res.CumVlm, 0)
	}

	t.Logf("GetAccountRateLimits() = %v", res)
}

func TestInfoAPI_GetL2BookSnapshot(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetL2BookSnapshot("BTC")
	if err != nil {
		t.Errorf("GetL2BookSnapshot() error = %v", err)
	}
	if res.Levels[0][0].Px <= 0 {
		t.Errorf("res.Levels[0][0].Px = %v, want > %v", res.Levels[0][0].Px, 0)
	}
	t.Logf("GetL2BookSnapshot() = %v", res)
}

func TestInfoAPI_GetCandleSnapshot(t *testing.T) {
	api := GetInfoAPI()
	startTime, endTime := GetDefaultTimeRange()
	res, err := api.GetCandleSnapshot("ETH", "1d", startTime, endTime)
	if err != nil {
		t.Errorf("GetCandleSnapshot() error = %v", err)
	}
	if len(*res) == 0 {
		t.Errorf("GetCandleSnapshot() len = %v, want > %v", res, 0)
	}
	if (*res)[0].Open <= 0 {
		t.Errorf("*res)[0].Open  = %v, want > %v", (*res)[0].Open, 0)
	}
	t.Logf("GetCandleSnapshot() = %v", res)
}

func TestInfoAPI_GetMeta(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetMeta()
	if err != nil {
		t.Errorf("GetMeta() error = %v", err)
	}
	t.Logf("GetMeta() = %v", res)
	if res.Universe[0].Name != "BTC" {
		t.Errorf("GetMeta() doesnt return %v, want %v", res.Universe[0].Name, "BTC")
	}
}

func TestInfoAPI_GetUserState(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetAccountState()
	if err != nil {
		t.Errorf("GetUserState() error = %v", err)
	}
	if res.Withdrawable == 0 {
		t.Errorf("GetUserState.Withdrawable = %v, want > %v", res.Withdrawable, 0)
	}
	if res.CrossMarginSummary.AccountValue == 0 {
		t.Errorf("GetUserState.AccountValue = %v, want > %v", res.CrossMarginSummary.AccountValue, 0)
	}
	t.Logf("GetUserState() = %v", res)
}

func TestInfoAPI_GetAccountOpenOrders(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetAccountOpenOrders()
	if err != nil {
		t.Errorf("GetAccountOpenOrders() error = %v", err)
	}
	if len(*res) == 0 {
		t.Errorf("GetAccountOpenOrders() len = %v, want > %v", res, 0)
	}
	t.Logf("GetAccountOpenOrders() = %v", res)
}

func TestInfoAPI_GetAccountFundingUpdates(t *testing.T) {
	api := GetInfoAPI()
	startTime, endTime := GetDefaultTimeRange()
	res, err := api.GetAccountFundingUpdates(startTime, endTime)
	if err != nil {
		t.Errorf("GetAccountFundingUpdates() error = %v", err)
	}
	if len(*res) == 0 {
		t.Errorf("GetAccountFundingUpdates() len = %v, want > %v", res, 0)
	}
	t.Logf("GetAccountFundingUpdates() = %v", res)
}

func TestInfoAPI_GetHistoricalFundingRates(t *testing.T) {
	api := GetInfoAPI()
	startTime, endTime := GetDefaultTimeRange()
	res, err := api.GetHistoricalFundingRates("BTC", startTime, endTime)
	if err != nil {
		t.Errorf("GetHistoricalFundingRates() error = %v", err)
	}
	if len(*res) == 0 {
		t.Errorf("GetHistoricalFundingRates()  len = %v, want > %v", res, 0)
	}
	t.Logf("GetHistoricalFundingRates() = %v", res)
}

func TestInfoAPI_GetAccountNonFundingUpdates(t *testing.T) {
	api := GetInfoAPI()
	startTime, endTime := GetDefaultTimeRange()
	res, err := api.GetAccountNonFundingUpdates(startTime, endTime)
	if err != nil {
		t.Errorf("GetAccountNonFundingUpdates() error = %v", err)
	}
	if len(*res) == 0 {
		t.Errorf("GetAccountNonFundingUpdates() len = %v, want > %v", res, 0)
	}
	// find first deposit
	for _, update := range *res {
		if update.Delta.Type == "deposit" {
			// check that usdc is in the deposit
			if update.Delta.Usdc == 0 {
				t.Errorf("update.Delta.Usdc = %v, want > %v", update.Delta.Amount, 0)
			}
		}
		if update.Delta.Type == "withdrawal" {
			if update.Delta.Usdc == 0 {
				t.Errorf("update.Delta.Usdc = %v, want > %v", update.Delta.Amount, 0)
			}
			if update.Delta.Nonce == 0 {
				t.Errorf("update.Delta.Nonce = %v, want > %v", update.Delta.Nonce, 0)
			}
			if update.Delta.Fee == 0 {
				t.Errorf("update.Delta.Fee = %v, want > %v", update.Delta.Fee, 0)
			}
		}
		if update.Delta.Type == "spotGenesis" {
			if update.Delta.Token == "" {
				t.Errorf("update.Delta.Token = %v", update.Delta.Amount)
			}
			if update.Delta.Amount == 0 {
				t.Errorf("update.Delta.Amount = %v, want > %v", update.Delta.Amount, 0)
			}
		}
		if update.Delta.Type == "accountClassTransfer" {
			if update.Delta.Usdc == 0 {
				t.Errorf("update.Delta.Usdc = %v, want > %v", update.Delta.Amount, 0)
			}
		}
	}
	t.Logf("GetAccountNonFundingUpdates() = %v", res)
}

func TestInfoAPI_GetAccountWithdrawals(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetAccountWithdrawals()
	if err != nil {
		t.Errorf("GetAccountWithdrawals() error = %v", err)
	}
	if len(*res) == 0 {
		t.Errorf("GetAccountWithdrawals() len = %v, want > %v", res, 0)
	}
	for _, withdrawal := range *res {
		if withdrawal.Amount == 0 {
			t.Errorf("withdrawal.Amount = %v, want > %v", withdrawal.Amount, 0)
		}
	}
	t.Logf("GetAccountWithdrawals() = %v", res)
}

func TestInfoAPI_GetAccountDeposits(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetAccountDeposits()
	if err != nil {
		t.Errorf("GetAccountDeposits() error = %v", err)
	}
	if len(*res) == 0 {
		t.Errorf("GetAccountDeposits() len = %v, want > %v", res, 0)
	}
	for _, deposit := range *res {
		if deposit.Amount == 0 {
			t.Errorf("deposit.Amount = %v, want > %v", deposit.Amount, 0)
		}
	}
	t.Logf("GetAccountDeposits() = %v", res)
}

func TestInfoAPI_GetMarketPx(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetMartketPx("BTC")
	if err != nil {
		t.Errorf("GetMartketPx() error = %v", err)
	}
	if res < 10_000 {
		t.Errorf("GetMartketPx() = %v, want > %v", res, 10_000)
	}
	t.Logf("GetMartketPx() = %v", res)
}

func TestInfoAPI_BuildMetaMap(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.BuildMetaMap()
	if err != nil {
		t.Errorf("BuildMetaMap() error = %v", err)
	}
	if len(res) == 0 {
		t.Errorf("BuildMetaMap() = %v, want > %v", res, 0)
	}
	// check BTC, ETH in map
	if _, ok := res["BTC"]; !ok {
		t.Errorf("BuildMetaMap() = %v, want %v", res, "BTC")
	}
	if _, ok := res["ETH"]; !ok {
		t.Errorf("BuildMetaMap() = %v, want %v", res, "ETH")
	}
	t.Logf("BuildMetaMap() = %v", res)
}

func TestInfoAPI_BuildSpotMetaMap(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.BuildSpotMetaMap()
	if err != nil {
		t.Errorf("BuildSpotMetaMap() error = %v", err)
	}
	if len(res) == 0 {
		t.Errorf("BuildSpotMetaMap() = %v, want > %v", res, 0)
	}
	// check PURR, HYPE in map
	if _, ok := res["PURR"]; !ok {
		t.Errorf("BuildSpotMetaMap() = %v, want %v", res, "PURR")
	}
	if _, ok := res["HYPE"]; !ok {
		t.Errorf("BuildSpotMetaMap() = %v, want %v", res, "HYPE")
	}
	t.Logf("map(PURR) = %+v", res["PURR"])
	t.Logf("BuildSpotMetaMap() = %+v", res)
}

func TestInfoAPI_GetSpotMeta(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetSpotMeta()
	if err != nil {
		t.Errorf("GetSpotMeta() error = %v", err)
	}
	if len(res.Tokens) == 0 {
		t.Errorf("GetSpotMeta() = %v, want > %v", res, 0)
	}
	t.Logf("GetSpotMeta() = %v", res)
}

func TestInfoAPI_GetAllSpotPrices(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetAllSpotPrices()
	if err != nil {
		t.Errorf("GetAllSpotPrices() error = %v", err)
	}
	if len(*res) == 0 {
		t.Errorf("GetAllSpotPrices() = %v, want > %v", res, 0)
	}
	t.Logf("GetAllSpotPrices() = %+v", res)
}

func TestInfoAPI_GetSpotMarketPx(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetSpotMarketPx("HYPE")
	if err != nil {
		t.Errorf("GetSpotMarketPx() error = %v", err)
	}
	if res < 0 {
		t.Errorf("GetSpotMarketPx() = %v, want > %v", res, 0)
	}
	t.Logf("GetSpotMarketPx(HYPE) = %v", res)
}

func TestInfoAPI_GetUserStateSpot(t *testing.T) {
	api := GetInfoAPI()
	res, err := api.GetAccountStateSpot()
	if err != nil {
		t.Errorf("GetUserStateSpot() error = %v", err)
	}
	if len(res.Balances) == 0 {
		t.Errorf("GetUserStateSpot() = %v, want > %v", res, 0)
	}
	t.Logf("GetUserStateSpot() = %+v", res)
}
