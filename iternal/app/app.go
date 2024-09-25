package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LaughG33k/chatWSService/iternal/client/mongo"
	rediscl "github.com/LaughG33k/chatWSService/iternal/client/redis"
	"github.com/LaughG33k/chatWSService/iternal/handler"
	messagesrepository "github.com/LaughG33k/chatWSService/iternal/repository/messagesRepository"
	"github.com/LaughG33k/chatWSService/pkg"

	"github.com/redis/go-redis/v9"
)

func Run() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpServer, err := initHttpServer(
		"127.0.0.1",
		"8081",
		time.Second*45,
		time.Second*50,
		time.Second*45,
		1000,
		1000,
	)

	if err != nil {
		pkg.Log.Fatal(err)
		return
	}

	redisCliet := rediscl.NewClient(ctx, &redis.Options{
		Addr:           "127.0.0.1:6379",
		DB:             0,
		MaxRetries:     5,
		MaxActiveConns: 1000,
		PoolSize:       1000,
		MaxIdleConns:   200,
		DialTimeout:    10 * time.Second,
		ReadTimeout:    1 * time.Minute,
		WriteTimeout:   45 * time.Second,
		PoolTimeout:    1 * time.Minute,
	}, 3*time.Second)

	redisCliet.Start()

	tmMongo, cancM := context.WithTimeout(ctx, 45*time.Second)
	defer cancM()

	mongoClient, err := mongo.NewMongoClient(tmMongo, mongo.MongoClientConfig{
		Host:                 "127.0.0.1",
		Port:                 "27017",
		Db:                   "messages",
		BulkWriteTimeSleep:   3 * time.Second,
		HealthCheakTimeSleep: 15 * time.Second,
		ReconectTimeSleep:    10 * time.Second,
		RecconectAttempts:    6,
		OperationTimeout:     45 * time.Second,
		MaxPoolSize:          250,
		RetryWrites:          true,
		RetryReads:           true,
	})

	if err != nil {
		pkg.Log.Fatal(err)
		return
	}

	msgRepo := messagesrepository.NewRepository(mongoClient, "messages")

	wsChatHandler := handler.NewWsChatHandler(ctx, redisCliet, msgRepo)

	wsChatHandler.StartHandler()

	if err := httpServer.ListenAndServe(); err != nil {
		log.Panic(err)
	}

	<-ctx.Done()

	tm, canc := context.WithTimeout(context.Background(), 40*time.Second)
	defer canc()
	pkg.C.Close(tm)

}
