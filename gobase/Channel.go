package gobase

import (
	"time"
	"sync"
)

const(
	ChannelStatus_START="ChannelStatus_START"
	ChannelStatus_PAUSE="ChannelStatus_PAUSE"
	ChannelStatus_STOP="ChannelStatus_STOP"
)

type Channel struct {
	Id int64
	Name string
	Status string
	Description string
	Pipelines *sync.Map
	GmtCreate time.Time
	GmtModified time.Time
}

type Pipeline struct{
	Id int64
	ChannelId int64
	Name string
	Description string
	SelectNodes []*Node
	ExtractNodes []*Node
	LoadNodes []*Node
	Pairs []*DataMediaPair
	GmtCreate time.Time
	GmtModified time.Time
	Parameters *PipelineParameter
}

const(
	RoundRbin="RoundRbin"
	/** 随机 */
	Random="Random"
	/** Stick */
	Stick="Stick"
	Canal="Canal"
	Eromanga="Eromanga"
)
type PipelineParameter struct {
	pipelineId int64
	parallelism int
	lbAlgorithm string
	home bool
	canalAddr string
	canalUserName string
	canalPasswd string
	canalFilter string
	selectorMode string
	destinationName string
	mainstemClientId int
	mainstemBatchsize int32
	batchTimeout int64
	extractPoolSize int
	loadPoolSize int
	channelInfo string
	dryRun bool
}

func NewPipeParameter(pipelineId int64,parallelism int,canalAddr string,canalUserName string,canalPasswd string,canalFilter string,destinationName string,mainstemClientId int,mainstemBatchsize int32,batchTimeout int64,extractPoolSize int,loadPoolSize int,channelInfo string,dryRun bool) *PipelineParameter{
	param:=new(PipelineParameter)
	param.pipelineId=pipelineId
	param.parallelism=parallelism
	param.lbAlgorithm=Random
	param.home=false
	param.canalAddr=canalAddr
	param.canalUserName=canalUserName
	param.canalPasswd=canalPasswd
	param.canalFilter=canalFilter
	param.selectorMode=Canal
	param.destinationName=destinationName
	param.mainstemClientId=mainstemClientId
	param.mainstemBatchsize=mainstemBatchsize
	param.batchTimeout=batchTimeout
	param.extractPoolSize=extractPoolSize
	param.loadPoolSize=loadPoolSize
	param.channelInfo=channelInfo
	param.dryRun=dryRun
	return param
}