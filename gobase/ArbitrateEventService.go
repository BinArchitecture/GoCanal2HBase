package gobase

import "time"

type ProcessEventData struct {
	pipelineId int64
	processId int64
	startTime time.Time
	endTime time.Time
	firstTime time.Time
	batchId int64
	number int64
	size int64
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

type TerminEventData struct{
	proceEvdata *ProcessEventData
	terminType string
	code string
	desc string
	currNid int64
}

type EtlEventData struct {
	proceEvdata *ProcessEventData
	currNid int
	nextNid int
	desc interface{}
}

type ArbitrateEvent interface {
	Await(pipelineId int64) *EtlEventData
	Single(data *EtlEventData)
}

type ArbitrateEventService interface {
	//MainStemEvent() MainStemArbitrateEvent
	SelectEvent() ArbitrateEvent
	ExtractEvent() ArbitrateEvent
	TransformEvent() ArbitrateEvent
	LoadEvent() ArbitrateEvent
	TerminEvent() ArbitrateEvent
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

