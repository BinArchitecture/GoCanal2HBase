package gobase

import "sync/atomic"

type AggregationItem struct{
	startTime int64
	size int64
	endTime int64
	number int64
}

type StageAggregation struct{
	bufferSize int
	indexMask int
	table []*AggregationItem
	sequence int64
}

func NewStageAggregation(buffSize int) *StageAggregation{
	agg:=new(StageAggregation)
	agg.bufferSize=buffSize
	agg.indexMask=agg.bufferSize-1
	agg.sequence=-1
	agg.table=make([]*AggregationItem,buffSize)
	return agg
}

func (self *StageAggregation) getIndex(sequcnce int64) int{
	return int(self.sequence) & self.indexMask
}

func (self *StageAggregation) Push(aggregation *AggregationItem){
	seq:=atomic.AddInt64(&self.sequence,1)
	self.table[self.getIndex(seq)]=aggregation
}

func (self *StageAggregation) Count() int64{
	return self.sequence
}
//平均处理时间
func (self *StageAggregation) Histogram() string{
	return ""
}

//public String histogram() {
//Long costs = 0L;
//Long items = 0L;
//Long max = 0L;
//Long min = Long.MAX_VALUE;
//Long tps = 0L;
//Long tpm = 0L;
//Long avg = 0L;
//
//Long lastTime = 0L;
//Long tpsCount = 0L;// 记录每秒的请求数，临时变量
//Long tpsTotal = 0L;// 总tps数多少
//Long tpsSecond = 0L;// 多少秒中有数据
//
//Long tpmCount = 0L; // 记录每分钟的请求数，临时变量
//Long tpmTotal = 0L; // 总tps数多少
//Long tpmMinute = 0L;// 多少分钟有数据
//for (int i = 0; i < table.length; i++) {
//AggregationItem aggregation = table[i];
//if (aggregation != null) {
//Long cost = aggregation.getEndTime() - aggregation.getStartTime();
//items += 1;
//costs += cost;
//if (cost > max) {
//max = cost;
//}
//if (cost < min) {
//min = cost;
//}
//
//if (lastTime != 0) {
//if (lastTime > aggregation.getEndTime() - ONE_SECOND) {// 说明在同一秒
//tpsCount++;
//} else {
//tpsTotal += tpsCount;
//tpsSecond++;
//tpsCount = 0L;
//}
//
//if (lastTime > aggregation.getEndTime() - ONE_MINUTE) {// 说明在同一分钟
//tpmCount++;
//} else {
//tpmTotal += tpmCount;
//tpmMinute++;
//tpmCount = 0L;
//}
//
//}
//
//lastTime = aggregation.getEndTime();
//}
//}
//// 设置一下最后一批tps/m统计信息
//tpsTotal += tpsCount;
//tpsSecond++;
//tpsCount = 0L;
//tpmTotal += tpmCount;
//tpmMinute++;
//tpmCount = 0L;
//
//if (items != 0) {
//avg = costs / items;
//}
//
//if (tpsSecond != 0) {
//tps = tpsTotal / tpsSecond;
//}
//
//if (tpmMinute != 0) {
//tpm = tpmTotal / tpmMinute;
//}
//
//if (min == Long.MAX_VALUE) {
//min = 0L;
//}
//
//return String.format(HISTOGRAM_FORMAT, new Object[] { sequence.get() + 1, items, max, min, avg, tps, tpm });
//}