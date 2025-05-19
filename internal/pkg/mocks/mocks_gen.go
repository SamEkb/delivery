package mocks

//go:generate mockery --dir=../../core/ports --disable-version-string --with-expecter --name UnitOfWork --output ./ --filename uow_mock.go
//go:generate mockery --dir=../../core/ports --disable-version-string --with-expecter --name CourierRepository --output=. --filename courier_repository_mock.go
//go:generate mockery --dir=../../core/ports --disable-version-string --with-expecter --name OrderRepository --output=. --filename order_repository_mock.go
//go:generate mockery --dir=../../core/ports --disable-version-string --with-expecter --name GeoServiceClient --output ./ --filename geo_mock.go
