package hyperliquid

import (
	"math"
	"os"
	"testing"
	"time"
)

func GetExchangeAPI() *ExchangeAPI {
	exchangeAPI := NewExchangeAPI(true)
	if GLOBAL_DEBUG {
		exchangeAPI.SetDebugActive()
	}
	TEST_ADDRESS := os.Getenv("TEST_ADDRESS")
	TEST_PRIVATE_KEY := os.Getenv("TEST_PRIVATE_KEY")
	err := exchangeAPI.SetPrivateKey(TEST_PRIVATE_KEY)
	if err != nil {
		panic(err)
	}
	exchangeAPI.SetAccountAddress(TEST_ADDRESS)
	return exchangeAPI
}

func TestExchangeAPI_Endpoint(testing *testing.T) {
	exchangeAPI := GetExchangeAPI()
	res := exchangeAPI.Endpoint()
	if res != "/exchange" {
		testing.Errorf("Endpoint() = %v, want %v", res, "/exchange")
	}
}

func TestExchangeAPI_AccountAddress(testing *testing.T) {
	exchangeAPI := GetExchangeAPI()
	res := exchangeAPI.AccountAddress()
	TARGET_ADDRESS := os.Getenv("TEST_ADDRESS")
	if res != TARGET_ADDRESS {
		testing.Errorf("AccountAddress() = %v, want %v", res, TARGET_ADDRESS)
	}
}

func TestExchangeAPI_isMainnet(testing *testing.T) {
	exchangeAPI := GetExchangeAPI()
	res := exchangeAPI.IsMainnet()
	if res != true {
		testing.Errorf("isMainnet() = %v, want %v", res, true)
	}
}

func TestExchangeAPI_UpdateLeverage(testing *testing.T) {
	exchangeAPI := GetExchangeAPI()
	_, err := exchangeAPI.UpdateLeverage("ETH", true, 20)
	if err != nil {
		testing.Errorf("UpdateLeverage() error = %v", err)
	}
	// Set incorrect leverage 2000
	_, err = exchangeAPI.UpdateLeverage("ETH", true, 2000)
	if err == nil {
		testing.Errorf("UpdateLeverage() error = %v", err)
	} else if err.Error() != "Invalid leverage value" {
		testing.Errorf("UpdateLeverage() error = %v expected Invalid leverage value", err)
	}
	testing.Logf("UpdateLeverage() = %v", err)
}

func TestExchangeAPI_MarketOpen(testing *testing.T) {
	exchangeAPI := GetExchangeAPI()
	size := -0.01
	coin := "ETH"
	res, err := exchangeAPI.MarketOrder(coin, size, nil)
	if err != nil {
		testing.Errorf("MakeOpen() error = %v", err)
	}
	testing.Logf("MakeOpen() = %v", res)
	avgPrice := res.Response.Data.Statuses[0].Filled.AvgPx
	if avgPrice == 0 {
		testing.Errorf("res.Response.Data.Statuses[0].Filled.AvgPx = %v", avgPrice)
	}
	totalSize := res.Response.Data.Statuses[0].Filled.TotalSz
	if totalSize != math.Abs(size) {
		testing.Errorf("res.Response.Data.Statuses[0].Filled.TotalSz = %v", totalSize)
	}
	time.Sleep(2 * time.Second) // wait to execute order
	accountState, err := exchangeAPI.infoAPI.GetUserState(exchangeAPI.AccountAddress())
	if err != nil {
		testing.Errorf("GetAccountState() error = %v", err)
	}
	positionOpened := false
	positionCorrect := false
	for _, position := range accountState.AssetPositions {
		if position.Position.Coin == coin {
			positionOpened = true
		}
		if position.Position.Coin == coin && position.Position.Szi == size {
			positionCorrect = true
		}
	}
	if !positionOpened {
		testing.Errorf("Position not found: %v", accountState.AssetPositions)
	}
	if !positionCorrect {
		testing.Errorf("Position not correct: %v", accountState.AssetPositions)
	}
	testing.Logf("GetAccountState() = %v", accountState)
	time.Sleep(5 * time.Second) // wait to execute order
}

