package repositories

import (
	"context"
	"errors"
	"pedido-service/internal/core/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PedidoRepository struct {
	collection *mongo.Collection
}

func NovoPedidoRepository(client *mongo.Client, dbName string) *PedidoRepository {
	collection := client.Database(dbName).Collection("pedidos")
	return &PedidoRepository{
		collection: collection,
	}
}

func (r *PedidoRepository) Criar(ctx context.Context, pedido *domain.Pedido) error {
	_, err := r.collection.InsertOne(ctx, pedido)
	return err
}

func (r *PedidoRepository) BuscarPorID(ctx context.Context, id string) (*domain.Pedido, error) {
	var pedido domain.Pedido
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&pedido)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pedido, nil
}

func (r *PedidoRepository) Listar(ctx context.Context) ([]*domain.Pedido, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pedidos []*domain.Pedido
	if err = cursor.All(ctx, &pedidos); err != nil {
		return nil, err
	}

	return pedidos, nil
}

func (r *PedidoRepository) ListarPorStatus(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"status": status}, options.Find().SetSort(bson.M{"created_at": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pedidos []*domain.Pedido
	if err = cursor.All(ctx, &pedidos); err != nil {
		return nil, err
	}

	return pedidos, nil
}

func (r *PedidoRepository) ListarPorCliente(ctx context.Context, clienteID string) ([]*domain.Pedido, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"cliente_id": clienteID}, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pedidos []*domain.Pedido
	if err = cursor.All(ctx, &pedidos); err != nil {
		return nil, err
	}

	return pedidos, nil
}

func (r *PedidoRepository) Atualizar(ctx context.Context, pedido *domain.Pedido) error {
	pedido.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": pedido.ID},
		bson.M{"$set": pedido},
	)
	if err == mongo.ErrNoDocuments {
		return errors.New("pedido não encontrado")
	}
	return err
}
