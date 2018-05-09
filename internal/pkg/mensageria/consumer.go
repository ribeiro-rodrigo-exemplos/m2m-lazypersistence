package mensageria

import (
	"bytes"
	"encoding/json"
	"log"
	"m2m-lazypersistence/internal/pkg/config"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

// Consumer - Consumidor de mensagens do Rabbitmq
type Consumer struct {
	Host          string
	Port          int
	User          string
	Password      string
	PrefetchCount int
	Queues        []config.QueueConfig
	connection    connection
}

// Listener - chamado quando uma mensagem for recebida
type Listener func(request RequestPersistence)

type connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Connect - conecta consumidor ao rabbitmq
func (c *Consumer) Connect(listener Listener) {
	url := "amqp://" + c.User + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/"
	conn, err := amqp.Dial(url)

	if err != nil {
		log.Println("Erro na conexão com o rabbitmq -", err)
		go c.reconnect(listener)
		return
	}

	log.Println("Conexão estabelecida com o rabbitmq")

	c.connection = connection{conn: conn}

	c.startCloseListener(listener)

	c.openChannels(listener)
}

// Disconnect - desconecta consumidor do rabbitmq
func (c *Consumer) Disconnect() {
	if c.connection.channel != nil {
		c.connection.conn.Close()
		c.connection.channel = nil
	}
}

func (c *Consumer) openChannels(listener Listener) {

	channel := c.openChannel()

	for _, queue := range c.Queues {

		messages := c.createConsumer(channel, queue)

		log.Println("Ouvindo mensagens da fila", queue.Name)

		go func() {
			for m := range messages {
				message := Message{delivery: m}
				buffer := bytes.NewReader(m.Body)
				decoder := json.NewDecoder(buffer)
				decoder.UseNumber()
				err := decoder.Decode(&message.Payload)

				if err != nil {
					log.Printf("Erro ao deserializar mensagem %s\n", err)
					message.Discard()
					continue
				}

				request := RequestPersistence{Message: message, Headers: m.Headers}
				listener(request)
			}
		}()

		c.connection.channel = channel
	}
}

func (c *Consumer) openChannel() *amqp.Channel {

	channel, err := c.connection.conn.Channel()

	if err != nil {
		log.Fatal("Erro ao criar canal no rabbitmq -", err)
	}

	err = channel.Qos(
		c.PrefetchCount,
		0,
		false,
	)

	if err != nil {
		log.Fatal("Erro ao configurar prefetch do canal -", err)
	}

	return channel
}

func (c *Consumer) createConsumer(channel *amqp.Channel, queueCfg config.QueueConfig) (messages <-chan amqp.Delivery) {

	var args amqp.Table

	if queueCfg.DlqExchange != "" {
		args = amqp.Table{
			"x-dead-letter-exchange":    queueCfg.DlqExchange,
			"x-dead-letter-routing-key": queueCfg.DlqRoutingKey,
		}
		c.createExchange(
			channel,
			queueCfg.DlqExchange,
			queueCfg.DlqExchangeType,
			nil,
		)
	} else {
		args = nil
	}

	queue, err := channel.QueueDeclare(
		queueCfg.Name,
		queueCfg.Durable,
		false,
		false,
		false,
		args,
	)

	if err != nil {
		log.Fatalf("Erro ao criar fila %s - %s", queueCfg.Name, err)
	}

	if queueCfg.Exchange != "" {

		c.createExchange(
			channel,
			queueCfg.Exchange,
			queueCfg.ExchangeType,
			args,
		)

		err = channel.QueueBind(
			queue.Name,
			queueCfg.RoutingKey,
			queueCfg.Exchange,
			false,
			nil,
		)

		if err != nil {
			log.Fatalf("Erro ao realizar bind entre o exchange %s e a fila %s - %s",
				queueCfg.Exchange, queue.Name, err)
		}
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
		log.Fatalf("Erro ao criar consumidor para a fila %s - %s", queue.Name, err)
	}

	return
}

func (c *Consumer) startCloseListener(listener Listener) {

	channelClose := make(chan *amqp.Error)
	c.connection.conn.NotifyClose(channelClose)

	go func() {
		m := <-channelClose
		log.Println(m)
		c.reconnect(listener)
	}()
}

func (c *Consumer) reconnect(listener Listener) {
	c.Disconnect()
	time.Sleep(time.Second * 5)
	c.Connect(listener)
}

func (c *Consumer) createExchange(channel *amqp.Channel, name, kind string, args amqp.Table) {
	err := channel.ExchangeDeclare(
		name,
		kind,
		true,
		false,
		false,
		false,
		args,
	)

	if err != nil {
		log.Fatalf("Erro ao criar exchange %s no rabbitmq - %s", name, err)
	}
}
