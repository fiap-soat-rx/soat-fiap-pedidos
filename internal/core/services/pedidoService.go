package services

import (
	"context"
	"errors"
	"fmt"
	"pedido-service/internal/core/domain"
	"pedido-service/internal/core/ports"
	"pedido-service/pkg/clients"

	"github.com/google/uuid"
)

type PedidoService struct {
	pedidoRepository ports.PedidoRepository
	produtoClient    *clients.ProdutoClient
	clienteClient    *clients.ClienteClient
}

func NovoPedidoService(
	pedidoRepository ports.PedidoRepository,
	produtoClient *clients.ProdutoClient,
	clienteClient *clients.ClienteClient,
) *PedidoService {
	return &PedidoService{
		pedidoRepository: pedidoRepository,
		produtoClient:    produtoClient,
		clienteClient:    clienteClient,
	}
}

func (s *PedidoService) CriarPedido(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {
	if clienteID != nil && *clienteID != "" {
		if err := s.clienteClient.ValidarCliente(ctx, *clienteID); err != nil {
			return nil, fmt.Errorf("erro ao validar cliente: %w", err)
		}
	}

	for i, item := range itens {
		produto, err := s.produtoClient.ValidarProduto(ctx, item.ProdutoID)
		if err != nil {
			return nil, err
		}

		itens[i].Nome = produto.Nome
		itens[i].Preco = produto.Preco
	}

	id := uuid.New().String()
	pedido, err := domain.NovoPedido(id, clienteID, itens)
	if err != nil {
		return nil, err
	}

	err = s.pedidoRepository.Criar(ctx, pedido)
	if err != nil {
		return nil, err
	}

	return pedido, nil
}

func (s *PedidoService) BuscarPedidoPorID(ctx context.Context, id string) (*domain.Pedido, error) {
	return s.pedidoRepository.BuscarPorID(ctx, id)
}

func (s *PedidoService) ListarPedidos(ctx context.Context) ([]*domain.Pedido, error) {
	return s.pedidoRepository.Listar(ctx)
}

func (s *PedidoService) ListarPedidosPorStatus(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error) {
	if !domain.IsStatusValido(status) {
		return nil, errors.New("status inválido")
	}
	return s.pedidoRepository.ListarPorStatus(ctx, status)
}

func (s *PedidoService) ListarPedidosPorCliente(ctx context.Context, clienteID string) ([]*domain.Pedido, error) {
	return s.pedidoRepository.ListarPorCliente(ctx, clienteID)
}

func (s *PedidoService) AtualizarStatusPedido(ctx context.Context, id string, status domain.StatusPedido) error {
	if !domain.IsStatusValido(status) {
		return errors.New("status inválido")
	}

	pedido, err := s.pedidoRepository.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if pedido == nil {
		return errors.New("pedido não encontrado")
	}

	if !s.isTransicaoValida(pedido.Status, status) {
		return fmt.Errorf("transição de status inválida: de %s para %s", pedido.Status, status)
	}

	pedido.AtualizarStatus(status)
	return s.pedidoRepository.Atualizar(ctx, pedido)
}

func (s *PedidoService) isTransicaoValida(statusAtual, novoStatus domain.StatusPedido) bool {
	transicoesValidas := map[domain.StatusPedido][]domain.StatusPedido{
		domain.StatusRecebido: {
			domain.StatusEmPreparacao,
			domain.StatusCancelado,
		},
		domain.StatusEmPreparacao: {
			domain.StatusPronto,
			domain.StatusCancelado,
		},
		domain.StatusPronto: {
			domain.StatusFinalizado,
		},
		domain.StatusFinalizado: {},
		domain.StatusCancelado:  {},
	}

	statusPermitidos, existe := transicoesValidas[statusAtual]
	if !existe {
		return false
	}

	for _, statusPermitido := range statusPermitidos {
		if statusPermitido == novoStatus {
			return true
		}
	}

	return false
}
