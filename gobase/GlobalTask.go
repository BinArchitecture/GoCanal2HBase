package gobase

const (//StageType
	SELE   = "SELECT"
	EXTR  = "EXTRACT"
	TRNS = "TRANSFORM"
	LOAD      = "LOAD"
)

type GlobalTask struct {
	ArbitrateEventService *ArbitrateEventService
	RowDataPipeDelegate *RowDataPipeDelegate
	StageAggregationCollector *StageAggregationCollector
	ConfigClientService *ConfigClientService
	ArbitrateFactory *ArbitrateFactory
}

func NewGlobalTask() *GlobalTask{
	task:=new(GlobalTask)
	arbitrateService:=NewArbitrateEventService(task)
	task.ArbitrateEventService=arbitrateService
	rowDelegate:=NewRowDataPipeDelegate()
	task.RowDataPipeDelegate=rowDelegate
	stageCollector:=NewStageAggregationCollector()
	task.StageAggregationCollector=stageCollector
	configService:=NewConfigClientService()
	task.ConfigClientService=configService
	task.ArbitrateFactory=NewArbitrateFactory()
	return task
}

func (self *GlobalTask) Shutdown(){
}

