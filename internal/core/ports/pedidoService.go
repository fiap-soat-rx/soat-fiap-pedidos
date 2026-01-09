package ports

import (
	"context"
	"pedido-service/internal/core/domain"
)

type PedidoService interface {
	CriarPedido(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error)
	BuscarPedidoPorID(ctx context.Context, id string) (*domain.Pedido, error)
	ListarPedidos(ctx context.Context) ([]*domain.Pedido, error)
	ListarPedidosPorStatus(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error)
	ListarPedidosPorCliente(ctx context.Context, clienteID string) ([]*domain.Pedido, error)
	AtualizarStatusPedido(ctx context.Context, id string, status domain.StatusPedido) error
}

