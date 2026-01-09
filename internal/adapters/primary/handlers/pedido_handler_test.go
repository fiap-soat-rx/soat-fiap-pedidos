package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pedido-service/internal/adapters/primary/dtos"
	"pedido-service/internal/core/domain"
	"pedido-service/pkg/clients"
	"github.com/gorilla/mux"
)

// Mock service implementing ports.PedidoService
type MockPedidoService struct {
	CriarFunc               func(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error)
	BuscarPorIDFunc         func(ctx context.Context, id string) (*domain.Pedido, error)
	ListarFunc              func(ctx context.Context) ([]*domain.Pedido, error)
	ListarPorStatusFunc     func(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error)
	ListarPorClienteFunc    func(ctx context.Context, clienteID string) ([]*domain.Pedido, error)
	AtualizarStatusFunc     func(ctx context.Context, id string, status domain.StatusPedido) error
}

func (m *MockPedidoService) CriarPedido(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {
	if m.CriarFunc != nil {
		return m.CriarFunc(ctx, clienteID, itens)
	}
	return nil, nil
}
func (m *MockPedidoService) BuscarPedidoPorID(ctx context.Context, id string) (*domain.Pedido, error) {
	if m.BuscarPorIDFunc != nil {
		return m.BuscarPorIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *MockPedidoService) ListarPedidos(ctx context.Context) ([]*domain.Pedido, error) {
	if m.ListarFunc != nil {
		return m.ListarFunc(ctx)
	}
	return nil, nil
}
func (m *MockPedidoService) ListarPedidosPorStatus(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error) {
	if m.ListarPorStatusFunc != nil {
		return m.ListarPorStatusFunc(ctx, status)
	}
	return nil, nil
}
func (m *MockPedidoService) ListarPedidosPorCliente(ctx context.Context, clienteID string) ([]*domain.Pedido, error) {
	if m.ListarPorClienteFunc != nil {
		return m.ListarPorClienteFunc(ctx, clienteID)
	}
	return nil, nil
}
func (m *MockPedidoService) AtualizarStatusPedido(ctx context.Context, id string, status domain.StatusPedido) error {
	if m.AtualizarStatusFunc != nil {
		return m.AtualizarStatusFunc(ctx, id, status)
	}
	return nil
}

// Mock Pagamento client
type MockPagamentoClient struct {
	CriarFunc func(ctx context.Context, pedidoID string, valor float64, tipo string) (*clients.Pagamento, error)
}

func (m *MockPagamentoClient) CriarPagamento(ctx context.Context, pedidoID string, valor float64, tipo string) (*clients.Pagamento, error) {
	if m.CriarFunc != nil {
		return m.CriarFunc(ctx, pedidoID, valor, tipo)
	}
	return nil, nil
}

func TestCriarPedido_InvalidJSON(t *testing.T) {
	svc := &MockPedidoService{}
	h := NovoPedidoHandler(svc, nil)

	req := httptest.NewRequest(http.MethodPost, "/pedidos", bytes.NewBufferString("{invalid"))
	w := httptest.NewRecorder()

	h.CriarPedido(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid json, got %d", w.Result().StatusCode)
	}
}

func TestCriarPedido_ServiceError(t *testing.T) {
	svc := &MockPedidoService{CriarFunc: func(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {
		return nil, errors.New("service error")
	}}
	h := NovoPedidoHandler(svc, nil)

	payload := dtos.PedidoRequest{Itens: []dtos.ItemPedidoRequest{{ProdutoID: "p1", Quantidade: 1}}}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/pedidos", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.CriarPedido(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 when service returns error, got %d", w.Result().StatusCode)
	}
}

func TestCriarPedido_MissingPrice_ReturnsBadRequest(t *testing.T) {
	svc := &MockPedidoService{CriarFunc: func(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {
		return &domain.Pedido{ID: "id1"}, nil
	}}
	h := NovoPedidoHandler(svc, nil)

	payload := dtos.PedidoRequest{Itens: []dtos.ItemPedidoRequest{{ProdutoID: "p1", Quantidade: 1}}}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/pedidos", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.CriarPedido(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 when price is missing, got %d", w.Result().StatusCode)
	}
}

func TestListarPedidos_ServiceError(t *testing.T) {
	svc := &MockPedidoService{ListarFunc: func(ctx context.Context) ([]*domain.Pedido, error) {
		return nil, errors.New("list error")
	}}
	h := NovoPedidoHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/pedidos", nil)
	w := httptest.NewRecorder()

	h.ListarPedidos(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 when list service errors, got %d", w.Result().StatusCode)
	}
}

