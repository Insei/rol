package dtos

type PaginatedListDto[DtoType any] struct {
	Page  int
	Size  int
	Total int64
	Items *[]DtoType
}
