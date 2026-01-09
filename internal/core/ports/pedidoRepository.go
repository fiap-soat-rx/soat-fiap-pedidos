package ports

import (
	"context"
	"pedido-service/internal/core/domain"
)

type PedidoRepository interface {
	Criar(ctx context.Context, pedido *domain.Pedido) error
	BuscarPorID(ctx context.Context, id string) (*domain.Pedido, error)
	Listar(ctx context.Context) ([]*domain.Pedido, error)
	ListarPorStatus(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error)
	ListarPorCliente(ctx context.Context, clienteID string) ([]*domain.Pedido, error)
	Atualizar(ctx context.Context, pedido *domain.Pedido) error
}

