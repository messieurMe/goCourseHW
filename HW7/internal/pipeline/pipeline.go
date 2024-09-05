package pipeline

import (
	"context"
	"errors"
	"hw7/internal/model"
	"hw7/internal/semaphore"
	"sync"
)

const fanLimit = 10

type OrderPipeline interface {
	Start(ctx context.Context, actions model.OrderActions, orders <-chan model.OrderInitialized, processed chan<- model.OrderProcessFinished)
}

type OrderPipelineImplementation struct{}

func NewOrderPipelineImplementation() *OrderPipelineImplementation {
	return &OrderPipelineImplementation{}
}

type Utils struct {
	ctx       context.Context
	actions   *model.OrderActions
	semaphore *semaphore.Semaphore
}

func (o *OrderPipelineImplementation) Start(
	ctx context.Context,
	actions *model.OrderActions,
	orders <-chan model.OrderInitialized,
	processed chan<- model.OrderProcessFinished) {

	go func() {
		semaphoreChan := make(chan struct{}, fanLimit)
		utils := &Utils{
			ctx:       ctx,
			actions:   actions,
			semaphore: &semaphore.Semaphore{SemC: semaphoreChan},
		}

		initToStarted := addPipelineStage[model.OrderInitialized, model.OrderProcessStarted](
			utils,
			orders,
			proceedInitToStart,
		)

		startedToFinishedExternal := addPipelineStage[model.OrderProcessStarted, model.OrderFinishedExternalInteraction](
			utils,
			initToStarted,
			processStartedToFinishedExternal,
		)

		finished := addPipelineStage[model.OrderFinishedExternalInteraction, model.OrderProcessFinished](
			utils,
			startedToFinishedExternal,
			processFinishedExternalToFinishProcess)

		for i := range finished {
			processed <- i
		}
		close(processed)
	}()
}

func addPipelineStage[T model.OrderStates, K model.OrderStates](
	pipelineUtils *Utils,
	orders <-chan T,
	body func(*Utils, *T, chan<- K),
) <-chan K {

	outCh := make(chan K)

	go func() {

		defer func() {
			close(outCh)
		}()

	label:
		for order := range orders {
			select {
			case <-pipelineUtils.ctx.Done():
				break label
			default:
				body(pipelineUtils, &order, outCh)
			}
		}
	}()
	return outCh
}

func proceedInitToStart(pipelineUtils *Utils, order *model.OrderInitialized, outCh chan<- model.OrderProcessStarted) {
	defer func() {
		if err := recover(); err != nil {
			outCh <- model.OrderProcessStarted{
				OrderInitialized: *order,
				OrderStates:      order.OrderStates,
				Error:            errors.New(err.(string)),
			}
		}
	}()

	pipelineUtils.actions.InitToStarted()
	outCh <- model.OrderProcessStarted{
		OrderInitialized: *order,
		OrderStates:      append(order.OrderStates, model.ProcessStarted),
		Error:            nil,
	}
}

func processStartedToFinishedExternal(pipelineUtils *Utils, order *model.OrderProcessStarted, outCh chan<- model.OrderFinishedExternalInteraction) {
	defer func() {
		if err := recover(); err != nil {
			outCh <- model.OrderFinishedExternalInteraction{
				OrderProcessStarted: *order,
				StorageID:           0,
				PickupPointID:       0,
				OrderStates:         order.OrderStates,
				Error:               errors.New(err.(string)),
			}
		}
	}()

	if order.Error != nil {
		outCh <- model.OrderFinishedExternalInteraction{
			OrderProcessStarted: *order,
			StorageID:           0,
			PickupPointID:       0,
			OrderStates:         order.OrderStates,
			Error:               order.Error,
		}
	} else {
		pipelineUtils.actions.StartedToFinishedExternalInteraction()
		startedOrder := order

		fanned := fanIn(*startedOrder, fanOut(startedOrder, pipelineUtils.semaphore))
		for i := range fanned {
			outCh <- i
		}
	}
}

func fanOut(startedOrder *model.OrderProcessStarted, semaphore *semaphore.Semaphore) *[]<-chan ExternalResult {
	tasks := [2]func(startedOrder model.OrderProcessStarted, output chan<- ExternalResult){
		retrieveWarehouseId,
		retrieveDropPointId,
	}

	inChannels := make([]<-chan ExternalResult, 0)

	for _, task := range tasks {
		task := task
		channel := make(chan ExternalResult)
		inChannels = append(inChannels, channel)
		semaphore.Acquire()

		go func(
			currentTask func(startedOrder model.OrderProcessStarted, output chan<- ExternalResult),
			currentStartedOrder *model.OrderProcessStarted,
			currentOutput chan<- ExternalResult,
		) {
			currentTask(*currentStartedOrder, currentOutput)
			semaphore.Release()
			close(currentOutput)
		}(task, startedOrder, channel)
	}
	return &inChannels
}

func retrieveWarehouseId(startedOrder model.OrderProcessStarted, output chan<- ExternalResult) {
	output <- &WarehouseIdHolder{(startedOrder.OrderInitialized.ProductID % 2) + 1}
}

func retrieveDropPointId(startedOrder model.OrderProcessStarted, output chan<- ExternalResult) {
	output <- &DropPointIdHolder{(startedOrder.OrderInitialized.ProductID % 3) + 1}
}

func fanIn(startedOrder model.OrderProcessStarted, channels *[]<-chan ExternalResult) <-chan model.OrderFinishedExternalInteraction {

	mergedChannels := make(chan ExternalResult)
	outputChannel := make(chan model.OrderFinishedExternalInteraction)

	wg := sync.WaitGroup{}

	for _, channel := range *channels {
		channel := channel
		wg.Add(1)
		go func() {
			for result := range channel {
				mergedChannels <- result
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(mergedChannels)
	}()

	go func() {
		result := FullExternalResult{}

		readyFlags := map[ExternalResultId]bool{}

		for externalResult := range mergedChannels {
			switch externalResult.getUniqueExternalResultId() {
			case WarehouseId:
				result.warehouseId = externalResult.(*WarehouseIdHolder).warehouseId
				readyFlags[WarehouseId] = true
			case DropPointId:
				result.dropPointId = externalResult.(*DropPointIdHolder).dropPointId
				readyFlags[DropPointId] = true
			}

			if len(readyFlags) == 2 {
				outputChannel <- model.OrderFinishedExternalInteraction{
					OrderProcessStarted: startedOrder,
					StorageID:           result.warehouseId,
					PickupPointID:       result.dropPointId,
					OrderStates:         append(startedOrder.OrderStates, model.FinishedExternalInteraction),
				}
			}
		}
		close(outputChannel)
	}()
	return outputChannel
}

func processFinishedExternalToFinishProcess(pipelineUtils *Utils, order *model.OrderFinishedExternalInteraction, outCh chan<- model.OrderProcessFinished) {
	defer func() {
		if err := recover(); err != nil {
			outCh <- model.OrderProcessFinished{
				OrderFinishedExternalInteraction: *order,
				OrderStates:                      order.OrderStates,
				Error:                            errors.New(err.(string)),
			}
		}
	}()

	if order.Error != nil {
		outCh <- model.OrderProcessFinished{
			OrderFinishedExternalInteraction: *order,
			OrderStates:                      order.OrderStates,
			Error:                            order.Error,
		}
	} else {
		pipelineUtils.actions.FinishedExternalInteractionToProcessFinished()
		outCh <- model.OrderProcessFinished{
			OrderFinishedExternalInteraction: *order,
			OrderStates:                      append(order.OrderStates, model.ProcessFinished),
		}
	}
}
