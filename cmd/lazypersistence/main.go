package main

import (
	"fmt"
	"m2m-lazypersistence/internal/pkg/mensageria"
)

func main() {
	consumidor := mensageria.Consumidor{
		Host:    "localhost",
		Porta:   5672,
		Usuario: "guest",
		Senha:   "guest",
	}

	defer consumidor.Desconectar()

	consumidor.Conectar(func(mensagem mensageria.Mensagem) {
		fmt.Println(mensagem.Get())
		mensagem.Confirmar()
	})

	foreaver := make(chan bool)
	<-foreaver
}
