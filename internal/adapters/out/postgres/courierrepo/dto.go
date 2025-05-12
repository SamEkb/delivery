package courierrepo

import (
	"github.com/google/uuid"
)

type CourierDto struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name          string
	Speed         int
	Location      LocationDTO        `gorm:"embedded;embeddedPrefix:location_"`
	StoragePlaces []*StoragePlaceDto `gorm:"foreignKey:CourierID;constraint:OnDelete:CASCADE;"`
}

type StoragePlaceDto struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	TotalVolume int
	OrderID     *uuid.UUID `gorm:"type:uuid;"`
	CourierID   uuid.UUID  `gorm:"type:uuid;index"`
}

type LocationDTO struct {
	X int
	Y int
}

func (CourierDto) TableName() string {
	return "couriers"
}

func (StoragePlaceDto) TableName() string {
	return "storage_places"
}
