package gobase


type GlobalTask struct {
	running bool
	pipelineId int64
	arbitrateEventService *ArbitrateEventService
	rowDataPipeDelegate *RowDataPipeDelegate
	stageAggregationCollector *StageAggregationCollector
}

