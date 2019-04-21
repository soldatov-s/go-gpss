// transaction
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
	PrintInfo()                      // Print info about transact
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
}

func NewTransaction(id int, pipe IPipeline) *Transaction {
	t := &Transaction{}
	t.id = id
	t.pipe = pipe
	t.born = pipe.GetModelTime()
	return t
}

func (t *Transaction) GetId() int {
	return t.id
}

func (t *Transaction) GetLife() int {
	return t.rip - t.born
}

func (t *Transaction) PrintInfo() {
	verbose := t.GetPipeline().IsVerbose()
	PrintlnVerbose(verbose, "Transaction Id:\t", t.GetId())
	PrintlnVerbose(verbose, "Borned:\t\t", t.born)
	PrintlnVerbose(verbose, "Advance time:\t", t.advance)
	PrintlnVerbose(verbose, "Transaction life:\t", t.GetPipeline().GetModelTime()-t.born)
	PrintlnVerbose(verbose, "Holder Name:\t", t.holderName)
	PrintlnVerbose(verbose, "Tiks:\t\t", t.ticks)
	PrintlnVerbose(verbose, "Time in queue:\t", t.timequeue)
	PrintlnVerbose(verbose)
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
