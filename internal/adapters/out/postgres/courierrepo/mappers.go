package courierrepo

import (
	"github.com/delivery/internal/core/domain/model/courier"
	"github.com/delivery/internal/core/domain/model/kernel"
)

func DomainToDto(courier *courier.Courier) CourierDto {
	return CourierDto{
		ID:            courier.ID(),
		Name:          courier.Name(),
		Speed:         courier.Speed(),
		Location:      LocationDTO{X: courier.Location().X(), Y: courier.Location().Y()},
		StoragePlaces: mapStoragePlaces(courier),
	}
}

func DtoToDomain(dto CourierDto) *courier.Courier {
	var storagePlaces []*courier.StoragePlace
	for _, dtoStoragePlace := range dto.StoragePlaces {
		storagePlace := courier.RestoreStoragePlace(
			dtoStoragePlace.ID,
			dtoStoragePlace.Name,
			dtoStoragePlace.TotalVolume,
			dtoStoragePlace.OrderID,
		)

		storagePlaces = append(storagePlaces, storagePlace)
	}
	location, _ := kernel.NewLocation(dto.Location.X, dto.Location.Y)
	return courier.RestoreCourier(dto.ID, dto.Name, dto.Speed, location, storagePlaces)
}

func mapStoragePlaces(courier *courier.Courier) []*StoragePlaceDto {
	storagePlacesDTO := make([]*StoragePlaceDto, 0, len(courier.StoragePlaces()))
	for _, storagePlace := range courier.StoragePlaces() {
		storagePlaceDTO := &StoragePlaceDto{
			ID:          storagePlace.ID(),
			OrderID:     storagePlace.OrderID(),
			Name:        storagePlace.Name(),
			TotalVolume: storagePlace.TotalVolume(),
			CourierID:   courier.ID(),
		}
		storagePlacesDTO = append(storagePlacesDTO, storagePlaceDTO)
	}
	return storagePlacesDTO
}
