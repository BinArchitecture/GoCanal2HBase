package main

import (
	"github.com/BinArchitecture/GoCanal2HBase/gobase"
	"time"
	"sync"
)

func main() {
	channel:=new(gobase.Channel)
	channel.Id=1
	channel.Name="gomysql2hbase"
	channel.Status=gobase.ChannelStatus_START
	channel.Description="go2hbase"
	channel.GmtCreate=time.Now()
	channel.GmtModified=time.Now()
	channel.Pipelines=new(sync.Map)
	pipeline:=new(gobase.Pipeline)
	pipeline.Id=1
	pipeline.Name="gopipe"
	pipeline.ChannelId=1
	pipeline.Description="gopipedesc"
	pipeline.GmtModified=time.Now()
	pipeline.GmtCreate=time.Now()
	node:=gobase.NewNode(1,"seltnode","10.6.30.109",11111,"seltnode")
	pipeline.SelectNodes=[]*gobase.Node{node}
	pipeline.ExtractNodes=[]*gobase.Node{node}
	pipeline.LoadNodes=[]*gobase.Node{node}
	pipeline.Pairs=make([]*gobase.DataMediaPair,1)
	pipeline.Parameters=gobase.NewPipeParameter(pipeline.Id,3,"10.6.30.109:11111","","","","example",1,1024,-1,5,5,"channelFuck",false)
	channel.Pipelines.Store(pipeline.Id,pipeline)

	nodeTask:=gobase.NewNodeTask(pipeline)
	otterController:=gobase.NewOtterController()
	otterController.Gltask.ConfigClientService.PutChannel(channel.Id,channel)
	otterController.PutPipe2Controller(pipeline.Id)
	otterController.StartPipeline(nodeTask)
}
