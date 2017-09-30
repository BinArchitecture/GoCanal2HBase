package gobase

import (
	"sync"
)

type ConfigClientService struct{
	channelCache *sync.Map //RefreshMemoryMirror<Long, Channel>
	channelMapping *sync.Map //Map<Long, Long>
}

func NewConfigClientService() *ConfigClientService{
	ccs:=new(ConfigClientService)
	ccs.channelCache=new(sync.Map)
	ccs.channelMapping=new(sync.Map)
	return ccs
}

func (self *ConfigClientService) PutChannel(channelId int64,channel *Channel) {
	if channel!=nil{
		self.channelCache.Store(channelId,channel)
		if channel.Pipelines!=nil{
			channel.Pipelines.Range(func(key, value interface{}) bool {
				pipeId:=key.(int64)
				self.channelMapping.Store(pipeId,channelId)
				return true
			})
		}
	}
}

func (self *ConfigClientService) FindChannel(channelId int64) *Channel{
	cc,ok:=self.channelCache.Load(channelId)
	if ok{
		chanel:=cc.(*Channel)
		return chanel
	}
	return nil
}

func (self *ConfigClientService) FindChannelByPipelineId(pipelineId int64) *Channel{
	pp,ok:=self.channelMapping.Load(pipelineId)
	if ok{
		chanelId:=pp.(int64)
		return self.FindChannel(chanelId)
	}
	return nil
}

func (self *ConfigClientService) FindPipeline(pipelineId int64) *Pipeline{
	chanel:=self.FindChannelByPipelineId(pipelineId)
	pp,ok:=chanel.Pipelines.Load(pipelineId)
	if ok{
		return pp.(*Pipeline)
	}
	return nil
}

func (self *ConfigClientService) BuildCanal(pipe *Pipeline) *CanalConnector{
	canalconn:=NewCanalConnector(pipe.Parameters.canalAddr,pipe.Parameters.canalUserName,pipe.Parameters.canalPasswd,pipe.Parameters.destinationName)
	return canalconn
}

