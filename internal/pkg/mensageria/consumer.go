package mensageria

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

const fila = "lazypersistence"

// Consumer - Consumidor de mensagens do Rabbitmq
type Consumer struct {
	Host       string
	Port       int
	User       string
	Password   string
	connection connection
}

// Listener - chamado quando uma mensagem for recebida
type Listener func(message Message)

type connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Connect - conecta consumidor ao rabbitmq
func (c *Consumer) Connect(listener Listener) {
	url := "amqp://" + c.User + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/"
	conn, err := amqp.Dial(url)

	if err != nil {
		log.Fatal("Erro na conexão com o rabbitmq")
	}

	log.Println("Conexão com o rabbitmq aberta")

	c.connection = connection{conn: conn}
	c.openChannel(listener)
}

// Disconnect - desconecta consumidor do rabbitmq
func (c *Consumer) Disconnect() {
	c.connection.channel.Close()
	c.connection.conn.Close()
}

func (c *Consumer) openChannel(listener Listener) {

	if c.connection.channel == nil {
		channel, err := c.connection.conn.Channel()
		if err != nil {
			log.Fatal("Erro ao criar canal no rabbitmq")
		}
		c.connection.channel = channel
	}

	messages := c.createConsumer()

	go func() {
		for m := range messages {
			message := Message{delivery: m}
			json.Unmarshal(m.Body, &message)
			listener(message)
		}
	}()
}

func (c *Consumer) createConsumer() (messages <-chan amqp.Delivery) {
	queue, err := c.connection.channel.QueueDeclare(
		fila,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Erro ao criar fila", fila)
	}

	messages, err = c.connection.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Erro ao criar consumidor para a fila", fila)
	}

	return
}
