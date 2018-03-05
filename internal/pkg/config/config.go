package config

// Config - configuracões da aplicação
type Config struct {
	MaxMessages int
	RabbitMQ    struct {
		Host     string
		Port     int
		User     string
		Password string
	}
}
