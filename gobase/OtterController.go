package gobase

import "sync"

type OtterController struct {
	controllers *sync.Map //Map<Long, Map<StageType, GlobalTask>>
	stageAggregationCollector *StageAggregationCollector
}

func (self *OtterController) StartPipeline(nodeTask *NodeTask){
	//pipelineId:=nodeTask.pipeline.id
}
