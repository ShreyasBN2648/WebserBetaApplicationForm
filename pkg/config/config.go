package config

import (
	"WebFormServer/pkg/constants"
	errorsCust "WebFormServer/pkg/errors"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Port  string
	Mongo Mongo
}

type Mongo struct {
	URL        string
	Database   string
	Collection string
	Timeout    int
}

func New(path string) (c Config, err error) {

	currentDir, _ := filepath.Abs(".")
	parentDir := filepath.Dir(currentDir)
	trueConfigPath := filepath.Join(parentDir, "config", "config.json")
	err = Read(trueConfigPath, &c)
	if err != nil {
		return
	}
	return
}

func Read(path string, conf interface{}) error {
	if reflect.ValueOf(conf).Kind() != reflect.Ptr {
		return errorsCust.ErrorInterfaceNotPointer
	}
	file, err := os.Open(path)
	if err != nil {
		log.Error().Stack().Caller().Err(err).Msgf("unable to open file: %s", path)
		return err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(conf)
	if err != nil {
		log.Error().Stack().Caller().Err(err).Msgf("unable to parse file: %s", path)
		return err
	}
	if err := translate(reflect.TypeOf(conf), reflect.ValueOf(conf)); err != nil {
		log.Error().Stack().Caller().Err(err).Msgf("unable to translate file: %s", path)
		return err
	}
	return nil
}

func translate(reflectType reflect.Type, reflectVal reflect.Value) error {
	switch reflectVal.Kind() {
	case reflect.Ptr:
		pointerVal := reflectVal.Elem()
		if !pointerVal.IsValid() {
			return nil
		}
		if err := translate(pointerVal.Type(), pointerVal); err != nil {
			return err
		}
	case reflect.Interface:
		interfaceVal := reflectVal.Elem()
		copyVal := reflect.New(interfaceVal.Type()).Elem()
		if err := translate(interfaceVal.Type(), interfaceVal); err != nil {
			return err
		}
		reflectVal.Set(copyVal)
	case reflect.Struct:
		for i := 0; i < reflectType.NumField(); i++ {
			if _, ok := reflectType.Field(i).Tag.Lookup("base64"); ok {
				structField := reflectVal.Field(i)
				original := structField.String()
				if original != constants.EmptyString {
					decoded, err := base64.StdEncoding.DecodeString(original)
					if err != nil {
						return err
					}
					structField.SetString(string(decoded))
				}
			} else {
				if err := translate(reflectVal.Field(i).Type(), reflectVal.Field(i)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
