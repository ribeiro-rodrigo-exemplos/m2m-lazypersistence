package config

// Config - configuracões da aplicação
type Config struct {
	MaxMessages      int    `default:"30"`
	RetentionSeconds int    `default:"30"`
	LogFile          string `required:"false"`
	RabbitMQ         struct {
		Host     string        `default:"localhost"`
		Port     int           `default:"5672"`
		User     string        `default:"guest"`
		Password string        `default:"guest"`
		Queues   []QueueConfig `required:"true"`
	}
	MongoDB struct {
		Host     string `default:"localhost"`
		Port     int    `default:"27017"`
		Database string `required:"true"`
	}
}

// QueueConfig - configuração das filas
type QueueConfig struct {
	Name            string
	Exchange        string
	ExchangeType    string
	RoutingKey      string
	DlqExchange     string
	DlqExchangeType string
	DlqRoutingKey   string
	Durable         bool
}
