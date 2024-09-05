package workerpool

import (
	"context"
	"hw7/internal/model"
	"hw7/internal/pipeline"
	"sync"
)

type OrderWorkerPool interface {
	StartWorkerPool(ctx context.Context, orders <-chan model.OrderInitialized, additionalActions model.OrderActions, workersCount int) <-chan model.OrderProcessFinished
}

type OrderWorkerPoolImplementation struct{}

func NewOrderWorkerPoolImplementation() *OrderWorkerPoolImplementation {
	return &OrderWorkerPoolImplementation{}
}

func (o *OrderWorkerPoolImplementation) StartWorkerPool(ctx context.Context, orders <-chan model.OrderInitialized, additionalActions model.OrderActions, workersCount int) <-chan model.OrderProcessFinished {
	resultChan := make(chan model.OrderProcessFinished)

	wg := sync.WaitGroup{}

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func() {
			workerResultChan := make(chan model.OrderProcessFinished)

			pipeline.NewOrderPipelineImplementation().Start(ctx, &additionalActions, orders, workerResultChan)

			for workerResult := range workerResultChan {
				resultChan <- workerResult
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}
