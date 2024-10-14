package server

type SearchRequest struct {
	Query string `json:"query" validate:"required"`
}

type DataRequest struct {
	ID string `json:"id" validate:"required"`
}
