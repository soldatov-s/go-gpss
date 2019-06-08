// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

// Parameter for modification
type Parameter struct {
	Name  string      // Name of parameter
	Value interface{} // Value of parameter
}

// Transaction struct for description of transaction
type Transaction struct {
	pipe       *Pipeline              // Pipeline
	parameters map[string]interface{} // Parameters of transaction
}

// NewTransaction create new transaction
func NewTransaction(pipe *Pipeline) *Transaction {
	t := &Transaction{
		pipe:       pipe,
		parameters: make(map[string]interface{}),
	}
	// Transact ID
	t.parameters["id"] = pipe.NewID()
	// Moment of borning
	t.parameters["born"] = pipe.ModelTime
	// Full time in advice state
	t.parameters["advance"] = 0
	// Time in queue at this moment
	t.parameters["timequeue"] = 0
	// Tiks for change state
	t.parameters["ticks"] = 0
	// Kill moment
	t.parameters["rip"] = 0
	// Holder object name
	t.parameters["holder"] = ""
	// Part id, for splitting. Default is "0/0". After splitting
	// may be "1/6" - first part of six parts or "5/6" - fifth part
	// of six parts
	t.parameters["part"] = 0
	// Number of parts
	t.parameters["parts"] = 0
	// ID of parent transaction, for splitting
	t.parameters["parent_id"] = 0
	return t
}

// Copy create copy of transact
func (t *Transaction) Copy() *Transaction {
	copyTr := &Transaction{}
	copyTr.pipe = t.pipe
	copyTr.parameters = make(map[string]interface{})
	for key, value := range t.parameters {
		copyTr.parameters[key] = value
	}
	return copyTr
}

// SetID set transact ID
func (t *Transaction) SetID(id int) {
	t.SetParameter("id", id)
}

// Get transact ID
func (t *Transaction) GetID() int {
	return t.GetIntParameter("id")
}

// GetLife get transact time of life, rip - born
func (t *Transaction) GetLife() int {
	return t.GetIntParameter("rip") - t.GetIntParameter("born")
}

// PrintInfo - print info about transact
func (t *Transaction) PrintInfo() {
	Log.Trace.Println("Transaction ID:\t", t.GetID(),
		"Borned:\t", t.GetIntParameter("born"),
		"Advance time:\t", t.GetAdvanceTime(),
		"Holder Name:\t", t.GetHolder(),
		"Tiks:\t\t", t.GetTicks(),
		"Time in queue:\t", t.GetQueueTime())
}

// SetTiсks - set ticks and increases advance value to same value.
func (t *Transaction) SetTiсks(interval int) {
	t.SetParameter("ticks", interval)
	t.SetParameter("advance", t.GetAdvanceTime()+interval)
}

// InqQueueTime - increment time in queue
func (t *Transaction) InqQueueTime() {
	t.SetParameter("timequeue", t.GetQueueTime()+1)
	t.SetParameter("advance", t.GetAdvanceTime()+1)
}

// GetTicks - get current value of ticks
func (t *Transaction) GetTicks() int {
	return t.GetIntParameter("ticks")
}

// IsTheEnd - is ticks value equal zero?
func (t *Transaction) IsTheEnd() bool {
	return bool(t.GetIntParameter("ticks") == 0)
}

// SetHolder - set holder of transact
func (t *Transaction) SetHolder(holderName string) {
	t.SetParameter("holder", holderName)
}

// GetHolder - get current holder of transact
func (t *Transaction) GetHolder() string {
	return t.GetStringParameter("holder")
}

// DecTiсks - decremet ticks, if ticks is less than zero, set ticks value to zero.
func (t *Transaction) DecTiсks() {
	ticks := t.GetTicks()
	ticks--
	if ticks < 0 {
		ticks = 0
	}
	t.SetParameter("ticks", ticks)
}

// Kill transact
func (t *Transaction) Kill() {
	t.SetParameter("rip", t.pipe.ModelTime)
}

// IsKilled - is transact killed?
func (t *Transaction) IsKilled() bool {
	return bool(t.GetIntParameter("rip") != 0)
}

// GetQueueTime - get current value of time in queue
func (t *Transaction) GetQueueTime() int {
	return t.GetIntParameter("timequeue")
}

// GetAdvanceTime - get full time in advice state
func (t *Transaction) GetAdvanceTime() int {
	return t.GetIntParameter("advance")
}

// GetPipeline - get pipeline for object
func (t *Transaction) GetPipeline() *Pipeline {
	return t.pipe
}

// ResetQueueTime - reset time in queue
func (t *Transaction) ResetQueueTime() {
	t.SetParameter("timequeue", 0)
}

// GetParts - get parts info
func (t *Transaction) GetParts() (int, int, int) {
	return t.GetIntParameter("part"),
		t.GetIntParameter("parts"),
		t.GetIntParameter("parent_id")
}

// SetParts - set parts info
func (t *Transaction) SetParts(part, parts, parent_id int) {
	t.SetParameters([]Parameter{
		{Name: "part", Value: part},
		{Name: "parts", Value: parts},
		{Name: "parent_id", Value: parent_id},
	})
}

// SetParameters - set parameters to transuct
func (t *Transaction) SetParameters(parameters []Parameter) {
	for _, v := range parameters {
		if v.Value != nil {
			t.parameters[v.Name] = v.Value
		} else {
			delete(t.parameters, v.Name)
		}
	}
}

// GetParameters - get all parameters of transact
func (t *Transaction) GetParameters() map[string]interface{} {
	return t.parameters
}

// SetParameter - set value of parameter
func (t Transaction) SetParameter(name string, value interface{}) {
	t.parameters[name] = value
}

// GetParameter - get parameter of transact by name
func (t *Transaction) GetParameter(name string) interface{} {
	return t.parameters[name]
}

// GetIntParameter - get int parameter of transact by name
func (t *Transaction) GetIntParameter(name string) int {
	return t.GetParameter(name).(int)
}

// GetIntParameter - get string parameter of transact by name
func (t *Transaction) GetStringParameter(name string) string {
	return t.GetParameter(name).(string)
}
