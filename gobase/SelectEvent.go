package gobase

import (
	"time"
)

type SelectMemoryArbitrateEvent struct{
	Gltask *GlobalTask
}

func NewSelectMemoryArbitrateEvent(gltask *GlobalTask) *SelectMemoryArbitrateEvent{
	event:=new(SelectMemoryArbitrateEvent)
	event.Gltask=gltask
	return event
}


func (self *SelectMemoryArbitrateEvent) Await(pipelineId int64) *EtlEventData{
	stageController:=self.Gltask.ArbitrateFactory.GetInstance(pipelineId)
	processId:=stageController.WaitForProcess(SELE)
	eventData:=new(EtlEventData)
	eventData.ProceEvdata.PipelineId=pipelineId
	eventData.ProceEvdata.ProcessId=processId
	eventData.ProceEvdata.StartTime=time.Now()
	//Long nid = ArbitrateConfigUtils.getCurrentNid();
	//eventData.setCurrNid(nid);
	//eventData.setNextNid(nid);
	return eventData
}

func (self *SelectMemoryArbitrateEvent) Single(data *EtlEventData){

}