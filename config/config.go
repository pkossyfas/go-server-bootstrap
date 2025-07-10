/*
Package config is responsible for loading the app configuration.
*/
package config

import (
	"os"
	"reflect"

	"github.com/pkossyfas/go-server-bootstrap/logger"
)

// GetConfig is the pointer to the struct holding the configuration.
var GetConfig *AppConfig

// AppConfig is the structure which holds all the static configuration of the application.
// Tags are being used to map the configuration variables with environmental variables.
type AppConfig struct {
	ServerPort         string `env:"APP_SERVER_PORT"           default:"8024"`
	ServerReadTimeout  string `env:"APP_SERVER_READ_TIMEOUT"   default:"5s"`
	ServerWriteTimeout string `env:"APP_SERVER_WRITE_TIMEOUT"  default:"5s"`
	ClientTimeout      string `env:"APP_CLIENT_TIMEOUT"        default:"5s"`
	ShutdownTimeout    string `env:"APP_SHUTDOWN_TIMEOUT"      default:"15s"`
	DBUser             string `env:"APP_DB_USER"               default:"postgres"`
	DBPassword         string `env:"APP_DB_PASSWORD"           default:"postgres"`
	DBHost             string `env:"APP_DB_HOST"               default:"localhost"`
	DBPort             string `env:"APP_DB_PORT"               default:"5432"`
	DBName             string `env:"APP_DB_NAME"               default:"postgres"`
}

// LoadAppConfig loads configuration from env variables if available.
// Use of reflect library and struct tags to get all the configuration variables which
// are mapped with an environmental variable. If env variable is not available the
// default value will be used.
func LoadAppConfig() {

	var c AppConfig
	st := reflect.ValueOf(&c).Elem()
	typeOfSt := st.Type()
	for i := 0; i < st.NumField(); i++ {
		field := typeOfSt.Field(i)
		if envVar, ok := field.Tag.Lookup("env"); ok {
			if envVarValue, ok := os.LookupEnv(envVar); ok {
				st.Field(i).SetString(envVarValue)
			} else {
				if envVarValue, ok := field.Tag.Lookup("default"); ok {
					st.Field(i).SetString(envVarValue)
				}
			}
		}
	}
	GetConfig = &c

	logger.Info("app config loaded")
}
