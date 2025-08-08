package entity

type ProductCreate struct {
	Name        string  `json:"name"`
	Size        int     `json:"size"`
	Type        string  `json:"type"`
	ImgUrl      string  `json:"img_url"`
	Price       float32 `json:"price"`
	Count       int     `json:"count"`
	Description string  `json:"description"`
	CategoryId  string  `json:"category_id"`
}

type ProductUpdate struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Size        int     `json:"size"`
	Type        string  `json:"type"`
	ImgUrl      string  `json:"img_url"`
	Price       float32 `json:"price"`
	Count       int     `json:"count"`
	Description string  `json:"description"`
	CategoryId  string  `json:"category_id"`
}

type ProductRes struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Size        int     `json:"size"`
	Type        string  `json:"type"`
	ImgUrl      string  `json:"img_url"`
	Price       float32 `json:"price"`
	Count       int     `json:"count"`
	Description string  `json:"description"`
	CategoryId  string  `json:"category_id"`
	CreatedAt   string  `json:"created_at"`
}

type ProductGetAllReq struct {
	Search     string `json:"search"`
	CategoryId string `json:"category_id"`
	Filter     Filter `json:"filter"`
}
type ProductGetAllRes struct {
	Products []ProductRes `json:"Products"`
	Count    int          `json:"count"`
}
