package cmd

import (
	"log"
	"time"

	"github.com/delivery/internal/adapters/in/http"
	"github.com/delivery/internal/adapters/in/jobs"
	consumer "github.com/delivery/internal/adapters/in/kafka"
	"github.com/delivery/internal/adapters/out/grpc/geo"
	producer "github.com/delivery/internal/adapters/out/kafka"
	"github.com/delivery/internal/adapters/out/postgres"
	"github.com/delivery/internal/core/application/eventhandlers"
	"github.com/delivery/internal/core/application/usecases/commands"
	"github.com/delivery/internal/core/application/usecases/queries"
	"github.com/delivery/internal/core/domain/model/order"
	"github.com/delivery/internal/core/domain/service"
	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/pkg/ddd"
	"gorm.io/gorm"
)

type CompositionRoot struct {
	config          *Config
	gormDb          *gorm.DB
	DomainServices  DomainServices
	Repositories    Repositories
	CommandHandlers CommandHandlers
	QueryHandlers   QueryHandlers
	Servers         Servers
	Jobs            Jobs
	KafkaConsumer   consumer.BasketConfirmedConsumer
	KafkaProducer   ports.OrderProducer
	EventHandler    ddd.EventHandler
	Mediatr         ddd.Mediatr
}

type DomainServices struct {
	DispatchService service.DispatchService
}

type Repositories struct {
	UnitOfWork        ports.UnitOfWork
	OrderRepository   ports.OrderRepository
	CourierRepository ports.CourierRepository
}

type CommandHandlers struct {
	AssignOrderCommandHandler   commands.AssignOrderHandler
	CreateOrderCommandHandler   commands.CreateOrderHandler
	CreateCourierCommandHandler commands.CreateCourierHandler
	MoveCourierCommandHandler   commands.MoveCourierHandler
}

type QueryHandlers struct {
	GetAllCouriersQueryHandler        queries.GetAllCouriersHandler
	GetNotCompletedOrdersQueryHandler queries.GetAllUncompletedOrdersHandler
}

type Servers struct {
	HttpServer *http.Server
}

type Jobs struct {
	AssignOrderJob jobs.AssignOrderJob
	MoveCourierJob jobs.MoveCourierJob
}

func NewCompositionRoot(config *Config, gormDb *gorm.DB) CompositionRoot {
	mediatr := ddd.NewMediatr()
	unitOfWork, err := postgres.NewUnitOfWork(gormDb, mediatr)
	if err != nil {
		log.Fatalf("failed to create unit of work: %v", err)
	}

	// Services
	dispatchService := service.NewDispatchService()

	// Repositories
	orderRepository := unitOfWork.OrderRepository()
	courierRepository := unitOfWork.CourierRepository()

	// Clients
	geoClient, err := geo.NewGeoClient(config.GeoServiceGrpcHost, 5*time.Second)
	if err != nil {
		log.Fatalf("failed to create geo service client: %v", err)
	}

	// Command Handlers
	createOrderCommandHandler, err := commands.NewAddCreateOrderHandler(unitOfWork, geoClient)
	if err != nil {
		log.Fatalf("failed to create create order command handler: %v", err)
	}

	assignOrderCommandHandler, err := commands.NewAssignOrderHandler(unitOfWork, dispatchService)
	if err != nil {
		log.Fatalf("failed to create assign order command handler: %v", err)
	}

	moveCourierCommandHandler, err := commands.NewMoveCourierHandler(unitOfWork)
	if err != nil {
		log.Fatalf("failed to create move courier command handler: %v", err)
	}

	createCourierCommandHandler, err := commands.NewCreateCourierHandler(unitOfWork)
	if err != nil {
		log.Fatalf("failed to create create courier command handler: %v", err)
	}

	// Queries
	getAllCouriersQueryHandler, err := queries.NewGetAllCouriersHandler(unitOfWork)
	if err != nil {
		log.Fatalf("failed to create get all couriers query handler: %v", err)
	}

	getNotCompletedOrdersQueryHandler, err := queries.NewGetAllUncompletedOrdersHandler(unitOfWork)
	if err != nil {
		log.Fatalf("failed to create get not completed orders query handler: %v", err)
	}

	// Jobs
	assignOrderJob, err := jobs.NewAssignOrderJob(assignOrderCommandHandler)
	if err != nil {
		log.Fatalf("failed to create assign order job: %v", err)
	}

	moveCourierJob, err := jobs.NewMoveCourierJob(moveCourierCommandHandler)
	if err != nil {
		log.Fatalf("failed to create move courier job: %v", err)
	}

	// Kafka Consumer
	kafkaConsumer, err := consumer.NewConsumer(
		[]string{config.KafkaHost},
		config.KafkaBasketConfirmedTopic,
		config.KafkaConsumerGroup,
		createOrderCommandHandler,
	)
	if err != nil {
		log.Fatalf("failed to create kafka consumer: %v", err)
	}

	// Kafka Producer
	kafkaProducer, err := producer.NewOrderStatusChangedProducer(
		[]string{config.KafkaHost},
		config.KafkaOrderChangedTopic,
	)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}

	//Handler
	handler, err := eventhandlers.NewOrderStatusChangedEventHandler(kafkaProducer)
	if err != nil {
		log.Fatalf("failed to create order status changed event handler: %v", err)
	}

	// Mediatr
	event := order.NewStatusChangedDomainEventWithoutData()
	mediatr.Subscribe(handler, event)

	// Servers
	httpServer, err := http.NewServer(
		assignOrderCommandHandler,
		createOrderCommandHandler,
		createCourierCommandHandler,
		getAllCouriersQueryHandler,
		getNotCompletedOrdersQueryHandler,
	)

	return CompositionRoot{
		config: config,
		gormDb: gormDb,
		DomainServices: DomainServices{
			DispatchService: dispatchService,
		},
		Repositories: Repositories{
			UnitOfWork:        unitOfWork,
			OrderRepository:   orderRepository,
			CourierRepository: courierRepository,
		},
		CommandHandlers: CommandHandlers{
			AssignOrderCommandHandler:   assignOrderCommandHandler,
			CreateOrderCommandHandler:   createOrderCommandHandler,
			CreateCourierCommandHandler: createCourierCommandHandler,
			MoveCourierCommandHandler:   moveCourierCommandHandler,
		},
		QueryHandlers: QueryHandlers{
			GetAllCouriersQueryHandler:        getAllCouriersQueryHandler,
			GetNotCompletedOrdersQueryHandler: getNotCompletedOrdersQueryHandler,
		},
		Jobs: Jobs{
			AssignOrderJob: *assignOrderJob,
			MoveCourierJob: *moveCourierJob,
		},
		Servers: Servers{
			HttpServer: httpServer,
		},
		KafkaConsumer: kafkaConsumer,
		KafkaProducer: kafkaProducer,
		EventHandler:  handler,
		Mediatr:       mediatr,
	}
}
