package mongo

import "go.mongodb.org/mongo-driver/mongo"

type transactionFunc func(mongo.SessionContext) (interface{}, error)
