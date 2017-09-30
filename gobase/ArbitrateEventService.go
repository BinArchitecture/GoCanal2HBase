package gobase

import (
	"time"
	"sync"
)

type ProcessEventData struct {
	PipelineId int64
	ProcessId  int64
	StartTime  time.Time
	EndTime    time.Time
	FirstTime  time.Time
	BatchId    int64
	Number     int64
	Size       int64
}

const(//TerminType
	/** 正常结束 */
	NORMAL="NORMAL"
	/** 警告信息 */
	WARNING="WARNING"
	/** 回滚对应同步 */
	ROLLBACK="ROLLBACK"
	/** 重新开始同步 */
	RESTART="RESTART"
	/** 关闭同步 */
	SHUTDOWN="SHUTDOWN"
)

type EtlEventData struct {
	ProceEvdata *ProcessEventData
	CurrNid     int
	NextNid     int
	Desc        interface{}
}

type ArbitrateEvent interface {
	Await(pipelineId int64) *EtlEventData
	Single(data *EtlEventData)
}

type ArbitrateFactory struct {
	cache *sync.Map  // 两层的Map接口，第一层为pipelineId，第二层为具体的资源类型class
}

func NewArbitrateFactory() *ArbitrateFactory{
	arbitrateFactory:=new(ArbitrateFactory)
	arbitrateFactory.cache=new(sync.Map)
	return arbitrateFactory
}

func (self *ArbitrateFactory) GetInstance(pipelineId int64) *MemoryStageController{
	con,ok:=self.cache.Load(pipelineId)
	if ok{
		return con.(*MemoryStageController)
	}
	return nil
}

func (self *ArbitrateFactory) Put(pipelineId int64,controller *MemoryStageController) {
	self.cache.Store(pipelineId,controller)
}

type ArbitrateEventService struct {
	//MainStemEvent() MainStemArbitrateEvent
	SelectEvent ArbitrateEvent
	ExtractEvent ArbitrateEvent
	TransformEvent ArbitrateEvent
	LoadEvent ArbitrateEvent
	TerminEvent ArbitrateEvent
}

func NewArbitrateEventService(gltask *GlobalTask) *ArbitrateEventService {
	service:=new(ArbitrateEventService)
	service.SelectEvent=NewSelectMemoryArbitrateEvent(gltask)
	service.ExtractEvent=NewExtractMemoryArbitrateEvent(gltask)
	service.TransformEvent=NewTransformMemoryArbitrateEvent(gltask)
	service.LoadEvent=NewLoadMemoryArbitrateEvent(gltask)
	service.TerminEvent=NewTerminMemoryArbitrateEvent(gltask)
	return service
}


type MainStemArbitrateEvent struct{
}

const(//Status
	OVERTAKE="OVERTAKE"
	TAKEING="TAKEING"
)

type MainStemEventData struct {
	nid int64
	status string //Status
	active bool
	pipelineId int64
}

//func (self *MainStemArbitrateEvent) Await(pipeId int64) error{
//
//}
//
//func (self *MainStemArbitrateEvent) Check(pipeId int64) bool{
//
//}
//
//func (self *MainStemArbitrateEvent) Release(pipeId int64) bool{
//
//}
//
//func (self *MainStemArbitrateEvent) Single(data *MainStemEventData) bool{
//
//}

