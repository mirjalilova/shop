package entity

type CategoryCreate struct {
	Name string `json:"name"`
}

type CategoryUpdate struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CategoryRes struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type CategoryGetAllRes struct {
	Categorys []CategoryRes `json:"Categorys"`
	Count     int           `json:"count"`
}
