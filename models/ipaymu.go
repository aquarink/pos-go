package models

type PaymentChannel struct {
	Code        string    `json:"Code"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	Channels    []Channel `json:"Channels"`
}

type Channel struct {
	Code                   string         `json:"Code"`
	Name                   string         `json:"Name"`
	Description            string         `json:"Description"`
	Logo                   string         `json:"Logo"`
	PaymentInstructionsDoc string         `json:"PaymentInstructionsDoc"`
	FeatureStatus          string         `json:"FeatureStatus"`
	HealthStatus           string         `json:"HealthStatus"`
	TransactionFee         TransactionFee `json:"TransactionFee"`
}

type TransactionFee struct {
	ActualFee     float64 `json:"ActualFee"`
	ActualFeeType string  `json:"ActualFeeType"`
	AdditionalFee float64 `json:"AdditionalFee"`
}

type PaymentChannelsResponse struct {
	Status  int              `json:"Status"`
	Success bool             `json:"Success"`
	Message string           `json:"Message"`
	Data    []PaymentChannel `json:"Data"`
}

type Method struct {
	ID            string   `json:"$id"`
	Name          string   `json:"name"`
	Tipe          string   `json:"type"`
	Icon          []string `json:"icon"`
	TrxFeePercent float32  `json:"trx_fee_percent"`
	TrxFeeDenom   float32  `json:"trx_fee_denom"`
	TrxTax        float32  `json:"trx_tax"`
	Status        string   `json:"status"`
	CreatedAt     string   `json:"$createdAt"`
	UpdatedAt     string   `json:"$updatedAt"`
}

const (
	TypeBank   = "bank"
	TypeWallet = "wallet"
)

func IsValidTypeMethod(tipe string) bool {
	switch tipe {
	case TypeBank, TypeWallet:
		return true
	default:
		return false
	}
}
