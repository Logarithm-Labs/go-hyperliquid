package hyperliquid

import (
	"log"
	"math"
	"os"
	"testing"
	"time"
)

func GetExchangeAPI() *ExchangeAPI {
	exchangeAPI := NewExchangeAPI(false)
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

func TestExchangeAPI_Endpoint(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	res := exchangeAPI.Endpoint()
	if res != "/exchange" {
		t.Errorf("Endpoint() = %v, want %v", res, "/exchange")
	}
}

func TestExchangeAPI_AccountAddress(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	res := exchangeAPI.AccountAddress()
	TARGET_ADDRESS := os.Getenv("TEST_ADDRESS")
	if res != TARGET_ADDRESS {
		t.Errorf("AccountAddress() = %v, want %v", res, TARGET_ADDRESS)
	}
}

func TestExchangeAPI_isMainnet(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	res := exchangeAPI.IsMainnet()
	if res != false {
		t.Errorf("isMainnet() = %v, want %v", res, true)
	}
}

func TestExchageAPI_TestMetaIsNotEmpty(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	meta := exchangeAPI.meta
	if meta == nil {
		t.Errorf("Meta() = %v, want not nil", meta)
	}
	if len(meta) == 0 {
		t.Errorf("Meta() = %v, want not empty", meta)
	}
	t.Logf("Meta() = %+v", meta)
}

func TestExchangeAPI_UpdateLeverage(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	_, err := exchangeAPI.UpdateLeverage("ETH", true, 20)
	if err != nil {
		t.Errorf("UpdateLeverage() error = %v", err)
	}
	// Set incorrect leverage 2000
	_, err = exchangeAPI.UpdateLeverage("ETH", true, 2000)
	if err == nil {
		t.Errorf("UpdateLeverage() error = %v", err)
	} else if err.Error() != "Invalid leverage value" {
		t.Errorf("UpdateLeverage() error = %v expected Invalid leverage value", err)
	}
	t.Logf("UpdateLeverage() = %v", err)
}

func TestExchangeAPI_MarketOpen(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	size := -0.01
	coin := "ETH"
	res, err := exchangeAPI.MarketOrder(coin, size, nil)
	if err != nil {
		t.Errorf("MakeOpen() error = %v", err)
	}
	t.Logf("MakeOpen() = %v", res)
	avgPrice := res.Response.Data.Statuses[0].Filled.AvgPx
	if avgPrice == 0 {
		t.Errorf("res.Response.Data.Statuses[0].Filled.AvgPx = %v", avgPrice)
	}
	totalSize := res.Response.Data.Statuses[0].Filled.TotalSz
	if totalSize != math.Abs(size) {
		t.Errorf("res.Response.Data.Statuses[0].Filled.TotalSz = %v", totalSize)
	}
	time.Sleep(2 * time.Second) // wait to execute order
	accountState, err := exchangeAPI.infoAPI.GetUserState(exchangeAPI.AccountAddress())
	if err != nil {
		t.Errorf("GetAccountState() error = %v", err)
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
		t.Errorf("Position not found: %v", accountState.AssetPositions)
	}
	if !positionCorrect {
		t.Errorf("Position not correct: %v", accountState.AssetPositions)
	}
	t.Logf("GetAccountState() = %v", accountState)
	time.Sleep(5 * time.Second) // wait to execute order
}

func TestExchangeAPI_MarketClose(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	res, err := exchangeAPI.ClosePosition("ETH")
	if err != nil {
		t.Errorf("MakeClose() error = %v", err)
	}
	t.Logf("MakeClose() = %v", res)
}

func TestExchangeAPI_LimitOrder(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	size := 100.1234
	coin := "PNUT"
	px := 0.154956
	res, err := exchangeAPI.LimitOrder(TifGtc, coin, size, px, false)
	if err != nil {
		t.Errorf("MakeLimit() error = %v", err)
	}
	t.Logf("MakeLimit() = %v", res)
}

func TestExchangeAPI_CancelAllOrders(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	res, err := exchangeAPI.CancelAllOrders()
	if err != nil {
		t.Errorf("CancelAllOrders() error = %v", err)
	}
	t.Logf("CancelAllOrders() = %v", res)
}

func TestExchangeAPI_CreateLimitOrderAndCancelOrderByCloidt(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	size := -0.01
	coin := "BTC"
	px := 105000.0
	cloid := "0x1234567890abcdef1234567890abcdef"
	res, err := exchangeAPI.LimitOrder(TifGtc, coin, size, px, false, cloid)
	if err != nil {
		t.Errorf("MakeLimit() error = %v", err)
	}
	t.Logf("MakeLimit() = %v", res)
	openOrders, err := exchangeAPI.infoAPI.GetOpenOrders(exchangeAPI.AccountAddress())
	if err != nil {
		t.Errorf("GetAccountOpenOrders() error = %v", err)
	}
	t.Logf("GetAccountOpenOrders() = %v", openOrders)
	orderOpened := false
	var orderCloid string
	for _, order := range *openOrders {
		t.Logf("Order: %+v", order)
		if order.Coin == coin && order.Sz == -size && order.LimitPx == px {
			orderOpened = true
			orderCloid = order.Cloid
			break
		}
	}
	if !orderOpened {
		t.Errorf("Order not found: %v", openOrders)
	}
	time.Sleep(5 * time.Second) // wait to execute order
	cancelRes, err := exchangeAPI.CancelOrderByCloid(coin, orderCloid)
	if err != nil {
		t.Errorf("CancelOrderByCloid() error = %v", err)
	}
	t.Logf("CancelOrderByCloid() = %v", cancelRes)
}

func TestExchangeAPI_CreateLimitOrderAndCancelOrderByOid(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	size := 0.1
	coin := "BTC"
	px := 85000.0
	res, err := exchangeAPI.LimitOrder(TifGtc, coin, size, px, false)
	if err != nil {
		t.Errorf("MakeLimit() error = %v", err)
	}
	t.Logf("MakeLimit() = %v", res)
	openOrders, err := exchangeAPI.infoAPI.GetOpenOrders(exchangeAPI.AccountAddress())
	if err != nil {
		t.Errorf("GetAccountOpenOrders() error = %v", err)
	}
	t.Logf("GetAccountOpenOrders() = %v", openOrders)
	orderOpened := false
	var orderOid int64
	for _, order := range *openOrders {
		t.Logf("Order: %+v", order)
		if order.Coin == coin && order.Sz == size && order.LimitPx == px {
			orderOpened = true
			orderOid = order.Oid
			break
		}
	}
	if !orderOpened {
		t.Errorf("Order not found: %v", openOrders)
	}
	time.Sleep(5 * time.Second) // wait to execute order
	cancelRes, err := exchangeAPI.CancelOrderByOID(coin, int(orderOid))
	if err != nil {
		t.Errorf("CancelOrderByOid() error = %v", err)
	}
	t.Logf("CancelOrderByOid() = %v", cancelRes)
}

func TestExchangeAPI_TestModifyOrder(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	size := 0.005
	coin := "ETH"
	px := 2000.0
	res, err := exchangeAPI.LimitOrder(TifGtc, coin, size, px, false)
	if err != nil {
		t.Errorf("MakeLimit() error = %v", err)
	}
	t.Logf("MakeLimit() = %v", res)
	openOrders, err := exchangeAPI.infoAPI.GetOpenOrders(exchangeAPI.AccountAddress())
	if err != nil {
		t.Errorf("GetAccountOpenOrders() error = %v", err)
	}
	t.Logf("GetAccountOpenOrders() = %v", openOrders)
	orderOpened := false
	for _, order := range *openOrders {
		if order.Coin == coin && order.Sz == size && order.LimitPx == px {
			orderOpened = true
			break
		}
	}
	log.Printf("Order ID: %v", res.Response.Data.Statuses[0].Resting.OrderID)
	if !orderOpened {
		t.Errorf("Order not found: %+v", openOrders)
	}
	time.Sleep(5 * time.Second) // wait to execute order
	// modify order
	newPx := 2500.0
	orderType := OrderType{
		Limit: &LimitOrderType{
			Tif: TifGtc,
		},
	}
	orderID := res.Response.Data.Statuses[0].Resting.OrderID

	modifyOrderRequest := OrderRequest{
		OrderID:    &orderID,
		Coin:       coin,
		Sz:         size,
		LimitPx:    newPx,
		OrderType:  orderType,
		IsBuy:      true,
		ReduceOnly: false,
	}
	modifyRes, err := exchangeAPI.BulkModifyOrders([]OrderRequest{modifyOrderRequest})
	if err != nil {
		t.Errorf("ModifyOrder() error = %v", err)
	}
	t.Logf("ModifyOrder() = %+v", modifyRes)
	time.Sleep(5 * time.Second) // wait to execute order
	cancelRes, err := exchangeAPI.CancelAllOrders()
	if err != nil {
		t.Errorf("CancelAllOrders() error = %v", err)
	}
	t.Logf("CancelAllOrders() = %v", cancelRes)
}

func TestExchangeAPI_TestMultipleMarketOrder(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	testCases := []struct {
		coin string
		size float64
	}{
		{"BTC", 0.001},
		{"BTC", -0.001},
		{"ETH", 0.12},
		{"ETH", -0.12},
		{"INJ", 21.1},
		{"INJ", -21.1},
		{"PNUT", 100.122},
		{"PNUT", -100.1},
		{"ADA", 100.123456},
		{"ADA", -100.123456},
	}
	for _, tc := range testCases {
		t.Run(tc.coin, func(t *testing.T) {
			res, err := exchangeAPI.MarketOrder(tc.coin, tc.size, nil)
			if err != nil {
				t.Errorf("MarketOrder() error = %v", err)
			}
			t.Logf("MarketOrder() = %v", res)
		})

	}
}

func TestExchangeAPI_TestIncorrectOrderSize(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	size := 0.1
	coin := "ADA"
	res, err := exchangeAPI.MarketOrder(coin, size, nil)
	if err != nil {
		t.Errorf("MarketOrder() error = %v", err)
	}
	if res.Response.Data.Statuses[0].Error != "Order has zero size." {
		t.Errorf("MarketOrder() error = %s but expected %s", res.Response.Data.Statuses[0].Error, "Order has zero size.")
	}
}

func TestExchangeAPI_TestClosePositionByMarket(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	size := -1.0
	coin := "ETH"
	res, err := exchangeAPI.MarketOrder(coin, size, nil)
	if err != nil {
		t.Errorf("MarketOrder() error = %v", err)
	}
	t.Logf("MarketOrder() = %v", res)
	time.Sleep(5 * time.Second)
	// close position with big size
	closeRes, err := exchangeAPI.MarketOrder(coin, -size, nil)
	if err != nil {
		t.Errorf("MarketOrder() error = %v", err)
	}
	t.Logf("MarketOrder() = %v", closeRes)
	accountState, err := exchangeAPI.infoAPI.GetUserState(exchangeAPI.AccountAddress())
	if err != nil {
		t.Errorf("GetAccountState() error = %v", err)
	}
	// check that there is no opened position
	if len(accountState.AssetPositions) != 0 {
		t.Errorf("Account has opened positions: %v", accountState.AssetPositions)
	}
}

// Test Mainnet Only
func TestExchangeAPI_TestWithdraw(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	withdrawAmount := 20.0
	stateBefore, err := exchangeAPI.infoAPI.GetUserState(exchangeAPI.AccountAddress())
	if err != nil {
		t.Errorf("GetAccountState() error = %v", err)
	}
	t.Logf("GetAccountState() = %v", stateBefore)
	balanceBefore := stateBefore.Withdrawable
	if balanceBefore < withdrawAmount {
		t.Errorf("Insufficient balance: %v", stateBefore)
	}
	accountAddress := exchangeAPI.AccountAddress() // withdraw to the same address
	res, err := exchangeAPI.Withdraw(accountAddress, withdrawAmount)
	if err != nil {
		t.Errorf("Withdraw() error = %v", err)
	}
	t.Logf("Withdraw() = %v", res)
}

func TestExchageAPI_TestMarketOrderSpot(t *testing.T) {
	exchangeAPI := GetExchangeAPI()
	size := 0.81
	coin := "HYPE"
	res, err := exchangeAPI.MarketOrderSpot(coin, size, nil)
	if err != nil {
		t.Errorf("MakeOpen() error = %v", err)
	}
	t.Logf("MakeOpen() = %v", res)
	avgPrice := res.Response.Data.Statuses[0].Filled.AvgPx
	if avgPrice == 0 {
		t.Errorf("res.Response.Data.Statuses[0].Filled.AvgPx = %v", avgPrice)
	}
}
