package router

import (
	"context"
	"fmt"
	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"good-template-go/conf"
	ginSwaggerDocs "good-template-go/docs"
	"good-template-go/pkg/errors"
	handlers "good-template-go/pkg/handler"
	"good-template-go/pkg/middlewares"
	"good-template-go/pkg/repo"
	"good-template-go/pkg/service"
	"good-template-go/pkg/utils/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"net/http"
	"strings"
	"time"
)

type kafkaProcessHandler interface {
	KafkaProcess(ctx context.Context, message kafka.Message) error
}

func ApplicationV1Router(router *gin.Engine, db *gorm.DB) {
	// Router
	routerV1 := router.Group("/api/v1")

	// Init repo
	todoRepo := repo.NewRepoTodo(db)

	// Init service
	todoService := service.NewTodoService(todoRepo)

	// Init handler
	migrateHandler := handlers.NewMigrationHandler(db)
	todoHandler := handlers.NewTodoHandler(todoService)

	// APIs

	// Internal apis
	internalRoutes := routerV1.Group("/internal", middlewares.AuthAdminJWTMiddleware())
	{
		internalRoutes.POST("/migrate-public", migrateHandler.MigratePublic)
		internalRoutes.POST("/migrate-log", migrateHandler.MigrateLog)
	}

	// Todo apis
	todoRoutes := routerV1.Group("/todo", middlewares.AuthManagerJWTMiddleware())
	{
		todoRoutes.POST("/create", middlewares.AuthManagerJWTMiddleware(), todoHandler.Create)
		todoRoutes.POST("/get-one/:id", middlewares.AuthManagerJWTMiddleware(), todoHandler.Get)
		todoRoutes.POST("/get-list", middlewares.AuthManagerJWTMiddleware(), todoHandler.List)
		todoRoutes.POST("/update/:id", middlewares.AuthManagerJWTMiddleware(), todoHandler.Update)
		todoRoutes.POST("/delete/:id", middlewares.AuthManagerJWTMiddleware(), todoHandler.Delete)
	}

	// Swagger
	ginSwaggerDocs.SwaggerInfo.Host = conf.GetConfig().SwaggerHost
	ginSwaggerDocs.SwaggerInfo.Title = conf.GetConfig().AppName
	ginSwaggerDocs.SwaggerInfo.BasePath = routerV1.BasePath()
	ginSwaggerDocs.SwaggerInfo.Version = "v1"
	ginSwaggerDocs.SwaggerInfo.Schemes = []string{"http", "https"}

	routerV1.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.PersistAuthorization(true),
	))

	// Kafka Handler
	//kafkaHandlersTodo := kafkaHandlers.NewTodoKafkaHandlers(todoRepo)
	//
	//_ = map[string]kafkaProcessHandler{
	//	utils.TodoTopicPrefix + model.Todo{}.TableName(): kafkaHandlersTodo,
	//}

	//if len(kafkaTopic) > 0 {
	//	// call migrate
	//	err := migrateHandler.MigrateDatabase()
	//	if err != nil {
	//		log.Fatalf("error migrating tables: %v", err)
	//	}
	//}

	//for topic, handler := range kafkaTopic {
	//	fmt.Printf("Fetching message for topic: %s \n", topic)
	//	go KafkaConsumer(topic, handler)
	//}
}

func NewRoute() {
	// Log
	logger.Init(conf.GetConfig().AppName)
	logger.DefaultLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		PadLevelText:     true,
		ForceQuote:       true,
		QuoteEmptyFields: true,
	})

	// GetDB
	db := GetDBPostgres()

	// Cors
	router := gin.Default()
	router.Use(limit.MaxAllowed(200))
	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{"*"}
	router.Use(cors.New(configCors))

	//
	router.Use(errors.ErrorHandlerMiddleware)
	ApplicationV1Router(router, db)
	startServer(router)
}

func startServer(router http.Handler) {
	s := &http.Server{
		Addr:           ":" + conf.GetConfig().Port,
		Handler:        router,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Server running on:", conf.GetConfig().Host+":"+conf.GetConfig().Port)
	if err := s.ListenAndServe(); err != nil {
		_ = fmt.Errorf("fatal error description: %s", strings.ToLower(err.Error()))
		panic(err)
	}
}

func GetDBPostgres() *gorm.DB {
	dsn := postgres.Open(fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable connect_timeout=5",
		conf.GetConfig().DBHost,
		conf.GetConfig().DBPort,
		conf.GetConfig().DBUser,
		conf.GetConfig().DBName,
		conf.GetConfig().DBPass,
	))
	db, err := gorm.Open(dsn, &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{
			SingularTable: true,
			//TablePrefix:   conf.GetConfig().DBSchema + ".",
		},
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		log.Fatalf("error opening connection to database: %v", err)
	}

	conn, err := db.DB()
	if err != nil {
		log.Fatalf("error initializing database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err = conn.PingContext(ctx); err != nil {
		log.Fatalf("error opening connection to database: %v", err)
	}
	log.Println("Postgres connected!")

	return db
}

func KafkaConsumer(topic string, handlers kafkaProcessHandler) {
	brokers := strings.Split(conf.GetConfig().KafkaBrokers, ";")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		Partition: 0,
		GroupID:   conf.GetConfig().AppName,
		GroupBalancers: []kafka.GroupBalancer{
			kafka.RangeGroupBalancer{},
			kafka.RoundRobinGroupBalancer{},
		},
	})

	var (
		m   kafka.Message
		err error
	)

	for {
		m, err = r.FetchMessage(context.Background())
		if err != nil {
			break
		}

		// commit empty message
		ctx := context.Background()
		if string(m.Value) == "" {
			err = r.CommitMessages(ctx, m)
		} else {

			// process kafka message
			if err = handlers.KafkaProcess(ctx, m); err != nil {
				message := fmt.Sprintf("error process message at offset %d, topic %s: %v", m.Offset, m.Topic, err)
				log.Printf(message)
				go func() {
					w := kafka.Writer{
						Addr:                   kafka.TCP(fmt.Sprintf("%s:%s", conf.GetConfig().KafkaHost, conf.GetConfig().KafkaPort)),
						Topic:                  "error_log",
						AllowAutoTopicCreation: true,
					}

					err = w.WriteMessages(context.Background(), kafka.Message{
						Key:   []byte("error"),
						Value: []byte(message),
					})
					if err != nil {
						log.Printf("failed to log error at offset %d, topic %s: %v:", m.Offset, m.Topic, err)
					}
					if err = w.Close(); err != nil {
						log.Printf("failed to close writer: %v", err)
					}
				}()

				continue
			} else {
				err = r.CommitMessages(ctx, m)
			}
		}
		if err != nil {
			log.Fatalf("failed to commit message at offset %d, topic %s: %v:", m.Offset, m.Topic, err)
		}
	}

	if err = r.Close(); err != nil {
		log.Fatalf("failed to close reader: %v", err)
	}
}
