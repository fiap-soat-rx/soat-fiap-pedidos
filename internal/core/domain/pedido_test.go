package domain

import (
	"testing"
)

func TestNovoPedido_ValidAndInvalid(t *testing.T) {
	itens := []ItemPedido{{ProdutoID: "p1", Nome: "P1", Preco: 2.5, Quantidade: 2}}
	p, err := NovoPedido("id1", nil, itens)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ValorTotal != 5.0 {
		t.Fatalf("expected total 5.0, got %v", p.ValorTotal)
	}

	// no items
	_, err = NovoPedido("id", nil, []ItemPedido{})
	if err == nil {
		t.Fatalf("expected error for empty items")
	}

	// item with zero quantity
	_, err = NovoPedido("id", nil, []ItemPedido{{ProdutoID: "p", Quantidade: 0, Preco: 1}})
	if err == nil {
		t.Fatalf("expected error for zero quantity")
	}

	// item with zero price
	_, err = NovoPedido("id", nil, []ItemPedido{{ProdutoID: "p", Quantidade: 1, Preco: 0}})
	if err == nil {
		t.Fatalf("expected error for zero price")
	}
}

func TestCalcularValorTotal(t *testing.T) {
	p := &Pedido{Itens: []ItemPedido{{ProdutoID: "p1", Preco: 1.5, Quantidade: 2}, {ProdutoID: "p2", Preco: 2, Quantidade: 3}}}
	p.CalcularValorTotal()
	if p.ValorTotal != 1.5*2+2*3 {
		t.Fatalf("unexpected total: %v", p.ValorTotal)
	}
}

func TestAtualizarStatusAndIsStatusValido(t *testing.T) {
	p := &Pedido{Status: StatusRecebido}
	p.AtualizarStatus(StatusEmPreparacao)
	if p.Status != StatusEmPreparacao {
		t.Fatalf("status not updated")
	}

	if !IsStatusValido(StatusFinalizado) {
		t.Fatalf("expected status finalizado to be valid")
	}
	if IsStatusValido(StatusPedido("X")) {
		t.Fatalf("expected unknown status to be invalid")
	}
}
