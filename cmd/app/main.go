package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/delivery/internal/generated/servers"
	_ "github.com/lib/pq"

	"github.com/delivery/cmd"
	"github.com/delivery/internal/adapters/out/postgres/courierrepo"
	"github.com/delivery/internal/adapters/out/postgres/orderrepo"
	"github.com/delivery/internal/pkg/errs"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config := getConfigs()

	connectionString, err := makeConnectionString(
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPassword,
		config.DbName,
		config.DbSslMode)
	if err != nil {
		log.Fatal(err.Error())
	}

	crateDbIfNotExists(config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPassword,
		config.DbName,
		config.DbSslMode)
	gormDb := mustGormOpen(connectionString)
	mustAutoMigrate(gormDb)

	compositionRoot := cmd.NewCompositionRoot(
		config,
		gormDb,
	)

	startKafkaConsumer(compositionRoot)
	//startCronJobs(compositionRoot)
	startWebServer(compositionRoot, config.HttpPort)
}

func getConfigs() *cmd.Config {
	return &cmd.Config{
		HttpPort:                  goDotEnvVariable("HTTP_PORT"),
		DbHost:                    goDotEnvVariable("DB_HOST"),
		DbPort:                    goDotEnvVariable("DB_PORT"),
		DbUser:                    goDotEnvVariable("DB_USER"),
		DbPassword:                goDotEnvVariable("DB_PASSWORD"),
		DbName:                    goDotEnvVariable("DB_NAME"),
		DbSslMode:                 goDotEnvVariable("DB_SSLMODE"),
		GeoServiceGrpcHost:        goDotEnvVariable("GEO_SERVICE_GRPC_HOST"),
		KafkaHost:                 goDotEnvVariable("KAFKA_HOST"),
		KafkaConsumerGroup:        goDotEnvVariable("KAFKA_CONSUMER_GROUP"),
		KafkaBasketConfirmedTopic: goDotEnvVariable("KAFKA_BASKET_CONFIRMED_TOPIC"),
		KafkaOrderChangedTopic:    goDotEnvVariable("KAFKA_ORDER_CHANGED_TOPIC"),
	}
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func crateDbIfNotExists(host string, port string, user string,
	password string, dbName string, sslMode string) {
	dsn, err := makeConnectionString(host, port, user, password, "postgres", sslMode)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer db.Close()

	//_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	//if err != nil {
	//	log.Printf("Ошибка создания БД (возможно, уже существует): %v", err)
	//}
}

func makeConnectionString(host string, port string, user string,
	password string, dbName string, sslMode string) (string, error) {
	if host == "" {
		return "", errs.NewValueIsRequiredError("host")
	}
	if port == "" {
		return "", errs.NewValueIsRequiredError("port")
	}
	if user == "" {
		return "", errs.NewValueIsRequiredError("user")
	}
	if password == "" {
		return "", errs.NewValueIsRequiredError("password")
	}
	if dbName == "" {
		return "", errs.NewValueIsRequiredError("dbName")
	}
	if sslMode == "" {
		return "", errs.NewValueIsRequiredError("sslMode")
	}
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		host,
		port,
		user,
		password,
		dbName,
		sslMode), nil
}

func mustGormOpen(connectionString string) *gorm.DB {
	pgGorm, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN:                  connectionString,
			PreferSimpleProtocol: true,
		},
	), &gorm.Config{})
	if err != nil {
		log.Fatalf("connection to postgres through gorm\n: %s", err)
	}
	return pgGorm
}

func mustAutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&courierrepo.CourierDto{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	err = db.AutoMigrate(&courierrepo.StoragePlaceDto{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	err = db.AutoMigrate(&orderrepo.OrderDTO{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

}

func startWebServer(compositionRoot cmd.CompositionRoot, port string) {
	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy")
	})

	servers.RegisterHandlers(e, compositionRoot.Servers.HttpServer)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
}

func startCronJobs(compositionRoot cmd.CompositionRoot) {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddJob("* * * * * *", &compositionRoot.Jobs.AssignOrderJob)
	if err != nil {
		log.Fatalf("failed to add assign order job: %v", err)
	}
	_, err = c.AddJob("* * * * * *", &compositionRoot.Jobs.MoveCourierJob)
	if err != nil {
		log.Fatalf("failed to add move courier job: %v", err)
	}

	c.Start()
}

func startKafkaConsumer(compositionRoot cmd.CompositionRoot) {
	go func() {
		if err := compositionRoot.KafkaConsumer.Consume(); err != nil {
			log.Fatalf("Kafka consumer error: %v", err)
		}
	}()
}
