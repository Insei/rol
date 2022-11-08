package tests

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
	"testing"
	"time"
)

/////////////////////////////////////////////////////////////////
////////////////Entities and their fields section////////////////
/////////////////////////////////////////////////////////////////

type TestEntityFields struct {
	String       string
	SecondString string
	SearchString string
	Yesterday    time.Time
	Tomorrow     time.Time
	NullDate     *time.Time
	NullableDate *time.Time
	Number       int
	Iterator     int
}

func (f TestEntityFields) Equals(fields TestEntityFields) bool {
	if fields.NullDate != f.NullDate {
		return false
	}
	if fields.NullableDate != f.NullableDate {
		return false
	}
	if fields.String != f.String {
		return false
	}
	if fields.SearchString != f.SearchString {
		return false
	}
	if fields.SecondString != f.SecondString {
		return false
	}
	if fields.Number != f.Number {
		return false
	}
	if fields.Tomorrow != f.Tomorrow {
		return false
	}
	if fields.Yesterday != f.Yesterday {
		return false
	}
	return true
}

func GetTestEntityFields() TestEntityFields {
	nullableDate := time.Now()
	//DO NOT CHANGE!
	//IF YOU NEED ADD SOME NEW VALUES FOR NEW TEST ADD NEW FIELDS
	//AND DO NOT FORGET IMPLEMENT CHECK IN Equals() method
	return TestEntityFields{
		String:       "empty",
		SecondString: "second",
		SearchString: "search string with example text",
		NullDate:     nil,
		NullableDate: &nullableDate,
		Yesterday:    time.Now().AddDate(0, 0, -1),
		Tomorrow:     time.Now().AddDate(0, 0, 1),
		Number:       100,
		Iterator:     0,
	}
}

func GetTestEntityFieldsWithIteration(iteration int) TestEntityFields {
	nullableDate := time.Now()
	//DO NOT CHANGE!
	//IF YOU NEED ADD SOME NEW VALUES FOR NEW TEST ADD NEW FIELDS
	//AND DO NOT FORGET IMPLEMENT CHECK IN Equals() method
	return TestEntityFields{
		String:       "empty",
		SecondString: "second",
		SearchString: "search string with example text",
		NullDate:     nil,
		NullableDate: &nullableDate,
		Yesterday:    time.Now().AddDate(0, 0, -1),
		Tomorrow:     time.Now().AddDate(0, 0, 1),
		Number:       100,
		Iterator:     iteration,
	}
}

type BaseTestEntity[IDType comparable] struct {
	domain.Entity[IDType]
	TestEntityFields
}

func (e BaseTestEntity[IDType]) GetDataFields() TestEntityFields {
	return e.TestEntityFields
}

type TestEntityUUID struct {
	BaseTestEntity[uuid.UUID]
}

type TestEntityInt struct {
	BaseTestEntity[int]
}

type ITestEntity[IDType comparable] interface {
	interfaces.IEntityModel[IDType]
	GetDataFields() TestEntityFields
}

/////////////////////////////////////////////////////////////////
/////////////////// Generic Tester interface ////////////////////
/////////////////////////////////////////////////////////////////

type IGenericRepositoryTester interface {
	SetTesting(testing *testing.T)
	Dispose() error
	_0Count_Zero()
	_1InsertAndDelete()
	_2DeleteAll()
	_3GetByID_Found()
	_4GetByID_NotFound()
	_5GetByIDExtended_Found()
	_6GetByIDExtended_NotFound()
	_7GetByIDExtended_WithEmptyConditions()
	_8GetByIDExtended_WithNotExistedFieldCondition()
	_9GetByIDExtended_WithExistedFieldCondition_Found()
	_10GetByIDExtended_WithExistedFieldCondition_NotFound()
	_11GetByIDExtended_WithMultipleExistedFieldCondition_Found()
	_12GetByIDExtended_WithMultipleExistedFieldCondition_NotFound()
	_13GetByIDExtended_WithORMultipleExistedFieldCondition_NotFound()
	_14GetByIDExtended_WithORMultipleExistedFieldCondition_Found()
	_15Delete()
	_16IsExist_Exist()
	_17IsExist_NotExist()
	_18IsExist_WithEmptyConditions_Exist()
	_19IsExist_WithNotExistedFieldCondition_Error()
	_20IsExist_WithExistedFieldCondition_Exist()
	_21IsExist_WithExistedFieldCondition_NotExist()
	_22IsExist_WithMultipleExistedFieldCondition_Exist()
	_23IsExist_WithMultipleExistedFieldCondition_NotExist()
	_24IsExist_WithORMultipleExistedFieldCondition_NotExist()
	_25IsExist_WithORMultipleExistedFieldCondition_Exist()
	_26IsExist_WithExistedNullFieldCondition_Exist()
	_27IsExist_WithExistedNullFieldCondition_NotExist()
	_28IsExist_WithExistedNullableFieldCondition_Exist()
	_29IsExist_WithExistedNullableFieldCondition_NotExist()
	_30Count_All()
	_31Count_WithNotExistedFieldCondition_Error()
	_32Count_ById_OneElem()
	_33Count_ById_TwoElems()
	_34Count_ByIteratorRange()
}

