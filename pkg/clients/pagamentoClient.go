package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Pagamento struct {
	ID         string  `json:"id"`
	PedidoID   string  `json:"pedido_id"`
	Valor      float64 `json:"valor"`
	Tipo       string  `json:"tipo"`
	Status     string  `json:"status"`
	QRCodeData string  `json:"qr_code_data,omitempty"`
	ExternalID string  `json:"external_id,omitempty"`
	ExpiresAt  *string `json:"expires_at,omitempty"`
}

type CriarPagamentoRequest struct {
	PedidoID string  `json:"pedido_id"`
	Valor    float64 `json:"valor"`
	Tipo     string  `json:"tipo"`
}

type PagamentoClient struct {
	baseURL    string
	httpClient *http.Client
}

func NovoPagamentoClient(baseURL string) *PagamentoClient {
	return &PagamentoClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *PagamentoClient) CriarPagamento(ctx context.Context, pedidoID string, valor float64, tipo string) (*Pagamento, error) {
	url := fmt.Sprintf("%s/api/v1/pagamentos", c.baseURL)

	reqBody := CriarPagamentoRequest{
		PedidoID: pedidoID,
		Valor:    valor,
		Tipo:     tipo,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro ao criar pagamento: status %d, body: %s", resp.StatusCode, string(body))
	}

	var pagamento Pagamento
	if err := json.NewDecoder(resp.Body).Decode(&pagamento); err != nil {
		return nil, err
	}

	return &pagamento, nil
}
