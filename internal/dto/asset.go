package dto

type AssetStat struct {
	Title    string `json:"title"`
	Price    int    `json:"price"`
	Cashflow int    `json:"cashflow"`
	AssetID  int    `json:"asset_id"`
}

type AssetCard struct {
	Title    string `json:"title"`
	Descr    string `json:"descr"`
	Price    int    `json:"price"`
	Cashflow int    `json:"cashflow"`
}
