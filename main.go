package main

import (
	"flag"
	"time"
)

const (
	DefaultWaiterCount  = 1
	DefaultBaristaCount = 3

	DefaultCafeName = "Chanbucks"
)

var (
	waiterCount  int
	baristaCount int

	cafeName string
)

func init() {
	flag.IntVar(&waiterCount, "waiter", DefaultWaiterCount, "number of waiters")
	flag.IntVar(&baristaCount, "barista", DefaultBaristaCount, "number of baristas")

	flag.StringVar(&cafeName, "name", DefaultCafeName, "name of the cafe")
}

func run() {
	nameChan := make(chan string, 1)
	numChan := make(chan int, 2)
	cafeChan := make(chan *Cafe, 1)

	go NewCafe(nameChan, numChan, cafeChan)

	numChan <- waiterCount
	numChan <- baristaCount

	nameChan <- cafeName

	// receive cafe
	cafe := <-cafeChan

	cafe.Open()

	// start taking order from customers
	go cafe.TakeOrder()

	time.Sleep(time.Second * 10)

	cafe.Close()
}

func main() {
	flag.Parse()

	run()

}
