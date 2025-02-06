package hyperliquid

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func ToTypedSig(r [32]byte, s [32]byte, v byte) RsvSignature {
	return RsvSignature{
		R: hexutil.Encode(r[:]),
		S: hexutil.Encode(s[:]),
		V: v,
	}
}

func ArrayAppend(data []byte, toAppend []byte) []byte {
	return append(data, toAppend...)
}

func HexToBytes(addr string) []byte {
	if strings.HasPrefix(addr, "0x") {
		fAddr := strings.Replace(addr, "0x", "", 1)
		b, _ := hex.DecodeString(fAddr)
		return b
	} else {
		b, _ := hex.DecodeString(addr)
		return b
	}
}

func OrderWiresToOrderAction(orders []OrderWire, grouping Grouping) PlaceOrderAction {
	return PlaceOrderAction{
		Type:     "order",
		Grouping: grouping,
		Orders:   orders,
	}
}

func OrderRequestToWire(req OrderRequest, meta map[string]AssetInfo, isSpot bool) OrderWire {
	info := meta[req.Coin]
	var assetId, maxDecimals int
	if isSpot {
		// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/asset-ids
		assetId = info.AssetId + 10000
		maxDecimals = SPOT_MAX_DECIMALS
	} else {
		assetId = info.AssetId
		maxDecimals = PERP_MAX_DECIMALS
	}
	return OrderWire{
		Asset:      assetId,
		IsBuy:      req.IsBuy,
		LimitPx:    RoundOrderPrice(req.LimitPx, info.SzDecimals, maxDecimals),
		SizePx:     RoundOrderSize(req.Sz, info.SzDecimals),
		ReduceOnly: req.ReduceOnly,
		OrderType:  OrderTypeToWire(req.OrderType),
		Cloid:      req.Cloid,
	}
}

func ModifyOrderRequestToWire(req ModifyOrderRequest, meta map[string]AssetInfo, isSpot bool) ModifyOrderWire {
	info := meta[req.Coin]
	var assetId, maxDecimals int
	if isSpot {
		// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/asset-ids
		assetId = info.AssetId + 10000
		maxDecimals = SPOT_MAX_DECIMALS
	} else {
		assetId = info.AssetId
		maxDecimals = PERP_MAX_DECIMALS
	}
	return ModifyOrderWire{
		OrderId: req.OrderId,
		Order: OrderWire{
			Asset:      assetId,
			IsBuy:      req.IsBuy,
			LimitPx:    RoundOrderPrice(req.LimitPx, info.SzDecimals, maxDecimals),
			SizePx:     RoundOrderSize(req.Sz, info.SzDecimals),
			ReduceOnly: req.ReduceOnly,
			OrderType:  OrderTypeToWire(req.OrderType),
		},
	}
}

func OrderTypeToWire(orderType OrderType) OrderTypeWire {
	if orderType.Limit != nil {
		return OrderTypeWire{
			Limit: &LimitOrderType{
				Tif: orderType.Limit.Tif,
			},
			Trigger: nil,
		}
	} else if orderType.Trigger != nil {
		return OrderTypeWire{
			Trigger: &TriggerOrderType{
				TpSl:      orderType.Trigger.TpSl,
				TriggerPx: orderType.Trigger.TriggerPx,
				IsMarket:  orderType.Trigger.IsMarket,
			},
			Limit: nil,
		}
	}
	return OrderTypeWire{}
}

// Format the float with custom decimal places, default is 6 (perp), 8 (spot).
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/tick-and-lot-size
func FloatToWire(x float64, maxDecimals int, szDecimals int) string {
	bigf := big.NewFloat(x)
	var maxDecSz uint
	intPart, _ := bigf.Int64()
	intSize := len(strconv.FormatInt(intPart, 10))
	if intSize >= maxDecimals {
		maxDecSz = 0
	} else {
		maxDecSz = uint(maxDecimals - intSize)
	}
	x, _ = bigf.Float64()
	rounded := fmt.Sprintf("%.*f", maxDecSz, x)
	if strings.Contains(rounded, ".") {
		for strings.HasSuffix(rounded, "0") {
			rounded = strings.TrimSuffix(rounded, "0")
		}
	}
	if strings.HasSuffix(rounded, ".") {
		rounded = strings.TrimSuffix(rounded, ".")
	}
	return rounded
}

// To sign raw messages via EIP-712
func StructToMap(strct any) (res map[string]interface{}, err error) {
	a, err := json.Marshal(strct)
	if err != nil {
		return map[string]interface{}{}, err
	}
	json.Unmarshal(a, &res)
	return res, nil
}

// Round the order size to the nearest tick size
func RoundOrderSize(x float64, szDecimals int) string {
	newX := math.Round(x*math.Pow10(szDecimals)) / math.Pow10(szDecimals)
	// TODO: add rounding
	return big.NewFloat(newX).Text('f', szDecimals)
}

// Round the order price to the nearest tick size
func RoundOrderPrice(x float64, szDecimals int, maxDecimals int) string {
	maxSignFigures := 5
	allowedDecimals := maxDecimals - szDecimals
	numberOfDigitsInIntegerPart := len(strconv.Itoa(int(x)))
	if numberOfDigitsInIntegerPart >= maxSignFigures {
		return RoundOrderSize(x, 0)
	}
	allowedSignFigures := maxSignFigures - numberOfDigitsInIntegerPart
	if x < 1 {
		text := RoundOrderSize(x, allowedDecimals)
		startSignFigures := false
		for i := 2; i < len(text); i++ {
			if text[i] == '0' && !startSignFigures {
				continue
			}
			startSignFigures = true
			allowedSignFigures--
			if allowedSignFigures == 0 {
				return text[:i+1]
			}
		}
		return text
	} else {
		return RoundOrderSize(x, min(allowedSignFigures, allowedDecimals))
	}
}
