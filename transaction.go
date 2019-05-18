// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

type ITransaction interface {
	SetID(int)                                   // Set transact ID
	GetId() int                                  // Get transact ID
	GetLife() int                                // Get transact time of life, rip - born
	SetTiсks(interval int)                       // Set advance ticks
	DecTiсks()                                   // Decrement ticks
	GetTicks() int                               // Get current value of ticks
	IsTheEnd() bool                              // Is ticks value equal zero?
	SetHolderName(holderName string)             // Set holder of transact
	GetHolderName() string                       // Get current holder of transact
	InqQueueTime()                               // Increment time in queue
	GetQueueTime() int                           // Get current value of time in queue
	ResetQueueTime()                             // Reset time in queue
	GetAdvanceTime() int                         // Get full time in advice state
	Kill()                                       // Kill transact
	IsKilled() bool                              // Is transact killed?
	GetPipeline() IPipeline                      // Get pipeline for object
	SetParts(part, parts, parent_id int)         // Set parts info
	GetParts() (int, int, int)                   // Get parts info
	SetParameters(parameters []Parameter)        // Set parameters to transuct
	GetParameters() map[string]interface{}       // Get all parameters of trunsact
	SetParameter(name string, value interface{}) // Set value of parameter
	GetParameter(name string) interface{}        // Get parameter of trunsact by name
	PrintInfo()                                  // Print info about transact
	Copy() ITransaction                          // Create copy of transact
}

// Struct for splitting
type Parts struct {
	part      int // Part id
	parts     int // Number of parts
	parent_id int // ID of parent transaction, for splitting
}

type Transaction struct {
	pipe  IPipeline // Pipeline
	parts Parts     /* For splitting. Default is "0/0". After splitting
	may be "1/6" - first part of six parts or "5/6" - fifth part of six parts */
	parameters map[string]interface{} // Parameters of transaction
}

func NewTransaction(pipe IPipeline) ITransaction {
	t := &Transaction{}
	t.parameters = make(map[string]interface{})
	t.SetParameters([]Parameter{
		{Name: "id", Value: pipe.GetIDNewTransaction()}, // Transact ID
		{Name: "born", Value: pipe.GetModelTime()},      // Moment of borning
		{Name: "advance", Value: 0},                     // Full time in advice state
		{Name: "timequeue", Value: 0},                   // Time in queue at this moment
		{Name: "ticks", Value: 0},                       // Tiks for change state
		{Name: "rip", Value: 0},                         // Kill moment
		{Name: "holderName", Value: ""},                 // Holder object name
	})
	t.pipe = pipe

	t.parts = Parts{0, 0, 0}
	return t
}

func (t *Transaction) Copy() ITransaction {
	copy_t := &Transaction{}
	copy_t.pipe = t.pipe
	copy_t.parts = t.parts
	copy_t.parameters = make(map[string]interface{})
	for key, value := range t.parameters {
		copy_t.parameters[key] = value
	}
	return copy_t
}

func (t *Transaction) SetID(id int) {
	t.SetParameter("id", id)
}

func (t *Transaction) GetId() int {
	return t.parameters["id"].(int)
}

func (t *Transaction) GetLife() int {
	return t.GetIntParameter("rip") - t.GetIntParameter("born")
}

func (t *Transaction) PrintInfo() {
	trace := t.GetPipeline().GetLogger().GetTrace()
	trace.Println("Transaction Id:\t", t.GetId(),
		"Borned:\t", t.GetIntParameter("born"),
		"Advance time:\t", t.GetIntParameter("advance"),
		"Transaction life:\t", t.GetPipeline().GetModelTime()-t.GetIntParameter("born"),
		"Holder Name:\t", t.GetStringParameter("holderName"),
		"Tiks:\t\t", t.GetIntParameter("ticks"),
		"Time in queue:\t", t.GetIntParameter("timequeue"))
}

// Set ticks and increases advance value to same value.
func (t *Transaction) SetTiсks(interval int) {
	t.SetParameter("ticks", interval)
	t.SetParameter("advance", t.GetIntParameter("advance")+interval)
}

func (t *Transaction) InqQueueTime() {
	t.SetParameter("timequeue", t.GetIntParameter("timequeue")+1)
	t.SetParameter("advance", t.GetIntParameter("advance")+1)
}

func (t *Transaction) GetTicks() int {
	return t.GetIntParameter("ticks")
}

func (t *Transaction) IsTheEnd() bool {
	return bool(t.GetIntParameter("ticks") == 0)
}

func (t *Transaction) SetHolderName(holderName string) {
	t.SetParameter("holderName", holderName)
}

func (t *Transaction) GetHolderName() string {
	return t.GetParameter("holderName").(string)
}

// Decremet ticks. If ticks is less than zero, set ticks value to zero.
func (t *Transaction) DecTiсks() {
	ticks := t.GetIntParameter("ticks")
	ticks--
	if ticks < 0 {
		ticks = 0
	}
	t.SetParameter("ticks", ticks)
}

func (t *Transaction) Kill() {
	t.SetParameter("rip", t.GetPipeline().GetModelTime())
}

func (t *Transaction) IsKilled() bool {
	return bool(t.GetIntParameter("rip") != 0)
}

func (t *Transaction) GetQueueTime() int {
	return t.GetIntParameter("timequeue")
}

func (t *Transaction) GetAdvanceTime() int {
	return t.GetIntParameter("advance")
}

func (t *Transaction) GetPipeline() IPipeline {
	return t.pipe
}

func (t *Transaction) ResetQueueTime() {
	t.SetParameter("timequeue", 0)
}

func (t *Transaction) GetParts() (int, int, int) {
	return t.parts.part, t.parts.parts, t.parts.parent_id
}

func (t *Transaction) SetParts(part, parts, parent_id int) {
	t.parts = Parts{part, parts, parent_id}
}

func (t *Transaction) SetParameters(parameters []Parameter) {
	for _, v := range parameters {
		if v.Value != nil {
			t.parameters[v.Name] = v.Value
		} else {
			delete(t.parameters, v.Name)
		}
	}
}

func (t *Transaction) GetParameters() map[string]interface{} {
	return t.parameters
}

func (t Transaction) SetParameter(name string, value interface{}) {
	t.parameters[name] = value
}

func (t *Transaction) GetParameter(name string) interface{} {
	return t.parameters[name]
}

func (t *Transaction) GetIntParameter(name string) int {
	return t.GetParameter(name).(int)
}

func (t *Transaction) GetStringParameter(name string) string {
	return t.GetParameter(name).(string)
}
