package gobase

import (
	"time"
)

type TerminEventData struct{
	ProceEvdata *ProcessEventData
	TerminType  string
	Code        string
	Desc        string
	CurrNid     int64
}

type TerminMemoryArbitrateEvent struct{
	Gltask *GlobalTask
}

func NewTerminMemoryArbitrateEvent(gltask *GlobalTask) *TerminMemoryArbitrateEvent{
	event:=new(TerminMemoryArbitrateEvent)
	event.Gltask=gltask
	return event
}

func (self *TerminMemoryArbitrateEvent) Await(pipelineId int64) *EtlEventData{
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

func (self *TerminMemoryArbitrateEvent) Single(data *EtlEventData){

}