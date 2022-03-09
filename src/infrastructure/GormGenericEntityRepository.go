package infrastructure

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/domain/entities"
	"time"
)

type GormGenericEntityRepository struct {
	Db     *gorm.DB
	logger *logrus.Logger
}

func NewGormGenericEntityRepository(dialector gorm.Dialector, log *logrus.Logger) (*GormGenericEntityRepository, error) {
	db, err := gorm.Open(dialector, &gorm.Config{})
	db.AutoMigrate(
		&entities.EthernetSwitch{},
		&entities.EthernetSwitchPort{},
	)
	return &GormGenericEntityRepository{
		Db:     db,
		logger: log,
	}, err
}

func generateOrderString(orderBy string, orderDirection string) string {
	order := ""
	if len(orderBy) > 0 {
		order = orderBy
		if len(orderDirection) > 0 {
			order = order + " " + orderDirection
		}
	}
	if len(order) < 1 {
		order = "id"
	}
	return order
}

func (ger *GormGenericEntityRepository) GetList(entities interface{}, orderBy string, orderDirection string, page int, size int, query string, args ...interface{}) (int64, error) {
	offset, count := int64((page-1)*size), int64(0)
	entityModel, orderString := mappers.GetEmptyEntityFromArrayType(entities), generateOrderString(orderBy, orderDirection)
	gormQuery := ger.Db.Model(entityModel).Order(orderString)
	date, _ := time.Parse("2006-01-02", "1999-01-01")
	gormQuery = gormQuery.Where("deleted_at < ?", date)
	if len(query) > 0 {
		if len(args) > 0 {
			gormQuery = gormQuery.Where(query, args)
		} else {
			gormQuery = gormQuery.Where(query)
		}
	}
	gormQuery = gormQuery.Count(&count)
	if count < offset {
		offset = 0
	}
	return count, gormQuery.Offset(int(offset)).Limit(size).Find(entities).Error
}

func (ger *GormGenericEntityRepository) GetAll(entities interface{}) error {
	date, _ := time.Parse("2006-01-02", "1999-01-01")
	return ger.Db.Preload(clause.Associations).Where("deleted_at < ?", date).Find(entities).Error
}

func (ger *GormGenericEntityRepository) GetById(entity interfaces.IEntityModel, id uint) error {
	date, _ := time.Parse("2006-01-02", "1999-01-01")
	return ger.Db.Preload(clause.Associations).Where("deleted_at < ?", date).First(entity, id).Error
}

func (ger *GormGenericEntityRepository) Update(entity interfaces.IEntityModel) error {
	return ger.Db.Save(entity).Error
}

func (ger *GormGenericEntityRepository) Insert(entity interfaces.IEntityModel) (uint, error) {
	if err := ger.Db.Create(entity).Error; err != nil {
		return 0, err
	}
	return entity.GetId(), nil
}

func (ger *GormGenericEntityRepository) Delete(entity interfaces.IEntityModel) error {
	return ger.Db.Select(clause.Associations).Delete(entity).Error
}
