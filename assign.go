// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

// Modify Transaction Parameters of Active Transaction
type Assign struct {
	BaseObj
	// Parameters for modification
	parameters []Parameter
}

// Parameter for modification
type Parameter struct {
	Name  string      // Name of parameter
	Value interface{} // Value of parameter
}

// Creates new Assign.
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

func (obj *Assign) AppendTransact(transact ITransaction) bool {
	transact.PrintInfo()
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			transact.SetParameters(obj.parameters)
			obj.GetLogger().GetTrace().Println("Append transact ", transact.GetId(), " to Assign")
			return true
		}
	}
	return false
}

func (obj *Assign) PrintReport() {
	return
}
