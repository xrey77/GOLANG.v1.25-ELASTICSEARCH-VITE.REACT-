package dto

type ProductsItem struct {
	Category       string `json:"category"`
	Descriptions   string `json:"descriptons"`
	Qty            string `json:"qty"`
	Unit           string `json:"unit"`
	Costprice      string `json:"costprice"`
	Sellprice      string `json:"sellprice"`
	Saleprice      string `json:"saleprice"`
	Productpicture string `json:"productpicture"`
	Alertstocks    string `json:"alertstocks"`
	Criticalstocks string `json:"criticalstocks"`
}
