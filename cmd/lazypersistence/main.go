package main

import (
	"flag"
	"log"
	"m2m-lazypersistence/internal/app"
	"m2m-lazypersistence/internal/pkg/config"
	"os"

	"github.com/jinzhu/configor"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const productionEnv = "PRODUCTION"

func configLog(config config.Config, environmentFlag string) {

	if environmentFlag == productionEnv || os.Getenv("M2M-ENVIRONMENT") == productionEnv {
		log.SetOutput(&lumberjack.Logger{
			Filename:   config.LogFile,
			MaxSize:    100,
			MaxBackups: 14,
			MaxAge:     28,
		})

	} else {
		log.SetOutput(os.Stdout)
	}
}

func loadConfig(configLocation string) config.Config {
	config := new(config.Config)
	configor.Load(config, configLocation)
	return *config
}

func loadFlags() (string, string) {
	configLocation := flag.String("config-location", "./configs/config.json", "a string")
	environment := flag.String("m2m-environment", "DEVELOPMENT", "a string")
	flag.Parse()

	return *configLocation, *environment
}

func main() {

	configLocation, environmentFlag := loadFlags()

	config := loadConfig(configLocation)
	configLog(config, environmentFlag)

	app.Bootstrap(config)

	foreaver := make(chan bool)
	<-foreaver
}
