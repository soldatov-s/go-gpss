[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/soldatov-s/go-gpss)
[![Build Status](https://travis-ci.org/soldatov-s/go-gpss.svg?branch=master)](https://travis-ci.org/soldatov-s/go-gpss)
  
Readme
========

go-gpss is a framework based on conception GPSS (General Purpose Simulation System).
It help quick and easy build model for simulation modeling
 
It include today few blocks:
- Generator - Generator transaction
- Advance - Advance block used for simulation waiting/holding time
- Queue - Queue of transaction
- Facility - Any facility
- Hole - Hole in which fall in transactions
All blocks need to add in Pipeline and than start simulation.
For generate random values used pseudo-random generation function from math/rand.
After simulation you can print report about simulation.

# Example
Barbershop: random client go to Barbershop every 18 minutes with deviation 6 minutes.
We have only one barber. Barber spends for each client 16 minutes with deviation
4 minutes. How many people will get a haircut per day? How many people will be in queue?
How long does one haircut last?

```Golang
p := NewPipeline("Barbershop", true)
g := NewGenerator("Clients", 18, 6, 0, 0, nil)
q := NewQueue("Chairs")
f := NewFacility("Master", 16, 4)
h := NewHole("Out")
p.Append(g, q)
p.Append(q, f)
p.Append(f, h)
p.Append(h)
p.Start(480)
<-p.Done
p.PrintReport()
```
In report we will see next information (may be diferent values, becouse 
timing was randomized):

```bash
Pipeline name " Barbershop "
Simulation time 480
Object name " Chairs "
Max content 1
Total entries 26
Zero entries 11
Persent zero entries 42.31%
In queue 0
Average time/trans 2.58
Average time/trans without zero entries 4.47

Object name " Clients "
Generated 26

Object name " Master "
Average advance 16.46
Average utilization 89.17
Number entries 26.00
Transact 26 in facility

Object name " Out "
Killed 25
Average advance 16.56
Average life 19.44
```
We have served 25 client. 26 served at the end of the simulation. No more 1 
client in queue. 11 client not waiting in queue (42.31 percent), but are 
immediately served. Average waiting time in queue 4.47 minutes. Barber busy at 
89.17 percent.

# TODO
- Extend list of blocks
- Extend report
- Clear code
- Make tests
