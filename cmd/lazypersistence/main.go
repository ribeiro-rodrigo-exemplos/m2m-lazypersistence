package main

import "m2m-lazypersistence/internal/pkg/manager"

func main() {

	manager.Init()

	foreaver := make(chan bool)
	<-foreaver
}
