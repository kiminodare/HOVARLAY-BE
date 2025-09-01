package dtoHistory

type GetHistoriesQuery struct {
	Page  int `json:"page" validate:"min=0"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}
