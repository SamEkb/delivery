package orderrepo

import (
	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/delivery/internal/core/domain/model/order"
)

func DomainToDto(order *order.Order) OrderDTO {
	return OrderDTO{
		ID:        order.ID(),
		CourierID: order.CourierID(),
		Location:  order.Location(),
		Volume:    order.Volume(),
		Status:    order.Status(),
	}
}

func DtoToDomain(dto OrderDTO) *order.Order {
	var aggregate *order.Order
	location, _ := kernel.NewLocation(dto.Location.X(), dto.Location.Y())
	aggregate = order.RestoreOrder(dto.ID, dto.CourierID, location, dto.Volume, dto.Status)
	return aggregate
}
