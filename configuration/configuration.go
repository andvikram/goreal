package configuration

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

const (
	// DevEnv ...
	DevEnv = "development"
	// TestEnv ...
	TestEnv = "test"
	// ProdEnv ...
	ProdEnv = "production"
)

type conf struct {
	Environment  string
	Origin       string
	Datasink     string
	DatasinkURL  string
	AppScheme    string
	AppHost      string
	AppPort      string
	LogLevel     string
	LogFilePath  string
	CertFilePath string
	KeyFilePath  string
}

var (
	// Config ...
	Config  *conf
	once    sync.Once
	pathSep = string(os.PathSeparator)
)

// Initialize ...
func Initialize() {
	once.Do(func() {
		Config = new(conf)
		initConfig(viper.GetString("CONFIGDIRPATH"))
		initVars()
	})
}

func initVars() {
	Config.Environment = viper.GetString("ENVIRONMENT")
	Config.Origin = viper.GetString("ORIGIN")
	Config.Datasink = viper.GetString("DATASINK")
	Config.DatasinkURL = viper.GetString("DATASINK_URL")
	Config.AppScheme = viper.GetString("APPLICATION.SCHEME")
	Config.AppHost = viper.GetString("APPLICATION.HOST")
	Config.AppPort = viper.GetString("APPLICATION.PORT")
	Config.LogLevel = viper.GetString("APPLICATION.LOG_LEVEL")
	Config.LogFilePath = viper.GetString("APPLICATION.LOG_FILE_PATH")

	Config.CertFilePath = viper.GetString("SECURE.CERT_PATH")
	Config.KeyFilePath = viper.GetString("SECURE.KEY_PATH")

	switch {
	case Config.Environment == "":
		panic("Environment variable ENVIRONMENT not set")
	case Config.Origin == "":
		panic("Environment variable ORIGIN not set")
	case Config.Datasink == "":
		panic("Environment variable DATASINK not set")
	case Config.DatasinkURL == "":
		panic("Environment variable DATASINK_URL not set")
	case Config.AppScheme == "":
		panic("Environment variable APPLICATION.SCHEME not set")
	case Config.AppHost == "":
		panic("Environment variable APPLICATION.HOST not set")
	case Config.LogLevel == "":
		panic("Environment variable APPLICATION.LOG_LEVEL not set")
	case Config.LogFilePath == "":
		panic("Environment variable APPLICATION.LOG_FILE_PATH not set")
	case Config.AppScheme == "https" && Config.CertFilePath == "":
		panic("AppScheme is https, but environment variable SECURE.CERT_PATH is not set")
	case Config.AppScheme == "https" && Config.KeyFilePath == "":
		panic("AppScheme is https, but environment variable SECURE.KEY_PATH is not set")
	}
}

func initConfig(configDirPath string) {
	viper.AutomaticEnv() //env var will overwrite keys in yml file
	viper.SetConfigType("yaml")
	viper.SetConfigName("goreal-config")

	viper.AddConfigPath(configDirPath)

	// if config file does not exist, don't attempt to read it
	if _, err := os.Stat(configDirPath + "/goreal-config.yml"); os.IsNotExist(err) {
		return
	}

	err := viper.ReadInConfig()
	if viper.ReadInConfig() != nil {
		panic(fmt.Errorf("fatal error in config file: %s", err))
	}
}
