package entity

type ShoesCreate struct {
	Name        string   `json:"name"`
	Size        []int    `json:"size"`
	Color       []string `json:"color"`
	ImgUrl      []string `json:"img_url"`
	Price       float32  `json:"price"`
	Description string   `json:"description"`
	CategoryId  string   `json:"category_id"`
}

type ShoesUpdate struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Size        []int    `json:"size"`
	Color       []string `json:"color"`
	ImgUrl      []string `json:"img_url"`
	Price       float32  `json:"price"`
	Description string   `json:"description"`
	CategoryId  string   `json:"category_id"`
}

type ShoesRes struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Size        []int    `json:"size"`
	Color       []string `json:"color"`
	ImgUrl      []string `json:"img_url"`
	Price       float32  `json:"price"`
	Description string   `json:"description"`
	CategoryId  string   `json:"category_id"`
	CreatedAt   string   `json:"created_at"`
}

type ShoesGetAllRes struct {
	Shoess []ShoesRes `json:"shoess"`
	Count  int        `json:"count"`
}
