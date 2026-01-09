package dtos

import (
	"pedido-service/internal/core/domain"
	"time"
)

type PedidoRequest struct {
	ClienteID *string             `json:"cliente_id,omitempty"`
	Itens     []ItemPedidoRequest `json:"itens"`
}

type ItemPedidoRequest struct {
	ProdutoID  string `json:"produto_id"`
	Quantidade int    `json:"quantidade"`
	Observacao string `json:"observacao,omitempty"`
}

type PedidoResponse struct {
	ID         string               `json:"id"`
	ClienteID  *string              `json:"cliente_id,omitempty"`
	Itens      []ItemPedidoResponse `json:"itens"`
	ValorTotal float64              `json:"valor_total"`
	Status     string               `json:"status"`
	CreatedAt  time.Time            `json:"created_at"`
}

type ItemPedidoResponse struct {
	ProdutoID  string  `json:"produto_id"`
	Nome       string  `json:"nome"`
	Preco      float64 `json:"preco"`
	Quantidade int     `json:"quantidade"`
	Observacao string  `json:"observacao,omitempty"`
}

func (r *PedidoRequest) ToDomain() (*domain.Pedido, error) {
	itens := make([]domain.ItemPedido, len(r.Itens))
	for i, item := range r.Itens {
		itens[i] = domain.ItemPedido{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
			Observacao: item.Observacao,
		}
	}

	pedido, err := domain.NovoPedido("", r.ClienteID, itens)
	if err != nil {
		return nil, err
	}

	return pedido, nil
}

func FromDomain(p *domain.Pedido) *PedidoResponse {
	itens := make([]ItemPedidoResponse, len(p.Itens))
	for i, item := range p.Itens {
		itens[i] = ItemPedidoResponse{
			ProdutoID:  item.ProdutoID,
			Nome:       item.Nome,
			Preco:      item.Preco,
			Quantidade: item.Quantidade,
			Observacao: item.Observacao,
		}
	}

	return &PedidoResponse{
		ID:         p.ID,
		ClienteID:  p.ClienteID,
		Itens:      itens,
		ValorTotal: p.ValorTotal,
		Status:     string(p.Status),
		CreatedAt:  p.CreatedAt,
	}
}
