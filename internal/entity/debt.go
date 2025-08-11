package entity

type DebtLogCreate struct {
	UserID string `json:"user_id"`
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
}

type DebtLogUpdateBody struct {
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
	Status string `json:"status"`
}

type DebtLogUpdate struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
	Status string `json:"status"`
}

type DebtLogRes struct {
	Id        string `json:"id"`
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	Amount    int    `json:"amount"`
	Reason    string `json:"reason"`
	Status    string `json:"status"`
	TakenTime string `json:"taken_time"`
	GivenTime string `json:"given_time"`
}

type DebtLogGetAllRes struct {
	DebtLogs []DebtLogRes `json:"debt_logs"`
	Count    int          `json:"count"`
}
