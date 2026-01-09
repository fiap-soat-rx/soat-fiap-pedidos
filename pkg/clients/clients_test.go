package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPagamentoClient_CriarPagamento_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Pagamento{ID: "pay1", PedidoID: "p1", Valor: 10})
	}))
	defer srv.Close()

	cli := NovoPagamentoClient(srv.URL)
	p, err := cli.CriarPagamento(context.Background(), "p1", 10, "PIX")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.ID != "pay1" {
		t.Fatalf("expected pay1 id")
	}
}

func TestPagamentoClient_CriarPagamento_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	cli := NovoPagamentoClient(srv.URL)
	_, err := cli.CriarPagamento(context.Background(), "p1", 10, "PIX")
	if err == nil {
		t.Fatalf("expected error on invalid json")
	}
}
func TestPagamentoClient_CriarPagamento_NonCreated(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("err"))
	}))
	defer srv.Close()

	cli := NovoPagamentoClient(srv.URL)
	_, err := cli.CriarPagamento(context.Background(), "p1", 10, "PIX")
	if err == nil {
		t.Fatalf("expected error on non-created status")
	}
}

func TestProdutoClient_BuscarPorID_Various(t *testing.T) {
	// success
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Produto{ID: "prod1", Nome: "X", Preco: 5, Disponivel: true})
	}))
	defer srv.Close()

	pc := NovoProdutoClient(srv.URL)
	p, err := pc.BuscarPorID(context.Background(), "prod1")
	if err != nil || p == nil || p.ID != "prod1" {
		t.Fatalf("expected product returned")
	}

	// not found
	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv2.Close()
	pc2 := NovoProdutoClient(srv2.URL)
	p2, err2 := pc2.BuscarPorID(context.Background(), "p2")
	if err2 != nil {
		t.Fatalf("expected no error on 404, got %v", err2)
	}
	if p2 != nil {
		t.Fatalf("expected nil product on 404")
	}

	// server error
	srv3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("boom"))
	}))
	defer srv3.Close()
	pc3 := NovoProdutoClient(srv3.URL)
	_, err3 := pc3.BuscarPorID(context.Background(), "p3")
	if err3 == nil {
		t.Fatalf("expected error on non-200 status")
	}
}

func TestProdutoClient_ValidarProduto_Errors(t *testing.T) {
	// not found
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	pc := NovoProdutoClient(srv.URL)
	_, err := pc.ValidarProduto(context.Background(), "p1")
	if err == nil {
		t.Fatalf("expected error when produto not found")
	}

	// unavailable
	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Produto{ID: "p", Disponivel: false})
	}))
	defer srv2.Close()
	pc2 := NovoProdutoClient(srv2.URL)
	_, err2 := pc2.ValidarProduto(context.Background(), "p")
	if err2 == nil {
		t.Fatalf("expected error when produto not available")
	}
}
func TestClienteClient_BuscarPorID_And_Validar(t *testing.T) {
	// not found
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	cc := NovoClienteClient(srv.URL)
	c, err := cc.BuscarPorID(context.Background(), "c1")
	if err != nil {
		t.Fatalf("expected no error on 404, got %v", err)
	}
	if c != nil {
		t.Fatalf("expected nil cliente on 404")
	}

	// success and validar
	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Cliente{ID: "c1", Nome: "name"})
	}))
	defer srv2.Close()
	cc2 := NovoClienteClient(srv2.URL)
	err2 := cc2.ValidarCliente(context.Background(), "c1")
	if err2 != nil {
		t.Fatalf("expected no error on validar cliente success")
	}

	// invalid json
	srv3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not-json"))
	}))
	defer srv3.Close()
	cc3 := NovoClienteClient(srv3.URL)
	_, err3 := cc3.BuscarPorID(context.Background(), "c2")
	if err3 == nil {
		t.Fatalf("expected error when response invalid json")
	}
}
