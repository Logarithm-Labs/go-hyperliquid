package hyperliquid

// Base request for /info
type InfoRequest struct {
	User      string `json:"user,omitempty"`
	Type      string `json:"type"`
	Oid       string `json:"oid,omitempty"`
	Coin      string `json:"coin,omitempty"`
	StartTime int64  `json:"startTime,omitempty"`
	EndTime   int64  `json:"endTime,omitempty"`
}

type UserStateRequest struct {
	User string `json:"user"`
	Type string `json:"type"`
}

type Role string

const (
	RoleUser       Role = "user"
	RoleAgent      Role = "agent"
	RoleVault      Role = "vault"
	RoleSubAccount Role = "subAccount"
	RoleMissing    Role = "missing"
)

// UserRole represents a user role.
// Role can be one of "user", "agent", "vault",  "subAccount" or "missing"
type UserRole struct {
	Role Role `json:"role"`
	Data struct {
		Master string `json:"master,omitempty"`
	} `json:"data,omitempty"`
}

// IsVaultOrSubAccount checks if the role is either Vault or SubAccount.
func (r Role) IsVaultOrSubAccount() bool {
	return r == RoleVault || r == RoleSubAccount
}

// Asset represents an asset.
type Asset struct {
	Name         string `json:"name"`
	SzDecimals   int    `json:"szDecimals"`
	MaxLeverage  int    `json:"maxLeverage"`
	OnlyIsolated bool   `json:"onlyIsolated"`
}

type UserState struct {
	Withdrawable               float64         `json:"withdrawable,string"`
	CrossMaintenanceMarginUsed float64         `json:"crossMaintenanceMarginUsed,string"`
	AssetPositions             []AssetPosition `json:"assetPositions"`
	CrossMarginSummary         MarginSummary   `json:"crossMarginSummary"`
	MarginSummary              MarginSummary   `json:"marginSummary"`
	Time                       int64           `json:"time"`
}

type UserRoleResponse struct {
	UserRole UserRole `json:"userRole"`
}

type AssetPosition struct {
	Position Position `json:"position"`
	Type     string   `json:"type"`
}

type Position struct {
	Coin           string   `json:"coin"`
	EntryPx        float64  `json:"entryPx,string"`
	Leverage       Leverage `json:"leverage"`
	LiquidationPx  float64  `json:"liquidationPx,string"`
	MarginUsed     float64  `json:"marginUsed,string"`
	PositionValue  float64  `json:"positionValue,string"`
	ReturnOnEquity float64  `json:"returnOnEquity,string"`
	Szi            float64  `json:"szi,string"`
	UnrealizedPnl  float64  `json:"unrealizedPnl,string"`
	MaxLeverage    int      `json:"maxLeverage"`
	CumFunding     struct {
		AllTime   float64 `json:"allTime,string"`
		SinceOpne float64 `json:"sinceOpen,string"`
		SinceChan float64 `json:"sinceChange,string"`
	} `json:"cumFunding"`
}
type UserStateSpot struct {
	Balances []SpotAssetPosition `json:"balances"`
}

// SpotAssetPosition represents an asset position.
//
//	SpotAssetPosition{"coin": "USDC", "token": 0, "hold": "0.0", "total": "14.625485", "entryNtl": "0.0"}
type SpotAssetPosition struct {
	Coin     string  `json:"coin"`
	Token    int     `json:"token"`
	Hold     float64 `json:"hold,string"`
	Total    float64 `json:"total,string"`
	EntryNtl float64 `json:"entryNtl,string"`
}

type Order struct {
	Children         []any   `json:"children,omitempty"`
	Cloid            string  `json:"cloid,omitempty"`
	Coin             string  `json:"coin"`
	IsPositionTpsl   bool    `json:"isPositionTpsl,omitempty"`
	IsTrigger        bool    `json:"isTrigger,omitempty"`
	LimitPx          float64 `json:"limitPx,string,omitempty"`
	Oid              int64   `json:"oid"`
	OrderType        string  `json:"orderType,omitempty"`
	OrigSz           float64 `json:"origSz,string,omitempty"`
	ReduceOnly       bool    `json:"reduceOnly,omitempty"`
	Side             string  `json:"side"`
	Sz               float64 `json:"sz,string,omitempty"`
	Tif              string  `json:"tif,omitempty"`
	Timestamp        int64   `json:"timestamp"`
	TriggerCondition string  `json:"triggerCondition,omitempty"`
	TriggerPx        float64 `json:"triggerPx,string,omitempty"`
}

type Leverage struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

