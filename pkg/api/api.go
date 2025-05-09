package api

import (
	"WebFormServer/pkg/config"
	"WebFormServer/pkg/models"
	"WebFormServer/pkg/mongo"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var l *zerolog.Logger

func init() {
	logger := log.With().Str("pkg", "api").Logger()
	l = &logger
}

type Api struct {
	Serv *http.ServeMux
	m    mongo.Saver
}

func New(conf config.Config) (*Api, error) {
	mongoSaver, err := mongo.New(&conf.Mongo)
	if err != nil {
		return nil, err
	}
	return &Api{
		Serv: http.NewServeMux(),
		m:    mongoSaver,
	}, err
}

func (a *Api) FormHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		l.Error().Str("Userinfo Recorded", err.Error())
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	userInfo := models.UserInfo{
		Name:           r.FormValue("name"),
		Email:          r.FormValue("email"),
		Age:            r.FormValue("age-range"),
		Location:       r.FormValue("location"),
		Gpu:            r.FormValue("gpu"),
		Cpu:            r.FormValue("cpu"),
		SubmissionDate: time.Now().UTC().Format("2006-01-02"),
	}
	a.m.StoreBetaApplications(userInfo)
	currentDir, _ := filepath.Abs(".")
	parentDir := filepath.Dir(currentDir)
	filePath := filepath.Join(parentDir, "cmd", "staticHTML/submitted.html")
	invalidMethod := filepath.Join(parentDir, "cmd", "staticHTML/invalidMethod.html")

	content, err := os.ReadFile(filePath)
	if err != nil {
		l.Error().Str("Error reading file", err.Error())
		return
	}
	if r.Method == http.MethodGet {
		invMethod, err := os.ReadFile(invalidMethod)
		if err != nil {
			l.Error().Str("Error reading file", err.Error())
			return
		}
		fmt.Fprintf(w, "%s", invMethod)
		return
	}
	htmlString := string(content)
	fmt.Fprintf(w, "%s", htmlString)
}

func FormHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not supported", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "HomePage")
}
