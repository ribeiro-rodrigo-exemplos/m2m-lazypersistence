package dispatcher

import (
	"fmt"
	"m2m-lazypersistence/internal/pkg/mensageria"
)

// Dispatch - Salva as mensagens no mongodb
func Dispatch(repository map[string][]mensageria.Message) {
	fmt.Println("signal - gravando dados no mongo")
}

func copy(repository map[string][]mensageria.Message) map[string][]mensageria.Message {
	newRepository := make(map[string][]mensageria.Message)

	for chave, valor := range repository {
		newRepository[chave] = valor
	}

	return newRepository
}