func TestExchangeAPI_LimitOrderAndCancel(testing *testing.T) {
	exchangeAPI := GetExchangeAPI()
	size := 0.01
	coin := "ETH"
	px := 2000.0
	res, err := exchangeAPI.LimitOrder(TifGtc, coin, size, px, false)
	if err != nil {
		testing.Errorf("MakeLimit() error = %v", err)
	}
	testing.Logf("MakeLimit() = %v", res)
	openOrders, err := exchangeAPI.infoAPI.GetOpenOrders(exchangeAPI.AccountAddress())
	if err != nil {
		testing.Errorf("GetAccountOpenOrders() error = %v", err)
	}
	testing.Logf("GetAccountOpenOrders() = %v", openOrders)
	orderOpened := false
	for _, order := range *openOrders {
		if order.Coin == coin && order.Sz == size && order.LimitPx == px {
			orderOpened = true
			break
		}
	}
	if !orderOpened {
		testing.Errorf("Order not found: %v", openOrders)
	}
	time.Sleep(5 * time.Second) // wait to execute order
	cancelRes, err := exchangeAPI.CancelAllOrders()
	if err != nil {
		testing.Errorf("CancelAllOrders() error = %v", err)
	}
	testing.Logf("CancelAllOrders() = %v", cancelRes)
}

func TestExchangeAPI_CancelAllOrders(testing *testing.T) {
	exchangeAPI := GetExchangeAPI()
	res, err := exchangeAPI.CancelAllOrders()
	if err != nil {
		testing.Errorf("CancelAllOrders() error = %v", err)
	}
	testing.Logf("CancelAllOrders() = %v", res)
}

func TestExchangeAPI_MarketClose(testing *testing.T) {
	exchangeAPI := GetExchangeAPI()
	res, err := exchangeAPI.ClosePosition("ETH")
	if err != nil {
		testing.Errorf("MakeClose() error = %v", err)
	}
	testing.Logf("MakeClose() = %v", res)
}

func TestExchangeAPI_TestWithdraw(testing *testing.T) {
	exchangeAPI := GetExchangeAPI()
	withdrawAmount := 10.0
	stateBefore, err := exchangeAPI.infoAPI.GetUserState(exchangeAPI.AccountAddress())
	if err != nil {
		testing.Errorf("GetAccountState() error = %v", err)
	}
	testing.Logf("GetAccountState() = %v", stateBefore)
	balanceBefore := stateBefore.Withdrawable
	if balanceBefore < withdrawAmount {
		testing.Errorf("Insufficient balance: %v", stateBefore)
	}
	accountAddress := exchangeAPI.AccountAddress() // withdraw to the same address
	res, err := exchangeAPI.Withdraw(accountAddress, withdrawAmount)
	if err != nil {
		testing.Errorf("Withdraw() error = %v", err)
	}
	testing.Logf("Withdraw() = %v", res)
	time.Sleep(30 * time.Second) // wait to execute order
	stateAfter, err := exchangeAPI.infoAPI.GetUserState(exchangeAPI.AccountAddress())
	if err != nil {
		testing.Errorf("GetAccountState() error = %v", err)
	}
	testing.Logf("GetAccountState() = %v", stateAfter)
	balanceAfter := stateAfter.Withdrawable
	if balanceAfter >= balanceBefore {
		testing.Errorf("Balance not updated: %v", stateAfter)
	}
}

func TestExchageAPI_TestMarketOrderSpot(testing *testing.T) {
	exchangeAPI := GetExchangeAPI()
	size := 1600.0
	coin := "YEETI"
	res, err := exchangeAPI.MarketOrderSpot(coin, size, nil)
	if err != nil {
		testing.Errorf("MakeOpen() error = %v", err)
	}
	testing.Logf("MakeOpen() = %v", res)
	avgPrice := res.Response.Data.Statuses[0].Filled.AvgPx
	if avgPrice == 0 {
		testing.Errorf("res.Response.Data.Statuses[0].Filled.AvgPx = %v", avgPrice)
	}
}