func TestBuscarPedidoPorID_NotFound(t *testing.T) {
	svc := &MockPedidoService{BuscarPorIDFunc: func(ctx context.Context, id string) (*domain.Pedido, error) {
		return nil, nil
	}}
	h := NovoPedidoHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/pedidos/id1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "id1"})
	w := httptest.NewRecorder()

	h.BuscarPedidoPorID(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 when pedido not found, got %d", w.Result().StatusCode)
	}
}

func TestAtualizarStatusPedido_InvalidJSON(t *testing.T) {
	svc := &MockPedidoService{}
	h := NovoPedidoHandler(svc, nil)

	req := httptest.NewRequest(http.MethodPut, "/pedidos/id1/status", bytes.NewBufferString("{invalid"))
	req = mux.SetURLVars(req, map[string]string{"id": "id1"})
	w := httptest.NewRecorder()

	h.AtualizarStatusPedido(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid json, got %d", w.Result().StatusCode)
	}
}

func TestAtualizarStatusPedido_ServiceError(t *testing.T) {
	svc := &MockPedidoService{AtualizarStatusFunc: func(ctx context.Context, id string, status domain.StatusPedido) error {
		return errors.New("update error")
	}}
	h := NovoPedidoHandler(svc, nil)

	payload := AtualizarStatusRequest{Status: domain.StatusEmPreparacao}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/pedidos/id1/status", bytes.NewBuffer(b))
	req = mux.SetURLVars(req, map[string]string{"id": "id1"})
	w := httptest.NewRecorder()

	h.AtualizarStatusPedido(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 when update service errors, got %d", w.Result().StatusCode)
	}
}

func TestAtualizarStatusPedido_Success(t *testing.T) {
	svc := &MockPedidoService{AtualizarStatusFunc: func(ctx context.Context, id string, status domain.StatusPedido) error {
		return nil
	}}
	h := NovoPedidoHandler(svc, nil)

	payload := AtualizarStatusRequest{Status: domain.StatusEmPreparacao}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/pedidos/id1/status", bytes.NewBuffer(b))
	req = mux.SetURLVars(req, map[string]string{"id": "id1"})
	w := httptest.NewRecorder()

	h.AtualizarStatusPedido(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on update success, got %d", w.Result().StatusCode)
	}
}

func TestCheckoutPedido_PaymentError(t *testing.T) {
	pedido := &domain.Pedido{ID: "id1", ValorTotal: 10}
	svc := &MockPedidoService{CriarFunc: func(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {
		return pedido, nil
	}}

	pagSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	}))
	defer pagSrv.Close()

	pagCli := clients.NovoPagamentoClient(pagSrv.URL)
	h := NovoPedidoHandler(svc, pagCli)

	payload := dtos.PedidoRequest{Itens: []dtos.ItemPedidoRequest{{ProdutoID: "p1", Quantidade: 1}}}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/pedidos/checkout", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.CheckoutPedido(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 when payment fails, got %d", w.Result().StatusCode)
	}
}

func TestCheckoutPedido_Success(t *testing.T) {
	pedido := &domain.Pedido{ID: "id1", ValorTotal: 10, Status: domain.StatusRecebido}
	svc := &MockPedidoService{CriarFunc: func(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {
		return pedido, nil
	}}

	pagSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		ex := "2026-01-01T00:00:00Z"
		json.NewEncoder(w).Encode(clients.Pagamento{ID: "pay1", PedidoID: "id1", Valor: 10, QRCodeData: "qrcode", ExternalID: "ext", ExpiresAt: &ex})
	}))
	defer pagSrv.Close()

	pagCli := clients.NovoPagamentoClient(pagSrv.URL)
	h := NovoPedidoHandler(svc, pagCli)

	payload := dtos.PedidoRequest{Itens: []dtos.ItemPedidoRequest{{ProdutoID: "p1", Quantidade: 1}}}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/pedidos/checkout", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.CheckoutPedido(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 on checkout success, got %d", w.Result().StatusCode)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["pagamento_id"] != "pay1" {
		t.Fatalf("expected pagamento_id pay1, got %v", resp["pagamento_id"])
	}
}