/////////////////////////////////////////////////////////////////
///////////////// Generic Tester implementation /////////////////
/////////////////////////////////////////////////////////////////

type GenericRepositoryTester[IDType comparable, EntityType ITestEntity[IDType]] struct {
	ctx  context.Context
	t    *testing.T
	repo interfaces.IGenericRepository[IDType, EntityType]
}

func NewGenericRepositoryTester[IDType comparable, EntityType ITestEntity[IDType]](repo interfaces.IGenericRepository[IDType, EntityType]) (*GenericRepositoryTester[IDType, EntityType], error) {
	tester := &GenericRepositoryTester[IDType, EntityType]{}
	tester.repo = repo
	tester.ctx = context.Background()
	return tester, nil
}

func (t *GenericRepositoryTester[IDType, EntityType]) DataCleanUp() {
	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("CreatedAt", "!=", time.Now().AddDate(1, 1, 1))
	err := t.repo.DeleteAll(t.ctx, queryBuilder)
	if err != nil {
		panic("tests can't work correctly, database cleanup is not work")
	}
}

func (t *GenericRepositoryTester[IDType, EntityType]) SetTesting(testing *testing.T) {
	t.t = testing
}

func (t *GenericRepositoryTester[IDType, EntityType]) Dispose() error {
	return t.repo.Dispose()
}

func (t *GenericRepositoryTester[IDType, EntityType]) CreatePredefinedTestEntityWithIteration(iteration int) EntityType {
	entity := new(EntityType)
	var entityObj interface{} = entity
	if _, ok := entityObj.(interfaces.IEntityModel[uuid.UUID]); ok {
		entityObj = TestEntityUUID{BaseTestEntity: BaseTestEntity[uuid.UUID]{
			TestEntityFields: GetTestEntityFieldsWithIteration(iteration),
		}}
		return entityObj.(EntityType)
	}
	if _, ok := entityObj.(interfaces.IEntityModel[int]); ok {
		entityObj = TestEntityInt{BaseTestEntity: BaseTestEntity[int]{
			TestEntityFields: GetTestEntityFieldsWithIteration(iteration),
		}}
		return entityObj.(EntityType)
	}
	t.t.Error("Selected entity is not implement IEntityModel[IDType] interface")
	return *entity
}

func (t *GenericRepositoryTester[IDType, EntityType]) CreatePredefinedTestEntity() EntityType {
	return t.CreatePredefinedTestEntityWithIteration(0)
}

func (t *GenericRepositoryTester[IDType, EntityType]) GetNotExistedId() IDType {
	id := new(IDType)
	var idObj interface{} = *id
	switch idObj.(type) {
	case int, int8, int16, int32, int64:
		idObj = 9999999
		return idObj.(IDType)
	case uuid.UUID:
		idObj = uuid.New()
		return idObj.(IDType)
	default:
		t.t.Error("GetNotExistedId: id type is not supported")
	}
	return *id
}

func (t *GenericRepositoryTester[IDType, EntityType]) Create3PredefinedTestEntities() []EntityType {
	//add a few entities
	savedEntityOne, err := t.repo.Insert(t.ctx, t.CreatePredefinedTestEntity())
	assert.NoError(t.t, err)
	savedEntityTwo, err := t.repo.Insert(t.ctx, t.CreatePredefinedTestEntity())
	assert.NoError(t.t, err)
	savedEntityThree, err := t.repo.Insert(t.ctx, t.CreatePredefinedTestEntity())
	assert.NoError(t.t, err)
	return []EntityType{savedEntityOne, savedEntityTwo, savedEntityThree}
}

