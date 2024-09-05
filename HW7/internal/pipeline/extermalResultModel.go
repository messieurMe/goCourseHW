package pipeline

type ExternalResultId string

const (
	WarehouseId ExternalResultId = "WHID"
	DropPointId ExternalResultId = "DPID"
)

type ExternalResult interface {
	getUniqueExternalResultId() ExternalResultId
}

type WarehouseIdHolder struct {
	warehouseId int
}

func (w *WarehouseIdHolder) getUniqueExternalResultId() ExternalResultId {
	return WarehouseId
}

type DropPointIdHolder struct {
	dropPointId int
}

func (w *DropPointIdHolder) getUniqueExternalResultId() ExternalResultId {
	return DropPointId
}

type FullExternalResult struct {
	warehouseId int
	dropPointId int
}
