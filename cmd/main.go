package main

import (
	"WebFormServer/pkg/api"
	"WebFormServer/pkg/config"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var l *zerolog.Logger

const (
	configPath = "config/config.json"
)

func init() {
	logger := log.With().
		Str("pkg", "api").
		Logger()
	l = &logger
}

func main() {
	c, err := config.New(getConfigPath())
	if err != nil {
		l.Fatal().
			Str("error", err.Error()).
			Msg("unable to start server")
	}
	a, err := api.New(c)
	if err != nil {
		l.Fatal().
			Str("error", err.Error()).
			Msg("unable to start server")
	}
	formServer := http.FileServer(http.Dir("./staticHTML"))
	http.Handle("/", formServer)
	http.HandleFunc("/form", a.FormHandler)
	// a.Serv.Handle("/", formServer)
	// a.Serv.HandleFunc("/form", a.FormHandler)
	l.Print("Starting server at port 8000")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", c.Port), nil); err != nil {
		l.Fatal().
			Str("error", err.Error()).
			Msg("unable to start server")
	}
}

func getConfigPath() string {
	if ConfigPath, ok := os.LookupEnv("CONFIG_PATH"); ok {
		return ConfigPath
	}
	return configPath
}
