package _select

import (
	"time"
	"github.com/BinArchitecture/GoCanal2HBase/gobase"
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
	base *gobase.GlobalTask
}

func NewSelectTask(base *gobase.GlobalTask,addr string,uname string,passwd string,destination string) *SelectTask{
	task:=new(SelectTask)
	task.base=base
	task.canalconn=NewCanalConnector(addr,uname,passwd,destination)
	task.needCheck=false
	task.isStart=false
	task.rversion=0
	task.lastResetTime=time.Now()
	task.batchBuffer=make(chan *BatchTermin)
	task.canStartSelector=make(chan bool)
	return task
}

func (self *SelectTask) Run(){
	for self.base.Running{

	}
}
type BatchTermin struct{
	batchId int64
	processId int64
	needWait bool
}
