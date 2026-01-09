package dtos

import (
	"testing"

	"pedido-service/internal/core/domain"
)

func TestPedidoRequest_ToDomain_ErrorWhenPriceMissing(t *testing.T) {
	req := &PedidoRequest{ClienteID: nil, Itens: []ItemPedidoRequest{{ProdutoID: "p1", Quantidade: 1}}}
	_, err := req.ToDomain()
	if err == nil {
		t.Fatalf("expected error when price missing from itens")
	}
}

func TestFromDomain_Converts(t *testing.T) {
	it := domain.ItemPedido{ProdutoID: "p1", Nome: "Prod", Preco: 5, Quantidade: 2, Observacao: "obs"}
	p := &domain.Pedido{ID: "id", ClienteID: nil, Itens: []domain.ItemPedido{it}, ValorTotal: 10}
	r := FromDomain(p)
	if r.ID != "id" {
		t.Fatalf("expected id mapped")
	}
	if len(r.Itens) != 1 || r.Itens[0].Preco != 5 {
		t.Fatalf("expected item price 5")
	}
}
