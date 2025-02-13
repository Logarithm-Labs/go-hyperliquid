package hyperliquid

import (
	"os"
	"sync"
	"testing"
)

func GetHyperliquidAPI() *Hyperliquid {
	hl := NewHyperliquid(&HyperliquidClientConfig{
		IsMainnet:      false,
		AccountAddress: os.Getenv("TEST_ADDRESS"),
		PrivateKey:     os.Getenv("TEST_PRIVATE_KEY"),
	})
	if GLOBAL_DEBUG {
		hl.infoAPI.SetDebugActive()
		hl.ExchangeAPI.SetDebugActive()
	}
	return hl
}

func TestHyperliquid_CheckFieldsConsistency(t *testing.T) {
	hl := GetHyperliquidAPI()
	if hl.ExchangeAPI.baseEndpoint != "/exchange" {
		t.Errorf("baseEndpoint = %v, want %v", hl.ExchangeAPI.baseEndpoint, "/exchange")
	}
	if hl.InfoAPI.baseEndpoint != "/info" {
		t.Errorf("baseEndpoint = %v, want %v", hl.InfoAPI.baseEndpoint, "/info")
	}
	var apiUrl string
	if hl.ExchangeAPI.IsMainnet() {
		apiUrl = MAINNET_API_URL
	} else {
		apiUrl = TESTNET_API_URL
	}
	if hl.InfoAPI.baseURL != apiUrl {
		t.Errorf("baseUrl = %v, want %v", hl.InfoAPI.baseURL, apiUrl)
	}
	if hl.InfoAPI.Debug != hl.ExchangeAPI.Debug {
		t.Errorf("debug = %v, want %v", hl.InfoAPI.Debug, hl.ExchangeAPI.Debug)
	}
	savedAddress := hl.AccountAddress()
	newAddress := "0x1234567890"
	hl.SetAccountAddress(newAddress)
	if hl.InfoAPI.AccountAddress() != newAddress {
		t.Errorf("InfoAPI.AccountAddress = %v, want %v", hl.InfoAPI.AccountAddress(), newAddress)
	}
	if hl.ExchangeAPI.AccountAddress() != newAddress {
		t.Errorf("ExchangeAPI.AccountAddress = %v, want %v", hl.ExchangeAPI.AccountAddress(), newAddress)
	}
	if hl.AccountAddress() != newAddress {
		t.Errorf("gl.AccountAddress = %v, want %v", hl.AccountAddress(), newAddress)
	}
	hl.SetAccountAddress(savedAddress)
}

func TestHyperliquid_MakeSomeTradingLogic(t *testing.T) {
	client := GetHyperliquidAPI()

	// Make limit order
	res1, err := client.LimitOrder(TifGtc, "ETH", 0.01, 1234.1, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("LimitOrder(TifIoc, ETH, 0.01, 1234.1, false): %v", res1)

	res2, err := client.LimitOrder(TifGtc, "ETH", 0.01, 1200.1, true)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("LimitOrder(TifGtc, ETH, 0.01, 1200.1, true): %v", res2)

	res3, err := client.LimitOrder(TifGtc, "ETH", -0.01, 5000.1, true)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("LimitOrder(TifGtc, ETH, -0.01, 5000.1, true): %v", res3)

	res4, err := client.LimitOrder(TifGtc, "ETH", 0.01, 1234.1, false, "0x1234567890abcdef1234567890abcdef")
	if err != nil {
		if err != nil {
			t.Errorf("Error: %v", err)
		}
	}
	t.Logf("LimitOrder(TifIoc, ETH, 0.01, 1234.1, false, 0x1234567890abcdef1234567890abcdef): %v", res4)

	// Get all ordres
	res5, err := client.GetAccountOpenOrders()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("GetAccountOpenOrders(): %v", res5)

	// Close all orders
	res6, err := client.CancelAllOrders()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("CancelAllOrders(): %v", res6)

	// Make market order
	res7, err := client.MarketOrder("ETH", 0.01, nil)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("MarketOrder(ETH, 0.01, nil): %v", res7)

	// Close position
	res8, err := client.ClosePosition("ETH")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("ClosePosition(ETH): %v", res8)

	// Get account balance
	res9, err := client.GetAccountState()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("GetAccountState(): %v", res9)
}

func TestHyperliquid_MarketOrder(t *testing.T) {
	client := GetHyperliquidAPI()
	order, err := client.MarketOrder("ADA", 100, nil)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("MarketOrder(ADA, 100, nil): %+v", order)
}

func TestHyperliquid_LimitOrder(t *testing.T) {
	client := GetHyperliquidAPI()
	order1, err := client.LimitOrder(TifGtc, "BTC", 0.01, 70000, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("LimitOrder(TifGtc, BTC, 0.01, 70000, false): %+v", order1)
	order2, err := client.LimitOrder(TifGtc, "BTC", -0.01, 120000, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("LimitOrder(TifGtc, BTC, -0.01, 120000, false): %+v", order2)
}

func TestHyperliquid_GoLimitOrders(t *testing.T) {
	client := GetHyperliquidAPI()
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			order, err := client.LimitOrder(TifGtc, "BTC", 0.001, 60000, false)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			t.Logf("LimitOrder(TifGtc, BTC, 0.01, 70000, false): %+v", order)
		}(wg)
	}
	wg.Wait()
}
