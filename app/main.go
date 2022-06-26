package main

import (
	"fmt"
	"io"
	"os"
	"storage-api/app/helpers"
	"storage-api/app/services"
	httpDelivery "storage-api/media/delivery/http"
	cloudinaryrepo "storage-api/media/repository/cloudinary"
	googlebucket "storage-api/media/repository/google_bucket"
	mongoRepo "storage-api/media/repository/mongo"
	mediaUsecase "storage-api/media/usecase"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	_ = godotenv.Load()

	// logger
	logMaxSize, _ := strconv.Atoi(os.Getenv("LOG_MAX_SIZE"))
	if logMaxSize == 0 {
		logMaxSize = 50 //default 50 megabytes
	}
	lumberjackLog := &lumberjack.Logger{
		Filename:  "server.log",
		MaxSize:   logMaxSize,
		LocalTime: true,
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(io.MultiWriter(lumberjackLog, os.Stdout))
}

func main() {
	router := gin.New()
	gin.DefaultWriter = logrus.StandardLogger().Writer()
	timeout, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	t := time.Duration(timeout) * time.Second

	// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
	goenv := os.Getenv("GO_ENV")
	if goenv == "" {
		goenv = "production"
	}
	if err := sentry.Init(sentry.ClientOptions{
		Environment: goenv,
		Release:     "v1.0.0",
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	mongoDB := services.ConnectMongo(t)
	mr := mongoRepo.NewMongoRepo(mongoDB)
	gbr := googlebucket.NewGoogleBucketRepo()
	cr := cloudinaryrepo.NewCloudinaryRepo()
	mediaUsecase := mediaUsecase.NewUsecaseMedia(mr, gbr, cr, t)

	router.Use(helpers.RequestLogger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(helpers.CustomRecovery())

	router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	prefix := "/storage"
	httpDelivery.NewHandler(router, mediaUsecase, prefix)

	router.Run(":" + os.Getenv("PORT"))
}
