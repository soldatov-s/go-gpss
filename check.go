// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
)

// Check is object for check parameters of transact or any another parameters
// sumulation model. Check maybe have two results: true ir false. You can replace
// Checking Handler with your handler.
type Check struct {
	BaseObj
	// Checking handler
	HandleChecking HandleCheckingFunc
	// Destination object in case false result of checking
	falseObj IBaseObj
	// Parameters of transact for checking
	parameters []Parameter
	// For counting true result checking
	cnt_true int
	// For counting false result checking
	cnt_false int
}

type HandleCheckingFunc func(obj *Check, transact ITransaction) bool

func Checking(obj *Check, transact ITransaction) bool {
	res := true
	for _, v := range obj.parameters {
		parameter := transact.GetParameterByName(v.Name)
		if parameter != v.Value {
			res = bool(res && false)
		}
	}
	return res
}

func NewCheck(name string, hndl HandleCheckingFunc, falseObj IBaseObj, parameters ...Parameter) *Check {
	obj := &Check{parameters: parameters, falseObj: falseObj}
	obj.name = name
	if hndl != nil {
		obj.HandleChecking = hndl
	} else {
		obj.HandleChecking = Checking
	}
	return obj
}

func (obj *Check) AppendTransact(transact ITransaction) bool {
	transact.PrintInfo()
	obj.GetLogger().GetTrace().Println("Append transact ", transact.GetId(), " to Check")
	if !obj.HandleChecking(obj, transact) {
		obj.cnt_false++
		if obj.falseObj != nil {
			if obj.falseObj.AppendTransact(transact) {
				return true
			} else {
				return false
			}
		}
		return false
	}
	obj.cnt_true++
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			return true
		}
	}
	return false
}

func (obj *Check) PrintReport() {
	obj.BaseObj.PrintReport()
	fmt.Printf("Check result true %d\tCheck result false %d\n\n", obj.cnt_true, obj.cnt_false)
	return
}
