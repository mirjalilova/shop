package entity

type BucketItemCreate struct {
	UserID  string `json:"user_id"`
	ProductID string `json:"product_id"`
	Count     int    `json:"count"`
}

type BucketItemUpdateBody struct {
	Count int     `json:"count"`
	Price float32 `json:"price"`
}

type BucketItemUpdate struct {
	Id    string  `json:"id"`
	Count int     `json:"count"`
	Price float32 `json:"price"`
}

type BucketItemRes struct {
	Id           string  `json:"id"`
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductSize  int     `json:"product_size"`
	ProductType  string  `json:"product_type"`
	ProductPrice float32 `json:"product_price"`
	ImgUrl       string  `json:"img_url"`
	Count        int     `json:"count"`
	Price        float32 `json:"price"`
}

type BucketRes struct {
	Buskets    []BucketItemRes `json:"Buskets"`
	TotalPrice float32         `json:"total_price"`
}
