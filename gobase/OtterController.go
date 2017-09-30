package gobase

import (
	"sync"
)

type AbstractTask interface {
	StartUp()
	Shutdown()
}

type OtterController struct {
	controllers *sync.Map //Map<Long, Map<StageType, GlobalTask>>
	stageAggregationCollector *StageAggregationCollector
	Gltask *GlobalTask
}

func NewOtterController() *OtterController{
	controller:=new(OtterController)
	controller.controllers=new(sync.Map)
	controller.stageAggregationCollector=new(StageAggregationCollector)
	controller.Gltask=NewGlobalTask()
	return controller
}

func (self *OtterController) PutPipe2Controller(pipelineId int64){
	mm:=new(sync.Map)
	self.controllers.Store(pipelineId,mm)
}

func (self *OtterController) StartPipeline(nodeTask *NodeTask){
	pipelineId:=nodeTask.pipeline.Id
	tmMap,ok:=self.controllers.Load(pipelineId)
	if ok{
		mapTask:=tmMap.(*sync.Map)
		for i,stage:=range nodeTask.stageType{
			event:=nodeTask.taskEvent[i]
			if event==CRE{
				self.startTask(pipelineId,mapTask,stage)
			}else{
				self.stopTask(mapTask,stage)
			}
		}
	}
}

func (self *OtterController) startTask(pipelineId int64,tasks *sync.Map,stageType string) {
	var task AbstractTask
	switch stageType {
	case SELE:
		task=NewSelectTask(self.Gltask,pipelineId)
	case EXTR:
	case TRNS:
	case LOAD:
	}
	if task!=nil{
		task.StartUp()
		tasks.Store(stageType,task)
	}
}

func (self *OtterController) stopTask(tasks *sync.Map,stageType string) {
	tt,ok:=tasks.Load(stageType)
	if ok{
		task:=tt.(*GlobalTask)
		if task!=nil{
			task.Shutdown()
		}
	}
}