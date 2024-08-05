package models

type Packages struct {
	ID                string `json:"$id"`
	Name              string `json:"name"`
	Price             int    `json:"price"`
	MerchantAvailable int    `json:"merchant"`
	CashierAvailable  int    `json:"cashier"`
	CategoryAvailable int    `json:"category"`
	ProductAvailable  int    `json:"product"`
	TableAvailable    int    `json:"table"`
	Description       string `json:"Description"`
	CreatedAt         string `json:"$createdAt"`
	UpdatedAt         string `json:"$updatedAt"`
}

type Owner struct {
	ID                string `json:"$id"`
	OwnerId           string `json:"owner_id"`
	OwnerName         string `json:"owner_name"`
	PackageId         string `json:"package_id"`
	PackageName       string `json:"package_name"`
	MerchantAvailable int    `json:"max_merchant"`
	CashierAvailable  int    `json:"max_cashier"`
	CategoryAvailable int    `json:"max_category"`
	ProductAvailable  int    `json:"max_product"`
	TableAvailable    int    `json:"max_table"`
	CreatedAt         string `json:"$createdAt"`
	UpdatedAt         string `json:"$updatedAt"`
}
