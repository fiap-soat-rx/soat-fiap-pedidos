package routes

import (
	"net/http"
	"pedido-service/internal/adapters/primary/handlers"

	"github.com/gorilla/mux"
)

func ConfigurarRotas(r *mux.Router, pedidoHandler *handlers.PedidoHandler, healthHandler *handlers.HealthHandler) {
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/health", healthHandler.HealthCheck).Methods(http.MethodGet)

	api.HandleFunc("/checkout", pedidoHandler.CheckoutPedido).Methods(http.MethodPost)
	api.HandleFunc("/pedidos", pedidoHandler.CriarPedido).Methods(http.MethodPost)
	api.HandleFunc("/pedidos", pedidoHandler.ListarPedidos).Methods(http.MethodGet)
	api.HandleFunc("/pedidos/{id}", pedidoHandler.BuscarPedidoPorID).Methods(http.MethodGet)
	api.HandleFunc("/pedidos/{id}/status", pedidoHandler.AtualizarStatusPedido).Methods(http.MethodPatch)
}
