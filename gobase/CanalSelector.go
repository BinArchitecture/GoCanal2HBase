package gobase

import (
	"github.com/golang/glog"
	"time"
	canal_entry "github.com/BinArchitecture/GoCanal2HBase/com_alibaba_otter_canal_protocol/entry"
)

type CanalSelector struct {
	canalConnector *CanalConnector
	pipeline *Pipeline
	running bool
}

type OtterMessage struct{
	Id int64
	Datas []*EtlEventData
}

func (self *CanalSelector) Connect(){
	self.canalConnector.Connect()
}

func (self *CanalSelector) Subscribe(){
	filter:=self.pipeline.Parameters.canalFilter
	self.canalConnector.Subscribe(filter)
}

func (self *CanalSelector) Selector() *OtterMessage{
	batchSize:=self.pipeline.Parameters.mainstemBatchsize
	batchTimeout:=self.pipeline.Parameters.batchTimeout
	var msg *Message
	var err error
	if batchTimeout<0{
		for self.running{
			msg,err=self.canalConnector.GetWithoutAck(int(batchSize))
			if err!=nil{
				glog.Error(err)
				return nil
			}
			if msg==nil||msg.id==-1{
				time.Sleep(time.Second)
				glog.Infoln("no data from canal now")
			}else{
				break
			}
		}
	}else{
		//canaltimeout
	}
	evenDataList:=parse(self.pipeline,msg.entries)
	otterMsg:=new(OtterMessage)
	otterMsg.Id=msg.id
	otterMsg.Datas=evenDataList
	return otterMsg
}

//此处省去回环数据过滤操作
func parse(pipeline *Pipeline,entrys []*canal_entry.Entry) []*EtlEventData{
	for _,entry:=range entrys{
		switch *entry.EntryType {
		case canal_entry.EntryType_TRANSACTIONBEGIN:
		case canal_entry.EntryType_ROWDATA:
		case canal_entry.EntryType_TRANSACTIONEND:
		}

	}
	return nil
}