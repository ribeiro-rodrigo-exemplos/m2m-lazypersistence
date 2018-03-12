package main

import (
	"flag"
	"fmt"
	"log"
	"m2m-lazypersistence/internal/app"
	"m2m-lazypersistence/internal/pkg/config"
	"os"

	"github.com/jinzhu/configor"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const productionEnv = "PRODUCTION"

func configLog(config config.Config, environmentFlag string) {

	if environmentFlag == productionEnv || os.Getenv("ENVIRONMENT") == productionEnv {
		log.SetOutput(&lumberjack.Logger{
			Filename:   config.LogFile,
			MaxSize:    1,
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
	configLocation := flag.String("config-location", "./configs/config.yml", "a string")
	environment := flag.String("environment", "DEV", "a string")
	flag.Parse()

	return *configLocation, *environment
}

func main() {

	configLocation, environmentFlag := loadFlags()

	config := loadConfig(configLocation)
	configLog(config, environmentFlag)
	fmt.Println(config)
	app.Bootstrap(config)

	foreaver := make(chan bool)
	<-foreaver
}
