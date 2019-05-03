[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/soldatov-s/go-gpss)
[![Build Status](https://travis-ci.org/soldatov-s/go-gpss.svg?branch=master)](https://travis-ci.org/soldatov-s/go-gpss)
[![Coverage Status](http://codecov.io/github/soldatov-s/go-gpss/coverage.svg?branch=master)](http://codecov.io/github/soldatov-s/go-gpss?branch=master)
[![codebeat badge](https://codebeat.co/badges/c344eca5-937c-4f96-9d7e-f518d6d2e4e5)](https://codebeat.co/projects/github-com-soldatov-s-go-gpss-master)
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

# Example 1
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

Full source [example1](examples/example1/main.go).  

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

# Example 1.1
Same as Example 1, but used bifacility. Bifacility is a component that inlcude
two parts - in_elemet and out_element. Between theise elements we can insert 
any anower component, for example advance. This allow build more complex models.
```Golang
p := NewPipeline("Barbershop", false)
g := NewGenerator("Clients", 18, 6, 0, 0, nil)
q := NewQueue("Chairs")
f_in, f_out := NewBifacility("Master")
a := NewAdvance("Master work", 16, 4)
h := NewHole("Out")
p.Append(g, q)
p.Append(q, f_in)
p.Append(f_in, a)
p.Append(a, f_out)
p.Append(f_out, h)
p.Append(h)
p.Start(480)
<-p.Done
p.PrintReport() 
```

Full source [example1](examples/example1/main.go).  

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
p := NewPipeline("Cafe Simulation", false)
g := NewGenerator("Visitors", 18, 6, 0, 0, nil)
q := NewQueue("Visitors queue")
orders_f := NewFacility("Order Acceptance", 5, 3)
split := NewSplit("Split orders", 1, 1, nil)
barista_q := NewQueue("Queue of orders to barista")
barista_f := NewFacility("Barista", 5, 2)
cook_q := NewQueue("Queue of orders to cook")
cook_f := NewFacility("Cook", 10, 5)
aggregate := NewAggregate("Aggregate orders")
h := NewHole("Out")
p.Append(g, q)
p.Append(q, orders_f)
p.Append(orders_f, split)
p.Append(split, barista_q, cook_q)
p.Append(barista_q, barista_f)
p.Append(cook_q, cook_f)
p.Append(barista_f, aggregate)
p.Append(cook_f, aggregate)
p.Append(aggregate, h)
p.Append(h)
p.Start(480)
<-p.Done
p.PrintReport()
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
Simulation of Restaurant. We have restaurant with 24 tables, 8 waiters,
2 hostes, 4 cooks and 2 barmans. Random client go to restaurant every 
10 minutes with deviation 5 minutes. If queue to restaurant more than 6, clients 
go out. Hostes spends for each client 5 minutes with deviation 3 minutes. 
Waiters spends for client for each client 5 minutes with deviation 
3 minutes (both for take orders and giving dishes).
Сlient can request 4 dishes and one drink. Barman spends 4 minutes with deviation 2 
minutes for one order. Cooks spends 7..15 minutes with deviation 3..5 minutes for one 
order. Clients spends for eating 45 minutes with deviation 10 minutes. And
clients spends for payment 5 minutes with deviation 2 minutes.
How many people can be served in a restaurant? What will be the employment of 
the staff?
This example will be use Assign and Check blocks.  
<p align="center">
  <a href="/images/pic04.jpg">
  <img src="/images/pic04.jpg" width="400" height="650" alt="Pic03"/>
  </a>
  <br /> 
  <b>Pic 04 - Restaurant simulation</b>
</p>
Full source [example4](examples/example4/main.go).

# Fixes
- Fixed report, ordered by id in Pipeline
- Fixed HoldedTransactID in facility, zeroing after removing transact
- Fixed Generator, Modificator is used now

# TODO
- Extend list of blocks
- Extend report
- Clear code
- Make tests
