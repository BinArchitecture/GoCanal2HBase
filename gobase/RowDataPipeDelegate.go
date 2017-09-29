package gobase

import "sync"

const(//PipeDataType
	/** 数据库 */
	DB_BATCH="db_batch"
	/** 附件记录 */
	FILE_BATCH="file_batch"
	/** mq记录 */
	MQ_BATCH="mq_batch"
	/** cache记录 */
	CACHE_BATCH="cache_batch"
)
type RowDataPipeDelegate struct {
	rowDataMemoryPipe *RowDataMemoryPipe
	RowDataRpcPipe *RowDataRpcPipe
}

type MemoryPipeKey struct {
	identity *Identity
	time int64
	dataType string //PipeDataType
}

type Identity struct {
	channelId int64
	pipelineId int64
	processId int64
}

type RowDataMemoryPipe struct {
	cache *sync.Map //Map<MemoryPipeKey, DbBatch>
	timeout int
	retry int
}

func NewRowDataMemoryPipe() *RowDataMemoryPipe{
	pipe:=new(RowDataMemoryPipe)
	pipe.cache=new(sync.Map)
	pipe.timeout=60000
	pipe.retry=3
	return pipe
}

func (self *RowDataMemoryPipe) put(batch RowBatch) *MemoryPipeKey {
	key := new(MemoryPipeKey)
	key.identity = batch.identity
	key.dataType = DB_BATCH
	self.cache.Store(key, batch)
	return key
}

type RowDataRpcPipe struct {

}


