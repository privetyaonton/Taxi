package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxi/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	client *mongo.Client
	cfg    *config.Config
}

type log struct {
	Level  string
	Caller string
	Msg    string
	Method string
	Uuid   string
	Err    string
	Time   string
}

func New(cfg *config.Config) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.GetMongoUrl()))
	if err != nil {
		return nil, fmt.Errorf("connect failed: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return &Mongo{client, cfg}, nil
}

func (m *Mongo) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("disconnect failed: %w", err)
	}
	return nil
}

func (m *Mongo) Write(p []byte) (n int, err error) {
	var logs log
	err = json.Unmarshal(p, &logs)
	if err != nil {
		return 0, fmt.Errorf("unmarshal failed: %w", err)
	}

	logger := m.client.Database(m.cfg.MONGO_DB_USERNAME).Collection("logs")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = logger.InsertOne(ctx, bson.M{
		"level":  logs.Level,
		"caller": logs.Caller,
		"msg":    logs.Msg,
		"method": logs.Method,
		"uuid":   logs.Uuid,
		"error":  logs.Err,
		"time":   logs.Time,
	})
	if err != nil {
		return
	}

	return len(p), nil
}
