package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Cliente struct {
	ID       string `json:"id"`
	Nome     string `json:"Nome"`
	CPF      string `json:"cpf"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
}

type ClienteClient struct {
	baseURL    string
	httpClient *http.Client
}

func NovoClienteClient(baseURL string) *ClienteClient {
	return &ClienteClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *ClienteClient) BuscarPorID(ctx context.Context, id string) (*Cliente, error) {
	url := fmt.Sprintf("%s/api/v1/clientes/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro ao buscar cliente: %s", string(body))
	}

	var cliente Cliente
	if err := json.NewDecoder(resp.Body).Decode(&cliente); err != nil {
		return nil, err
	}

	return &cliente, nil
}

func (c *ClienteClient) ValidarCliente(ctx context.Context, clienteID string) error {
	cliente, err := c.BuscarPorID(ctx, clienteID)
	if err != nil {
		return err
	}

	if cliente == nil {
		return errors.New("cliente não encontrado")
	}

	return nil
}
