package generator

import (
	"context"
	"hw7/internal/model"
)

type OrderGenerator interface {
	GenerateOrdersStream(ctx context.Context, orders []model.OrderInitialized) <-chan model.OrderInitialized
}

type OrderGeneratorImplementation struct{}

func NewOrderGeneratorImplementation() *OrderGeneratorImplementation {
	return &OrderGeneratorImplementation{}
}

func (o *OrderGeneratorImplementation) GenerateOrdersStream(ctx context.Context, orders []model.OrderInitialized) <-chan model.OrderInitialized {
	onii := make(chan model.OrderInitialized)

	go func() {
		for _, v := range orders {
			select {
			case <-ctx.Done():
				break
			default:
				onii <- v
			}
		}
		close(onii)
	}()

	return onii
}
