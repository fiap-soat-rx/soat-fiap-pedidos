package configs

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerPort          string
	MongoURI            string
	MongoDBName         string
	ProdutoServiceURL   string
	ClienteServiceURL   string
	PagamentoServiceURL string
	LogLevel            string
}

func LoadConfig() *Config {
	serverPort := getEnvOrPanic("SERVER_PORT", "8083")
	mongoURI := getEnvOrPanic("MONGO_URI", "mongodb://localhost:27017")
	mongoDBName := getEnvOrPanic("MONGO_DB_NAME", "pedido_db")
	produtoServiceURL := getEnvOrPanic("PRODUTO_SERVICE_URL", "http://localhost:8082")
	clienteServiceURL := getEnvOrPanic("CLIENTE_SERVICE_URL", "http://localhost:8081")
	pagamentoServiceURL := getEnvOrPanic("PAGAMENTO_SERVICE_URL", "http://localhost:8084")
	logLevel := getEnv("LOG_LEVEL", "info")

	return &Config{
		ServerPort:          serverPort,
		MongoURI:            mongoURI,
		MongoDBName:         mongoDBName,
		ProdutoServiceURL:   produtoServiceURL,
		ClienteServiceURL:   clienteServiceURL,
		PagamentoServiceURL: pagamentoServiceURL,
		LogLevel:            logLevel,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvOrPanic(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if defaultValue != "" {
		return defaultValue
	}
	panic(fmt.Sprintf("Variável de ambiente obrigatória não definida: %s", key))
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, strconv.FormatBool(defaultValue))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
