package presenters

import (
	"errors"
	"testing"
	"time"

	"pedido-service/internal/core/domain"
)

func TestPedidoPresenter_PresentAndList(t *testing.T) {
	p := &domain.Pedido{ID: "id1", ValorTotal: 10, Status: domain.StatusRecebido, CreatedAt: time.Now()}
	pr := NewPedidoPresenter().Present(p)
	if pr.ID != "id1" {
		t.Fatalf("expected id1, got %s", pr.ID)
	}

	list := NewPedidoPresenter().PresentList([]*domain.Pedido{p})
	if len(list) != 1 {
		t.Fatalf("expected list length 1, got %d", len(list))
	}
}

func TestPedidoPresenter_PresentErrorAndSuccess(t *testing.T) {
	pres := NewPedidoPresenter()
	status, errResp := pres.PresentError(errors.New("some error"))
	if status != 400 {
		t.Fatalf("expected 400 for error, got %d", status)
	}
	m, ok := errResp.(map[string]string)
	if !ok || m["error"] == "" {
		t.Fatalf("expected error map with message")
	}

	status2, success := pres.PresentSuccess("ok", nil)
	if status2 != 200 {
		t.Fatalf("expected 200 for success, got %d", status2)
	}
	if success == nil {
		t.Fatalf("expected non-nil success response")
	}
}
