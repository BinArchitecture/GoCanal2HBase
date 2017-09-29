package gobase

type Setl interface {
	fetch() []interface{}
	handle() error
}

type MemoryStageController struct {

}