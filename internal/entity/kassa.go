package entity

type KassaEvent struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Size      int     `json:"size"`
	ImgURL    *string `json:"img_url"`
	Count     int     `json:"count"`
	Price     float64 `json:"price"`
}

type Kassa struct {
	Items map[string]*KassaEvent `json:"items"`
}

type Formalize struct {
	ProductID string  `json:"product_id"`
	Count     int     `json:"count"`
}
