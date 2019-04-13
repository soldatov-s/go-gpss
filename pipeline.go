// pipeline
package gpss

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
)

type IPipeline interface {
	Append(obj IBaseObj, src ...IBaseObj) // Append  object to pipeline
	Delete(obj IBaseObj)                  // Delete object from pipeline
	Start(value int)                      // Start simulation
	Stop()                                // Stop simulation
	IsVerbose() bool                      // Show verbose info
	GetSimTime() int                      // Get Simulation time
	GetModelTime() int                    // Get current model time
	PrintReport()                         // Print report
}

type Pipeline struct {
	name      string              // Pipeline name
	objects   map[string]IBaseObj // Maps of objects
	modelTime int                 // Current Model Time
	Done      chan struct{}       // Chan for done
	simTime   int                 // Simulation time
	verbose   bool
}

// Create new Pipeline
func NewPipeline(name string, verbose bool) *Pipeline {
	p := &Pipeline{}
	p.objects = make(map[string]IBaseObj)
	p.name = name
	p.Done = make(chan struct{})
	p.modelTime = 0
	p.verbose = verbose
	return p
}

// Append object to pipeline. Src is multiple sources of transact for appended
// object.
func (p *Pipeline) Append(obj IBaseObj, dst ...IBaseObj) {
	obj.SetDst(dst)
	obj.SetPipeline(p)
	p.objects[obj.GetName()] = obj
}

// Delete object from pipeline
func (p *Pipeline) Delete(obj IBaseObj) {
	o := obj.(IBaseObj)
	delete(p.objects, o.GetName())
}

// Print list of objects ib pipeline
func (p *Pipeline) PrintObjects() {
	keys := make([]string, 0, len(p.objects))
	for k := range p.objects {
		keys = append(keys, k)
	}
	fmt.Println("Pipeline ", p.name)
	for _, k := range keys {
		fmt.Println("Key:", k, "Value:", reflect.TypeOf(p.objects[k]))
	}
}

// Start simulation
func (p *Pipeline) Start(value int) {
	var wg sync.WaitGroup

	p.simTime = value
	go func() {
		for {
			select {
			case <-p.Done:
				return
			default:
				if p.verbose {
					fmt.Println("\nModelTime ", p.modelTime)
				}
				wg.Add(len(p.objects))
				for _, o := range p.objects {
					o.HandleTransacts(&wg)
				}
				wg.Wait()
				if p.modelTime++; p.modelTime == value {
					p.Stop()
				}
			}
		}
	}()
}

// Stop simulation
func (p *Pipeline) Stop() {
	close(p.Done)
}

func (p *Pipeline) IsVerbose() bool {
	return p.verbose
}

func (p *Pipeline) PrintReport() {
	fmt.Println("Pipeline name \"", p.name, "\"")
	fmt.Println("Simulation time", p.modelTime)
	sortedKeys := make([]string, 0, len(p.objects))
	for k := range p.objects {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		p.objects[k].PrintReport()
	}
}

func (p *Pipeline) GetSimTime() int {
	return p.simTime
}

func (p *Pipeline) GetModelTime() int {
	return p.modelTime
}
