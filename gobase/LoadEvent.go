package gobase


type LoadMemoryArbitrateEvent struct {
	Gltask *GlobalTask
}
func NewLoadMemoryArbitrateEvent(gltask *GlobalTask) *LoadMemoryArbitrateEvent{
	event:=new(LoadMemoryArbitrateEvent)
	event.Gltask=gltask
	return event
}

func (self *LoadMemoryArbitrateEvent) Await(pipelineId int64) *EtlEventData {
	return nil
}

func (self *LoadMemoryArbitrateEvent) Single(data *EtlEventData) {

}
