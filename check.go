// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

type Check struct {
	BaseObj
	HandleChecking HandleCheckingFunc
	falseObj       IBaseObj
	parameters     []Parameter
}

type HandleCheckingFunc func(obj *Check, transact ITransaction) bool

func Checking(obj *Check, transact ITransaction) bool {
	res := true
	for _, v := range obj.parameters {
		parameter := transact.GetParameterByName(v.name)
		if parameter != v.value {
			res = bool(res && false)
		}
	}
	return res
}

func NewCheck(name string, hndl HandleCheckingFunc, falseObj IBaseObj, parameters ...Parameter) *Check {
	obj := &Check{parameters: parameters, falseObj: falseObj}
	if hndl != nil {
		obj.HandleChecking = hndl
	} else {
		obj.HandleChecking = Checking
	}
	return obj
}

func (obj *Check) AppendTransact(transact ITransaction) bool {
	transact.PrintInfo()
	if !obj.HandleChecking(obj, transact) {
		if obj.falseObj != nil {
			if obj.falseObj.AppendTransact(transact) {
				return true
			} else {
				return false
			}
		}
	}
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			obj.GetLogger().GetTrace().Println("Append transact ", transact.GetId(), " to Check")
			return true
		}
	}
	return false
}

func (obj *Check) PrintReport() {
	return
}
