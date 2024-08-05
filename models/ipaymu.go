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
