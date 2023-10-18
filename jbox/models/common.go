package models

const (
	OrderByAsc  = "asc"
	OrderByDesc = "desc"
)

type OrderOption struct {
	By   string
	Sort string
}

type PaginationOption struct {
	PageNum  int
	PageSize int
}
