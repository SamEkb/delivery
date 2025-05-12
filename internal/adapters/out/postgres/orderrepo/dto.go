package orderrepo

import (
	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
)

type OrderDTO struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey"`
	CourierID *uuid.UUID      `gorm:"type:uuid;index"`
	Location  kernel.Location `gorm:"embedded;embeddedPrefix:location_"`
	Volume    int
	Status    order.Status `gorm:"type:varchar(20)"`
}

type LocationDTO struct {
	X int
	Y int
}

func (OrderDTO) TableName() string {
	return "orders"
}
