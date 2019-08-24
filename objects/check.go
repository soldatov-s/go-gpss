// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
)

// Check compares parameters of Transaction or any another parameters of sumulation
// model, and controls the destination of the Active Transaction based on the
// result of the comparison.
type Check struct {
	BaseObj
	// Function for checking
	HandleChecking HandleCheckingFunc
	// Destination object in case false result of checking
	falseObj IBaseObj
	// Parameters of transact for checking
	parameters []Parameter
	// For counting true result checking
	cntTrue int
	// For counting false result checking
	cntFalse int
}

// HandleCheckingFunc is a checking function signature
type HandleCheckingFunc func(obj *Check, transact *Transaction) bool

// Checking is default function for checking
func Checking(obj *Check, transact *Transaction) bool {
	for _, v := range obj.parameters {
		parameter := transact.GetParameter(v.Name)
		if parameter != v.Value {
			return false
		}
	}
	return true
}

// NewCheck creates new Check.
// name - name of object; hndl - function for checking; falseObj - destination
// of the Active Transaction in case false result of checking; parameters -
// parameters for checking
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

// AppendTransact append transact to object
func (obj *Check) AppendTransact(transact *Transaction) bool {
	obj.BaseObj.AppendTransact(transact)
	if !obj.HandleChecking(obj, transact) {
		obj.cntFalse++
		if obj.falseObj != nil {
			if obj.falseObj.AppendTransact(transact) {
				return true
			}
		}
		return false
	}
	obj.cntTrue++
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			return true
		}
	}
	return false
}

// Report - print report about object
func (obj *Check) Report() {
	obj.BaseObj.Report()
	fmt.Printf("Check result true %d\tCheck result false %d\n\n", obj.cntTrue, obj.cntFalse)
	return
}
