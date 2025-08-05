package dto

type AssetStat struct {
	Title    string `json:"title"`
	Price    int    `json:"price"`
	Cashflow int    `json:"cashflow"`
}

type AssetCard struct {
	Title    string `json:"title"`
	Descr    string `json:"descr"`
	Price    int    `json:"price"`
	Cashflow int    `json:"cashflow"`
}