type MarginSummary struct {
	AccountValue    float64 `json:"accountValue,string"`
	TotalMarginUsed float64 `json:"totalMarginUsed,string"`
	TotalNtlPos     float64 `json:"totalNtlPos,string"`
	TotalRawUsd     float64 `json:"totalRawUsd,string"`
}

type SpotMeta struct {
	Universe []struct {
		Tokens      []int  `json:"tokens"`
		Name        string `json:"name"`
		Index       int    `json:"index"`
		IsCanonical bool   `json:"isCanonical"`
	} `json:"universe"`
	Tokens []struct {
		Name        string `json:"name"`
		SzDecimals  int    `json:"szDecimals"`
		WeiDecimals int    `json:"weiDecimals"`
		Index       int    `json:"index"`
		TokenID     string `json:"tokenId"`
		IsCanonical bool   `json:"isCanonical"`
		EvmContract any    `json:"evmContract"`
		FullName    any    `json:"fullName"`
	} `json:"tokens"`
}

type Meta struct {
	Universe []Asset `json:"universe"`
}

type OrderFill struct {
	Cloid         string       `json:"cloid"`
	ClosedPnl     float64      `json:"closedPnl,string"`
	Coin          string       `json:"coin"`
	Crossed       bool         `json:"crossed"`
	Dir           string       `json:"dir"`
	Fee           float64      `json:"fee,string"`
	FeeToken      string       `json:"feeToken"`
	Hash          string       `json:"hash"`
	Oid           int          `json:"oid"`
	Px            float64      `json:"px,string"`
	Side          string       `json:"side"`
	StartPosition string       `json:"startPosition"`
	Sz            float64      `json:"sz,string"`
	Tid           int64        `json:"tid"`
	Time          int64        `json:"time"`
	Liquidation   *Liquidation `json:"liquidation"`
}

type Context struct {
	DayNtlVlm    string   `json:"dayNtlVlm"`
	Funding      string   `json:"funding"`
	ImpactPxs    []string `json:"impactPxs"`
	MarkPx       string   `json:"markPx"`
	MidPx        string   `json:"midPx"`
	OpenInterest string   `json:"openInterest"`
	OraclePx     string   `json:"oraclePx"`
	Premium      string   `json:"premium"`
	PrevDayPx    string   `json:"prevDayPx"`
}

type HistoricalFundingRate struct {
	Coin        string `json:"coin"`
	FundingRate string `json:"fundingRate"`
	Premium     string `json:"premium"`
	Time        int64  `json:"time"`
}

type L2BookSnapshot struct {
	Coin   string `json:"coin"`
	Time   int64  `json:"time"`
	Levels [][]struct {
		Px float64 `json:"px,string"`
		Sz float64 `json:"sz,string"`
		N  int     `json:"n"`
	} `json:"levels"`
}

type CandleSnapshotSubRequest struct {
	Coin      string `json:"coin"`
	Interval  string `json:"interval"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
}

type CandleSnapshotRequest struct {
	Type string                   `json:"type"`
	Req  CandleSnapshotSubRequest `json:"req"`
}

type CandleSnapshot struct {
	CloseTime int64   `json:"t"`
	OpenTime  int64   `json:"T"`
	Symbol    string  `json:"s"`
	Interval  string  `json:"i"`
	Open      float64 `json:"o,string"`
	Close     float64 `json:"c,string"`
	High      float64 `json:"h,string"`
	Low       float64 `json:"l,string"`
	Volume    float64 `json:"v,string"`
	N         int     `json:"n"`
}

type NonFundingUpdate struct {
	Hash  string          `json:"hash"`
	Time  int64           `json:"time"`
	Delta NonFundingDelta `json:"delta"`
}

type FundingUpdate struct {
	Hash  string       `json:"hash"`
	Time  int64        `json:"time"`
	Delta FundingDelta `json:"delta"`
}

type RatesLimits struct {
	CumVlm        float64 `json:"cumVlm,string"`
	NRequestsUsed int     `json:"nRequestsUsed"`
	NRequestsCap  int     `json:"nRequestsCap"`
}

type SpotMetaAndAssetCtxsResponse [2]interface{} // Array of exactly 2 elements

type Market struct {
	PrevDayPx         string `json:"prevDayPx,omitempty"`
	DayNtlVlm         string `json:"dayNtlVlm,omitempty"`
	MarkPx            string `json:"markPx,omitempty"`
	MidPx             string `json:"midPx,omitempty"`
	CirculatingSupply string `json:"circulatingSupply,omitempty"`
	Coin              string `json:"coin,omitempty"`
	TotalSupply       string `json:"totalSupply,omitempty"`
	DayBaseVlm        string `json:"dayBaseVlm,omitempty"`
}
