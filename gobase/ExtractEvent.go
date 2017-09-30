package gobase

type ExtractMemoryArbitrateEvent struct {
	Gltask *GlobalTask
}

func NewExtractMemoryArbitrateEvent(gltask *GlobalTask) *ExtractMemoryArbitrateEvent{
	event:=new(ExtractMemoryArbitrateEvent)
	event.Gltask=gltask
	return event
}

func (self *ExtractMemoryArbitrateEvent) Await(pipelineId int64) *EtlEventData{
	return nil
}

func (self *ExtractMemoryArbitrateEvent) Single(data *EtlEventData){

}
