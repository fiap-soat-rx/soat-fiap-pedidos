package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pedido-service/internal/core/domain"
	"pedido-service/pkg/clients"
)

// Mock repository
type MockPedidoRepository struct {
	CriarFunc         func(ctx context.Context, pedido *domain.Pedido) error
	BuscarPorIDFunc   func(ctx context.Context, id string) (*domain.Pedido, error)
	ListarFunc        func(ctx context.Context) ([]*domain.Pedido, error)
	ListarPorStatusFunc func(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error)
	ListarPorClienteFunc func(ctx context.Context, clienteID string) ([]*domain.Pedido, error)
	AtualizarFunc     func(ctx context.Context, pedido *domain.Pedido) error
}

func (m *MockPedidoRepository) Criar(ctx context.Context, pedido *domain.Pedido) error {
	if m.CriarFunc != nil {
		return m.CriarFunc(ctx, pedido)
	}
	return nil
}
func (m *MockPedidoRepository) BuscarPorID(ctx context.Context, id string) (*domain.Pedido, error) {
	if m.BuscarPorIDFunc != nil {
		return m.BuscarPorIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *MockPedidoRepository) Listar(ctx context.Context) ([]*domain.Pedido, error) {
	if m.ListarFunc != nil {
		return m.ListarFunc(ctx)
	}
	return nil, nil
}
func (m *MockPedidoRepository) ListarPorStatus(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error) {
	if m.ListarPorStatusFunc != nil {
		return m.ListarPorStatusFunc(ctx, status)
	}
	return nil, nil
}
func (m *MockPedidoRepository) ListarPorCliente(ctx context.Context, clienteID string) ([]*domain.Pedido, error) {
	if m.ListarPorClienteFunc != nil {
		return m.ListarPorClienteFunc(ctx, clienteID)
	}
	return nil, nil
}
func (m *MockPedidoRepository) Atualizar(ctx context.Context, pedido *domain.Pedido) error {
	if m.AtualizarFunc != nil {
		return m.AtualizarFunc(ctx, pedido)
	}
	return nil
}

// helper to create a produto http test server
func produtoServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(handler))
}

// helper to create a cliente http test server
func clienteServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(handler))
}


func TestCriarPedido_ClienteValidationFails(t *testing.T) {
	repo := &MockPedidoRepository{}

	cliSrv := clienteServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer cliSrv.Close()

	prodSrv := produtoServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(clients.Produto{ID: "p1", Nome: "Prod", Preco: 1.0, Disponivel: true})
	})
	defer prodSrv.Close()

	prodClient := clients.NovoProdutoClient(prodSrv.URL)
	cliClient := clients.NovoClienteClient(cliSrv.URL)

	service := NovoPedidoService(repo, prodClient, cliClient)

	_, err := service.CriarPedido(context.Background(), ptr("cliente1"), []domain.ItemPedido{{ProdutoID: "p1", Quantidade: 1}})
	if err == nil {
		t.Fatalf("expected error when cliente validation fails")
	}
}

func TestCriarPedido_ProdutoValidationFails(t *testing.T) {
	repo := &MockPedidoRepository{}

	prodSrv := produtoServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer prodSrv.Close()

	prodClient := clients.NovoProdutoClient(prodSrv.URL)
	service := NovoPedidoService(repo, prodClient, nil)

	_, err := service.CriarPedido(context.Background(), nil, []domain.ItemPedido{{ProdutoID: "p1", Quantidade: 1}})
	if err == nil {
		t.Fatalf("expected error when produto validation fails")
	}
}

func TestCriarPedido_Success(t *testing.T) {
	called := false
	repo := &MockPedidoRepository{CriarFunc: func(ctx context.Context, pedido *domain.Pedido) error {
		called = true
		return nil
	}}

	prodSrv := produtoServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(clients.Produto{ID: "p1", Nome: "Prod", Preco: 2.5, Disponivel: true})
	})
	defer prodSrv.Close()

	cliSrv := clienteServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(clients.Cliente{ID: "c1"})
	})
	defer cliSrv.Close()

	prodClient := clients.NovoProdutoClient(prodSrv.URL)
	cliClient := clients.NovoClienteClient(cliSrv.URL)
	service := NovoPedidoService(repo, prodClient, cliClient)

	pedido, err := service.CriarPedido(context.Background(), ptr("c1"), []domain.ItemPedido{{ProdutoID: "p1", Quantidade: 2}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected repository Criar to be called")
	}
	if pedido.ValorTotal != 5.0 {
		t.Fatalf("expected total 5.0, got %v", pedido.ValorTotal)
	}
	if pedido.Itens[0].Nome != "Prod" {
		t.Fatalf("expected item nome to be filled by produto client")
	}
}

func TestListarPedidosPorStatus_Invalid(t *testing.T) {
	repo := &MockPedidoRepository{}
	service := NovoPedidoService(repo, nil, nil)
	_, err := service.ListarPedidosPorStatus(context.Background(), domain.StatusPedido("X"))
	if err == nil {
		t.Fatalf("expected error for invalid status")
	}
}

func TestAtualizarStatusPedido_InvalidStatus(t *testing.T) {
	repo := &MockPedidoRepository{}
	service := NovoPedidoService(repo, nil, nil)
	err := service.AtualizarStatusPedido(context.Background(), "id", domain.StatusPedido("X"))
	if err == nil {
		t.Fatalf("expected error for invalid status")
	}
}

func TestAtualizarStatusPedido_NotFound(t *testing.T) {
	repo := &MockPedidoRepository{BuscarPorIDFunc: func(ctx context.Context, id string) (*domain.Pedido, error) {
		return nil, nil
	}}
	service := NovoPedidoService(repo, nil, nil)
	err := service.AtualizarStatusPedido(context.Background(), "id", domain.StatusEmPreparacao)
	if err == nil {
		t.Fatalf("expected error when pedido not found")
	}
}

func TestAtualizarStatusPedido_InvalidTransition(t *testing.T) {
	repo := &MockPedidoRepository{BuscarPorIDFunc: func(ctx context.Context, id string) (*domain.Pedido, error) {
		p := &domain.Pedido{ID: id, Status: domain.StatusRecebido}
		return p, nil
	}}
	service := NovoPedidoService(repo, nil, nil)
	err := service.AtualizarStatusPedido(context.Background(), "id", domain.StatusFinalizado)
	if err == nil {
		t.Fatalf("expected error for invalid transition")
	}
}

func TestAtualizarStatusPedido_Success(t *testing.T) {
	updated := false
	repo := &MockPedidoRepository{
		BuscarPorIDFunc: func(ctx context.Context, id string) (*domain.Pedido, error) {
			p := &domain.Pedido{ID: id, Status: domain.StatusRecebido}
			return p, nil
		},
		AtualizarFunc: func(ctx context.Context, pedido *domain.Pedido) error {
			updated = true
			return nil
		},
	}
	service := NovoPedidoService(repo, nil, nil)
	err := service.AtualizarStatusPedido(context.Background(), "id", domain.StatusEmPreparacao)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updated {
		t.Fatalf("expected Atualizar to be called")
	}
}

// helpers
func ptr(s string) *string { return &s }
