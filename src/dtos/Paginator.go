package dtos

type Paginator struct {
	CurrentPage int
	PageSize    int
	ItemsCount  int
	Items       interface{}
}
