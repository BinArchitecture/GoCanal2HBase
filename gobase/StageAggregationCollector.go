package gobase

import (
	"sync"
)

const (//StageType
	SELE   = "SELECT"
	EXTR  = "EXTRACT"
	TRNS = "TRANSFORM"
	LOAD      = "LOAD"
)

type StageAggregationCollector struct{
	collector *sync.Map //Map<Long, Map<StageType, StageAggregation>>
	profiling bool
}

func NewStageAggregationCollector() *StageAggregationCollector{
	collector:=new(StageAggregationCollector)
	collector.profiling=true
	//bufferSize:=1024
	collector.collector=new(sync.Map)
	return collector
}

func (self *StageAggregationCollector) push(pipelineId int64,stageType string ,aggregationItem AggregationItem){
	tmpMap,_:=self.collector.Load(pipelineId)
	if tmpMap!=nil{
		tmpmm:=tmpMap.(sync.Map)
		tmpq,_:=tmpmm.Load(stageType)
		if tmpq!=nil{

		}
	}


}
