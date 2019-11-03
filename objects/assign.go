// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package objects

// Assign - modify Transaction Parameters of Active Transaction
type Assign struct {
	BaseObj
	// Parameters for modification
	parameters []Parameter
}

// NewAssign creates new Assign.
// name - name of object
// parameters - parameters for assign.
// Example:
// Parameter{name: "param1_name", value: param1_value},
// Parameter{name: "param2_name", value: param2_value} ...
func NewAssign(name string, parameters ...Parameter) *Assign {
	obj := &Assign{parameters: parameters}
	obj.name = name
	return obj
}

// AppendTransact append transact to object
func (obj *Assign) AppendTransact(transact *Transaction) bool {
	obj.BaseObj.AppendTransact(transact)
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			transact.SetParameters(obj.parameters)
			return true
		}
	}
	return false
}

// Report - print report about object
func (obj *Assign) Report() {}
