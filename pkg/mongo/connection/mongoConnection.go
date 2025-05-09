package connection

import (
	"context"
	"encoding/base64"
	"net/url"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	name     string
	host     string
	username string
	Url      *url.URL
	Database *mongo.Database
	Logger   zerolog.Logger
}

type Config struct {
	Database string `json:"database"`
	URL      string `json:"url" base64:"true"`
	UserName string `json:"username" base64:"true"`
	Password string `json:"password" base64:"true"`
}

func (ds *Mongo) Connect(c Config) error {
	ds.Logger = log.With().Str("connection", c.Database).Logger()
	decodedURL, err := base64.StdEncoding.DecodeString(c.URL)
	if err != nil {
		ds.Logger.Error().Err(err).Send()
		return err
	}
	uri, err := url.Parse(string(decodedURL))
	if err != nil {
		ds.Logger.Error().Err(err).Send()
		return err
	}

	ds.Url = uri
	opts := options.Client()
	opts.ApplyURI(ds.Url.String())
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		ds.Logger.Error().Err(err).Send()
		return err
	}
	ds.Database = client.Database(c.Database)
	ds.name, ds.host, ds.username = ds.Database.Name(), ds.Url.Hostname(), ds.Url.User.Username()
	return nil
}
