package config

// Config - configuracões da aplicação
type Config struct {
	RabbitMQ struct {
		Host     string
		Port     int
		User     string
		Password string
	}
}
