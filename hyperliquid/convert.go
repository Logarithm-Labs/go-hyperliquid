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
func HexToInt(hexString string) (*big.Int, error) {
	value := new(big.Int)
	if len(hexString) > 1 && hexString[:2] == "0x" {
		hexString = hexString[2:]
	}
	_, success := value.SetString(hexString, 16)
	if !success {
		return nil, fmt.Errorf("invalid hexadecimal string: %s", hexString)
	}
	return value, nil
}

func IntToHex(value *big.Int) string {
	return "0x" + value.Text(16)
}

func OrderWiresToOrderAction(orders []OrderWire, grouping Grouping) PlaceOrderAction {
	return PlaceOrderAction{
		Type:     "order",
		Grouping: grouping,
		Orders:   orders,
	}
}

func (req *OrderRequest) isSpot() bool {
	return strings.ContainsAny(req.Coin, "@-")
}

// ToWire (OrderRequest) converts an OrderRequest to an OrderWire using the provided metadata.
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/asset-ids
func (req *OrderRequest) ToWireMeta(meta map[string]AssetInfo) OrderWire {
	info := meta[req.Coin]
	return req.ToWire(info)
}

// ToModifyWire converts an OrderRequest to a ModifyOrderWire using the provided AssetInfo.
func (req *OrderRequest) ToModifyWire(info AssetInfo) ModifyOrderWire {
	return ModifyOrderWire{
		OrderID: *req.OrderID,
		Order:   req.ToWire(info),
	}
}

// ToWire converts an OrderRequest to an OrderWire using the provided AssetInfo.
func (req *OrderRequest) ToWire(info AssetInfo) OrderWire {
	var assetID = info.AssetID
	var maxDecimals = PERP_MAX_DECIMALS
	if req.isSpot() {
		assetID = info.AssetID + 10000 // https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/asset-ids
		maxDecimals = SPOT_MAX_DECIMALS
	}
	return OrderWire{
		Asset:      assetID,
		IsBuy:      req.IsBuy,
		LimitPx:    PriceToWire(req.LimitPx, maxDecimals, info.SzDecimals),
		SizePx:     SizeToWire(req.Sz, info.SzDecimals),
		ReduceOnly: req.ReduceOnly,
		OrderType:  OrderTypeToWire(req.OrderType),
		Cloid:      req.Cloid,
	}
}

// OrderTypeToWire converts an OrderType to an OrderTypeWire.
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

/**
 * FloatToWire converts a float64 to a string representation following Hyperliquid's decimal rules.
 * FloatToWire converts a float64 to a string representation following Hyperliquid's decimal rules.
 *
 * The conversion adheres to market-specific decimal place constraints:
 * - Perpetual markets: Maximum 6 decimal places
 * - Spot markets: Maximum 8 decimal places
 *
 * The function dynamically adjusts decimal precision based on:
 * 1. Integer part magnitude
 * 2. Maximum allowed decimals (maxDecimals)
 * 3. Size decimal precision (szDecimals)
 *
 * Output formatting:
 * - Removes trailing zeros
 * - Trims unnecessary decimal points
 * - Maintains tick size precision requirements
 *
 * @see https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/tick-and-lot-size
 */
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

// fastPow10 returns 10^exp as a float64. For our purposes exp is small.
func pow10(exp int) float64 {
	var res float64 = 1
	for i := 0; i < exp; i++ {
		res *= 10
	}
	return res
}

// PriceToWire converts a price value to its string representation per Hyperliquid rules.
// It enforces:
//   - At most 5 significant figures,
//   - And no more than (maxDecimals - szDecimals) decimal places.
//
// Integer prices are returned as is.
func PriceToWire(x float64, maxDecimals, szDecimals int) string {
	// If the price is an integer, return it without decimals.
	if x == math.Trunc(x) {
		return strconv.FormatInt(int64(x), 10)
	}

	// Rule 1: The tick rule – maximum decimals allowed is (maxDecimals - szDecimals).
	allowedTick := maxDecimals - szDecimals

	// Rule 2: The significant figures rule – at most 5 significant digits.
	var allowedSig int
	if x >= 1 {
		// Count digits in the integer part.
		digits := int(math.Floor(math.Log10(x))) + 1
		allowedSig = 5 - digits
		if allowedSig < 0 {
			allowedSig = 0
		}
	} else {
		// For x < 1, determine the effective exponent.
		exponent := int(math.Ceil(-math.Log10(x)))
		allowedSig = 4 + exponent
	}

	// Final allowed decimals is the minimum of the tick rule and the significant figures rule.
	allowedDecimals := allowedTick
	if allowedSig < allowedDecimals {
		allowedDecimals = allowedSig
	}
	if allowedDecimals < 0 {
		allowedDecimals = 0
	}

	// Round the price to allowedDecimals decimals.
	factor := pow10(allowedDecimals)
	rounded := math.Round(x*factor) / factor

	// Format the number with fixed precision.
	s := strconv.FormatFloat(rounded, 'f', allowedDecimals, 64)
	// Only trim trailing zeros if the formatted string contains a decimal point.
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}

// SizeToWire converts a size value to its string representation,
// rounding it to exactly szDecimals decimals.
// Integer sizes are returned without decimals.
func SizeToWire(x float64, szDecimals int) string {
	// Return integer sizes without decimals.
	if szDecimals == 0 {
		return strconv.FormatInt(int64(x), 10)
	}
	// Return integer sizes directly.
	if x == math.Trunc(x) {
		return strconv.FormatInt(int64(x), 10)
	}

	// Round the size value to szDecimals decimals.
	factor := pow10(szDecimals)
	rounded := math.Round(x*factor) / factor

	// Format with fixed precision then trim any trailing zeros and the decimal point.
	s := strconv.FormatFloat(rounded, 'f', szDecimals, 64)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
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

// RoundOrderSize rounds the order size to the nearest tick size
func RoundOrderSize(x float64, szDecimals int) string {
	newX := math.Round(x*math.Pow10(szDecimals)) / math.Pow10(szDecimals)
	// TODO: add rounding
	return big.NewFloat(newX).Text('f', szDecimals)
}

// RoundOrderPrice rounds the order price to the nearest tick size
func RoundOrderPrice(x float64, szDecimals int, maxDecimals int) string {
	maxSignFigures := 5
	allowedDecimals := maxDecimals - szDecimals
	numberOfDigitsInIntegerPart := len(strconv.Itoa(int(x)))
	if numberOfDigitsInIntegerPart >= maxSignFigures {
		return RoundOrderSize(x, 0)
	}
	allowedSignFigures := maxSignFigures - numberOfDigitsInIntegerPart
	if x >= 1 {
		return RoundOrderSize(x, min(allowedSignFigures, allowedDecimals))
	}

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
}
