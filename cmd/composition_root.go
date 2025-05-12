package cmd

import (
	"log"

	"github.com/delivery/internal/adapters/out/postgres"
	"github.com/delivery/internal/core/domain/service"
	"github.com/delivery/internal/core/ports"
	"gorm.io/gorm"
)

type CompositionRoot struct {
	config *Config
	gormDb *gorm.DB
}

func NewCompositionRoot(config *Config, gormDb *gorm.DB) CompositionRoot {
	return CompositionRoot{
		config: config,
		gormDb: gormDb,
	}
}

func (c *CompositionRoot) NewDispatcherService() service.DispatchService {
	return service.NewDispatchService()
}

func (c *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	unitOfWork, err := postgres.NewUnitOfWork(c.gormDb)
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}
	return unitOfWork
}
