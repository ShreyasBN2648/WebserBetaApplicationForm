package mongo

import (
	"WebFormServer/pkg/config"
	"WebFormServer/pkg/constants"
	errorsCust "WebFormServer/pkg/errors"
	"WebFormServer/pkg/models"
	"WebFormServer/pkg/mongo/connection"
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoPak = "mongo"
)

type Saver interface {
	StoreBetaApplications(userInfo models.UserInfo) (docId string, err error)
}

var l *zerolog.Logger

type Mongo struct {
	db         *connection.Mongo
	URL        string
	Database   string
	Collection string
	Timeout    time.Duration
}

func init() {
	logger := log.With().Str("pkg", mongoPak).Logger()
	l = &logger
}

func New(conf *config.Mongo) (*Mongo, error) {
	db, err := initDatabase(*conf)
	if err != nil {
		return &Mongo{}, err
	}
	timeout := time.Duration(conf.Timeout) * time.Millisecond
	return &Mongo{
		db:         db,
		Database:   conf.Database,
		Collection: conf.Collection,
		Timeout:    timeout,
	}, nil
}

func initDatabase(conf config.Mongo) (db *connection.Mongo, err error) {
	c := connection.Config{
		Database: conf.Database,
		URL:      conf.URL,
	}
	db = &connection.Mongo{}
	err = db.Connect(c)
	if err != nil {
		db.Logger.Fatal().Err(err).Send()
	}
	return db, nil
}

func (m *Mongo) StoreBetaApplications(userInfo models.UserInfo) (docId string, err error) {
	start := time.Now()
	defer func(start time.Time) {
		var e *zerolog.Event
		if err != nil {
			e = l.Error().Err(err)
		} else {
			l.Debug()
		}
		e.Str(constants.Func, "StoreBetaApplications").
			Int64(constants.Duration, time.Since(start).Milliseconds()).Send()
	}(start)
	var result *mongo.InsertOneResult
	result, err = m.db.Database.Client().Database(m.Database).Collection(m.Collection).InsertOne(context.Background(), userInfo, &options.InsertOneOptions{})
	if err != nil {
		err = errorsCust.ErrInsertingBetaUserApplication
		return
	}
	docId = result.InsertedID.(primitive.ObjectID).Hex()
	return
}
