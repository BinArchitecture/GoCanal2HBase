package gobase

import "time"

const(
	CRE="CREATE"
	DEL="DELETE"
	NodeStart="NodeStart"
	NodeStop="NodeStop"
)

type NodeTask struct {
	pipeline *Pipeline
	stageType []string
	taskEvent []string
	shutdown bool
}

func NewNodeTask(pipeline *Pipeline) *NodeTask{
	nodetask:=new(NodeTask)
	nodetask.pipeline=pipeline
	nodetask.stageType=[]string{SELE,EXTR,TRNS,LOAD}
	nodetask.taskEvent=[]string{CRE,CRE,CRE,CRE}
	nodetask.shutdown=false
	return nodetask
}

type Node struct {
	Id int64
	Name string
	Ip string
	Port int
	NodeStatus string
	Description string
	GmtCreate time.Time
	GmtModified time.Time
	//private NodeParameter     parameters       = new NodeParameter(); // node对应参数信
}

func NewNode(id int64,name string,ip string,port int,desc string) *Node{
	node:=new(Node)
	node.Id=id
	node.Name=name
	node.Ip=ip
	node.Port=port
	node.Description=desc
	node.NodeStatus=NodeStart
	node.GmtCreate=time.Now()
	node.GmtModified=time.Now()
	return node
}