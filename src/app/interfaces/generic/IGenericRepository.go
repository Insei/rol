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
	// 	*[]EntityType - pointer to array of the entities.
	//	error - if an error occurred, otherwise nil
	GetList(orderBy string, orderDirection string, page int, size int, queryBuilder interfaces.IQueryBuilder) (*[]EntityType, error)
	//Count
	// Get count of entities with filtering
	//Params
	//	queryBuilder - QueryBuilder is repo.NewQueryBuilder()
	//Return
	//	int64 - count of entities
	//	error - if an error occurred, otherwise nil
	Count(queryBuilder interfaces.IQueryBuilder) (int64, error)
	//NewQueryBuilder
	//	Get QueryBuilder
	//Return
	//	IQueryBuilder pointer to object that implements IQueryBuilder interface for this repository
	NewQueryBuilder() interfaces.IQueryBuilder
	//GetById
	//	Get entity by ID from repository.
	//Params
	//	id - entity id
	//Return
	//  *EntityType - pointer to the entity.
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
	//  entity - pointer to the entity for update.
	//Return
	//	uint - entity id
	//	error - if an error occurred, otherwise nil
	Insert(entity EntityType) (uint, error)
	//Delete
	//	Mark an entity as deleted in the repository.
	//Params
	//  id - id of the entity for delete.
	//Return
	//	error - if an error occurred, otherwise nil
	Delete(id uint) error
}
