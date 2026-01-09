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

type Produto struct {
	ID         string  `json:"id"`
	Nome       string  `json:"Nome"`
	Descricao  string  `json:"descricao"`
	Preco      float64 `json:"preco"`
	Categoria  string  `json:"categoria"`
	Disponivel bool    `json:"disponivel"`
}

type ProdutoClient struct {
	baseURL    string
	httpClient *http.Client
}

func NovoProdutoClient(baseURL string) *ProdutoClient {
	return &ProdutoClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *ProdutoClient) BuscarPorID(ctx context.Context, id string) (*Produto, error) {
	url := fmt.Sprintf("%s/api/v1/produtos/%s", c.baseURL, id)

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
		return nil, fmt.Errorf("erro ao buscar produto: %s", string(body))
	}

	var produto Produto
	if err := json.NewDecoder(resp.Body).Decode(&produto); err != nil {
		return nil, err
	}

	return &produto, nil
}

func (c *ProdutoClient) ValidarProduto(ctx context.Context, produtoID string) (*Produto, error) {
	produto, err := c.BuscarPorID(ctx, produtoID)
	if err != nil {
		return nil, err
	}

	if produto == nil {
		return nil, errors.New("produto não encontrado")
	}

	if !produto.Disponivel {
		return nil, errors.New("produto não disponível")
	}

	return produto, nil
}