func TestBuscarPedidoPorID_Success(t *testing.T) {
	svc := &MockPedidoService{BuscarPorIDFunc: func(ctx context.Context, id string) (*domain.Pedido, error) {
		return &domain.Pedido{ID: id, Status: domain.StatusRecebido}, nil
	}}
	h := NovoPedidoHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/pedidos/id1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "id1"})
	w := httptest.NewRecorder()

	h.BuscarPedidoPorID(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when found, got %d", w.Result().StatusCode)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["id"] != "id1" {
		t.Fatalf("expected id id1, got %v", resp["id"])
	}
}


func TestCheckoutPedido_InvalidJSON(t *testing.T) {
	svc := &MockPedidoService{}
	h := NovoPedidoHandler(svc, nil)

	req := httptest.NewRequest(http.MethodPost, "/pedidos/checkout", bytes.NewBufferString("{invalid"))
	w := httptest.NewRecorder()

	h.CheckoutPedido(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid json in checkout, got %d", w.Result().StatusCode)
	}
}

func TestListarPedidos_FilterAndSort(t *testing.T) {
	p1 := &domain.Pedido{ID: "1", Status: domain.StatusEmPreparacao, CreatedAt: time.Now().Add(-time.Hour)}
	p2 := &domain.Pedido{ID: "2", Status: domain.StatusRecebido, CreatedAt: time.Now()}
	p3 := &domain.Pedido{ID: "3", Status: domain.StatusFinalizado, CreatedAt: time.Now().Add(-2 * time.Hour)}
	svc := &MockPedidoService{ListarFunc: func(ctx context.Context) ([]*domain.Pedido, error) {
		return []*domain.Pedido{p2, p3, p1}, nil
	}}
	h := NovoPedidoHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/pedidos", nil)
	w := httptest.NewRecorder()

	h.ListarPedidos(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on list success, got %d", w.Result().StatusCode)
	}

	var resp []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp) != 2 { // p3 should be filtered out
		t.Fatalf("expected 2 items after filtering finalizado, got %d", len(resp))
	}
	if resp[0]["id"] != "1" {
		t.Fatalf("expected p1 first due to priority, got %v", resp[0]["id"])
	}
}

func TestListarPedidos_ByStatusParam(t *testing.T) {
	p := &domain.Pedido{ID: "s1", Status: domain.StatusEmPreparacao}
	svc := &MockPedidoService{ListarPorStatusFunc: func(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error) {
		return []*domain.Pedido{p}, nil
	}}
	h := NovoPedidoHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/pedidos?status=EM_PREPARACAO", nil)
	w := httptest.NewRecorder()

	h.ListarPedidos(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on list by status, got %d", w.Result().StatusCode)
	}
}

func TestListarPedidos_ByClienteParam(t *testing.T) {
	p := &domain.Pedido{ID: "c1", Status: domain.StatusRecebido}
	svc := &MockPedidoService{ListarPorClienteFunc: func(ctx context.Context, clienteID string) ([]*domain.Pedido, error) {
		return []*domain.Pedido{p}, nil
	}}
	h := NovoPedidoHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/pedidos?cliente_id=abc", nil)
	w := httptest.NewRecorder()

	h.ListarPedidos(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on list by cliente, got %d", w.Result().StatusCode)
	}
}

func TestBuscarPedidoPorID_ServiceError(t *testing.T) {
	svc := &MockPedidoService{BuscarPorIDFunc: func(ctx context.Context, id string) (*domain.Pedido, error) {
		return nil, errors.New("db error")
	}}
	h := NovoPedidoHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/pedidos/id1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "id1"})
	w := httptest.NewRecorder()

	h.BuscarPedidoPorID(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 when service errors, got %d", w.Result().StatusCode)
	}
}

func TestAtualizarStatusPedido_ResponseMessage(t *testing.T) {
	svc := &MockPedidoService{AtualizarStatusFunc: func(ctx context.Context, id string, status domain.StatusPedido) error {
		return nil
	}}
	h := NovoPedidoHandler(svc, nil)

	payload := AtualizarStatusRequest{Status: domain.StatusEmPreparacao}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/pedidos/id1/status", bytes.NewBuffer(b))
	req = mux.SetURLVars(req, map[string]string{"id": "id1"})
	w := httptest.NewRecorder()

	h.AtualizarStatusPedido(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on update success, got %d", w.Result().StatusCode)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["message"] == nil {
		t.Fatalf("expected success message in response")
	}
}

