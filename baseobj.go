// baseobj
package gpss

import (
	"fmt"
	"sync"
)

type IBaseObj interface {
	GetName() string                    // Get object name
	SetDst([]IBaseObj)                  // Set dst for object
	GetDst() []IBaseObj                 // Get dst for object
	SetPipeline(pipe IPipeline)         // Set pipeline for object
	GetPipeline() IPipeline             // Get pipeline for object
	AppendTransact(ITransaction) bool   // Append transact to object
	HandleTransacts(wg *sync.WaitGroup) // Handle all transacts of object
	PrintReport()                       // Print report
}

type BaseObj struct {
	name    string
	objTime int
	dst     []IBaseObj
	pipe    IPipeline
	tb      ITransactTable
}

func NewBaseObj(name string) *BaseObj {
	obj := &BaseObj{}
	obj.Init(name)
	return obj
}

func (obj *BaseObj) Init(name string) {
	obj.name = name
	obj.tb = NewTransactTable()
}

func (obj *BaseObj) GetName() string {
	return obj.name
}

func (obj *BaseObj) SetDst(dst []IBaseObj) {
	obj.dst = dst
}

func (obj *BaseObj) GetDst() []IBaseObj {
	return obj.dst
}

func (obj *BaseObj) SetPipeline(pipe IPipeline) {
	obj.pipe = pipe
}

func (obj *BaseObj) GetPipeline() IPipeline {
	return obj.pipe
}

func (obj *BaseObj) GetTransactTable() ITransactTable {
	return obj.tb
}

func (obj *BaseObj) AppendTransact(t ITransaction) bool {
	obj.tb.Push(t)
	return true
}

func (obj *BaseObj) PrintReport() {
	fmt.Println("Object name \"", obj.name, "\"")
}
