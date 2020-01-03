[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/soldatov-s/go-gpss)
[![Build Status](https://travis-ci.org/soldatov-s/go-gpss.svg?branch=master)](https://travis-ci.org/soldatov-s/go-gpss)
[![Coverage Status](http://codecov.io/github/soldatov-s/go-gpss/coverage.svg?branch=master)](http://codecov.io/github/soldatov-s/go-gpss?branch=master)
[![codebeat badge](https://codebeat.co/badges/c344eca5-937c-4f96-9d7e-f518d6d2e4e5)](https://codebeat.co/projects/github-com-soldatov-s-go-gpss-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/soldatov-s/go-gpss)](https://goreportcard.com/report/github.com/soldatov-s/go-gpss)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
  
Readme
========

go-gpss is a framework based on conception GPSS (General Purpose Simulation System).
It help quick and easy build model for simulation modeling
 
It include today few blocks:
- Generator - sequentially generates Transactions
- Advance - delays the progress of a Transaction for a specified amount of simulated time
- Queue - Queue of Transactions
- Facility - facility entity with Advance in it
- Bifacility - as Facility, but without Advance in it, it present in two parts, first for takes ownership of a Facility, second for release ownership of a Facility
- Split - creates assembly set of sub-transactions of a Transaction
- Aggregate - aggregate multiple sub-transactions in Transaction
- Check - compares parameters of Transaction or any another parameters of simulation model, and controls the destination of the Active Transaction based on the result of the comparison
- Assign - modify Transaction Parameters of Active Transaction 
- Count - counts all Transactions which pass through the block, it present in two parts, first for increment Count value, second for decrement Count value
- Hole - Hole in which fall in Transactions

Active Transaction is a Transaction in current block.  
All blocks need to add in Pipeline and than start simulation. For generate random 
values used pseudo-random generation function from math/rand. After simulation 
you can print report about simulation.

# The difference between version 0.2 and 0.1
The new version supports a simpler construction of a simulation pipeline.
You can add objects to the pipeline one by one.

Old code
```Golang
p := NewPipeline("Barbershop", false)
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

New code
```Golang
p := objects.NewPipeline("Barbershop").
	AddObject(objects.NewGenerator("Clients", 18, 6, 0, 0, nil)).
	AddObject(objects.NewQueue("Chairs")).
	AddObject(objects.NewFacility("Master", 16, 4)).
	AddObject(objects.NewHole("Out"))
p.Start(480)
<-p.Done
p.Report()
```

You can link objects, for example, link barista-facility to barista-queue:

```Golang
baristaQ.LinkObject(baristaF)
```
That is mean baristaQ->baristaF

# Example 1.1
Barbershop: random client go to Barbershop every 18 minutes with deviation 6 minutes.
We have only one barber. Barber spends for each client 16 minutes with deviation
4 minutes. How many people will get a haircut per day? How many people will be in queue?
How long does one haircut last?

<p align="center">
  <img src="/images/pic01.jpg" width="300" height="400" alt="Pic01"/>
  <br /> 
  <b>Pic 01 - Barbershop simulation</b>
</p>

```Golang
// Build pipeline
// Generator -> Queue -> Facility -> Hole
p := objects.NewPipeline("Barbershop").
	AddObject(objects.NewGenerator("Clients", 18, 6, 0, 0, nil)).
	AddObject(objects.NewQueue("Chairs")).
	AddObject(objects.NewFacility("Master", 16, 4)).
	AddObject(objects.NewHole("Out"))
// Start simulation
p.Start(480)

<-p.Done
p.Report()
```

Full source [example1](examples/example1.1/main.go).  

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

Full source [example1.1](examples/example1.1/main.go).

# Example 1.2
Same as Example 1.1, but used bifacility. Bifacility is a component that inlcude
two parts - in_elemet and out_element. Between theise elements we can insert 
any anower component, for example advance. This allow build more complex models.
```Golang
// Build pipeline
// Generator -> Queue -> ...
p := objects.NewPipeline("Barbershop", doneHandler).
	AddObject(objects.NewGenerator("Clients", 18, 6, 0, 0, nil)).
	AddObject(objects.NewQueue("Chairs"))

// Create BiFacility
fIN, fOUT := objects.NewBifacility("Master")

// Build pipeline
// ... -> BiFacilityIn -> Advance -> BeFacilityOut -> Hole
p.
	AddObject(fIN).
	AddObject(objects.NewAdvance("Master work", 16, 4)).
	AddObject(fOUT).
	AddObject(objects.NewHole("Out"))
// Start simulation
p.Start(480)
<-p.Done
p.Report() 
```

Full source [example1.2](examples/example1.2/main.go).  

**Important**, 
The advance and facility components count the average for all transactions that 
entered into it. Bifacility component count the average value ​​only for 
transactions that are entered _and_ exited from the component.

# Example 2
Office with two WC and 10 employees. Employees go to the WC every 90 minutes 
with a deviation of 60 minutes. The way to the WC lasts 5 minutes with a 
deviation of 3 minutes. A worker takes a WC room for 15 minutes with a deviation 
of 10 minutes. The way from the WC lasts 5 minutes with a deviation of 3 minutes.
How many people will be in queue? How long will the worker wait for the toilet 
to be empty? How much time will each toilet be occupied?

<p align="center">
  <img src="/images/pic02.jpg" width="400" height="650" alt="Pic03"/>
  <br /> 
  <b>Pic 02 - WC simulation</b>
</p>

```Golang
p := objects.NewPipeline("Water Closet Simulation", doneHandler).
	AddObject(objects.NewGenerator("Office", 0, 0, 0, 10, nil)).
	AddObject(objects.NewAdvance("Wanted to use the toilet", 90, 60)).
	AddObject(objects.NewAdvance("Path to WC", 5, 3)).
	AddObject(objects.NewQueue("Queue to the WC")).
	AddObject(objects.NewFacility("WC1", 15, 10), objects.NewFacility("WC2", 15, 10)).
	AddObject(objects.NewAdvance("Path from WC", 5, 3)).
	Loop("Wanted to use the toilet")
p.Start(540)
<-p.Done
p.Report()
```

Full source [example2](examples/example2/main.go).  
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

# Example 3
Small cafeteria with one barista and cook. Random client go to cafeteria every 
18 minutes with deviation 6 minutes. Cashier spends for each client 5 minutes 
with deviation 3 minutes. Сlient can request coffee/tea and/or burger/piece of 
cake (one or two positions in order). Barista spends 5 minutes with deviation 2 
minutes for one order. Cook spends 10 minutes with deviation 5 minutes for one 
order.
How many people can be served in a cafe? How many people will be in queue?

<p align="center">
  <img src="/images/pic03.jpg" width="400" height="650" alt="Pic03"/>
  <br /> 
  <b>Pic 03 - Cafe simulation</b>
</p>

```Golang
p := objects.NewPipeline("Cafe Simulation", doneHandler).
	AddObject(objects.NewGenerator("Visitors", 18, 6, 0, 0, nil)).
	AddObject(objects.NewQueue("Visitors queue")).
	AddObject(objects.NewFacility("Order Acceptance", 5, 3)).
	AddObject(objects.NewSplit("Split orders", 1, 1, nil))
baristaQ := objects.NewQueue("Queue of orders to barista")
cookQ := objects.NewQueue("Queue of orders to cook")
p.AddObject(baristaQ, cookQ)
baristaF := objects.NewFacility("Barista", 5, 2)
cookF := objects.NewFacility("Cook", 10, 5)
baristaQ.LinkObject(baristaF)
cookQ.LinkObject(cookF)
aggregate := objects.NewAggregate("Aggregate orders")
baristaF.LinkObject(aggregate)
cookF.LinkObject(aggregate)
aggregate.LinkObject(objects.NewHole("Out"))
p.Start(480)
<-p.Done
p.Report()
```

Full source [example3](examples/example3/main.go).  
In report we will see next information (may be diferent values, becouse 
timing was randomized):

```bash
Pipeline name " Cafe Simulation "
Simulation time 480
Object name " Visitors "
Generated 26

Object name " Visitors queue "
Max content 0
Total entries 26
Zero entries 26
Persent zero entries 100.00%
Current contents 0
Average content 0.00
Average time/trans 0.00

Object name " Order Acceptance "
Average advance 4.81
Average utilization 26.04%
Number entries 26.00
Transact 26 in facility

Object name " Split orders "
Average split 1.48

Object name " Queue of orders to barista "
Max content 0
Total entries 23
Zero entries 23
Persent zero entries 100.00%
Current contents 0
Average content 0.00
Average time/trans 0.00

Object name " Queue of orders to cook "
Max content 0
Total entries 14
Zero entries 14
Persent zero entries 100.00%
Current contents 0
Average content 0.00
Average time/trans 0.00

Object name " Barista "
Average advance 5.48
Average utilization 26.25%
Number entries 23.00
Facility is empty

Object name " Cook "
Average advance 8.71
Average utilization 25.42%
Number entries 14.00
Facility is empty

Object name " Aggregate orders "
Number of aggregated transact 25.00

Object name " Out "
Killed 25
Average advance 12.12
Average life 12.72
```

We have served 24 client. 25 served at the end of the simulation. No clients in 
queue. Barista busy at 14.58 percent, cook busy at 32.71 percent.

# Example 4
Simulation of Restaurant. We have restaurant with 24 tables and staff: 
8 waiters, 2 hostes, 4 cooks and 2 barmans. Random visitors enter to restaurant 
every 10 minutes with deviation 5 minutes. If queue to restaurant contains 
more than 6 people, visitors leave restaurant without waiting free tables. 
Hostess spend for each visitor 5 minutes with deviation 3 minutes. Waiters spend 
for each visitor 5 minutes with deviation 3 minutes (both for take orders and giving dishes).
Visitor can request 4 dishes and one drink. Barman spend 4 for one drink minutes 
with deviation 2 minutes for one drink. Cooks spend 7...15 minutes with deviation 
3...5 minutes for one dish. Visitor eats one dish for 45 minutes with deviation 
10 minutes. And visitors spend for payment 5 minutes with deviation 2 minutes.

How many people can to serve a restaurant?
How many empty tables in restaurant?
Are there many or few staff in the restaurant?
This example will be use Assign and Check blocks.  
<p align="center">
  <img src="/images/pic04.jpg" width="400" height="650" alt="Pic03"/>
  <br /> 
  <b>Pic 04 - Restaurant simulation</b>
</p>

Full source [example4](examples/example4/main.go).  

On the basis of the [report](examples/example4/report.txt) the following conclusions can be made (there may be different values, because the time was randomized):
- two tables are not used (23 and 24) and a quarter of the tables have very low utilization
- the restaurant served 29 visitors and none left without visiting to the restaurant
- all visitors did not wait in queue
- at the end of the simulation, 12 visitors have already received some of the dishes and are waiting for the remaining
- a cook 1 and a cook 4 have a very large utilization (91.46%, 88.33%)
- a barman 2 has a very low utilization (1.67%)
- half of the waiters have a very small utilization
- a hostess 2 has a very low utilization (9.38%)

We can say that our restaurant very big and have many staff. Or we opened a restaurant in a bad place, and there are too few visitors.  
I tested model when clients enter to restaurant every 5 minutes with deviation 3 minutes and utilization rate began to grow.
If it necessary developer can change how will be splits orders between cooks or add an order split between waiters for uniform distribution of orders.

# Fixes
- Fixed report, ordered by id in Pipeline
- Fixed HoldedTransactID in facility, zeroing after removing transact
- Fixed Generator, Modificator is used now
- Fixed examples after linter

# TODO
- Extend list of blocks
- Extend report
- Clear code
- Make tests

# Features of go-gpss
- pipeline as block for pipline
- generate schem based on pipeline
- generate report in the form of a table
- visualizate simulation
