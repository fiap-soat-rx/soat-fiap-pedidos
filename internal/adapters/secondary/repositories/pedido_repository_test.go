package repositories

import (
	"context"
	"testing"
)

func TestPedidoRepository_PanicOnNilCollection(t *testing.T) {
	r := &PedidoRepository{collection: nil}
	// Criar should panic due to nil collection
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic when collection is nil on Criar")
			}
		}()
		_ = r.Criar(context.Background(), nil)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic when collection is nil on BuscarPorID")
			}
		}()
		_, _ = r.BuscarPorID(context.Background(), "id")
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic when collection is nil on Listar")
			}
		}()
		_, _ = r.Listar(context.Background())
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic when collection is nil on Atualizar")
			}
		}()
		_ = r.Atualizar(context.Background(), nil)
	}()
}
