// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

type ITransaction interface {
	GetId() int                      // Get transact ID
	GetLife() int                    // Get transact time of life, rip - born
	SetTiсks(interval int)           // Set advance ticks
	DecTiсks()                       // Decrement ticks
	GetTicks() int                   // Get current value of ticks
	IsTheEnd() bool                  // Is ticks value equal zero?
	SetHolderName(holderName string) // Set holder of transact
	GetHolderName() string           // Get current holder of transact
	InqQueueTime()                   // Increment time in queue
	GetQueueTime() int               // Get current value of time in queue
	ResetQueueTime()                 // Reset time in queue
	GetAdvanceTime() int             // Get full time in advice state
	Kill()                           // Kill transact
	IsKilled() bool                  // Is transact killed?
	GetPipeline() IPipeline          // Get pipeline for object
	SetParts(part, parts int)        // Set parts info
	GetParts() (int, int)            // Get parts info
	PrintInfo()                      // Print info about transact
	Copy() ITransaction              // Create copy of transact
}

// Struct for splitting
type Parts struct {
	part  int // Part id
	parts int // Number of parts
}

type Transaction struct {
	id         int       // Transact ID
	born       int       // Moment of borning
	rip        int       // Kill moment
	advance    int       // Full time in advice state
	ticks      int       // Tiks for change state
	holderName string    // Holder object name
	timequeue  int       // Time in queue at this moment
	pipe       IPipeline // Pipeline
	parts      Parts     /* For splitting. Default is "0/0". After splitting
	may be "1/6" - first part of six parts or "5/6" - fifth part of six parts */
}

func NewTransaction(id int, pipe IPipeline) *Transaction {
	t := &Transaction{}
	t.id = id
	t.pipe = pipe
	t.born = pipe.GetModelTime()
	t.parts = Parts{0, 0}
	return t
}

func (t *Transaction) Copy() ITransaction {
	copy_t := &Transaction{}
	copy_t.id = t.id
	copy_t.pipe = t.pipe
	copy_t.born = t.born
	copy_t.parts = t.parts
	return copy_t
}

func (t *Transaction) GetId() int {
	return t.id
}

func (t *Transaction) GetLife() int {
	return t.rip - t.born
}

func (t *Transaction) PrintInfo() {
	trace := t.GetPipeline().GetLogger().GetTrace()
	trace.Println("Transaction Id:\t", t.GetId(), "Borned:\t", t.born,
		"Advance time:\t", t.advance, "Transaction life:\t",
		t.GetPipeline().GetModelTime()-t.born, "Holder Name:\t", t.holderName,
		"Tiks:\t\t", t.ticks, "Time in queue:\t", t.timequeue)
}

// Set ticks and increases advance value to same value.
func (t *Transaction) SetTiсks(interval int) {
	t.ticks = interval
	t.advance += interval
}

func (t *Transaction) InqQueueTime() {
	t.timequeue++
}

func (t *Transaction) GetTicks() int {
	return t.ticks
}

func (t *Transaction) IsTheEnd() bool {
	return bool(t.ticks == 0)
}

func (t *Transaction) SetHolderName(holderName string) {
	t.holderName = holderName
}

func (t *Transaction) GetHolderName() string {
	return t.holderName
}

// Decremet ticks. If ticks is less than zero, set ticks value to zero.
func (t *Transaction) DecTiсks() {
	t.ticks--
	if t.ticks < 0 {
		t.ticks = 0
	}
}

func (t *Transaction) Kill() {
	t.rip = t.GetPipeline().GetModelTime()
}

func (t *Transaction) IsKilled() bool {
	return bool(t.rip != 0)
}

func (t *Transaction) GetQueueTime() int {
	return t.timequeue
}

func (t *Transaction) GetAdvanceTime() int {
	return t.advance
}

func (t *Transaction) GetPipeline() IPipeline {
	return t.pipe
}

func (t *Transaction) ResetQueueTime() {
	t.timequeue = 0
}

func (t *Transaction) GetParts() (int, int) {
	return t.parts.part, t.parts.parts
}

func (t *Transaction) SetParts(part, parts int) {
	t.parts = Parts{part, parts}
}
