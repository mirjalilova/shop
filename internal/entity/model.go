package entity

type Filter struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type MultilingualField struct {
	Uz string `json:"uz" example:"Uzbek"`
	Ru string `json:"ru" example:"Русский"`
	En string `json:"en" example:"English"`
}

type ById struct {
	Id string `json:"id"`
}
