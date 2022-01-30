package main

import (
	"fmt"
	"os"
	"strconv"
)

type envConfig struct {
	httpAddr                          string
	usersGRPCAddr                     string
	mongoConnString                   string
	mongoDBName                       string
	mongoFriendsCollectionName        string
	mongoFriendRequestsCollectionName string
	kafkaBroker                       string
	kafkaTopic                        string
	kafkaGroupID                      string
	eventsTimeout                     int
}

func loadEnvConfig() (*envConfig, error) {
	var err error

	cfg := &envConfig{}

	cfg.httpAddr, err = lookupEnv("HTTP_SERVER_ADDR")
	if err != nil {
		return nil, err
	}

	cfg.usersGRPCAddr, err = lookupEnv("USERS_GRPC_ADDR")
	if err != nil {
		return nil, err
	}

	cfg.mongoDBName, err = lookupEnv("MONGO_DB_NAME")
	if err != nil {
		return nil, err
	}

	cfg.mongoFriendsCollectionName, err = lookupEnv("MONGO_FRIENDS_COLLECTION_NAME")
	if err != nil {
		return nil, err
	}

	cfg.mongoFriendRequestsCollectionName, err = lookupEnv("MONGO_FRIEND_REQUESTS_COLLECTION_NAME")
	if err != nil {
		return nil, err
	}

	cfg.mongoConnString, err = lookupEnv("MONGO_CONN_STRING")
	if err != nil {
		return nil, err
	}

	cfg.kafkaBroker, err = lookupEnv("KAFKA_BROKER")
	if err != nil {
		return nil, err
	}

	cfg.kafkaTopic, err = lookupEnv("KAFKA_TOPIC")
	if err != nil {
		return nil, err
	}

	cfg.kafkaGroupID, err = lookupEnv("KAFKA_GROUP_ID")
	if err != nil {
		return nil, err
	}

	timeoutString, err := lookupEnv("TIMEOUT")
	if err != nil {
		return nil, err
	}

	cfg.eventsTimeout, err = strconv.Atoi(timeoutString)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func lookupEnv(name string) (string, error) {
	const provideEnvErrorMsg = `please provide "%s" environment variable`

	val, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf(provideEnvErrorMsg, name)
	}

	return val, nil
}
