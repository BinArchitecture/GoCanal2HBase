package gobase

import (
	"time"
	"fmt"
	"github.com/golang/glog"
)

type SelectTask struct {
	isStart bool
	//statisticsClientService StatisticsClientService
	canalconn *CanalConnector
	batchBuffer chan *BatchTermin
	canStartSelector chan bool
	needCheck bool
	rversion int
	lastResetTime time.Time
	pipelineId int64
	pipeline *Pipeline
	Running bool
	base *GlobalTask
}

func NewSelectTask(base *GlobalTask,pipelineId int64) *SelectTask{
	task:=new(SelectTask)
	task.base=base
	pipe:=task.base.ConfigClientService.FindPipeline(pipelineId)
	task.pipeline=pipe
	task.canalconn=task.base.ConfigClientService.BuildCanal(pipe)
	task.needCheck=false
	task.isStart=false
	task.rversion=0
	task.lastResetTime=time.Now()
	task.batchBuffer=make(chan *BatchTermin)
	task.canStartSelector=make(chan bool)
	task.Running=true
	task.pipelineId=pipelineId
	return task
}

func (self *SelectTask) StartUp(){
	fmt.Printf("selectTask begin")
	pipe:=self.base.ConfigClientService.FindPipeline(self.pipelineId)
	if pipe==nil{
		glog.Error("pipe is null cannot start select task")
		return
	}
	conn:=NewCanalConnector(pipe.Parameters.canalAddr,pipe.Parameters.canalUserName,pipe.Parameters.canalPasswd,pipe.Parameters.destinationName)
	conn.Connect()
	conn.Subscribe(pipe.Parameters.canalFilter)
	for self.Running{
		//if (canStartSelector.state() == false) { // 是否出现异常
		//	// 回滚在出现异常的瞬间，拿出来的数据，因为otterSelector.selector()会循环，可能出现了rollback，其还未感知到
		//	rollback(gotMessage.getId());
		//	continue;
		//}
		msg,err:=conn.GetWithoutAck(1024)

		if err!=nil{
			glog.Error(err)
			return
		}
		if len(msg.entries)==0{
			self.batchBuffer<-&BatchTermin{msg.id,-1,false}
			continue
		}
		//etlEventData:=self.base.ArbitrateEventService.SelectEvent.Await(pipe.Id)
		//if (rversion.get() != startVersion) {// 说明存在过变化，中间出现过rollback，需要丢弃该数据
		//	logger.warn("rollback happend , should skip this data and get new message.");
		//	canStartSelector.get();// 确认一下rollback是否完成
		//	gotMessage = otterSelector.selector();// 这时不管有没有数据，都需要执行一次s/e/t/l
		//}

	}
}

func (self *SelectTask) Shutdown(){
	self.Running=false
}

type BatchTermin struct{
	batchId int64
	processId int64
	needWait bool
}
