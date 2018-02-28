package main

import (
	"m2m-lazypersistence/internal/app"
	"m2m-lazypersistence/internal/pkg/config"

	"github.com/jinzhu/configor"
)

func main() {

	config := new(config.Config)

	configor.Load(config, "./configs/config.yml")

	app.Bootstrap(*config)

	foreaver := make(chan bool)
	<-foreaver
}
