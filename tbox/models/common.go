package models

const (
	OrderByAsc  = "asc"
	OrderByDesc = "desc"
)

type OrderOption struct {
	By   string
	Type string
}

type PaginationOption struct {
	Page     int
	PageSize int
}
