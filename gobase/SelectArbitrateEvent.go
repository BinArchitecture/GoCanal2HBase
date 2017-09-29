package gobase

import "time"

type SelectMemoryArbitrateEvent struct{
}

func (self *SelectMemoryArbitrateEvent) Await(pipelineId int64) *EtlEventData{
	stageController:=NewMemoryStageController(pipelineId)
	processId:=stageController.WaitForProcess(SELE)
	eventData:=new(EtlEventData)
	eventData.proceEvdata.pipelineId=pipelineId
	eventData.proceEvdata.processId=processId
	eventData.proceEvdata.startTime=time.Now()
	//Long nid = ArbitrateConfigUtils.getCurrentNid();
	//eventData.setCurrNid(nid);
	//eventData.setNextNid(nid);
	return eventData
}

func (self *SelectMemoryArbitrateEvent) Single(data *EtlEventData){

}