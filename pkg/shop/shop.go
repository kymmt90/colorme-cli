package shop

type ShopResource struct {
	Shop Shop `json:"shop"`
}

type Shop struct {
	URL string `json:"url"`
}
