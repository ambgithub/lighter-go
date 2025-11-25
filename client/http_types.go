package client

const (
	CodeOK = 200
)

type ResultCode struct {
	Code    int32  `json:"code,example=200"`
	Message string `json:"message,omitempty"`
}

type NextNonce struct {
	ResultCode
	Nonce int64 `json:"nonce,example=722"`
}

type ApiKey struct {
	AccountIndex int64  `json:"account_index,example=3"`
	ApiKeyIndex  uint8  `json:"api_key_index,example=0"`
	Nonce        int64  `json:"nonce,example=722"`
	PublicKey    string `json:"public_key"`
}

type AccountApiKeys struct {
	ResultCode
	ApiKeys []*ApiKey `json:"api_keys"`
}

type TxHash struct {
	ResultCode
	TxHash string `json:"tx_hash,example=0x70997970C51812dc3A010C7d01b50e0d17dc79C8"`
}

type TransferFeeInfo struct {
	ResultCode
	TransferFee int64 `json:"transfer_fee_usdc"`
}

// ===== 新增: 用于 orderBookDetails 接口的结构体 =====

// MarketDetail 用于存储来自 orderBookDetails 接口的单个市场的详细信息。
type MarketDetail struct {
	Symbol        string `json:"symbol"`
	MarketID      uint8  `json:"market_id"`
	Status        string `json:"status"`
	SizeDecimals  int    `json:"size_decimals"`  // 数量（基础货币）的小数位数
	PriceDecimals int    `json:"price_decimals"` // 价格（报价货币）的小数位数
}

// OrderBookDetailsResponse 用于解析 /api/v1/orderBookDetails 接口的完整响应。
type OrderBookDetailsResponse struct {
	ResultCode
	OrderBookDetails []MarketDetail `json:"order_book_details"`
}

//11
