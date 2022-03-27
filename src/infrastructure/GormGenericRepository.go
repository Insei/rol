package infrastructure

import (
	"rol/app/interfaces"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormGenericRepository[EntityType interfaces.IEntityModel] struct {
	Db     *gorm.DB
	DbRoot *gorm.DB
	logger *logrus.Logger
}

func NewGormGenericRepository[EntityType interfaces.IEntityModel](db *gorm.DB, log *logrus.Logger) (*GormGenericRepository[EntityType], error) {
	model := new(EntityType)
	return &GormGenericRepository[EntityType]{
		DbRoot: db,
		Db:     db.Model(&model).Preload(clause.Associations),
		logger: log,
	}, nil
}

func (ger *GormGenericRepository[EntityType]) NewQueryBuilder() interfaces.IQueryBuilder {
	return NewGormQueryBuilder()
}

func (ger *GormGenericRepository[EntityType]) addQueryToGorm(gormQuery *gorm.DB, queryBuilder interfaces.IQueryBuilder) error {
	if queryBuilder != nil {
		query, err := queryBuilder.Build()
		if err != nil {
			return err
		}
		arrQuery := query.([]interface{})
		// TODO: We need more checks here
		switch arrQuery[0].(type) {
		case string:
			queryString := arrQuery[0].(string)
			queryArgs := make([]interface{}, 0)
			for i := 1; i < len(arrQuery); i++ {
				queryArgs = append(queryArgs, arrQuery[i])
			}
			gormQuery.Where(queryString, queryArgs...)
		}
	}
	return nil
}

func (ger *GormGenericRepository[EntityType]) GetList(orderBy string, orderDirection string, page int, size int, queryBuilder interfaces.IQueryBuilder) (*[]EntityType, error) {
	entities := &[]EntityType{}
	offset := int64((page - 1) * size)
	orderString := generateOrderString(orderBy, orderDirection)
	gormQuery := ger.Db.Order(orderString)
	err := ger.addQueryToGorm(gormQuery, queryBuilder)
	if err != nil {
		return entities, err
	}
	err = gormQuery.Offset(int(offset)).Limit(size).Find(entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (ger *GormGenericRepository[EntityType]) Count(queryBuilder interfaces.IQueryBuilder) (int64, error) {
	count := int64(0)
	model := new(EntityType)
	gormQuery := ger.DbRoot.Model(&model)
	err := ger.addQueryToGorm(gormQuery, queryBuilder)
	if err != nil {
		return 0, err
	}
	err = gormQuery.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ger *GormGenericRepository[EntityType]) GetById(id uint) (*EntityType, error) {
	entity := new(EntityType)
	err := ger.Db.First(entity, id).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (ger *GormGenericRepository[EntityType]) Update(entity *EntityType) error {
	return ger.Db.Save(entity).Error
}

func (ger *GormGenericRepository[EntityType]) Insert(entity EntityType) (uint, error) {
	if err := ger.Db.Create(&entity).Error; err != nil {
		return 0, err
	}
	return entity.GetId(), nil
}

func (ger *GormGenericRepository[EntityType]) Delete(id uint) error {
	entity := new(EntityType)
	return ger.Db.Select(clause.Associations).Delete(entity, id).Error
}
