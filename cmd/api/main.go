package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pedido-service/configs"
	"pedido-service/internal/adapters/primary/handlers"
	"pedido-service/internal/adapters/secondary/repositories"
	"pedido-service/internal/core/services"
	"pedido-service/internal/routes"
	"pedido-service/pkg/clients"
	"pedido-service/pkg/database"

	"github.com/gorilla/mux"
)

const (
	AppVersion = "1.0.0"
)

func main() {
	cfg := configs.LoadConfig()

	mongoClient, err := database.ConectarMongoDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	produtoClient := clients.NovoProdutoClient(cfg.ProdutoServiceURL)
	clienteClient := clients.NovoClienteClient(cfg.ClienteServiceURL)
	pagamentoClient := clients.NovoPagamentoClient(cfg.PagamentoServiceURL)

	pedidoRepository := repositories.NovoPedidoRepository(mongoClient, cfg.MongoDBName)

	pedidoService := services.NovoPedidoService(pedidoRepository, produtoClient, clienteClient)

	pedidoHandler := handlers.NovoPedidoHandler(pedidoService, pagamentoClient)
	healthHandler := handlers.NovoHealthHandler(AppVersion)

	router := mux.NewRouter()
	routes.ConfigurarRotas(router, pedidoHandler, healthHandler)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Pedido Service iniciado na porta %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Desligando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor: %v", err)
	}

	log.Println("Servidor encerrado")
}
