package gobase

const(//EventType
	Inser=0
	Delet=1
	Updat=2
)

type RowBatch struct {
	datas *EventData
	identity *Identity
}

type EventData struct {
	tableId int64
	tableName string
	schemaName string
	eventType int //EventType
	//变更数据的业务时间.
	executeTime int64
	//变更前的主键值,如果是insert/delete变更前和变更后的主键值是一样的.
	oldKeys []*EventColumn
	// 变更后的主键值,如果是insert/delete变更前和变更后的主键值是一样的.
	keys []*EventColumn
	//非主键的其他字段
	columns []*EventColumn
	//预计的size大小，基于binlog event的推算 default 1024
	size int
	//同步映射关系的id
	pairId int
	sql string

}

type EventColumn struct{
	index int
	columnType int
	columnName string
	columnValue string
	isNull bool
	isKey bool
	isUpdate bool
}