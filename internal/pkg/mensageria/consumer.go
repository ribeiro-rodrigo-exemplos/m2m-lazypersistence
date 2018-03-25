package mensageria

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

// Consumer - Consumidor de mensagens do Rabbitmq
type Consumer struct {
	Host       string
	Port       int
	User       string
	Password   string
	Queue      string
	Queues     []string
	connection connection
}

// Listener - chamado quando uma mensagem for recebida
type Listener func(request RequestPersistence)

type connection struct {
	conn     *amqp.Connection
	channels []*amqp.Channel
}

// Connect - conecta consumidor ao rabbitmq
func (c *Consumer) Connect(listener Listener) {
	url := "amqp://" + c.User + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/"
	conn, err := amqp.Dial(url)

	if err != nil {
		log.Fatal("Erro na conexão com o rabbitmq")
	}

	log.Println("Conexão estabelecida com o rabbitmq")

	c.connection = connection{conn: conn}
	c.openChannels(listener)
}

// Disconnect - desconecta consumidor do rabbitmq
func (c *Consumer) Disconnect() {
	for _, channel := range c.connection.channels {
		channel.Close()
	}
	c.connection.conn.Close()
	c.connection.channels = make([]*amqp.Channel, 0)
}

func (c *Consumer) openChannels(listener Listener) {

	for _, queueName := range c.Queues {

		channel := c.openChannel()
		messages := c.createConsumer(channel, queueName)

		log.Println("Ouvindo mensagens da fila", queueName)

		go func() {
			for m := range messages {
				message := Message{delivery: m}
				json.Unmarshal(m.Body, &message.Payload)
				request := RequestPersistence{Message: message, Headers: m.Headers}
				listener(request)
			}
		}()

		c.connection.channels = append(c.connection.channels, channel)
	}
}

func (c *Consumer) openChannel() *amqp.Channel {

	channel, err := c.connection.conn.Channel()
	if err != nil {
		log.Fatal("Erro ao criar canal no rabbitmq")
	}

	return channel
}

func (c *Consumer) createConsumer(channel *amqp.Channel, queueName string) (messages <-chan amqp.Delivery) {
	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Erro ao criar fila", queueName)
	}

	messages, err = channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Erro ao criar consumidor para a fila", queue.Name)
	}

	return
}
