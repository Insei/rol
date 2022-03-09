package generic

import (
	"rol/app/interfaces"
	"rol/dtos"
)

type IGenericEntityService interface {
	//GetList
	//	Get list of elements with filtering and pagination.
	//Params
	//	dtoArr - pointer to array of pointers to the entity DTO.
	//	search - string value to search
	//	orderBy - order field name
	//	orderDirection - Order direction, desc/asc
	//	page - number of the page
	//	pageSize - size of the page
	//Return
	//	paginator - struct with pagination info and items
	//	error - if an error occurred, otherwise nil
	GetList(dtoArr interface{}, search, orderBy, orderDirection string, page, pageSize int) (*dtos.Paginator, error)
	//GetAll
	//	Get all entities DTO from service.
	//Params
	//  dtoArr - pointer to array of pointers to the entities DTO.
	//Return
	//	error - if an error occurred, otherwise nil
	GetAll(dtoArr interface{}) error
	//GetById
	//	Get entity by ID from service.
	//Params
	//  dto - pointer to the entity DTO.
	//	id - entity id
	//Return
	//	error - if an error occurred, otherwise nil
	GetById(dto interfaces.IEntityDtoModel, id uint) error
	//Update
	//	Save the changes to the existing entity in the service.
	//Params
	//  updateDto - pointer to the entity DTO.
	//	id - entity id
	//Return
	//	error - if an error occurred, otherwise nil
	Update(updateDto interfaces.IEntityDtoModel, id uint) error
	//Create
	//	Add entity to the service.
	//Params
	//  createDto - pointer to the entity DTO.
	//Return
	//	uint - entity id
	//	error - if an error occurred, otherwise nil
	Create(createDto interfaces.IEntityDtoModel) (uint, error)
	//Delete
	//	Delete entity from the service.
	//Params
	//  dto - pointer to the entity DTO.
	//	id - entity id
	//Return
	//	error - if an error occurred, otherwise nil
	Delete(dto interfaces.IEntityDtoModel, id uint) error
}
