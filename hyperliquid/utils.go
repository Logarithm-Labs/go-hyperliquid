package hyperliquid

import (
	"crypto/rand"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// global nonce counter
var nonceCounter = time.Now().UnixMilli()

// Hyperliquid uses timestamps in milliseconds for nonce
// GetNonce returns a unique nonce that is always at least the current time in milliseconds.
// It ensures thread-safe updates using atomic operations.
func GetNonce() uint64 {
	now := time.Now().UnixMilli()
	for {
		// Load the current nonce value atomically.
		current := atomic.LoadInt64(&nonceCounter)

		// If the current time is greater than the stored nonce,
		// attempt to update the nonce to the current time.
		if current < now {
			if atomic.CompareAndSwapInt64(&nonceCounter, current, now) {
				return uint64(now)
			}
			// If the swap fails, retry.
			continue
		}

		// Otherwise, increment the nonce by one.
		newNonce := current + 1
		if atomic.CompareAndSwapInt64(&nonceCounter, current, newNonce) {
			return uint64(newNonce)
		}
	}
}

// Retruns a random cloid (Client Order ID)
func GetRandomCloid() string {
	buf := make([]byte, 16)
	// then we can call rand.Read.
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}
	return hexutil.Encode(buf)
}

// Calculate the slippage of a trade
func CalculateSlippage(isBuy bool, px float64, slippage float64) float64 {
	if isBuy {
		px = px * (1 + slippage)
	} else {
		px = px * (1 - slippage)
	}
	// Format the float with a precision of 6 significant figures
	pxStr := strconv.FormatFloat(px, 'g', 5, 64)
	// Convert the formatted string to a float
	pxFloat, err := strconv.ParseFloat(pxStr, 64)
	if err != nil {
		return px
	}
	// Round the float to 6 decimal places
	return pxFloat
}

func IsBuy(szi float64) bool {
	if szi > 0 {
		return true
	} else {
		return false
	}
}

// Get the slippage of a trade
// Returns the default slippage if the slippage is nil
func GetSlippage(sl *float64) float64 {
	slippage := DEFAULT_SLIPPAGE
	if sl != nil {
		slippage = *sl
	}
	return slippage
}

// Returns default time range of 90 days
// Returns the start time and end time in milliseconds
func GetDefaultTimeRange() (int64, int64) {
	endTime := time.Now().UnixMilli()
	startTime := time.Now().AddDate(0, 0, -90).UnixMilli()
	return startTime, endTime
}
