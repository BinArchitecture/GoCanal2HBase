package gobase

type TransformMemoryArbitrateEvent struct {
	Gltask *GlobalTask
}

func NewTransformMemoryArbitrateEvent(gltask *GlobalTask) *TransformMemoryArbitrateEvent{
	event:=new(TransformMemoryArbitrateEvent)
	event.Gltask=gltask
	return event
}


func (self *TransformMemoryArbitrateEvent) Await(pipelineId int64) *EtlEventData{
	return nil
}

func (self *TransformMemoryArbitrateEvent) Single(data *EtlEventData){

}
