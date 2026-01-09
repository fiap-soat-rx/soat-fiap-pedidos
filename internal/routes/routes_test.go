package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"pedido-service/internal/adapters/primary/handlers"
)

func TestConfigurarRotas_Health(t *testing.T) {
	r := mux.NewRouter()
	h := handlers.NovoHealthHandler("v1")
	// pass empty PedidoHandler (we won't call its routes in this test)
	ConfigurarRotas(r, &handlers.PedidoHandler{}, h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on health route, got %d", w.Result().StatusCode)
	}
}