func (t *GenericRepositoryTester[IDType, EntityType]) _0Count_Zero() {
	t.DataCleanUp()
	totalCount, err := t.repo.Count(t.ctx, nil)
	assert.NoError(t.t, err)
	assert.Equal(t.t, 0, totalCount)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _1InsertAndDelete() {
	t.DataCleanUp()
	newEntity := t.CreatePredefinedTestEntity()
	newEntityFields := newEntity.GetDataFields()
	savedEntity, err := t.repo.Insert(t.ctx, newEntity)
	assert.NoError(t.t, err)
	savedEntityFields := savedEntity.GetDataFields()
	assert.Equal(t.t, true, savedEntityFields.Equals(newEntityFields))
	totalCount, err := t.repo.Count(t.ctx, nil)
	assert.NoError(t.t, err)
	assert.Equal(t.t, 1, totalCount)
	err = t.repo.Delete(t.ctx, savedEntity.GetID())
	assert.NoError(t.t, err)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _2DeleteAll() {
	t.DataCleanUp()
	_ = t.Create3PredefinedTestEntities()
	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("CreatedAt", "!=", time.Now().AddDate(1, 1, 1))
	err := t.repo.DeleteAll(t.ctx, queryBuilder)
	assert.NoError(t.t, err)
	totalCount, err := t.repo.Count(t.ctx, nil)
	assert.NoError(t.t, err)
	assert.Equal(t.t, 0, totalCount)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _3GetByID_Found() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}
	//search
	getEntity, err := t.repo.GetByID(t.ctx, entities[1].GetID())
	assert.NoError(t.t, err)
	assert.Equal(t.t, getEntity.GetID(), entities[1].GetID())
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _4GetByID_NotFound() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	_, err := t.repo.GetByID(t.ctx, t.GetNotExistedId())
	assert.Error(t.t, err)
	assert.Equal(t.t, true, errors.As(err, errors.NotFound))
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _5GetByIDExtended_Found() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	//search
	getEntity, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), nil)
	assert.NoError(t.t, err)
	assert.Equal(t.t, getEntity.GetID(), entities[1].GetID())
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _6GetByIDExtended_NotFound() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	//search
	_, err := t.repo.GetByIDExtended(t.ctx, t.GetNotExistedId(), nil)
	assert.Error(t.t, err)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _7GetByIDExtended_WithEmptyConditions() {
	t.DataCleanUp()
	savedEntity, err := t.repo.Insert(t.ctx, t.CreatePredefinedTestEntity())
	assert.NoError(t.t, err)
	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	_, err = t.repo.GetByIDExtended(t.ctx, savedEntity.GetID(), queryBuilder)
	assert.NoError(t.t, err)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _8GetByIDExtended_WithNotExistedFieldCondition() {
	t.DataCleanUp()
	savedEntity, err := t.repo.Insert(t.ctx, t.CreatePredefinedTestEntity())
	assert.NoError(t.t, err)
	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("NotExistedField", "==", "empty")
	_, err = t.repo.GetByIDExtended(t.ctx, savedEntity.GetID(), queryBuilder)
	assert.Error(t.t, err)
	//make sure that error is not "not found" error
	assert.NotEqual(t.t, true, errors.As(err, errors.NotFound))
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _9GetByIDExtended_WithExistedFieldCondition_Found() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("SecondString", "==", "second")
	receivedEntity, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, receivedEntity.GetID(), entities[1].GetID())
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _10GetByIDExtended_WithExistedFieldCondition_NotFound() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("Number", "==", 1)
	_, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), queryBuilder)
	assert.Error(t.t, err)
	assert.Equal(t.t, true, errors.As(err, errors.NotFound))
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _11GetByIDExtended_WithMultipleExistedFieldCondition_Found() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("SecondString", "==", "second").
		Where("SearchString", "LIKE", "%example%")
	_, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.NotEqual(t.t, true, errors.As(err, errors.NotFound))
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _12GetByIDExtended_WithMultipleExistedFieldCondition_NotFound() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("SecondString", "!=", "second").
		Where("SearchString", "LIKE", "%example%")
	_, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), queryBuilder)
	assert.Error(t.t, err)
	assert.Equal(t.t, true, errors.As(err, errors.NotFound))
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _13GetByIDExtended_WithORMultipleExistedFieldCondition_NotFound() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Or("SecondString", "==", "notsecond").
		Or("String", "!=", "empty")
	_, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), queryBuilder)
	assert.Error(t.t, err)
	assert.Equal(t.t, true, errors.As(err, errors.NotFound))
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _14GetByIDExtended_WithORMultipleExistedFieldCondition_Found() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("SecondString", "==", "second").
		Or("String", "!=", "empty")
	_, err := t.repo.GetByIDExtended(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.NotEqual(t.t, true, errors.As(err, errors.NotFound))
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _15Delete() {
	t.DataCleanUp()
	savedEntity, err := t.repo.Insert(t.ctx, t.CreatePredefinedTestEntity())
	assert.NoError(t.t, err)
	err = t.repo.Delete(t.ctx, savedEntity.GetID())
	assert.NoError(t.t, err)
	_, err = t.repo.GetByID(t.ctx, savedEntity.GetID())
	assert.Error(t.t, err)
	assert.Equal(t.t, true, errors.As(err, errors.NotFound))
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _16IsExist_Exist() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	//checks existence
	exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), nil)
	assert.NoError(t.t, err)
	assert.Equal(t.t, true, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _17IsExist_NotExist() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	//checks existence
	exist, err := t.repo.IsExist(t.ctx, t.GetNotExistedId(), nil)
	assert.NoError(t.t, err)
	assert.Equal(t.t, false, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _18IsExist_WithEmptyConditions_Exist() {
	t.DataCleanUp()
	savedEntity, err := t.repo.Insert(t.ctx, t.CreatePredefinedTestEntity())
	assert.NoError(t.t, err)
	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	exist, err := t.repo.IsExist(t.ctx, savedEntity.GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, true, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _19IsExist_WithNotExistedFieldCondition_Error() {
	t.DataCleanUp()
	savedEntity, err := t.repo.Insert(t.ctx, t.CreatePredefinedTestEntity())
	assert.NoError(t.t, err)
	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("NotExistedField", "==", "empty")
	_, err = t.repo.IsExist(t.ctx, savedEntity.GetID(), queryBuilder)
	assert.Error(t.t, err)
	//make sure that error is not "not found" error
	assert.Equal(t.t, false, errors.As(err, errors.NotFound))
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _20IsExist_WithExistedFieldCondition_Exist() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("SecondString", "==", "second")
	exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, true, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _21IsExist_WithExistedFieldCondition_NotExist() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("Number", "==", 1)
	exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, false, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _22IsExist_WithMultipleExistedFieldCondition_Exist() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("SecondString", "==", "second").
		Where("SearchString", "LIKE", "%example%")
	exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, true, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _23IsExist_WithMultipleExistedFieldCondition_NotExist() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("SecondString", "!=", "second").
		Where("SearchString", "LIKE", "%example%")
	exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, false, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _24IsExist_WithORMultipleExistedFieldCondition_NotExist() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Or("SecondString", "==", "notsecond").
		Or("String", "!=", "empty")
	exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, false, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _25IsExist_WithORMultipleExistedFieldCondition_Exist() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Or("SecondString", "==", "second").
		Or("String", "!=", "empty")
	exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, true, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _26IsExist_WithExistedNullFieldCondition_Exist() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("NullDate", "==", nil)
	exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, true, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _27IsExist_WithExistedNullFieldCondition_NotExist() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("NullDate", "!=", nil)
	exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, false, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _28IsExist_WithExistedNullableFieldCondition_Exist() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("NullableDate", "!=", nil)
	exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, true, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _29IsExist_WithExistedNullableFieldCondition_NotExist() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("NullableDate", "==", nil)
	exist, err := t.repo.IsExist(t.ctx, entities[1].GetID(), queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, false, exist)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _30Count_All() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}

	count, err := t.repo.Count(t.ctx, nil)
	assert.NoError(t.t, err)
	assert.Equal(t.t, 3, count)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _31Count_WithNotExistedFieldCondition_Error() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}
	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("NotExistedField", "==", "empty")
	_, err := t.repo.Count(t.ctx, queryBuilder)
	assert.Error(t.t, err)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _32Count_ById_OneElem() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}
	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("Id", "==", entities[0].GetID())
	count, err := t.repo.Count(t.ctx, queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, 1, count)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _33Count_ById_TwoElems() {
	t.DataCleanUp()
	entities := t.Create3PredefinedTestEntities()
	if len(entities) != 3 {
		return
	}
	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("Id", "==", entities[0].GetID()).
		Or("Id", "==", entities[1].GetID())
	count, err := t.repo.Count(t.ctx, queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, 2, count)
	t.DataCleanUp()
}

func (t *GenericRepositoryTester[IDType, EntityType]) _34Count_ByIteratorRange() {
	var entitiies []EntityType
	for i := 0; i < 20; i++ {
		entity, err := t.repo.Insert(t.ctx, t.CreatePredefinedTestEntityWithIteration(i))
		assert.NoError(t.t, err)
		entitiies = append(entitiies, entity)
	}
	count, err := t.repo.Count(t.ctx, nil)
	assert.NoError(t.t, err)
	assert.Equal(t.t, 20, count)
	queryBuilder := t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("Iterator", ">", 4).
		Where("Iterator", "<", 11)
	countIter, err := t.repo.Count(t.ctx, queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, 6, countIter)
	queryBuilder = t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("Iterator", ">=", 4).
		Where("Iterator", "<=", 11)
	countIterFirst, err := t.repo.Count(t.ctx, queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, 8, countIterFirst)
	queryBuilder = t.repo.NewQueryBuilder(t.ctx)
	queryBuilder.Where("Iterator", "<=", 4).
		Or("Iterator", ">=", 11)
	countIterSecond, err := t.repo.Count(t.ctx, queryBuilder)
	assert.NoError(t.t, err)
	assert.Equal(t.t, 14, countIterSecond)
}
