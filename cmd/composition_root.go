package cmd

import "github.com/delivery/internal/core/domain/service"

type CompositionRoot struct {
}

func NewCompositionRoot() CompositionRoot {
	return CompositionRoot{}
}

func (c CompositionRoot) NewDispatcherService() service.DispatchService {
	return service.NewDispatchService()
}
