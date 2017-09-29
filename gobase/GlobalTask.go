package gobase

const (//StageType
	SELE   = "SELECT"
	EXTR  = "EXTRACT"
	TRNS = "TRANSFORM"
	LOAD      = "LOAD"
)

type GlobalTask struct {
	Running bool
	PipelineId int64
	ArbitrateEventService *ArbitrateEventService
	RowDataPipeDelegate *RowDataPipeDelegate
	StageAggregationCollector *StageAggregationCollector
}

func NewGlobalTask(pipelineId int64) *GlobalTask{
	task:=new(GlobalTask)
	task.PipelineId=pipelineId
	return task
}

