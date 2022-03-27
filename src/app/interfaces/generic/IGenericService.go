package generic

import (
	"rol/app/interfaces"
	"rol/dtos"
)

type IGenericService[DtoType interfaces.IEntityDtoModel,
	CreateDtoType interfaces.IEntityDtoModel,
	UpdateDtoType interfaces.IEntityDtoModel,
	EntityType interfaces.IEntityModel] interface {
	//GetList
	//	Get list of elements with search and pagination.
	//Params
	//	search - string value to search
	//	orderBy - order field name
	//	orderDirection - Order direction, desc/asc
	//	page - number of the page
	//	pageSize - size of the page
	//Return
	//	paginator - pointer to the struct with pagination info and entities
	//	error - if an error occurred, otherwise nil
	GetList(search, orderBy, orderDirection string, page, pageSize int) (*dtos.PaginatedListDto[DtoType], error)
	//GetById
	//	Get entity by ID from service.
	//Params
	//  dto - pointer to the entity DTO.
	//	id - entity id
	//Return
	//	error - if an error occurred, otherwise nil
	GetById(id uint) (*DtoType, error)
	//Update
	//	Save the changes to the existing entity in the service.
	//Params
	//  updateDto - pointer to the entity DTO.
	//	id - entity id
	//Return
	//	error - if an error occurred, otherwise nil
	Update(updateDto UpdateDtoType, id uint) error
	//Create
	//	Add entity to the service.
	//Params
	//  createDto - pointer to the entity DTO.
	//Return
	//	uint - entity id
	//	error - if an error occurred, otherwise nil
	Create(createDto CreateDtoType) (uint, error)
	//Delete
	//	Delete entity from the service.
	//Params
	//	id - entity id
	//Return
	//	error - if an error occurred, otherwise nil
	Delete(id uint) error
}
