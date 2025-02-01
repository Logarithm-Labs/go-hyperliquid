package hyperliquid

const GLOBAL_DEBUG = false // Default debug that is used in all tests

// Execution constants
const DEFAULT_SLIPPAGE = 0.005 // 0.5% default slippage
const SPOT_MAX_DECIMALS = 8    // Default decimals for spot
const PERP_MAX_DECIMALS = 6    // Default decimals for perp
var SZ_DECIMALS = 2            // Default decimals for usdc

// Signing constants
const HYPERLIQUID_CHAIN_ID = 1337
const VERIFYING_CONTRACT = "0x0000000000000000000000000000000000000000"
const ARBITRUM_CHAIN_ID = 42161
const ARBITRUM_TESTNET_CHAIN_ID = 421614
