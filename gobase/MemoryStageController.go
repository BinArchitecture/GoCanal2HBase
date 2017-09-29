package gobase

import (
	"sync"
	"sort"
)

type MemoryStageController struct {
	atomicMaxProcessId int64
	replys *sync.Map //Map<StageType, ReplyProcessQueue>
	progress *sync.Map //Map<Long, StageProgress>
	termins chan *TerminEventData
	nullProgress *StageProgress
	pipelineId int64
	stop bool
	initLock sync.RWMutex
}

func NewMemoryStageController(pipelineId int64) *MemoryStageController{
	cc:=new(MemoryStageController)
	cc.pipelineId=pipelineId
	cc.stop=false
	cc.replys=new(sync.Map)
	cc.progress=new(sync.Map)
	cc.termins=make(chan *TerminEventData,20)
	cc.nullProgress=new(StageProgress)
	return cc
}

func (self *MemoryStageController) WaitForProcess(stageType string) int64{
	var tmpQueue interface{}
	if stageType==SELE {
		var ok bool
		tmpQueue,ok=self.replys.Load(stageType)
		if !ok{
			self.initSelect()
		}
	}
	queue:=tmpQueue.(chan int64)
	processId:=<-queue
	if stageType==SELE{
		self.progress.Store(processId,self.nullProgress)
	}
	return processId
}

func (self *MemoryStageController) GetLastData(processId int64) *EtlEventData{
	tmp,ok:=self.progress.Load(processId)
	if ok{
		data:=tmp.(*StageProgress)
		return data.data
	}
	return nil
}

func (self *MemoryStageController) Destory(){
	self.initLock.Lock()
	defer self.initLock.Unlock()
	self.replys.Range(func(key, value interface{}) bool {
		self.replys.Delete(key)
		return true
	})
	self.progress.Range(func(key, value interface{}) bool {
		self.progress.Delete(key)
		return true
	})
}

func (self *MemoryStageController) ClearProgress(processId int64){
	self.progress.Delete(processId)
}

func (self *MemoryStageController) Termin(terminType string){
	self.progress.Range(func(key, value interface{}) bool {
		processId:=key.(int64)
		da:=value.(*StageProgress)
		eventData:=da.data
		data:=new(TerminEventData)
		data.proceEvdata.processId=processId
		data.proceEvdata.pipelineId=self.pipelineId
		data.terminType=terminType
		data.code="channel"
		data.desc=terminType
		if eventData!=nil{
			data.proceEvdata.batchId=eventData.proceEvdata.batchId
			data.currNid=int64(eventData.currNid)
			data.proceEvdata.startTime=eventData.proceEvdata.startTime
			data.proceEvdata.endTime=eventData.proceEvdata.endTime
			data.proceEvdata.firstTime=eventData.proceEvdata.firstTime
			data.proceEvdata.number=eventData.proceEvdata.number
			data.proceEvdata.size=eventData.proceEvdata.size
		}
		self.termins<-data
		self.progress.Delete(processId)
		return true
	})
	self.initSelect()
}

func (self *MemoryStageController) Single(stageType string,etlEventData *EtlEventData) bool{
	var result bool
	switch stageType {
	case SELE:
		result=handleSingle(self, etlEventData, stageType,EXTR)
	case EXTR:
		result=handleSingle(self, etlEventData, stageType,TRNS)
	case TRNS:
		result=handleSingle(self, etlEventData, stageType,LOAD)
	case LOAD:
		self.progress.Delete(etlEventData.proceEvdata.processId)
		//Object removed = progress.remove(etlEventData.getProcessId());
		//// 并不是立即触发，通知下一个最小的一个process启动
		//computeNextLoad();
		//// 一个process完成了，自动添加下一个process
		//if (removed != null) {
		//	replys.get(StageType.SELECT).offer(atomicMaxProcessId.incrementAndGet());
		//	result = true;
		//}
	}
	return result
}

func handleSingle(self *MemoryStageController, etlEventData *EtlEventData, stageType string,nextStageType string) (result bool) {
	result=false
	_, ok := self.progress.Load(etlEventData.proceEvdata.processId)
	if ok {
		stageProgress := new(StageProgress)
		stageProgress.data = etlEventData
		stageProgress.stageType = stageType
		self.progress.Store(etlEventData.proceEvdata.processId, stageProgress)
		if nextStageType!=LOAD{
			tmq, ok := self.replys.Load(nextStageType)
			if ok {
				extr := tmq.(chan int64)
				extr <- etlEventData.proceEvdata.processId
				result = true
			}
		}else{
			result = true
			proArray:=make([]int,1)
			self.progress.Range(func(key, value interface{}) bool {
				processId:=key.(int64)
				proArray=append(proArray,int(processId))
				return true
			})
			sort.Ints(proArray)
			if len(proArray)>0{
				minProid:=proArray[0]
				tmq, ok := self.replys.Load(nextStageType)
				if ok {
					extr := tmq.(chan int64)
					extr <- int64(minProid)
				}
			}
		}
	}
	return result
}

func (self *MemoryStageController) initSelect(){
	self.initLock.Lock()
	defer self.initLock.Unlock()
	tmpQue,ok:=self.replys.Load(SELE)
	if ok{
		queue:=tmpQue.(chan int64)
		queue<-0
	}
}

type StageProgress struct {
	stageType string
	data *EtlEventData
}