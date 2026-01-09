package presenters

import (
	"pedido-service/internal/adapters/primary/dtos"
	"pedido-service/internal/core/domain"
)

type PedidoPresenter struct{}

func NewPedidoPresenter() *PedidoPresenter {
	return &PedidoPresenter{}
}

func (p *PedidoPresenter) Present(pedido *domain.Pedido) *dtos.PedidoResponse {
	return dtos.FromDomain(pedido)
}

func (p *PedidoPresenter) PresentList(pedidos []*domain.Pedido) []*dtos.PedidoResponse {
	result := make([]*dtos.PedidoResponse, len(pedidos))
	for i, pedido := range pedidos {
		result[i] = p.Present(pedido)
	}
	return result
}

func (p *PedidoPresenter) PresentError(err error) (int, interface{}) {
	return 400, map[string]string{"error": err.Error()}
}

func (p *PedidoPresenter) PresentSuccess(message string, data interface{}) (int, interface{}) {
	return 200, map[string]interface{}{
		"message": message,
		"data":    data,
	}
}

