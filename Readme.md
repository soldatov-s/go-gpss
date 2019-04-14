[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/soldatov-s/go-gpss)
[![Build Status](https://travis-ci.org/soldatov-s/go-gpss.svg?branch=master)](https://travis-ci.org/soldatov-s/go-gpss)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
  
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

# Example 1
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

# Example 2
Office with two WC and 10 employees. Employees go to the WC every 90 minutes 
with a deviation of 60 minutes. The way to the WC lasts 5 minutes with a 
deviation of 3 minutes. A worker takes a WC room for 15 minutes with a deviation 
of 10 minutes. The way from the WC lasts 5 minutes with a deviation of 3 minutes.
How many people will be in queue? How long will the worker wait for the toilet 
to be empty? How much time will each toilet be occupied?

```Golang
p := NewPipeline("Water Closet Simulation", false)
g := NewGenerator("Office", 0, 0, 0, 10, nil)
a1 := NewAdvance("Wanted to use the toilet", 90, 60)
a2 := NewAdvance("Path to WC", 5, 3)
q := NewQueue("Queue to the WC")
f1 := NewFacility("WC1", 15, 10)
f2 := NewFacility("WC2", 15, 10)
a3 := NewAdvance("Path from WC", 5, 3)
p.Append(g, a1)
p.Append(a1, a2)
p.Append(a2, q)
p.Append(q, f1, f2)
p.Append(f1, a3)
p.Append(f2, a3)
p.Append(a3, a1)
p.Start(540)
<-p.Done
p.PrintReport()
```

In report we will see next information (may be diferent values, becouse 
timing was randomized):

```bash
Pipeline name " Water Closet Simulation "
Simulation time 540
Object name " Office "
Generated 10

Object name " Path from WC "
Average advance 5.21

Object name " Path to WC "
Average advance 4.83

Object name " Queue to the WC "
Max content 3
Total entries 45
Zero entries 26
Persent zero entries 57.78%
Current contents 0
Average content 0.30
Average time/trans 3.98
Average time/trans without zero entries 9.42

Object name " WC1 "
Average advance 14.79
Average utilization 65.74
Number entries 24.00
Transact 5 in facility

Object name " WC2 "
Average advance 13.24
Average utilization 51.48
Number entries 21.00
Transact 8 in facility

Object name " Wanted to use the toilet "
Average advance 86.13
```

Max count in queue 3 employes. Average waiting time in queue 9.42 minutes.
WC1 was occupied during 14.79 minutes. WC2 was occupied during 13.24 minutes.

# TODO
- Extend list of blocks
- Extend report
- Clear code
- Make tests
