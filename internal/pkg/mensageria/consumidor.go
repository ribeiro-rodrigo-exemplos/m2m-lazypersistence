package mensageria

import (
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

const fila = "lazypersistence"

// Consumidor de mensagens do Rabbitmq
type Consumidor struct {
	Host    string
	Porta   int
	Usuario string
	Senha   string
	conexao conexao
}

// Listener - chamado quando uma mensagem for recebida
type Listener func(mensagem Mensagem)

type conexao struct {
	conn  *amqp.Connection
	canal *amqp.Channel
}

// Conectar - conecta consumidor ao rabbitmq
func (c *Consumidor) Conectar(listener Listener) {
	url := "amqp://" + c.Usuario + ":" + c.Senha + "@" + c.Host + ":" + strconv.Itoa(c.Porta) + "/"
	conn, err := amqp.Dial(url)

	if err != nil {
		log.Fatal("Erro na conexão com o rabbitmq")
	}

	log.Print("Conexão com o rabbitmq aberta")

	c.conexao = conexao{conn: conn}
	c.abrirCanal(listener)
}

// Desconectar - desconecta consumidor do rabbitmq
func (c *Consumidor) Desconectar() {
	c.conexao.canal.Close()
	c.conexao.conn.Close()
}

func (c *Consumidor) abrirCanal(listener Listener) {

	if c.conexao.canal == nil {
		canal, err := c.conexao.conn.Channel()
		if err != nil {
			log.Fatal("Erro ao criar canal no rabbitmq")
		}
		c.conexao.canal = canal
	}

	mensagens, err := c.conexao.canal.Consume(
		fila,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Erro ao criar consumir no lazypersistence")
	}

	go func() {
		for mensagem := range mensagens {
			listener(Mensagem{delivery: mensagem})
		}
	}()
}
