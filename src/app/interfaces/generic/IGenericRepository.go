package generic

import (
	"rol/app/interfaces"
)

type IGenericRepository[EntityType interfaces.IEntityModel] interface {
	//GetList
	//	Get list of elements with filtering and pagination.
	//Params
	//  orderBy - order field name
	//	orderDirection - Order direction, desc/asc
	//	page - number of the page
	//	size - size of the page
	//	queryBuilder - QueryBuilder is repo.NewQueryBuilder()
	//Return
	// 	entities - pointer to array of the entities.
	//	error - if an error occurred, otherwise nil
	GetList(orderBy string, orderDirection string, page int, size int, queryBuilder interfaces.IQueryBuilder) (*[]EntityType, error)
	Count(queryBuilder interfaces.IQueryBuilder) (int64, error)
	NewQueryBuilder() interfaces.IQueryBuilder
	//GetById
	//	Get entity by ID from repository.
	//Params
	//	id - entity id
	//Return
	//  entity - pointer to the entity.
	//	error - if an error occurred, otherwise nil
	GetById(id uint) (*EntityType, error)
	//Update
	//	Save the changes to the existing entity in the repository.
	//Params
	//  entity - pointer to the entity.
	//Return
	//	error - if an error occurred, otherwise nil
	Update(entity *EntityType) error
	//Insert
	//	Add entity to the repository.
	//Params
	//  entity - the entity.
	//Return
	//	uint - entity id
	//	error - if an error occurred, otherwise nil
	Insert(entity EntityType) (uint, error)
	//Delete
	//	Mark an entity as deleted in the repository.
	//Params
	//  entity - pointer to the entity.
	//Return
	//	error - if an error occurred, otherwise nil
	Delete(id uint) error
}
