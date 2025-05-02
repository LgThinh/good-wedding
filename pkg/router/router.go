package router

import (
	"context"
	"fmt"
	limit "github.com/aviddiviner/gin-limit"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"good-wedding/conf"
	ginSwaggerDocs "good-wedding/docs"
	"good-wedding/pkg/errors"
	handlers "good-wedding/pkg/handler"
	"good-wedding/pkg/middlewares"
	"good-wedding/pkg/repo"
	"good-wedding/pkg/service"
	"good-wedding/pkg/utils/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"net/http"
	"strings"
	"time"
)

func ApplicationV1Router(router *gin.Engine, db *gorm.DB, s3Bucket *s3.S3) {
	// Router
	routerV1 := router.Group("/api/v1")

	// Init repo
	weddingRepo := repo.NewWeddingRepo(db, s3Bucket)

	// Init service
	weddingService := service.NewWeddingService(weddingRepo)

	// Init handler
	migrateHandler := handlers.NewMigrationHandler(db)
	weddingHandler := handlers.NewWeddingHandler(weddingService)

	// APIs

	// Internal apis
	internalRoutes := routerV1.Group("/internal", middlewares.AuthAdminJWTMiddleware())
	{
		internalRoutes.POST("/migrate-public", migrateHandler.MigratePublic)
		internalRoutes.POST("/migrate-log", migrateHandler.MigrateLog)
	}

	// Wedding apis
	weddingRoutes := routerV1.Group("/wedding")
	{
		weddingRoutes.POST("/upload-image", weddingHandler.UploadImage)
		weddingRoutes.POST("/upload-video", weddingHandler.UploadVideo)
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

	// Get s3
	s3Bucket := CreateS3Session()

	// Cors
	router := gin.Default()
	router.Use(limit.MaxAllowed(200))
	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{"*"}
	router.Use(cors.New(configCors))

	//
	router.Use(errors.ErrorHandlerMiddleware)
	ApplicationV1Router(router, db, s3Bucket)
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

func CreateS3Session() *s3.S3 {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf.GetConfig().AWSRegion),
		Credentials: credentials.NewStaticCredentials(conf.GetConfig().AWSAccessKey, conf.GetConfig().AWSSecretKey, ""),
	})
	if err != nil {
		log.Fatal("Unable to create AWS session: ", err)
	}

	log.Printf("AWS S3 connected!")
	return s3.New(sess)
}
