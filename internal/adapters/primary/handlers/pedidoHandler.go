package handlers

import (
	"encoding/json"
	"net/http"
	"pedido-service/internal/adapters/primary/dtos"
	"pedido-service/internal/adapters/primary/presenters"
	"pedido-service/internal/core/domain"
	"pedido-service/internal/core/ports"
	"pedido-service/pkg/clients"
	"sort"

	"github.com/gorilla/mux"
)

type PedidoHandler struct {
	pedidoService   ports.PedidoService
	pagamentoClient *clients.PagamentoClient
	presenter       *presenters.PedidoPresenter
}

func NovoPedidoHandler(pedidoService ports.PedidoService, pagamentoClient *clients.PagamentoClient) *PedidoHandler {
	return &PedidoHandler{
		pedidoService:   pedidoService,
		pagamentoClient: pagamentoClient,
		presenter:       presenters.NewPedidoPresenter(),
	}
}

type AtualizarStatusRequest struct {
	Status domain.StatusPedido `json:"status"`
}

func (h *PedidoHandler) CriarPedido(w http.ResponseWriter, r *http.Request) {
	var req dtos.PedidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	pedido, err := req.ToDomain()
	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	createdPedido, err := h.pedidoService.CriarPedido(r.Context(), pedido.ClienteID, pedido.Itens)
	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := h.presenter.Present(createdPedido)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *PedidoHandler) ListarPedidos(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	clienteID := r.URL.Query().Get("cliente_id")

	var pedidos []*domain.Pedido
	var err error

	switch {
	case status != "":
		pedidos, err = h.pedidoService.ListarPedidosPorStatus(r.Context(), domain.StatusPedido(status))
	case clienteID != "":
		pedidos, err = h.pedidoService.ListarPedidosPorCliente(r.Context(), clienteID)
	default:
		pedidos, err = h.pedidoService.ListarPedidos(r.Context())
	}

	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	var filteredPedidos []*domain.Pedido
	for _, p := range pedidos {
		if p.Status != domain.StatusFinalizado {
			filteredPedidos = append(filteredPedidos, p)
		}
	}

	sort.Slice(filteredPedidos, func(i, j int) bool {
		if filteredPedidos[i].Status == filteredPedidos[j].Status {
			return filteredPedidos[i].CreatedAt.Before(filteredPedidos[j].CreatedAt)
		}
		return getStatusPrioridade(filteredPedidos[i].Status) > getStatusPrioridade(filteredPedidos[j].Status)
	})

	response := h.presenter.PresentList(filteredPedidos)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PedidoHandler) BuscarPedidoPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	pedido, err := h.pedidoService.BuscarPedidoPorID(r.Context(), id)
	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	if pedido == nil {
		http.Error(w, "Pedido não encontrado", http.StatusNotFound)
		return
	}

	response := h.presenter.Present(pedido)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PedidoHandler) AtualizarStatusPedido(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req AtualizarStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := h.pedidoService.AtualizarStatusPedido(r.Context(), id, req.Status)
	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	status, response := h.presenter.PresentSuccess("Status do pedido atualizado com sucesso", nil)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func (h *PedidoHandler) CheckoutPedido(w http.ResponseWriter, r *http.Request) {
	var req dtos.PedidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Converter itens do DTO para domain sem preço (o serviço preencherá)
	itens := make([]domain.ItemPedido, len(req.Itens))
	for i, item := range req.Itens {
		itens[i] = domain.ItemPedido{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
			Observacao: item.Observacao,
			// Preço será preenchido pelo serviço ao validar o produto
		}
	}

	createdPedido, err := h.pedidoService.CriarPedido(r.Context(), req.ClienteID, itens)
	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	pagamento, err := h.pagamentoClient.CriarPagamento(r.Context(), createdPedido.ID, createdPedido.ValorTotal, "PIX")
	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"pedido_id":    createdPedido.ID,
		"status":       string(createdPedido.Status),
		"valor":        createdPedido.ValorTotal,
		"pagamento_id": pagamento.ID,
		"qr_code":      pagamento.QRCodeData,
		"external_id":  pagamento.ExternalID,
		"expires_at":   pagamento.ExpiresAt,
		"message":      "Pedido criado com sucesso. Realize o pagamento via PIX para prosseguir.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func getStatusPrioridade(status domain.StatusPedido) int {
	switch status {
	case domain.StatusPronto:
		return 3
	case domain.StatusEmPreparacao:
		return 2
	case domain.StatusRecebido:
		return 1
	default:
		return 0
	}
}
