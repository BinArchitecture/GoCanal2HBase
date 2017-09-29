package gobase

import "time"

type DataMediaPair struct{
	id int64
	pipelineId int64
	source *DataMedia
	target *DataMedia
	pullWeight int64
	pushWeight int64
	columnPairMode string
	columnPairs []*ColumnPair
	columnGroups []*ColumnGroup
	gmtCreate time.Time
	gmtModified  time.Time
}

type ColumnGroup struct {
	id int64
	columnPairs []*ColumnPair
	dataMediaPairId int64
	gmtCreate time.Time
	gmtModified  time.Time
}

type ColumnPair struct {
	id int64
	dataMediaPairId int64
	gmtCreate time.Time
	gmtModified  time.Time
	sourceColumn string
	targetColumn string
}


type DataMedia struct {
	id int64
	name string
	namespace string
	source *DataMediaSource
	encode string
	gmtCreate time.Time
	gmtModified  time.Time
}

const(
	DataMediaType_MYSQL="MYSQL"
	DataMediaType_ORACLE="ORACLE"
	columnPairMode_INCLUDE="CINCLUDE"
	columnPairMode_EXCLUDE="CEXCLUDE"
)

type DataMediaSource struct {
	id int64
	name string
	DataMediaType   string
	encode string
	gmtCreate time.Time
	gmtModified  time.Time
}