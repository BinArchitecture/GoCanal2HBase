package gobase

import "time"

type ProcessEventData struct {
	pipelineId int64
	processId int64
	startTime time.Time
	endTime time.Time
	firstTime time.Time
	batchId int64
	number int64
	size int64
}

type EtlEventData struct {
	proceEvdata *ProcessEventData
	currNid int
	nextNid int
	desc interface{}
}

type ArbitrateEvent interface {
	await(pipelineId int64) *EtlEventData
	single(data *EtlEventData)
}

type ArbitrateEventService interface {
	selectEvent() ArbitrateEvent
	extractEvent() ArbitrateEvent
	transformEvent() ArbitrateEvent
	loadEvent() ArbitrateEvent
	terminEvent() ArbitrateEvent
}
