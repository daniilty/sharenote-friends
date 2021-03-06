package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/daniilty/sharenote-friends/internal/core"
	"github.com/daniilty/sharenote-friends/internal/kafka"
	"github.com/daniilty/sharenote-friends/internal/mongo"
	"github.com/daniilty/sharenote-friends/internal/server"
	"github.com/daniilty/sharenote-friends/internal/users"
	schema "github.com/daniilty/sharenote-grpc-schema"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	exitCodeInitError = 2
)

func run() error {
	cfg, err := loadEnvConfig()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	conn, err := grpc.DialContext(ctx, cfg.usersGRPCAddr, grpc.WithInsecure())
	if err != nil {
		cancel()

		return err
	}

	mongoClient, err := mongo.Connect(context.Background(), cfg.mongoConnString)
	if err != nil {
		cancel()

		return err
	}

	db := mongoClient.Database(cfg.mongoDBName)
	friendsCollection := db.Collection(cfg.mongoFriendsCollectionName)
	friendRequestsCollection := db.Collection(cfg.mongoFriendRequestsCollectionName)

	d := mongo.NewDBImpl(db, friendRequestsCollection, friendsCollection)

	client := schema.NewUsersClient(conn)
	service := core.NewService(d, client)

	loggerCfg := zap.NewProductionConfig()

	logger, err := loggerCfg.Build()
	if err != nil {
		cancel()

		return err
	}

	httpServer := server.NewHTTP(cfg.httpAddr, logger.Sugar(), service)

	consumer := kafka.NewConsumerImpl(cfg.kafkaTopic, []string{cfg.kafkaBroker}, cfg.kafkaGroupID)

	usersHandler := users.NewEventsHandler(logger.Sugar(), time.Duration(cfg.eventsTimeout)*time.Second, d, consumer)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		httpServer.Run(ctx)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		usersHandler.Listen(ctx)
		wg.Done()
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-termChan
	cancel()
	wg.Wait()

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(exitCodeInitError)
	}
}
