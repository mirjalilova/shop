package entity

type OrderCreate struct {
	UserID      string   `json:"user_id"`
	BucketID    string   `json:"bucket_id"`
	Status      string   `json:"status"`
	Location    Location `json:"location"`
	Description string   `json:"description"`
	PaymentType string   `json:"payment_type"`
}

type OrderItemRes struct {
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

type OrderRes struct {
	Buskets    []OrderItemRes `json:"Buskets"`
	TotalPrice float32        `json:"total_price"`
}
