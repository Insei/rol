package domain

import (
	"time"

	"gorm.io/gorm"
)

//BaseEntity structure
type Entity[IDType comparable] struct {
	//	ID - entity identifier
	ID IDType `gorm:"primary_key;"`
	//	CreatedAt - time when the entity was created
	CreatedAt time.Time `gorm:"index"`
	//	UpdatedAt - time when the entity was updated
	UpdatedAt *time.Time `gorm:"index"`
	//	DeletedAt - time when the entity was deleted
	DeletedAt gorm.DeletedAt
}

//GetID of the entity
func (e Entity[IDType]) GetID() IDType {
	return e.ID
}

//GetCreatedAt time
func (e Entity[IDType]) GetCreatedAt() time.Time {
	return e.CreatedAt
}

//GetUpdatedAt time
func (e Entity[IDType]) GetUpdatedAt() *time.Time {
	return e.UpdatedAt
}

//GetDeletedAt time
func (e Entity[IDType]) GetDeletedAt() *time.Time {
	return &e.DeletedAt.Time
}
