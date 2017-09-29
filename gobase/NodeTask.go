package gobase

import "time"

const(
	CRE="CREATE"
	DEL="DELETE"
	NodeStart="NodeStart"
	NodeStop="NodeStop"
)

type NodeTask struct {
	pipeline Pipeline
	stageType []string
	taskEvent []string
	shutdown bool
}

type Pipeline struct{
	id int64
	channelId int64
	name string
	description string
	selectNodes []*Node
	extractNodes []*Node
	loadNodes []*Node
	pairs []*DataMediaPair
	gmtCreate time.Time
	gmtModified time.Time
}

type Node struct {
	id int64
	name string
	ip string
	port int64
	nodeStatus string
	description string
	gmtCreate time.Time
	gmtModified time.Time
	//private NodeParameter     parameters       = new NodeParameter(); // node对应参数信
}
